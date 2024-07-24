package goevo

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"
)

func setupNeatTestStuff(numIn, numOut int, useRecurrent bool) Population[*NeatGenotype] {
	counter := NewCounter()

	originalGt := NewNeatGenotype(counter, numIn, numOut, Sigmoid)
	originalGt.AddRandomSynapse(counter, 0.3, false)

	selec := &TournamentSelection[*NeatGenotype]{
		TournamentSize: 3,
	}

	r := 0.0
	if useRecurrent {
		r = 0.5
	}

	mut := &NeatMutationStd{
		StdNumNewSynapses:          1,
		StdNumNewRecurrentSynapses: r,
		StdNumNewNeurons:           0.5,
		StdNumMutateSynapses:       2,
		StdNumPruneSynapses:        0,
		StdNumMutateActivations:    0.5,
		StdNewSynapseWeight:        0.2,
		StdMutateSynapseWeight:     0.4,
		Counter:                    counter,
		PossibleActivations:        AllActivations,
		MaxHiddenNeurons:           3,
	}
	crs := &NeatCrossoverSimple{}
	reprod := NewTwoPhaseReproduction(crs, mut)

	var pop Population[*NeatGenotype] = NewSimplePopulation[*NeatGenotype](func() *NeatGenotype {
		gt := Clone(originalGt)
		gt.AddRandomSynapse(counter, 0.3, false)
		return gt
	}, 100, selec, reprod)
	return pop
}

// Test the genotype constructor makes a valid genotype that can be run
func TestNeatNewGenotype(t *testing.T) {
	c := NewCounter()
	g := NewNeatGenotype(c, 10, 5, Tanh)
	assertEq(t, g.NumInputNeurons(), 10, "inputs")
	assertEq(t, g.NumOutputNeurons(), 5, "outputs")
	p := g.Build()
	outs := p.Forward(make([]float64, 10))
	assertEq(t, len(outs), 5, "output length")
}

// Check we can save and load the genotype
func TestNeatSaving(t *testing.T) {
	counter := NewCounter()
	gt := NewNeatGenotype(counter, 3, 2, Tanh)
	gt.AddRandomSynapse(counter, 0.5, false)
	gt.AddRandomSynapse(counter, 0.5, false)
	gt.AddRandomSynapse(counter, 0.5, false)
	gt.AddRandomSynapse(counter, 0.5, false)
	gt.AddRandomNeuron(counter, Tanh, Relu, Sigmoid)
	gt.AddRandomNeuron(counter, Tanh, Relu, Sigmoid)
	gt.AddRandomNeuron(counter, Tanh, Relu, Sigmoid)
	gt.AddRandomSynapse(counter, 0.5, false)
	gt.AddRandomSynapse(counter, 0.5, false)
	gt.AddRandomSynapse(counter, 0.5, false)
	gt.AddRandomSynapse(counter, 0.5, false)

	input := []float64{1, 1, 1}
	originalOutput := gt.Build().Forward(input)

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(gt); err != nil {
		t.Fatal(err)
	}
	var loadedGt *NeatGenotype
	if err := json.NewDecoder(buf).Decode(&loadedGt); err != nil {
		t.Fatal(err)
	}
	loadedOutput := loadedGt.Build().Forward(input)
	if originalOutput[0] != loadedOutput[0] || originalOutput[1] != loadedOutput[1] {
		t.Fatalf("unmatching outputs: %v and %v", loadedOutput, originalOutput)
	}
}

// Randomly perform mutation operations on a genotype to check if it remains valid
func TestNeatGenotypeStressTest(t *testing.T) {
	counter := NewCounter()
	gt := NewNeatGenotype(counter, 5, 3, Sigmoid)
	if err := gt.Validate(); err != nil {
		t.Fatalf("error after creating genotype: %v", err)
	}

	for i := 0; i < 1000; i++ {
		var op string
		cachedGt := Clone(gt)
		if err := cachedGt.Validate(); err != nil {
			t.Fatalf("error after cloning genotype: %v\nORIGINAL:\n%v\nCLONED:\n%v", err, gt, cachedGt)
		}
		opId := rand.Intn(6)
		switch opId {
		case 0:
			op = "AddFwdSynapse"
			gt.AddRandomSynapse(counter, 0.5, false)
		case 1:
			op = "AddRecSynapse"
			gt.AddRandomSynapse(counter, 0.5, true)
		case 2:
			op = "RemoveSynapse"
			gt.RemoveRandomSynapse()
		case 3:
			op = "AddNeuron"
			gt.AddRandomNeuron(counter, Relu, Tanh, Sigmoid)
		case 4:
			op = "MutateSynapse"
			gt.MutateRandomSynapse(.3)
		case 5:
			op = "MutateActivation"
			gt.MutateRandomActivation(Relu, Tanh, Sigmoid)
		}
		if err := gt.Validate(); err != nil {
			t.Fatalf("error after performing %v op on genotype: %v\nBEFORE:\n%v\nAFTER:\n%v", err, op, cachedGt, gt)
		}
	}
}

func TestNeatXOR(t *testing.T) {
	testWithXORDataset(t, setupNeatTestStuff(3, 1, false), nil)
}

func TestNeatReccurrent(t *testing.T) {
	testWithRecurrentDataset(t, setupNeatTestStuff(1, 1, true), nil)
}
