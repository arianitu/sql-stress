package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
)

func TestFileSorting(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "test_sql_stress")
	if err != nil {
		t.Fatal(err)
	}
	filesToMake := [...]string{"test_2.sql", "test_1.sql", "test_3.sql", "test_10.sql", "test_100.sql"}
	expected := [...]string{"test_1.sql", "test_2.sql", "test_3.sql", "test_10.sql", "test_100.sql"}

	for _, name := range filesToMake {
		ioutil.WriteFile(path.Join(dir, name), []byte{0}, os.ModeAppend)
	}
	filesInOrder, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	sort.Sort(ByTime(filesInOrder))
	for index, fileInfo := range filesInOrder {
		if fileInfo.Name() != expected[index] {
			err := os.RemoveAll(dir)
			if err != nil {
				fmt.Println("Failed to cleanup temporary directory!")
			}
			t.Fatalf("Directory not sorted properly, got %v, expected %v", fileInfo.Name(), expected[index])
		}
	}

	err = os.RemoveAll(dir)
	if err != nil {
		fmt.Println("Failed to cleanup temporary directory!")
	}
}

func TestSemicolonSplit(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "test_sql_stress")
	if err != nil {
		t.Fatal(err)
	}
	expected := [...]string{"SELECT * FROM a", "SELECT * FROM b", "SELECT * FROM c"}
	file.WriteString(strings.Join(expected[:], ";") + ";")
	file.Seek(0, 0)

	scanner := bufio.NewScanner(file)
	scanner.Split(SemicolonSplit)
	i := 0
	for scanner.Scan() {
		v := scanner.Text()
		if v != expected[i] {
			t.Fatalf("expected %v got %v", expected[i], v)
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}
