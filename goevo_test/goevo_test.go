package goevo_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/JoshPattman/goevo"
	"github.com/JoshPattman/goevo/geno/arr"
	"github.com/JoshPattman/goevo/geno/neat"
	"github.com/JoshPattman/goevo/pop/simple"
	"github.com/JoshPattman/goevo/pop/speciated"
	"github.com/JoshPattman/goevo/selec/tournament"
)

func assertEq[T comparable](t *testing.T, a T, b T, name string) {
	if a != b {
		t.Fatalf("error in check '%s' (not equal): '%v' and '%v'", name, a, b)
	}
}

// Test the genotype constructor makes a valid genotype that can be run
func TestNewGenotype(t *testing.T) {
	c := goevo.NewCounter()
	g := neat.NewGenotype(c, 10, 5, goevo.Tanh)
	assertEq(t, g.NumInputNeurons(), 10, "inputs")
	assertEq(t, g.NumOutputNeurons(), 5, "outputs")
	p := g.Build()
	outs := p.Forward(make([]float64, 10))
	assertEq(t, len(outs), 5, "output length")
}

// Test training a new genotype on the XOR problem, fail if do not solve the problem
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

	fitness := func(f goevo.Forwarder) float64 {
		fitness := 0.0
		for i := range X {
			pred := f.Forward(X[i])
			e := pred[0] - Y[i][0]
			fitness -= math.Pow(math.Abs(e), 3)
		}
		return fitness
	}

	counter := goevo.NewCounter()

	originalGt := neat.NewGenotype(counter, 3, 1, goevo.Sigmoid)
	neat.AddRandomSynapse(originalGt, counter, 0.3, false)

	selec := &tournament.Selection[*neat.Genotype]{
		TournamentSize: 3,
	}

	mut := &neat.StdMutation{
		StdNumNewSynapses:       1,
		StdNumNewNeurons:        0.5,
		StdNumMutateSynapses:    2,
		StdNumPruneSynapses:     0,
		StdNumMutateActivations: 0.5,
		StdNewSynapseWeight:     0.2,
		StdMutateSynapseWeight:  0.4,
		Counter:                 counter,
		PossibleActivations:     goevo.AllActivations,
		MaxHiddenNeurons:        3,
	}
	crs := &neat.SimpleCrossoverStrategy{}
	reprod := goevo.NewTwoPhaseReproduction(crs, mut)

	var pop goevo.Population[*neat.Genotype] = simple.NewPopulation[*neat.Genotype](func() *neat.Genotype {
		gt := goevo.Clone(originalGt)
		neat.AddRandomSynapse(gt, counter, 0.3, false)
		return gt
	}, 100, selec, reprod)

	var maxFitness float64
	var maxGt *neat.Genotype
	debugging := false
	for gen := 0; gen < 5000; gen++ {
		maxFitness = math.Inf(-1)
		maxGt = nil
		for _, a := range pop.All() {
			a.Fitness = fitness(a.Genotype.Build())
			if a.Fitness > maxFitness {
				maxFitness = a.Fitness
				maxGt = a.Genotype
			}
		}
		if maxFitness > -0.4 {
			mut.StdNumPruneSynapses = 0.5
		} else {
			mut.StdNumPruneSynapses = 0
		}

		if maxFitness > -0.1 {
			break
		}
		// Below is only for debug
		if debugging && gen%100 == 0 {
			fmt.Printf("Generation %v: Max fitness %v\n", gen, maxFitness)
		}

		pop = pop.NextGeneration()
	}
	if debugging {
		maxPt := maxGt.Build()
		fmt.Println(maxPt.Forward([]float64{1}))
		fmt.Println(maxPt.Forward([]float64{1}))
		fmt.Println(maxPt.Forward([]float64{1}))
		fmt.Println(maxPt.Forward([]float64{1}))
		bs, _ := json.MarshalIndent(maxGt, "", "\t")
		fmt.Println(string(bs))
		func() {
			f, _ := os.Create("img.png")
			defer f.Close()
			png.Encode(f, maxGt.RenderImage(20, 10))
		}()
	}
	if maxFitness < -0.1 {
		t.Fatalf("XOR Failed to converge, ending with fitness %f", maxFitness)
	}
	if err := maxGt.Validate(); err != nil {
		t.Fatalf("final genotype was not valid: %v\nGenotype:\n%v", err, maxGt)
	}
}

