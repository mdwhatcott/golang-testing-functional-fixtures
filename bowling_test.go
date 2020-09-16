package example

import "testing"

func TestBowling_GutterGame(t *testing.T) {
	_TestBowling(t, RollMany(20, 0), AssertScore(0))
}
func TestBowling_AllOnes(t *testing.T) {
	_TestBowling(t, RollMany(20, 1), AssertScore(20))
}
func TestBowling_Spare(t *testing.T) {
	_TestBowling(t, RollSpare(), Roll(3), Roll(1), AssertScore(17))
}
func TestBowling_Strike(t *testing.T) {
	_TestBowling(t, RollStrike(), Roll(3), Roll(4), AssertScore(24))
}
func TestBowling_PerfectGame(t *testing.T) {
	_TestBowling(t, RollMany(12, allPins), AssertScore(300))
}

func Roll(pins int) BowlingFixtureOption {
	return RollMany(1, pins)
}
func RollSpare() BowlingFixtureOption {
	return RollMany(2, 5)
}
func RollStrike() BowlingFixtureOption {
	return Roll(allPins)
}
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
			this.Error(expected, actual)
		}
	}
}

func _TestBowling(t *testing.T, options ...BowlingFixtureOption) {
	t.Parallel()
	fixture := &BowlingFixture{T: t, game: new(BowlingGame)}
	for _, option := range options {
		option(fixture)
	}
}

type (
	BowlingFixtureOption func(this *BowlingFixture)
	BowlingFixture       struct {
		*testing.T
		game *BowlingGame
	}
)
