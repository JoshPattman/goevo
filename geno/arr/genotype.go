// Package arr provides a genotype type that is a slice of values.
// It also provides functions for reproduction using these genotypes.
package arr

import (
	"math/rand/v2"

	"github.com/JoshPattman/goevo"
)

var _ goevo.Cloneable = Genotype[int]{}

// Genotype is a genotype that is a slice of float64 values.
type Genotype[T any] struct {
	Values []T
}

// NewFloatGenotype creates a new genotype with the given size,
// where each value is a random float (of type T) with a normal distribution with the given standard deviation.
func NewFloatGenotype[T floatType](size int, std T) *Genotype[T] {
	gt := make([]T, size)
	for i := range gt {
		gt[i] = T(rand.NormFloat64()) * std
	}
	return &Genotype[T]{Values: gt}
}

// NewRuneGenotype creates a new genotype with the given size,
// where each value is a random rune from the given runeset.
func NewRuneGenotype(size int, runeset []rune) *Genotype[rune] {
	gt := make([]rune, size)
	for i := range gt {
		gt[i] = runeset[rand.N(len(runeset))]
	}
	return &Genotype[rune]{Values: gt}
}

// NewBoolGenotype creates a new genotype with the given size,
// where each value is a random boolean.
func NewBoolGenotype(size int) *Genotype[bool] {
	gt := make([]bool, size)
	for i := range gt {
		gt[i] = rand.Float64() < 0.5
	}
	return &Genotype[bool]{Values: gt}
}

// Clone returns a new genotype that is a copy of this genotype.
func (g Genotype[T]) Clone() any {
	clone := make([]T, len(g.Values))
	copy(clone, g.Values)
	return &Genotype[T]{Values: clone}
}
