package goevo

// Forwarder is an interface for a thing that has a Forward function, used to take a list of inputs and return a list of outputs.
// Phenotype, LayeredHyperPhenotype, and other HyperNEAT phenotypes implement this interface.
type Forwarder interface {
	Forward(inputs []float64) []float64
}