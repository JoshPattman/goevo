package goevo

var _ Reproduction[Float64sGenotype] = &FloatsReproduction{}

// FloatsReproduction is a reproduction strategy for Float64sGenotype.
// It performs crossover and mutation.
type FloatsReproduction struct {
	// The probability of mutating each locus
	MutateProbability float64
	// The standard deviation for the mutation
	MutateStd float64
}

// Reproduce creates a new genotype by crossing over and mutating the given genotypes.
func (r *FloatsReproduction) Reproduce(a, b Float64sGenotype) Float64sGenotype {
	child := a.CrossoverWith(b)
	child.Mutate(r.MutateProbability, r.MutateStd)
	return child
}
