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
	StdNumNewSynapses          float64
	StdNumNewRecurrentSynapses float64
	StdNumNewNeurons           float64
	StdNumMutateSynapses       float64
	StdNumPruneSynapses        float64
	StdNumMutateActivations    float64

	StdNewSynapseWeight    float64
	StdMutateSynapseWeight float64

	MaxHiddenNeurons int

	Counter             *Counter
	PossibleActivations []Activation
}

func (r *StdReproduction) Reproduce(a, b *Genotype) *Genotype {
	g := a.CrossoverWith(b)

	for i := 0; i < stdN(r.StdNewSynapseWeight); i++ {
		g.AddRandomSynapse(r.Counter, r.StdNewSynapseWeight, false)
	}
	for i := 0; i < stdN(r.StdNumNewRecurrentSynapses); i++ {
		g.AddRandomSynapse(r.Counter, r.StdNewSynapseWeight, true)
	}
	for i := 0; i < stdN(r.StdNumNewNeurons); i++ {
		if r.MaxHiddenNeurons < 0 || g.NumHiddenNeurons() < r.MaxHiddenNeurons {
			g.AddRandomNeuron(r.Counter, r.PossibleActivations...)
		}
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

type ProbReproduction struct {
	NewNeuronProbability           float64
	NewSynapseProbability          float64
	NewRecurrentSynapseProbability float64
	RemoveSynapseProbability       float64
	MutateSynapseProbability       float64
	MutateActivationProbability    float64
	SetSynapseZeroProbability      float64

	NewSynapseStd    float64
	MutateSynapseStd float64
	Activations      []Activation

	UseUnfitParentProbability float64

	Counter *Counter
}

func (r *ProbReproduction) Reproduce(a, b *Genotype) *Genotype {
	if rand.Float64() < r.UseUnfitParentProbability {
		a, b = b, a
	}
	g := a.CrossoverWith(b)

	if rand.Float64() < r.NewNeuronProbability {
		g.AddRandomNeuron(r.Counter, r.Activations...)
	}
	if rand.Float64() < r.NewSynapseProbability {
		g.AddRandomSynapse(r.Counter, r.NewSynapseStd, false)
	}
	if rand.Float64() < r.NewRecurrentSynapseProbability {
		g.AddRandomSynapse(r.Counter, r.NewSynapseStd, true)
	}
	if rand.Float64() < r.RemoveSynapseProbability {
		g.RemoveRandomSynapse()
	}
	if rand.Float64() < r.SetSynapseZeroProbability {
		g.ResetRandomSynapse()
	}
	if rand.Float64() < r.MutateSynapseProbability {
		g.MutateRandomSynapse(r.MutateSynapseStd)
	}
	if rand.Float64() < r.MutateActivationProbability {
		g.MutateRandomActivation(r.Activations...)
	}
	return g
}
