package arr

import (
	"math/rand"

	"github.com/JoshPattman/goevo"
)

var _ goevo.Mutation[*Genotype[float64]] = &StdMutationStrategy[float64]{}

type StdMutationStrategy[T floatType] struct {
	// The probability of mutating each locus
	MutateProbability T
	// The standard deviation for the mutation
	MutateStd T
}

func (s *StdMutationStrategy[T]) Mutate(gt *Genotype[T]) {
	for i := range gt.Values {
		if T(rand.Float64()) < s.MutateProbability {
			gt.Values[i] += T(rand.NormFloat64()) * s.MutateStd
		}
	}
}
