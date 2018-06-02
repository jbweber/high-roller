package main

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

var manyDice = regexp.MustCompile(`(\s*\d*?d\d+(?:\s*(?:\+|\-)\s*\d+)?)`)
var oneDice = regexp.MustCompile(`\s*(\d+)?d(\d+)(?:\s*(\+|\-)\s*(\d+))?`)
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type diceRoll struct {
	count int
	dice  int
	oper  string
	mod   int
}

func Parse(in string) diceRoll {
	r := oneDice.FindStringSubmatch(in)

	count, dice, oper, mod := 1, 20, "", 0

	// match + 4 captures
	if len(r) != 5 {
		return diceRoll{count, dice, oper, mod}
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

	return diceRoll{count, dice, oper, mod}
}

func ParseMany(in string) []diceRoll {
	// parse
	parsed := manyDice.FindAllStringSubmatch(in, -1)

	if len(parsed) == 0 {
		return []diceRoll{diceRoll{1, 20, "", 0}}
	}

	result := make([]diceRoll, len(parsed))
	for i, p := range parsed {
		if len(p) != 2 {
			result[i] = diceRoll{1, 20, "", 0}
		}
		roll := Parse(p[1])
		result[i] = roll
	}

	return result
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
