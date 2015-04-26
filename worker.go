package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Query struct {
	Id     string
	Query  string
	Values []interface{}
	Done   chan<- int64
}

func Worker(db *sql.DB, queryIn <-chan Query) {
	for query := range queryIn {
		startTime := time.Now()
		_, err := db.Exec(query.Query, query.Values...)
		elapsed := time.Since(startTime)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		query.Done <- elapsed.Nanoseconds()
	}
}
