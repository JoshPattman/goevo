package goevo

import "math/rand"

// Float64sGenotype is a genotype that is a slice of float64 values.
type Float64sGenotype []float64

// NewFloat64sGenotype creates a new genotype with the given size, where each value is a random float64 with a normal distribution with the given standard deviation.
func NewFloat64sGenotype(size int, std float64) Float64sGenotype {
	gt := make(Float64sGenotype, size)
	for i := range gt {
		gt[i] = rand.NormFloat64() * std
	}
	return gt
}

// Mutate modifies the genotype by adding a random value from a normal distribution to each value with the given probability and standard deviation.
func (f64s Float64sGenotype) Mutate(probPerLocus, mutStd float64) {
	for i := range f64s {
		if rand.Float64() < probPerLocus {
			f64s[i] += rand.NormFloat64() * mutStd
		}
	}
}

// CrossoverWith returns a new genotype that is a combination of this genotype and the other genotype.
func (f64s Float64sGenotype) CrossoverWith(other Float64sGenotype) Float64sGenotype {
	if len(f64s) != len(other) {
		panic("genotypes must have the same length")
	}
	child := make(Float64sGenotype, len(f64s))
	for i := range f64s {
		if rand.Float64() < 0.5 {
			child[i] = f64s[i]
		} else {
			child[i] = other[i]
		}
	}
	return child
}

// Clone returns a new genotype that is a copy of this genotype.
func (f64s Float64sGenotype) Clone() Float64sGenotype {
	clone := make(Float64sGenotype, len(f64s))
	copy(clone, f64s)
	return clone
}
