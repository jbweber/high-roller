package main

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

var oneDice = regexp.MustCompile(`\s*(\d+)?d(\d+)(?:\s*(\+|\-)?\s*(\d+))?`)

var diceRegex = regexp.MustCompile(`([0-9]*)d(\d+)`)
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func Parse(in string) (int, int, string, int) {
	r := oneDice.FindStringSubmatch(in)

	count, dice, oper, mod := 1, 20, "", 0

	// match + 4 captures
	if len(r) != 5 {
		return count, dice, oper, mod
	}

	cs := r[1]
	ds := r[2]
	oper = r[3]
	ms := r[4]

	count, err := strconv.Atoi(cs)
	if err != nil {
		count = 1
	}

	dice, err = strconv.Atoi(ds)
	if err != nil {
		dice = 20
	}

	mod, err = strconv.Atoi(ms)
	if err != nil {
		mod = 0
	}

	return count, dice, oper, mod
}

func Roll(min, max int) int {
	ra := r.Intn(max) // [0, max)
	return ra + min   // add min to give us values on the edges
}

func RollMany(count, max int) []int {
	results := make([]int, count)
	for i := 0; i < count; i++ {
		results[i] = Roll(1, max)
	}
	return results
}
