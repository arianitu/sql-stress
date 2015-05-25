package main

import (
	"testing"
)

func TestResolveValues(t *testing.T) {
	expectedStrings := []string{"string", "that", "should", "stay", "the", "same", "randString", "randIntInclusive"}
	s := &Step{}
	s.Init()
	for _, val := range expectedStrings {
		s.Values = append(s.Values, val)
	}
	vals, err := s.ResolveValues()
	if err != nil {
		t.Fatal(err)
	}
	for idx, val := range vals {
		if val != expectedStrings[idx] {
			t.Fatalf("expected %v, got %v", expectedStrings[idx], val)
		}
	}

	s = &Step{Values: []interface{}{23.25, "test", true}}
	vals, err = s.ResolveValues()
	if err != nil {
		t.Fatalf("expected ResolveValues to accept string, float64 and bool, but it did not because %v", err)
	}

	s = &Step{Values: []interface{}{[]string{"1", "2", "3"}}}
	vals, err = s.ResolveValues()
	if err == nil {
		t.Fatalf("expected ResolveValues to only accept string, float64 and bool, but it accepted []string")
	}

}

func TestResolveString(t *testing.T) {
	s := &Step{}
	s.Init()
	value, err := s.resolveString("randIntInclusive(10, 10)", 0)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := value.(int)
	if !ok {
		t.Fatal("Expected randIntInclusive to return an int!")
	}

	value, err = s.resolveString("randString(10, 15)", 0)
	if err != nil {
		t.Fatal(err)
	}
	_, ok = value.(string)
	if !ok {
		t.Fatal("Expected randString to return a string!")
	}

}

func TestResolveStringIncrementCount(t *testing.T) {
	s := &Step{}
	s.Init()

	value, err := s.resolveString("incrementingCount(10, 1)", 0)
	if err != nil {
		t.Fatal(err)
	}
	count, ok := value.(int64)
	if !ok {
		t.Fatal("Expected incrementingCount to return an int64!")
	}
	if count != 10 {
		t.Fatalf("Expected incrementingCount to return starting count %v got %v", 10, count)
	}

	value, err = s.resolveString("incrementingCount(10, 1)", 0)
	if err != nil {
		t.Fatal(err)
	}
	count, ok = value.(int64)
	if !ok {
		t.Fatal("Expected incrementingCount to return an int64!")
	}
	if count != 11 {
		t.Fatalf("Expected incrementingCount to increment by 1 using the same idx, starting value is %v and now is %v", 10, count)
	}

	value, err = s.resolveString("incrementingCount(50, 1)", 1)
	if err != nil {
		t.Fatal(err)
	}
	count, ok = value.(int64)
	if !ok {
		t.Fatal("Expected incrementingCount to return an int64!")
	}
	if count != 50 {
		t.Fatalf("Expected incrementingCount in a different idx to initialize to %v, but is %v", 50, value)
	}

	value, err = s.resolveString("incrementingCount(0, -1)", 2)
	if err != nil {
		t.Fatal(err)
	}
	count, ok = value.(int64)
	if !ok {
		t.Fatal("Expected incrementingCount to return an int64!")
	}
	if count != 0 {
		t.Fatalf("Expected incrementingCount to start at %v, is %v", 0, count)
	}

	value, err = s.resolveString("incrementingCount(0, -1)", 2)
	if err != nil {
		t.Fatal(err)
	}
	count, ok = value.(int64)
	if !ok {
		t.Fatal("Expected incrementingCount to return an int64!")
	}
	if count != -1 {
		t.Fatalf("Expected incrementingCount to decremement by 1 from %v and equal to %v, is %v", 0, -1, count)
	}
}


func BenchmarkResolveStringRandIntInclusive(b *testing.B) {
	s := &Step{}
	s.Init()
	
	for i := 0; i < b.N; i++ {
		s.resolveString("randIntInclusive(1, 5000000)", 0)
	}
}


func BenchmarkResolveStringRandString(b *testing.B) {
	s := &Step{}
	s.Init()
	
	for i := 0; i < b.N; i++ {
		s.resolveString("randString(1, 300)", 0)
	}
}

func BenchmarkResolveStringIncrementingCount(b *testing.B) {
	s := &Step{}
	s.Init()
	
	for i := 0; i < b.N; i++ {
		s.resolveString("incrementingCount(1, 1)", 0)
	}
}
