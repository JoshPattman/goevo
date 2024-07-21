package arr

import (
	"math/rand"

	"github.com/JoshPattman/goevo"
)

var _ goevo.Mutation[*Genotype[bool]] = &RandomBoolMutationStrategy{}

type RandomBoolMutationStrategy struct {
	// The probability of mutating each locus
	MutateProbability float64
}

func (s *RandomBoolMutationStrategy) Mutate(gt *Genotype[bool]) {
	for i := range gt.Values {
		if rand.Float64() < s.MutateProbability {
			gt.Values[i] = !gt.Values[i]
		}
	}
}
