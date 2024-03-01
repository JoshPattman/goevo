package goevo

import "math/rand"

type Float64sGenotype []float64

func NewFloat64sGenotype(size int, std float64) Float64sGenotype {
	gt := make(Float64sGenotype, size)
	for i := range gt {
		gt[i] = rand.NormFloat64() * std
	}
	return gt
}

func (f64s Float64sGenotype) Mutate(probPerLocus, mutStd float64) {
	for i := range f64s {
		if rand.Float64() < probPerLocus {
			f64s[i] += rand.NormFloat64() * mutStd
		}
	}
}

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

func (f64s Float64sGenotype) Clone() Float64sGenotype {
	clone := make(Float64sGenotype, len(f64s))
	copy(clone, f64s)
	return clone
}
