package goevo

import (
	"errors"
	"math/rand"
)

func MutateRandomSynapse(g *Genotype, stddev float64) {
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

// I don't knwo an efficient way to pick a deffo available random synapse to create so i just repeatedly try many times over
func AddRandomSynapse(counter Counter, g *Genotype, weightStddev float64, attempts int) error {
	if attempts == 0 {
		return errors.New("did not find new synapse slot within nuber of attempts")
	}
	nao := rand.Intn(len(g.Neurons) - g.NumOut)
	start := g.NumIn
	if start <= nao {
		start = nao + 1
	}
	nbo := start + rand.Intn(len(g.Neurons)-start)
	_, err := g.NewSynapse(counter, g.NeuronOrder[nao], g.NeuronOrder[nbo], rand.NormFloat64()*weightStddev)
	if err != nil {
		return AddRandomSynapse(counter, g, weightStddev, attempts-1)
	}
	return nil
}

func AddRandomNeuron(counter Counter, g *Genotype, activation Activation) error {
	if len(g.Synapses) == 0 {
		return errors.New("no synapses to create neuron on")
	}
	k := rand.Intn(len(g.Synapses))
	for sid := range g.Synapses {
		if k == 0 {
			g.NewNeuron(counter, sid, activation)
			return nil
		}
		k--
	}
	panic("unreachable")
}
