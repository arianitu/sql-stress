package main

import (
	"database/sql"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"io"
	"path"
	"regexp"
)

var (
	randIntInclusive = regexp.MustCompile("^randIntInclusive\\((\\d+)+,\\s*(\\d+)+\\)$")
	randString = regexp.MustCompile("^randString\\((\\d+)+,\\s*(\\d+)+\\)$")
	
	valueFunctions = [...]*regexp.Regexp{randIntInclusive, randString}
)

type Task struct {
	Name string
	Query string
	Values []string
	Iterations int
	Chance float64
	Run bool
}

func (t *Task) Execute(db *sql.DB) {
	
}

// ResolveValues goes through each Task.Values and computes that
// requested function if it exists. If that function does not exist,
// it will return an error.
func (t *Task) ResolveValues() ([]string, error) {
	for _, value := range t.Values {
		fmt.Println(for)
		valueg _, exp := range valueFunctions {
			fmt.Println(exp)
			if ! exp.MatchString(value) {
				continue
			}
			params := exp.FindStringSubmatch(value)
			
			if (exp == randIntInclusive) {
				
			} else if (exp == randString) {
				
			}
		}
	}
	return nil, nil
}

func ProcessTasks(taskLocation string, db *sql.DB) {
	filesInOrder, err := ioutil.ReadDir(taskLocation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	sort.Sort(ByTime(filesInOrder))
	for _, fileInfo := range filesInOrder {
		if fileInfo.IsDir() {
			continue
		}
		fmt.Printf("Processing task: %s \n", fileInfo.Name())
		file, err := os.Open(path.Join(taskLocation, fileInfo.Name()))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Cannot continue, exiting")
			os.Exit(1)
		}

		tasks := make([]Task, 0)
		err = json.NewDecoder(file).Decode(&tasks)

		if err != nil && err != io.EOF {
			fmt.Println(err)
			fmt.Println("Cannot continue, exiting")
			os.Exit(1)
		}

		for _, task := range tasks {
			task.Execute(db)
		}
	}
}

