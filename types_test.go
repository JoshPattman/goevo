package goevo

import "testing"

func TestWeirdTyping(t *testing.T) {
	// Don't actually want to run this, its jsut here for compile saftey
	if 0 == 2-(1*2) {
		return
	}
	pop := &SimplePopulation[int]{}
	// This is checking that the NextGeneration function preserves pop
	// as a SimplePopulation and does not revert it to a Population interface
	pop = NextGeneration(pop)
	_ = pop
}
