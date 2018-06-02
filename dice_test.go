package main

import (
	"testing"
)

func checkSlice(in []bool) bool {
	for _, v := range in {
		if !v {
			return false
		}
	}
	return true
}

func TestRoll(t *testing.T) {
	tt := []struct {
		name string
		dice int
	}{
		{"1d4", 4},
		{"1d6", 6},
		{"1d8", 8},
		{"1d10", 10},
		{"1d12", 12},
		{"1d20", 20},
		{"1d100", 100},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// ensure we get all expected values and none out of range
			s := make([]bool, tc.dice)
			d := false
			for d == false {
				v := Roll(1, tc.dice)
				if v < 1 || v > tc.dice {
					t.Fatalf("%d outside the range of 1 to %d", v, tc.dice)
				}
				s[v-1] = true
				d = checkSlice(s)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tt := []struct {
		value string
		roll  diceRoll
	}{
		{"", diceRoll{1, 20, "", 0}},
		{"d20", diceRoll{1, 20, "", 0}},
		{"1d20", diceRoll{1, 20, "", 0}},
		{" 1d20  ", diceRoll{1, 20, "", 0}},
		{" 2d20  ", diceRoll{2, 20, "", 0}},
		{"4d6", diceRoll{4, 6, "", 0}},
		{"1d8+1", diceRoll{1, 8, "+", 1}},
		{"1d8 + 1", diceRoll{1, 8, "+", 1}},
		{"   1d8 + 1	", diceRoll{1, 8, "+", 1}},
		{"1d8-1", diceRoll{1, 8, "-", 1}},
		{"1d8 - 1", diceRoll{1, 8, "-", 1}},
		{"   1d8 - 1	", diceRoll{1, 8, "-", 1}},
	}

	for _, tc := range tt {
		t.Run(tc.value, func(t *testing.T) {
			roll := Parse(tc.value)

			if roll.count != tc.roll.count {
				t.Errorf("%s: count expected %d, actual %d", tc.value, tc.roll.count, roll.count)
			}

			if roll.dice != tc.roll.dice {
				t.Errorf("%s: dice expected %d, actual %d", tc.value, tc.roll.dice, roll.dice)
			}

			if roll.oper != tc.roll.oper {
				t.Errorf("%s: oper expected %s, actual %s", tc.value, tc.roll.oper, roll.oper)
			}

			if roll.mod != tc.roll.mod {
				t.Errorf("%s: mod expected %d, actual %d", tc.value, tc.roll.mod, roll.mod)
			}
		})
	}
}

func TestParseMany(t *testing.T) {
	tt := []struct {
		value string
		rolls []diceRoll
	}{
		{"", []diceRoll{diceRoll{1, 20, "", 0}}},
		{"d20", []diceRoll{diceRoll{1, 20, "", 0}}},
		{"1d20", []diceRoll{diceRoll{1, 20, "", 0}}},
		{" 1d20  ", []diceRoll{diceRoll{1, 20, "", 0}}},
		{" 2d20  ", []diceRoll{diceRoll{2, 20, "", 0}}},
		{"4d6", []diceRoll{diceRoll{4, 6, "", 0}}},
		{"1d8+1", []diceRoll{diceRoll{1, 8, "+", 1}}},
		{"1d8 + 1", []diceRoll{diceRoll{1, 8, "+", 1}}},
		{"   1d8 + 1	", []diceRoll{diceRoll{1, 8, "+", 1}}},
		{"1d8-1", []diceRoll{diceRoll{1, 8, "-", 1}}},
		{"1d8 - 1", []diceRoll{diceRoll{1, 8, "-", 1}}},
		{"   1d8 - 1	", []diceRoll{diceRoll{1, 8, "-", 1}}},
		{"1d8+1 3d6", []diceRoll{diceRoll{1, 8, "+", 1}, diceRoll{3, 6, "", 0}}},
		{"   1d8 + 1    3d6			", []diceRoll{diceRoll{1, 8, "+", 1}, diceRoll{3, 6, "", 0}}},
		{"   1d8 + 1    3d6			", []diceRoll{diceRoll{1, 8, "+", 1}, diceRoll{3, 6, "", 0}}},
	}

	for _, tc := range tt {
		t.Run(tc.value, func(t *testing.T) {
			rolls := ParseMany(tc.value)

			if len(rolls) != len(tc.rolls) {
				t.Errorf("expected %d rolls got %d", len(tc.rolls), len(rolls))
			}

			for i, roll := range rolls {
				tcRoll := tc.rolls[i]

				if roll.count != tcRoll.count {
					t.Errorf("%s: count expected %d, actual %d", tc.value, tcRoll.count, roll.count)
				}

				if roll.dice != tcRoll.dice {
					t.Errorf("%s: dice expected %d, actual %d", tc.value, tcRoll.dice, roll.dice)
				}

				if roll.oper != tcRoll.oper {
					t.Errorf("%s: oper expected %s, actual %s", tc.value, tcRoll.oper, roll.oper)
				}

				if roll.mod != tcRoll.mod {
					t.Errorf("%s: mod expected %d, actual %d", tc.value, tcRoll.mod, roll.mod)
				}

			}
		})
	}

}
