package main

import (
	"testing"
)

func TestResolveValues(t *testing.T) {
	expectedStrings := []string{"string", "that", "should", "stay", "the", "same", "randString", "randIntInclusive"}
	s := &Step{}
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
	
	s = &Step{ Values: []interface{}{[]string{"1", "2", "3"}} }
	vals, err = s.ResolveValues()
	if err == nil {
		t.Fatalf("expected ResolveValues to only accept string, float64 and bool, but it accepted []string")
	}
	
}
