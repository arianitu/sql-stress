package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
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
	return rand.Intn(max-min+1) + min
}

func RandomIntExclusive(min, max int) int {
	return rand.Intn(max-min) + min
}

func resolveString(value string) (interface{}, error) {
	for _, exp := range valueFunctions {
		if !exp.MatchString(value) {
			continue
		}
		params := exp.FindStringSubmatch(value)

		if exp == randIntInclusive {
			min, err := strconv.Atoi(params[1])
			if err != nil {
				return nil, fmt.Errorf("First parameter of randIntIncusive must be an integer! Got: %v", params[1])
			}
			max, err := strconv.Atoi(params[2])
			if err != nil {
				return nil, fmt.Errorf("Second parameter of randIntIncusive must be an integer! Got: %v", params[2])
			}
			return RandomIntInclusive(min, max), nil
		} else if exp == randString {
			min, err := strconv.Atoi(params[1])
			if err != nil {
				return nil, fmt.Errorf("First parameter of randString must be an integer! Got: %v", params[1])
			}
			max, err := strconv.Atoi(params[2])
			if err != nil {
				return nil, fmt.Errorf("Second parameter of randString must be an integer! Got: %v", params[2])
			}
			return RandomString(min, max), nil
		}
	}
	return value, nil
}