// Check we can save and load the genotype
func TestSaving(t *testing.T) {
	counter := goevo.NewCounter()
	gt := neat.NewGenotype(counter, 3, 2, goevo.Tanh)
	neat.AddRandomSynapse(gt, counter, 0.5, false)
	neat.AddRandomSynapse(gt, counter, 0.5, false)
	neat.AddRandomSynapse(gt, counter, 0.5, false)
	neat.AddRandomSynapse(gt, counter, 0.5, false)
	neat.AddRandomNeuron(gt, counter, goevo.Tanh, goevo.Relu, goevo.Sigmoid)
	neat.AddRandomNeuron(gt, counter, goevo.Tanh, goevo.Relu, goevo.Sigmoid)
	neat.AddRandomNeuron(gt, counter, goevo.Tanh, goevo.Relu, goevo.Sigmoid)
	neat.AddRandomSynapse(gt, counter, 0.5, false)
	neat.AddRandomSynapse(gt, counter, 0.5, false)
	neat.AddRandomSynapse(gt, counter, 0.5, false)
	neat.AddRandomSynapse(gt, counter, 0.5, false)

	input := []float64{1, 1, 1}
	originalOutput := gt.Build().Forward(input)

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(gt); err != nil {
		t.Fatal(err)
	}
	var loadedGt *neat.Genotype
	if err := json.NewDecoder(buf).Decode(&loadedGt); err != nil {
		t.Fatal(err)
	}
	loadedOutput := loadedGt.Build().Forward(input)
	if originalOutput[0] != loadedOutput[0] || originalOutput[1] != loadedOutput[1] {
		t.Fatalf("unmatching outputs: %v and %v", loadedOutput, originalOutput)
	}
}

// Randomly perform mutation operations on a genotype to check if it remains valid
func TestGenotypeStressTest(t *testing.T) {
	counter := goevo.NewCounter()
	gt := neat.NewGenotype(counter, 5, 3, goevo.Sigmoid)
	if err := gt.Validate(); err != nil {
		t.Fatalf("error after creating genotype: %v", err)
	}

	for i := 0; i < 5000; i++ {
		var op string
		cachedGt := goevo.Clone(gt)
		if err := cachedGt.Validate(); err != nil {
			t.Fatalf("error after cloning genotype: %v\nORIGINAL:\n%v\nCLONED:\n%v", err, gt, cachedGt)
		}
		opId := rand.Intn(6)
		switch opId {
		case 0:
			op = "AddFwdSynapse"
			neat.AddRandomSynapse(gt, counter, 0.5, false)
		case 1:
			op = "AddRecSynapse"
			neat.AddRandomSynapse(gt, counter, 0.5, true)
		case 2:
			op = "RemoveSynapse"
			neat.RemoveRandomSynapse(gt)
		case 3:
			op = "AddNeuron"
			neat.AddRandomNeuron(gt, counter, goevo.Relu, goevo.Tanh, goevo.Sigmoid)
		case 4:
			op = "MutateSynapse"
			neat.MutateRandomSynapse(gt, 0.3)
		case 5:
			op = "MutateActivation"
			neat.MutateRandomActivation(gt, goevo.Relu, goevo.Tanh, goevo.Sigmoid)
		}
		if err := gt.Validate(); err != nil {
			t.Fatalf("error after performing %v op on genotype: %v\nBEFORE:\n%v\nAFTER:\n%v", err, op, cachedGt, gt)
		}
	}
}

