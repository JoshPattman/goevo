package goevo

import (
	"errors"
	"math/rand"
)

// Mutate the weight of a random synapse in `g` by sampling the normal distribution with standard deviation `stddev“
func (g *Genotype) MutateRandomSynapse(stddev float64) {
	if len(g.Synapses) == 0 {
		return
	}
	k := rand.Intn(len(g.Synapses))
	for _, s := range g.Synapses {
		if k == 0 {
			s.Weight += rand.NormFloat64() * stddev
			return
		}
		k--
	}
	panic("unreachable")
}

// Prune a random synapse from `g`. This has the capability to prun more than one synapse and neuron as it also removes redundant neurons and synapses.
func (g *Genotype) PruneRandomSynapse() {
	if len(g.Synapses) == 0 {
		return
	}
	k := rand.Intn(len(g.Synapses))
	for s := range g.Synapses {
		if k == 0 {
			g.PruneSynapse(s)
			return
		}
		k--
	}
	panic("unreachable")
}

// Add a new synapse to `g` with weight sampled from normal distribution with standard deviation `stddev`.
func (g *Genotype) AddRandomSynapse(counter *Counter, weightStddev float64) error {
	return addRandomSynapse(counter, g, weightStddev, false, 10)
}

// Add a new recurrent (backwards) synapse to `g` with weight sampled from normal distribution with standard deviation `stddev`.
func (g *Genotype) AddRandomRecurrentSynapse(counter *Counter, weightStddev float64) error {
	return addRandomSynapse(counter, g, weightStddev, true, 10)
}

// Helper function used by both synapse adding functions
func addRandomSynapse(counter *Counter, g *Genotype, weightStddev float64, isRecurrent bool, attempts int) error {
	if attempts == 0 {
		return errors.New("did not find new synapse slot within nuber of attempts")
	}
	nao := rand.Intn(len(g.Neurons) - g.NumOut)
	start := g.NumIn
	if start <= nao {
		start = nao + 1
	}
	nbo := start + rand.Intn(len(g.Neurons)-start)
	if isRecurrent {
		temp := nao
		nao = nbo
		nbo = temp
	}
	_, err := g.AddSynapse(counter, g.NeuronOrder[nao], g.NeuronOrder[nbo], rand.NormFloat64()*weightStddev)
	if err != nil {
		return addRandomSynapse(counter, g, weightStddev, isRecurrent, attempts-1)
	}
	return nil
}

// Add a neuron on a random synapse of `g` with activation function `activation`
func (g *Genotype) AddRandomNeuron(counter *Counter, activation Activation) error {
	if len(g.Synapses) == 0 {
		return errors.New("no synapses to create neuron on")
	}
	k := rand.Intn(len(g.Synapses))
	for sid := range g.Synapses {
		if k <= 0 {
			// Only create on non recurrent synapse
			of, _ := g.GetNeuronOrder(g.Synapses[sid].From)
			ot, _ := g.GetNeuronOrder(g.Synapses[sid].To)
			if of < ot {
				g.AddNeuron(counter, sid, activation)
				return nil
			}
		}
		k--
	}
	// If there are only recurrent synapses, this will be the result
	return errors.New("no synapses to create neuron on") //panic("unreachable")
}

// Choose a randomly selected activation from `activations`
func ChooseActivationFrom(activations []Activation) Activation {
	return activations[rand.Intn(len(activations))]
}
