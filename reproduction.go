package goevo

import (
	"math"
	"math/rand"
)

// Reproduction is an interface for the reproduction of two parents to create a child
type Reproduction[T any] interface {
	// Reproduce creates a new genotype from the two parents, where the first parent is fitter
	Reproduce(a, b T) T
}

func stdN(std float64) int {
	v := math.Abs(rand.NormFloat64() * std)
	if v > std*10 {
		v = std * 10 // Lets just cap this at 10 std to prevent any sillyness
	}
	return int(math.Round(v))
}
