// Package floatarr provides a genotype type that is a slice of float64 values.
// It also provides functions for reproduction using these genotypes.
package floatarr

import (
	"math/rand"

	"github.com/JoshPattman/goevo"
)

var _ goevo.Cloneable = Genotype{}

// Genotype is a genotype that is a slice of float64 values.
type Genotype struct {
	Values []float64
}

// NewGenotype creates a new genotype with the given size, where each value is a random float64 with a normal distribution with the given standard deviation.
func NewGenotype(size int, std float64) *Genotype {
	gt := make([]float64, size)
	for i := range gt {
		gt[i] = rand.NormFloat64() * std
	}
	return &Genotype{Values: gt}
}

// Clone returns a new genotype that is a copy of this genotype.
func (f64s Genotype) Clone() any {
	clone := make([]float64, len(f64s.Values))
	copy(clone, f64s.Values)
	return &Genotype{Values: clone}
}
