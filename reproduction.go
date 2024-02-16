package goevo

import (
	"math"
	"math/rand"
)

type Reproduction interface {
	// A is assumed to be the fitter parent
	Reproduce(a, b *Genotype) *Genotype
}

func stdN(std float64) int {
	v := math.Abs(rand.NormFloat64() * std)
	if v > std*10 {
		v = std * 10 // Lets just cap this at 10 std to prevent any sillyness
	}
	return int(math.Round(v))
}

type StdReproduction struct {
	StdNumNewSynapses       float64
	StdNumNewNeurons        float64
	StdNumMutateSynapses    float64
	StdNumPruneSynapses     float64
	StdNumMutateActivations float64

	StdNewSynapseWeight    float64
	StdMutateSynapseWeight float64

	Counter             *Counter
	PossibleActivations []Activation
}

func (r *StdReproduction) Reproduce(a, b *Genotype) *Genotype {
	g := a.CrossoverWith(b)

	for i := 0; i < stdN(r.StdNewSynapseWeight); i++ {
		g.AddRandomSynapse(r.Counter, r.StdNewSynapseWeight, false)
	}
	for i := 0; i < stdN(r.StdNumNewNeurons); i++ {
		g.AddRandomNeuron(r.Counter, r.PossibleActivations...)
	}
	for i := 0; i < stdN(r.StdNumMutateSynapses); i++ {
		g.MutateRandomSynapse(r.StdMutateSynapseWeight)
	}
	for i := 0; i < stdN(r.StdNumPruneSynapses); i++ {
		g.RemoveRandomSynapse()
	}
	for i := 0; i < stdN(r.StdNumMutateActivations); i++ {
		g.MutateRandomActivation(r.PossibleActivations...)
	}

	return g
}
