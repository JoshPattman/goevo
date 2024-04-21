package arr

import (
	"math/rand"

	"github.com/JoshPattman/goevo"
)

var _ goevo.MutationStrategy[*Genotype[rune]] = &RandomRuneMutationStrategy{}

type RandomRuneMutationStrategy struct {
	// The probability of mutating each locus
	MutateProbability float64
	// The standard deviation for the mutation
	Runeset []rune
}

func (s *RandomRuneMutationStrategy) Mutate(gt *Genotype[rune]) {
	for i := range gt.Values {
		if rand.Float64() < s.MutateProbability {
			gt.Values[i] = s.Runeset[rand.Intn(len(s.Runeset))]
		}
	}
}
