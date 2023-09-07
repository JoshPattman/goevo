# `goevo` - work-in-progress NEAT implementation in Golang
GoEVO is designed to be a fast but easy-to-understand package that implements the NEAT algorithm. It is still in development and has not had a major release yet, so stability is not guaranteed. If you find a bug or have any suggestions, please do raise an issue and i'll try to fix it. \
To learn more about the NEAT algorithm, here is the original paper: [Stanley, K. O., & Miikkulainen, R. (2002). Evolving neural networks through augmenting topologies. Evolutionary computation, 10(2), 99-127.](https://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf)
## Usage
### Creating and Modifying a `Genotype`
A Genotype is a bit like DNA - it encodes all the information to build the network.
```go
// Create a counter. This is used to keep track of new neurons and synapses
counter := goevo.NewAtomicCounter()

// Create the initial genotype. This is sort of like an agents DNA.
genotype := goevo.NewGenotype(counter, 2, 1, goevo.ActivationLinear, goevo.ActivationSigmoid)

// Add a synapse from first input to output with a weight of 0.5.
// We don't need to check the error because we know this is a valid modification
synapseID, _ := genotype.AddSynapse(counter, 0, 2, 0.5)
// Add a neuron on the synapse we just created
neuronID, secondSynapseID, _ := genotype.AddNeuron(counter, synapseID, goevo.ActivationReLU)
// Add a synapse between the second input and our new neuron with a weight of -0.5
genotype.AddSynapse(counter, 1, neuronID, -0.5)
```
### Visualising a `Genotype`
It is quite hard to deduce the topology of a genotype by looking at a list of its neurons and synapses. `goevo` supports drawing a picture of a genotype either to a `draw.Image` or a `png` or `jpg` file.
```go
vis := goevo.NewGenotypeVisualiser()
vis.DrawImageToPNGFile("example_1.png", genotype)
```
Below is the image in the generated file `example_1.png`. The green cirlces are input neurons, pink circles are hidden neurons, and yellow circles are output neurons. A blue line is a positive weight and a red line is a negative weight. The thicker the line, the stronger the weight.

<img src="README_ASSETS/example_1.png" width="400">

### Pruning Synapses
One way to prevent the networks getting too big is to prune synapses (delete synapses). Pruning will remove the given synapse, then remove all neurons and synapses that become redundant due to the pruning.
```go
// Prune the synapse that connects the hidden neuron to the output neuron. This makes the hidden neuron nedundant so it is therefor removed too, along with its other synapses.
genotypePrunedA := goevo.NewGenotypeCopy(genotype)
genotypePrunedA.PruneSynapse(secondSynapseID)
vis.DrawImageToPNGFile("example_2.png", genotypePrunedA)
// Prune the synapse that connects the first input to the hidden neuron
genotypePrunedB := goevo.NewGenotypeCopy(genotype)
genotypePrunedB.PruneSynapse(synapseID)
vis.DrawImageToPNGFile("example_3.png", genotypePrunedB)
```

Below is `example_2.png`. Because the removed synapse was the only synapse connecting the hidden node to an output node, the hidden node was removed with the synapse, along with its synapses.

<img src="README_ASSETS/example_2.png" width="400">

Below is `example_3.png`. The hidden node still had a connection going in and a connection going out after this pruning so it was not removed.

<img src="README_ASSETS/example_3.png" width="400">

### Using the `Genotype`
A Genotype cannot directly be used to convert an input into an output. Instead, it must first be converted into a Phenotype, which can be thought of a bit like compiling code into an executable. After a Phenotype is created, it can be used as many times as you like, but the neurons and synapses cannot be edited (to do this you have to modify the genotype and create a new phenotype).
```go
// Create a phenotype from the genotype
phenotype := goevo.NewPhenotype(genotype)
// Calculate the outputs given some inputs
fmt.Println(phenotype.Forward([]float64{0, 1}))
fmt.Println(phenotype.Forward([]float64{1, 0}))
// Make sure to clear any recurrent connections memory (in this case there are no recurrent connections but this is just an exmaple)
phenotype.ClearRecurrentMemory()
```
Output:
```
[0.5]
[0.6224593312018546]
```

### Saving and Loading `Genotype`
If you have just trained a genotype, you may wish to save it. Genotypes can be json marshalled und unmarshalled with go's built-in json parser.
```go
// Convert the genotype to a json []byte
jsBytes, _ := json.Marshal(genotype)
// Create an empty genotype and load the json []byte into it
genotypeLoaded := goevo.NewGenotypeEmpty()
json.Unmarshal(jsBytes, genotypeLoaded)
```

## Example - XOR
In this example, a population of agents attempts to create a genotype that can do XOR logic. Note that each part of the NEAT algorithm is run seperately, meaning that the training loop is very easy to customise to your desires.

