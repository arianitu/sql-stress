package main

import (
	"testing"
	"io/ioutil"
	"path"
	"os"
	"fmt"
	"sort"
)

func TestFileSorting(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "test_sql_stress")
	if err != nil {
		t.Error(err)
	}
	filesToMake := [...]string {"test_2.sql", "test_1.sql", "test_3.sql", "test_10.sql", "test_100.sql"}
	expected := [...]string {"test_1.sql", "test_2.sql", "test_3.sql", "test_10.sql", "test_100.sql"}

	for _, name := range filesToMake {
		fmt.Println(path.Join(dir, name))
		ioutil.WriteFile(path.Join(dir, name), []byte{0}, os.ModeAppend)
	}
	filesInOrder, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sort.Sort(ByTime(filesInOrder))
	for index, fileInfo := range filesInOrder {
		if fileInfo.Name() != expected[index] {
			err := os.RemoveAll(dir)
			if err != nil {
				fmt.Println("Failed to cleanup temporary directory!")
			}
			t.Errorf("Directory not sorted properly, got %v, expected %v", fileInfo.Name(), expected[index])
		}
	}
	
	err = os.RemoveAll(dir)
	if err != nil {
		fmt.Println("Failed to cleanup temporary directory!")
	}
}

