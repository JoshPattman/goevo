package floatarr

import (
	"math/rand"

	"github.com/JoshPattman/goevo"
)

var _ goevo.MutationStrategy[*Genotype] = &StdMutationStrategy{}

type StdMutationStrategy struct {
	// The probability of mutating each locus
	MutateProbability float64
	// The standard deviation for the mutation
	MutateStd float64
}

func (s *StdMutationStrategy) Mutate(gt *Genotype) {
	for i := range gt.Values {
		if rand.Float64() < s.MutateProbability {
			gt.Values[i] += rand.NormFloat64() * s.MutateStd
		}
	}
}
