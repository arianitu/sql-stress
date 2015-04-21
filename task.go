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
	"strconv"
	"math/rand"
	"sync"
)

var (
	randIntInclusive = regexp.MustCompile("^randIntInclusive\\((\\d+)+,\\s*(\\d+)+\\)$")
	randString = regexp.MustCompile("^randString\\((\\d+)+,\\s*(\\d+)+\\)$")
	valueFunctions = [...]*regexp.Regexp{randIntInclusive, randString}
)

type Task struct {
	Url string
	Parallel string
	Steps []Step
}

func (t *Task) Step(db *sql.DB, queryIn chan<- Query) {
	for _, step := range t.Steps {
		err := step.Execute(db, queryIn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

type Step struct {
	Name string
	Query string
	Values []string
	Iterations int
	Chance float64
	Run bool
}

func (s *Step) Execute(db *sql.DB, queryIn chan<- Query) error {

	fmt.Println("     " + s.Name)

	wg := &sync.WaitGroup{}
	wg.Add(s.Iterations)
	for i := 0; i < s.Iterations; i++ {
		values, err := s.ResolveValues()
		if err != nil {
			return err
		}
		queryIn <- Query{Query: s.Query, Values: values, WaitGroup: wg }
	}
	wg.Wait()
	
	return nil
}

// ResolveValues goes through each Task.Values and computes that
// requested function if it exists. If that function does not exist,
// it will return an error.
func (s *Step) ResolveValues() ([]interface{}, error) {
	values := make([]interface{}, 0)
	for _, value := range s.Values {
		
		for _, exp := range valueFunctions {
			if ! exp.MatchString(value) {
				continue
			}
			params := exp.FindStringSubmatch(value)
			
			if (exp == randIntInclusive) {
				min, err := strconv.Atoi(params[1])
				if err != nil {
					return nil, fmt.Errorf("First parameter of randIntIncusive must be an integer! Got: %v", params[1])
				}
				max, err := strconv.Atoi(params[2])
				if err != nil {
					return nil, fmt.Errorf("Second parameter of randIntIncusive must be an integer! Got: %v", params[2])
				}
				
				r := rand.Intn(max - min) + min
				values = append(values, r)
			} else if (exp == randString) {
				min, err := strconv.Atoi(params[1])
				if err != nil {
					return nil, fmt.Errorf("First parameter of randString must be an integer! Got: %v", params[1])
				}
				max, err := strconv.Atoi(params[2])
				if err != nil {
					return nil, fmt.Errorf("Second parameter of randString must be an integer! Got: %v", params[2])
				}
				values = append(values, RandomString(min, max))
				
			}
		}
	}
	return values, nil
}

func ProcessTasks(taskLocation string, db *sql.DB, queryIn chan<- Query) {
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

		var task Task
		err = json.NewDecoder(file).Decode(&task)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			fmt.Println("Cannot continue, exiting")
			os.Exit(1)
		}
		task.Step(db, queryIn)
	}
}

