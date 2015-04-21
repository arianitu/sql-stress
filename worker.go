package main

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	"sync"
)

type Query struct {
	Query string
	Values []interface{}
	WaitGroup *sync.WaitGroup
}

func Worker(db *sql.DB, queryIn <-chan Query, sink chan<- int64) {
	for query := range queryIn {
		startTime := time.Now()
		_, err := db.Exec(query.Query, query.Values...)
		elapsed := time.Since(startTime)
		
		if err != nil {
			fmt.Println(err);
			os.Exit(1)
		}
		
		sink <- elapsed.Nanoseconds()
		query.WaitGroup.Done()
	}
}

func Sink(sink <-chan int64) {
	
	for _ = range sink {
		// avg the results
	}
}

