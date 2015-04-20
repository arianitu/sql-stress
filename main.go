package main

// TODO: print out table sizes/index sizes at the end
// TODO: print a bunch of metrics

import (
	"database/sql"
	"fmt"
	"flag"
	"runtime"
	"time"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

func worker(db *sql.DB, queries <-chan string, results chan<-int64, done chan<- bool) {
	for query := range queries {
		startTime := time.Now()
		_, err := db.Exec(query)
		elapsed := time.Since(startTime)
		
		if err != nil {
			fmt.Println(err);
			continue
		}
		
		results <- elapsed.Nanoseconds()
	}
}

func sink(results <-chan int64) {
	// avg query time
	for time := range results {
		fmt.Println(time)
	}
}

func main() {
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
	if err != nil {
		fmt.Println(err);
		fmt.Println("Cannot continue, exiting")
		os.Exit(1)
		return
	}

	if (*runFixtures == 1) {
		ProcessFixtures(*fixtureLocation, db)
	}
	
	queries := make(chan string)
	results := make(chan int64)
	done := make(chan bool)
	for i := 0; i < *workers; i++ {
		go worker(db, queries, results, done)
	}
	ProcessTasks(*taskLocation, db)
	queries <- "SELECT NOW()"
	<- done
}