```go
// Define a counter (for counting new neurons and synapses), and a species counter (for new species)
counter, specCounter := goevo.NewAtomicCounter(), goevo.NewAtomicCounter()

// Define all of the possible activations that can be used
possibleActivations := []goevo.Activation{
    goevo.ActivationReLU,
    goevo.ActivationSigmoid,
    goevo.ActivationStep,
}

// This is our input and output
X := [][]float64{
    {0, 0},
    {0, 1},
    {1, 0},
    {1, 1},
}
Y := [][]float64{
    {0},
    {1},
    {1},
    {0},
}

// Our fitness function calculates the mean squared error between the predictions and the actual values.
fitness := func(f *goevo.Phenotype) float64 {
    loss := 0.0
    for i := range X {
        pred := f.Forward(X[i])
        e := pred[0] - Y[i][0]
        loss += e * e
    }
    return (1 - loss/4)
}

// Our reproduction function takes two parent genotypes (DNA) and creates a child genotype.
reproduction := func(g1, g2 *goevo.Genotype) *goevo.Genotype {
    // Perform crossover on the two parents to get a child genotype.
    g := goevo.NewGenotypeCrossover(g1, g2)
    // 10% of the time, add a random synapse to the child. We are not using any recurrent synapses
    if rand.Float64() > 0.9 {
        goevo.AddRandomSynapse(counter, g, 0.1, false, 5)
    }
    // 5% of the time, add a random neuron with one of our previously specified activation functions
    if rand.Float64() > 0.95 {
        goevo.AddRandomNeuron(counter, g, goevo.ChooseActivationFrom(possibleActivations))
    }
    // 5% of the time, delete a random synapse. This will then recursively delete neurons and synapses that have no inputs or outputs.
    if rand.Float64() > 0.95 {
        goevo.PruneRandomSynapse(g)
    }
    // 3 times, with a 15% chance each, mutate the weight of a random synapse
    for i := 0; i < 3; i++ {
        if rand.Float64() > 0.85 {
            goevo.MutateRandomSynapse(g, 0.1)
        }
    }
    return g
}

// The target population size will be 100
targetPopSize := 100
// Create a population of 100 agents that are all empty genotypes. An agent contains a genotype, fitness, and species ID
population := make([]*goevo.Agent, targetPopSize)
// Fill the population with the same initial genotype.
// It is very important that all genotypes start from one initial one, as this means their input and output nodes all have the same IDs.
// This immidiately makes crossover feasable.
initialGenotype := goevo.NewGenotype(counter, 2, 1, goevo.ActivationLinear, goevo.ActivationSigmoid)
for p := range population {
    population[p] = goevo.NewAgent(goevo.NewGenotypeCopy(initialGenotype))
}

// Distance threshold is the maximum distance two genotypes can be from each other before being considered different species.
// During the generational loop, we will change this to hit the target species number.
distanceThreshold := 1.0
// The target number of species we want in a generation
targetSpecies := 10

// Used to remember the best network up to this point
var bestNet *goevo.Genotype
var bestFitness float64

// For 500 generations
for gen := 0; gen < 500; gen++ {
    // Calculate the fitness of each agent, while keeping track of the best
    bestFitness = -math.MaxFloat64
    for _, a := range population {
        pheno := goevo.NewPhenotype(a.Genotype)
        a.Fitness = fitness(pheno)
        if a.Fitness > bestFitness {
            bestFitness = a.Fitness
            bestNet = a.Genotype
        }
    }
    // Split the population into species. We will use the genetic distance function to calculate the distance between two agents.
    // The weightings are the same as in the original paper.
    specPop := goevo.Speciate(specCounter, population, distanceThreshold, false, goevo.GeneticDistance(1, 0.4))
    // Change the distance threshold so that next generation, there should be closer to the target species number of species
    if len(specPop) > targetSpecies {
        distanceThreshold *= 1.2
    } else if len(specPop) < targetSpecies {
        distanceThreshold /= 1.2
    }
    // Calculate how many offspring each species is allowed
    allowedOffspring := goevo.CalculateOffspring(specPop, targetPopSize)
    // Using the previous speciated population, the number of allowed offspring, and ProbabilisticSelection, create a new generation
    population = goevo.Repopulate(specPop, allowedOffspring, reproduction, goevo.ProbabilisticSelection)

    // Print some info
    if gen%20 == 0 {
        fmt.Println("Gen", gen, ", most fit", bestFitness, ", num spec", len(specPop))
    }
}

// Print out the best network and its results
bestP := goevo.NewPhenotype(bestNet)
for i := range X {
    var pred []float64
    pred = bestP.Forward(X[i])
    fmt.Println("X", X[i], "YP", pred[0], "Y", Y[i][0])
}

// Draw the network to a file
vis := goevo.NewGenotypeVisualiser()
vis.DrawImageToPNGFile("xor.png", bestNet)
```

Here is the ouput from running this once:
```
X [0 0] YP 0.0024126496048012813 Y 0
X [0 1] YP 0.9999999999951001 Y 1
X [1 0] YP 0.9996491462464224 Y 1
X [1 1] YP 1.9287480105160926e-05 Y 0
```

The final network is rather large compared to the minimal network required for XOR. However, the algorithm can be made to create smaller networks by adding a penalty for large networks in the training function. This is a very good idea to do if you plan to use NEAT, as the whole purpose of it is to create optimal topologies.

<img src="README_ASSETS/xor.png" width="400">

## TODO
- Add a function to remove a neuron and re-route its synapses
- Add a function to change a neurons activation
- For the above two, add functions to randomly perform those actions
- Finish and merge the HyperNEAT branch. It currently somewhat works, but it does not have all of the features that I would like yet.
