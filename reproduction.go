package goevo

import (
	"math"
	"math/rand"
)

type ReproductionFunc func(*Genotype, *Genotype) *Genotype

func StdReproduction(counter *Counter, numNewSynapseStd, numNewNeuronsStd, numSynapseMutationStd, numPruneSynapsesStd float64, synapseMutStd, synapseGrowStd float64, newNeuronActivations []Activation) ReproductionFunc {
	return func(a, b *Genotype) *Genotype {
		child := CrossoverGenotypes(a, b)

		// Add synapses
		numNewSynapses := int(math.Round(rand.NormFloat64() * numNewSynapseStd))
		for i := 0; i < numNewSynapses; i++ {
			child.AddRandomSynapse(counter, synapseGrowStd)
		}

		// Mutate synapses
		numMutSynapses := int(math.Round(rand.NormFloat64() * numSynapseMutationStd))
		for i := 0; i < numMutSynapses; i++ {
			child.MutateRandomSynapse(synapseMutStd)
		}

		// Add neurons
		numNewNeurons := int(math.Round(rand.NormFloat64() * numNewNeuronsStd))
		for i := 0; i < numNewNeurons; i++ {
			child.AddRandomNeuron(counter, ChooseActivationFrom(newNeuronActivations))
		}

		// Prune Synapses
		numPruneSynapses := int(math.Round(rand.NormFloat64() * numPruneSynapsesStd))
		for i := 0; i < numPruneSynapses; i++ {
			child.PruneRandomSynapse()
		}

		return child
	}
}
