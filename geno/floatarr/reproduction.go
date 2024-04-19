package floatarr

import "github.com/JoshPattman/goevo"

var _ goevo.Reproduction[Genotype] = &FloatsReproduction{}

// FloatsReproduction is a reproduction strategy for Float64sGenotype.
// It performs crossover and mutation.
type FloatsReproduction struct {
	// The probability of mutating each locus
	MutateProbability float64
	// The standard deviation for the mutation
	MutateStd float64
}

// Reproduce creates a new genotype by crossing over and mutating the given genotypes.
func (r *FloatsReproduction) Reproduce(gs []Genotype) Genotype {
	if len(gs) != 2 {
		panic("floatarr: expected 2 parents")
	}
	a, b := gs[0], gs[1]
	child := a.CrossoverWith(b)
	child.Mutate(r.MutateProbability, r.MutateStd)
	return child
}

// NumParents returns 2, as this reproduction strategy requires 2 parents.
func (r *FloatsReproduction) NumParents() int {
	return 2
}
