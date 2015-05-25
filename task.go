package main

import (
	"runtime"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"sort"
	"sync"
	"time"
)

type Conn struct {
	Vendor string
	Url string
	Workers int
	MaxOpenConn int
}

type Task struct {
	Conn Conn
	Parallel string
	Skip   bool
	Steps    []Step
}

type Step struct {
	Name       string
	Query      string
	Values     []interface{}
	Iterations int
	Tables     []string
	Skip     bool
	Chance     float64
	Run        bool
}

func PrintTableInfo(db *sql.DB, table string) {
	s := MySQLTableSize{Db: db, Table: table}
	err := s.Init()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf(TabN(2)+"Table: %v\n", table)
	fmt.Printf(TabN(3)+"table size: %v MB, index size: %v MB, avg row size: %v bytes, rows: %v \n",
		s.GetTableSize()/1000000,
		s.GetIndexSize()/1000000,
		s.GetAvgRowSize(),
		s.GetRows())
}

func (t *Task) Step(db *sql.DB, queryIn chan<- Query) {
	for _, step := range t.Steps {
		if step.Skip {
			continue
		}
		err := step.Execute(db, queryIn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func (s *Step) Execute(db *sql.DB, queryIn chan<- Query) error {

	fmt.Println(TabN(1) + s.Name)

	// iterations default value is 1
	if s.Iterations <= 0 {
		s.Iterations = 1
	}
	
	wg := &sync.WaitGroup{}
	wg.Add(s.Iterations)

	sink := make(chan int64)

	var worst int64 = 0
	var best int64 = math.MaxInt64
	var totalTime int64 = 0
	go func() {
		for t := range sink {
			totalTime += t
			if t > worst {
				worst = t
			}
			if t < best {
				best = t
			}
			wg.Done()
		}
	}()

	for i := 0; i < s.Iterations; i++ {
		values, err := s.ResolveValues()
		if err != nil {
			return err
		}
		queryIn <- Query{Query: s.Query, Values: values, Done: sink}
	}
	wg.Wait()

	total := time.Duration(totalTime) * time.Nanosecond
	avgDuration := time.Duration(totalTime/int64(s.Iterations)) * time.Nanosecond
	bestDuration := time.Duration(best) * time.Nanosecond
	worstDuration := time.Duration(worst) * time.Nanosecond

	fmt.Printf(TabN(2)+"Avg: %v Worst: %v Best: %v Total: %v \n", avgDuration, worstDuration, bestDuration, total)
	fmt.Println("")
	for _, table := range s.Tables {
		PrintTableInfo(db, table)
	}

	return nil
}

// ResolveValues goes through each Task.Values and computes that
// requested function if it exists. If that function does not exist,
// it will return an error.
func (s *Step) ResolveValues() ([]interface{}, error) {
	values := make([]interface{}, 0)
	for _, anything := range s.Values {
		switch v := anything.(type) {
		case string:
			r, err := resolveString(v)
			if err != nil {
				return nil, err
			}
			values = append(values, r)
		case float64:
			values = append(values, v)
		case bool:
			values = append(values, v)
		default:
			return nil, fmt.Errorf("Value array in Task.step must be a string, float64, or bool")

		}
	}
	return values, nil
}

func ProcessTasks(settings *Settings, db *sql.DB, queryIn chan<- Query) {
	filesInOrder, err := ioutil.ReadDir(settings.TaskLocation)
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
		file, err := os.Open(path.Join(settings.TaskLocation, fileInfo.Name()))
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
		if task.Skip {
			continue
		}

		// The task specified a db, we'll use the task db instead of the default one.
		if task.Conn.Vendor != "" {
			workers := task.Conn.Workers
			if workers <= 0 {
				workers = settings.Workers
			}
			maxOpenConn := task.Conn.MaxOpenConn
			if maxOpenConn <= 0 {
				maxOpenConn = runtime.NumCPU()
			}
			taskQueryIn, taskDb, err := SpawnWorkers(task.Conn.Vendor, task.Conn.Url, workers, maxOpenConn)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Cannot continue, exiting")
				os.Exit(1)
			}
			
			task.Step(taskDb, taskQueryIn)
			// task.Step waits for all queries to complete before continuing, we're safe to close the channel
			close(taskQueryIn)
			taskDb.Close()
		} else {
			task.Step(db, queryIn)
		}
	}
}
