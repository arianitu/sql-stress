package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"runtime"
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

func SpawnWorkers(vendor string, url string, workers int) (chan<- Query, *sql.DB, error) {
	queryIn := make(chan Query)
	db, err := sql.Open(vendor, url)
	if err != nil {
		return nil, nil, err
	}

	db.SetMaxOpenConns(runtime.NumCPU())
	for i := 0; i < workers; i++ {
		go Worker(db, queryIn)
	}
	return queryIn, db, nil
}
