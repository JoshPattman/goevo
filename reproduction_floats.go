package goevo

var _ Reproduction[Float64sGenotype] = &FloatsReproduction{}

type FloatsReproduction struct {
	// The probability of mutating each locus
	MutateProbability float64
	// The standard deviation for the mutation
	MutateStd float64
}

// Reproduce implements Reproduction.
func (r *FloatsReproduction) Reproduce(a, b Float64sGenotype) Float64sGenotype {
	child := a.CrossoverWith(b)
	child.Mutate(r.MutateProbability, r.MutateStd)
	return child
}
