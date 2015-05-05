package main

import (
	"testing"
)

func TestRandStringMin(t *testing.T) {
	l := RandomString(10, 11)
	if len(l) < 10 {
		t.Errorf("Random string min length is wrong, got %v expected %v", len(l), 10)
	}
}

func TestRandIntInclusiveMin(t *testing.T) {
	l := RandomIntInclusive(10, 10)
	if l != 10 {
		t.Errorf("min length is wrong, got %v expected %v", l, 10)
	}	
}

func TestRandIntExclusiveMin(t *testing.T) {
	l := RandomIntExclusive(10, 11)
	if l != 10 {
		t.Errorf("min length is wrong, got %v expected at least %v", l, 10)
	}	
}

func TestTabN(t *testing.T) {
	tabs := TabN(5)
	expected := "\t\t\t\t\t"
	if tabs != expected {
		t.Errorf("expected to be %v, got %v", expected, tabs)
	}
}