// Test recurrent connections can evolve to remeber the sequence 0, 1, 1, 0 (with the same input each time)
func TestRecurrency(t *testing.T) {
	// We add a bias on the end of each input, which is always 1
	X := [][]float64{
		{1},
		{1},
		{1},
		{1},
	}
	Y := [][]float64{
		{0},
		{1},
		{1},
		{0},
	}

	fitness := func(f goevo.Forwarder) float64 {
		fitness := 0.0
		for i := range X {
			pred := f.Forward(X[i])
			e := pred[0] - Y[i][0]
			fitness -= math.Pow(math.Abs(e), 3)
		}
		return fitness
	}

	counter := goevo.NewCounter()

	originalGt := neat.NewGenotype(counter, 1, 1, goevo.Sigmoid)
	neat.AddRandomSynapse(originalGt, counter, 0.3, false)

	selec := &tournament.Selection[*neat.Genotype]{
		TournamentSize: 3,
	}

	mut := &neat.StdMutation{
		StdNumNewSynapses:          1,
		StdNumNewRecurrentSynapses: 0.5,
		StdNumNewNeurons:           0.5,
		StdNumMutateSynapses:       2,
		StdNumPruneSynapses:        0,
		StdNumMutateActivations:    0.5,
		StdNewSynapseWeight:        0.2,
		StdMutateSynapseWeight:     0.4,
		Counter:                    counter,
		PossibleActivations:        goevo.AllActivations,
		MaxHiddenNeurons:           3,
	}
	crs := &neat.SimpleCrossoverStrategy{}
	reprod := goevo.NewTwoPhaseReproduction(crs, mut)

	var pop goevo.Population[*neat.Genotype] = simple.NewPopulation[*neat.Genotype](func() *neat.Genotype {
		gt := goevo.Clone(originalGt)
		neat.AddRandomSynapse(gt, counter, 0.3, false)
		return gt
	}, 100, selec, reprod)

	var maxFitness float64
	var maxGt *neat.Genotype
	debugging := false
	for gen := 0; gen < 5000; gen++ {
		maxFitness = math.Inf(-1)
		maxGt = nil
		for _, a := range pop.All() {
			a.Fitness = fitness(a.Genotype.Build())
			if a.Fitness > maxFitness {
				maxFitness = a.Fitness
				maxGt = a.Genotype
			}
		}
		if maxFitness > -0.4 {
			mut.StdNumPruneSynapses = 0.5
		} else {
			mut.StdNumPruneSynapses = 0
		}

		if maxFitness > -0.1 {
			break
		}
		// Below is only for debug
		if debugging && gen%100 == 0 {
			fmt.Printf("Generation %v: Max fitness %v\n", gen, maxFitness)
		}

		pop = pop.NextGeneration()
	}
	if debugging {
		maxPt := maxGt.Build()
		fmt.Println(maxPt.Forward([]float64{1}))
		fmt.Println(maxPt.Forward([]float64{1}))
		fmt.Println(maxPt.Forward([]float64{1}))
		fmt.Println(maxPt.Forward([]float64{1}))
		bs, _ := json.MarshalIndent(maxGt, "", "\t")
		fmt.Println(string(bs))
		func() {
			f, _ := os.Create("img.png")
			defer f.Close()
			png.Encode(f, maxGt.RenderImage(20, 10))
		}()
	}
	if maxFitness < -0.1 {
		t.Fatalf("Recurrency Failed to converge, ending with fitness %f", maxFitness)
	}
	if err := maxGt.Validate(); err != nil {
		t.Fatalf("final genotype was not valid: %v\nGenotype:\n%v", err, maxGt)
	}
}

func TestFloatsGt(t *testing.T) {
	counter := goevo.NewCounter()
	mut := &arr.StdMutationStrategy[float64]{
		MutateProbability: 0.1,
		MutateStd:         0.05,
	}
	crs := &arr.KPointCrossoverStrategy[float64]{K: 2}
	reprod := goevo.NewTwoPhaseReproduction(crs, mut)
	selec := &tournament.Selection[*arr.Genotype[float64]]{
		TournamentSize: 3,
	}
	var pop goevo.Population[*arr.Genotype[float64]] = speciated.NewPopulation(
		counter,
		func() *arr.Genotype[float64] { return arr.NewFloatGenotype(10, 0.5) },
		5,
		20,
		0.1,
		2.5,
		selec,
		reprod,
	)
	// Fitness is max (0) when all the numbers sum to 10
	fitness := func(f *arr.Genotype[float64]) float64 {
		total := 0.0
		for i := range f.Values {
			total += f.Values[i]
		}
		return -math.Abs(10 - total)
	}
	var highestFitness float64
	for gen := 0; gen < 100; gen++ {
		highestFitness = math.Inf(-1)
		for _, a := range pop.All() {
			a.Fitness = fitness(a.Genotype)
			if a.Fitness > highestFitness {
				highestFitness = a.Fitness
			}
		}
		pop = pop.NextGeneration()
	}
	if highestFitness < -0.1 {
		t.Fatalf("Failed to converge, ending with fitness %f", highestFitness)
	}
}
