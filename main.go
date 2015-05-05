package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU())

	var workers = flag.Int("workers", 10, "The number of goroutines to execute queries on.")
	var runFixtures = flag.Int("run-fixtures", 1, "If we should run fixtures")
	var fixtureLocation = flag.String("fixture-location", "./fixtures", "The location of fixtures")
	var taskLocation = flag.String("task-location", "./tasks", "The location of tasks")
	var url = flag.String("url", "root:@/sql_stress_test", ` A database url. 
    mysql: username:password@localhost/dbname
    postgres: postgres://username:password@localhost/dbname
    sqlite: /some/location/test.db

`)
	
	var vendor = flag.String("vendor", "mysql", "The sql vendor. Possible values are: mysql, postgres, sqlite")
	flag.Parse()

	fmt.Printf("Running on %v workers \n", *workers)
	db, err := sql.Open(*vendor, *url)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot continue, exiting")
		os.Exit(1)
		return
	}
	
	db.SetMaxOpenConns(runtime.NumCPU())
	queryIn := make(chan Query)
	for i := 0; i < *workers; i++ {
		go Worker(db, queryIn)
	}

	if *runFixtures == 1 {
		ProcessFixtures(*fixtureLocation, db)
	}

	ProcessTasks(*taskLocation, db, queryIn)
	fmt.Println("Done!")
}
