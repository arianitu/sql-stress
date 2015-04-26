package main

import (
	"math/rand"
	"strings"
)

func TabN(n int) string {
	return strings.Repeat("\t", n)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func RandomString(min, max int) string {
	n := rand.Intn(max - min) + min
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
