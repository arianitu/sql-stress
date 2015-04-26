package main

// TODO: print out table sizes/index sizes at the end
// TODO: print a bunch of metrics

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	fmt.Printf("Running on %v workers \n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	var url = flag.String("url", "root:@/sql_stress_test", "A database url")
	var workers = flag.Int("workers", 10, "The number of workers to execute queries on.")
	var query = flag.String("query", "", "The query")
	var values = flag.String("values", "", "The values for the query separate by a coma.")
	var iterations = flag.Int("iterations", 1, "The number of times to run Query")
	var runFixtures = flag.Int("run-fixtures", 1, "If we should run fixtures")
	var table = flag.String("table", "", "The table we're running against. We currently only support 1 table at a time")
	var fixtureLocation = flag.String("fixture-location", "./fixtures", "The location of fixtures")
	var taskLocation = flag.String("task-location", "./tasks", "The location of tasks")

	fmt.Println(*query)
	fmt.Println(*values)
	fmt.Println(*iterations)
	fmt.Println(*table)

	flag.Parse()
	db, err := sql.Open("mysql", *url)
	db.SetMaxOpenConns(runtime.NumCPU())

	if err != nil {
		fmt.Println(err)
		fmt.Println("Cannot continue, exiting")
		os.Exit(1)
		return
	}

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
