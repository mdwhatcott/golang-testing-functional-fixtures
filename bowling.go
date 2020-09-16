package example

const (
	allPins   = 10
	allFrames = 10
	maxThrows = 21
)

type BowlingGame struct {
	score  int
	throw  int
	throws [maxThrows]int
}

func (this *BowlingGame) RecordRoll(pins int) {
	this.throws[this.throw] += pins
	this.throw++
}

func (this *BowlingGame) CalculateScore() int {
	this.throw = 0
	for frame := 0; frame < allFrames; frame++ {
		this.score += this.scoreThrowsInFrame()
		this.throw += this.advanceThrowToNextFrame()
	}
	return this.score
}

func (this *BowlingGame) scoreThrowsInFrame() int {
	if this.currentThrowIsStrike() {
		return this.strikeScore()
	} else if this.currentFrameIsSpare() {
		return this.spareScore()
	} else {
		return this.frameScore()
	}
}
func (this *BowlingGame) advanceThrowToNextFrame() int {
	if this.currentThrowIsStrike() {
		return 1
	} else {
		return 2
	}
}

func (this *BowlingGame) currentThrowIsStrike() bool {
	return this.at(0) == allPins
}
func (this *BowlingGame) currentFrameIsSpare() bool {
	return this.frameScore() == allPins
}

func (this *BowlingGame) strikeScore() int {
	return allPins + this.at(1) + this.at(2)
}
func (this *BowlingGame) spareScore() int {
	return allPins + this.at(2)
}
func (this *BowlingGame) frameScore() int {
	return this.at(0) + this.at(1)
}

func (this *BowlingGame) at(offset int) int {
	return this.throws[this.throw+offset]
}
