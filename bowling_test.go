package example

import (
	"strings"
	"testing"
)

func TestBowling(t *testing.T) {
	var cases = map[string][]BowlingFixtureOption{
		"gutter game": {
			RollMany(20, 0),
			AssertScore(0),
		},
		"all ones": {
			RollMany(20, 1),
			AssertScore(20),
		},
		"spare": {
			RollSpare(),
			Roll(3),
			Roll(1),
			AssertScore(17),
		},
		"strike": {
			RollStrike(),
			Roll(3),
			Roll(4),
			AssertScore(24),
		},
		"perfection": {
			RollMany(12, allPins),
			AssertScore(300),
		},
	}
	for name, options := range cases {
		t.Run(name, func(t *testing.T) {
			if strings.HasPrefix(name, "SKIP") {
				t.SkipNow()
			}
			fixture := &BowlingFixture{T: t, game: new(BowlingGame)}
			for _, option := range options {
				option(fixture)
			}
		})
	}
}
func Roll(pins int) BowlingFixtureOption { return RollMany(1, pins) }
func RollSpare() BowlingFixtureOption    { return RollMany(2, 5) }
func RollStrike() BowlingFixtureOption   { return Roll(allPins) }
func RollMany(times int, pins int) BowlingFixtureOption {
	return func(this *BowlingFixture) {
		for ; times > 0; times-- {
			this.game.RecordRoll(pins)
		}
	}
}
func AssertScore(expected int) BowlingFixtureOption {
	return func(this *BowlingFixture) {
		actual := this.game.CalculateScore()
		if !(actual == expected) {
			this.Helper()
			this.Error(expected, actual)
		}
	}
}

type BowlingFixtureOption func(this *BowlingFixture)
type BowlingFixture struct {
	*testing.T
	game *BowlingGame
}
