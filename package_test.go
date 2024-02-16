package goevo

import (
	"math"
	"testing"
)

func assertEq[T comparable](t *testing.T, a T, b T, name string) {
	if a != b {
		t.Fatalf("error in check '%s' (not equal): '%v' and '%v'", name, a, b)
	}
}

func TestNewGenotype(t *testing.T) {
	c := NewCounter()
	g := NewGenotype(c, 10, 5, Tanh)
	assertEq(t, g.numInputs, 10, "inputs")
	assertEq(t, g.numOutputs, 5, "outputs")
	p := g.Build()
	outs := p.Forward(make([]float64, 10))
	assertEq(t, len(outs), 5, "output length")
}

func TestXOR(t *testing.T) {
	// We add a bias on the end of each input, which is always 1
	X := [][]float64{
		{0, 0, 1},
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
	}
	Y := [][]float64{
		{0},
		{1},
		{1},
		{0},
	}

	fitness := func(f Forwarder) float64 {
		fitness := 0.0
		for i := range X {
			pred := f.Forward(X[i])
			e := pred[0] - Y[i][0]
			fitness -= math.Pow(math.Abs(e), 3)
		}
		return fitness
	}

	counter := NewCounter()

	originalGt := NewGenotype(counter, 3, 1, Sigmoid)
	originalGt.AddRandomSynapse(counter, 0.3, false)
	pop := NewSimplePopulation(func() *Genotype {
		gt := originalGt.Clone()
		gt.AddRandomSynapse(counter, 0.3, false)
		return gt
	}, 100)

	selec := &TournamentSelection{
		TournamentSize: 3,
	}
	reprod := &StdReproduction{
		StdNumNewSynapses:       1,
		StdNumNewNeurons:        0.5,
		StdNumMutateSynapses:    2,
		StdNumPruneSynapses:     0,
		StdNumMutateActivations: 0.5,
		StdNewSynapseWeight:     0.2,
		StdMutateSynapseWeight:  0.4,
		Counter:                 counter,
		PossibleActivations:     []Activation{Relu, Tanh, Sigmoid, Sin, Cos},
		MaxHiddenNeurons:        3,
	}
	var maxFitness float64
	for gen := 0; gen < 5000; gen++ {
		maxFitness = math.Inf(-1)
		//var maxGt *Genotype
		for _, a := range pop.Agents() {
			a.Fitness = fitness(a.Genotype.Build())
			if a.Fitness > maxFitness {
				maxFitness = a.Fitness
				//maxGt = a.Genotype
			}
		}
		if maxFitness > -0.4 {
			reprod.StdNumPruneSynapses = 0.5
		} else {
			reprod.StdNumPruneSynapses = 0
		}
		// Below is only for debug
		/*if gen%100 == 0 {
			fmt.Printf("Generation %v: Max fitness %v\n", gen, maxFitness)
			func() {
				f, _ := os.Create("img.png")
				defer f.Close()
				png.Encode(f, maxGt.Draw(20, 10))
			}()
		}*/

		pop = pop.NextGeneration(selec, reprod)
	}
	if maxFitness < -0.1 {
		t.Fatalf("XOR Failed to converge, ending with fitness %f", maxFitness)
	}
}
