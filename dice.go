package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

var diceRegex = regexp.MustCompile(`([0-9]*)d(\d+)`)
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func Parse(in string) (int, int) {
	result := diceRegex.FindStringSubmatch(in)

	if len(result) < 3 {
		return 1, 20
	}

	countStr := result[1]
	maxStr := result[2]

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 1
	}
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		count = 20
	}
	return count, max
}

func Roll(min, max int) int {
	ra := r.Intn(max) // [0, max)
	return ra + min   // add min to give us values on the edges
}

func RollMany(count, max int) []int {
	fmt.Printf("RollMany(%v, %v)\n", count, max)
	results := make([]int, count)
	for i := 0; i < count; i++ {
		results[i] = Roll(1, max)
	}
	return results
}
