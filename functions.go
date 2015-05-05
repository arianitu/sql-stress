package main

import (
	"math/rand"
	"strings"
	"regexp"
)

var (
	randIntInclusive = regexp.MustCompile("^randIntInclusive\\((\\d+)+,\\s*(\\d+)+\\)$")
	randString       = regexp.MustCompile("^randString\\((\\d+)+,\\s*(\\d+)+\\)$")
	valueFunctions   = [...]*regexp.Regexp{randIntInclusive, randString}
)

func TabN(n int) string {
	return strings.Repeat("\t", n)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(min, max int) string {
	n := rand.Intn(max-min) + min
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomIntInclusive(min, max int) int {
	return rand.Intn(max-min + 1) + min
}

func RandomIntExclusive(min, max int) int {
	return rand.Intn(max-min) + min
}
