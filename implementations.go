package goevo

// ================================== Utilities ==================================

// Generators
var _ Generator[float64] = NewGeneratorNormal(0.0, 0.0)
var _ Generator[rune] = NewGeneratorChoices([]rune("abcdefg"))

// Reproductions
var _ Reproduction[any] = &TwoPhaseReproduction[any]{}

// ================================== Genotypes ==================================

// Array genotypes
var _ Cloneable = &ArrayGenotype[int]{}
var _ Crossover[*ArrayGenotype[any]] = NewArrayCrossoverUniform[any]()
var _ Crossover[*ArrayGenotype[any]] = NewArrayCrossoverAsexual[any]()
var _ Crossover[*ArrayGenotype[any]] = NewArrayCrossoverKPoint[any](0)
var _ Mutation[*ArrayGenotype[float64]] = NewArrayMutationGeneratorAdd(NewGeneratorNormal(0.0, 0.0), 0.0)
var _ Mutation[*ArrayGenotype[bool]] = NewArrayMutationGeneratorReplace(NewGeneratorChoices([]bool{true, false}), 0.0)
var _ Mutation[*ArrayGenotype[bool]] = NewArrayMutationGenerator(NewGeneratorChoices([]bool{true, false}), func(old, new bool) bool { return old && new }, 0.0)

// Dense genotypes
var _ Cloneable = &DenseGenotype{}
var _ Forwarder = &DenseGenotype{}
var _ Crossover[*DenseGenotype] = &denseCrossoverUniform{}
var _ Mutation[*DenseGenotype] = &denseMutationUniform{}

// NEAT genotypes + phenotypes
var _ Cloneable = &NeatGenotype{}
var _ Buildable = &NeatGenotype{}
var _ Forwarder = &NeatPhenotype{}
var _ Crossover[*NeatGenotype] = &NeatCrossoverSimple{}
var _ Crossover[*NeatGenotype] = &NeatCrossoverAsexual{}
var _ Mutation[*NeatGenotype] = &NeatMutationStd{}

// ================================== Selections ==================================

// Elite selection
var _ Selection[any] = &EliteSelection[any]{}

// Tournament selection
var _ Selection[any] = &TournamentSelection[any]{}

// ================================== Populations ==================================

// Simple population
var _ Population[any] = &SimplePopulation[any]{}

// Speiated population
var _ Population[any] = &SpeciatedPopulation[any]{}

// Hill climber population
var _ Population[any] = &HillClimberPopulation[any]{}
