package goevo

import "math/rand"

func integerMidpoint(ia, ib int) int {
	d := (ib - ia) / 2
	if d == 0 {
		return ia + 1
	}
	return ia + d
}

// randRange exclusive
func randRange(i, a int) int {
	if a-i <= 0 {
		return i
	}
	return rand.Intn(a-i) + i
}
