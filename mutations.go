package goevo

import (
	"errors"
	"math/rand"
)

// Mutate a random synapses weight by normal distribution of standard deviation stddev
func MutateRandomSynapse(g *Genotype, stddev float64) {
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

// Add a random synapse with weight from normal distribution with standard deviation weightStddev
func AddRandomSynapse(counter Counter, g *Genotype, weightStddev float64, isRecurrent bool, attempts int) error {
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
		return AddRandomSynapse(counter, g, weightStddev, isRecurrent, attempts-1)
	}
	return nil
}

// Add a neuron on a random synapse
func AddRandomNeuron(counter Counter, g *Genotype, activation Activation) error {
	if len(g.Synapses) == 0 {
		return errors.New("no synapses to create neuron on")
	}
	k := rand.Intn(len(g.Synapses))
	for sid := range g.Synapses {
		if k == 0 {
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
	panic("unreachable")
}
