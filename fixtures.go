package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type ByTime []os.FileInfo

func (a ByTime) Len() int {
	return len(a)
}
func (a ByTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByTime) Less(i, j int) bool {
	formatString := "Failed to sort fixture directory. File %s is not in the format name_integertimestamp"

	nameA := strings.TrimSuffix(a[i].Name(), filepath.Ext(a[i].Name()))
	timeA, err := strconv.Atoi(strings.Split(nameA, "_")[1])
	if err != nil {
		panic(fmt.Sprintf(formatString, a[i].Name()))
	}
	nameB := strings.TrimSuffix(a[j].Name(), filepath.Ext(a[j].Name()))
	timeB, err := strconv.Atoi(strings.Split(nameB, "_")[1])
	if err != nil {
		panic(fmt.Sprintf(formatString, a[j].Name()))
	}
	return timeA < timeB
}

func SemicolonSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, ';'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we're done
	if atEOF {
		return 0, nil, nil
	}
	// Request more data.
	return 0, nil, nil
}

func ProcessFixtures(fixtureLocation string, db *sql.DB) {
	filesInOrder, err := ioutil.ReadDir(fixtureLocation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sort.Sort(ByTime(filesInOrder))
	for _, fileInfo := range filesInOrder {
		if fileInfo.IsDir() {
			continue
		}
		fmt.Printf("Processing fixture:%s \n", fileInfo.Name())
		file, err := os.Open(path.Join(fixtureLocation, fileInfo.Name()))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Cannot continue, exiting")
			os.Exit(1)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(SemicolonSplit)
		for scanner.Scan() {
			_, err := db.Exec(scanner.Text())
			if err != nil {
				fmt.Println(err)
				fmt.Println("Cannot continue, exiting")
				os.Exit(1)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Invalid input: %s", err)
		}

	}
}

