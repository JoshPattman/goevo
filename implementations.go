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
var _ Crossover[*ArrayGenotype[any]] = &ArrayCrossoverUniform[any]{}
var _ Crossover[*ArrayGenotype[any]] = &ArrayCrossoverAsexual[any]{}
var _ Crossover[*ArrayGenotype[any]] = &ArrayCrossoverKPoint[any]{}
var _ Mutation[*ArrayGenotype[bool]] = &ArrayMutationRandomBool{}
var _ Mutation[*ArrayGenotype[rune]] = &ArrayMutationRandomRune{}
var _ Mutation[*ArrayGenotype[float64]] = &ArrayMutationStd[float64]{}

// Dense genotypes
var _ Cloneable = &DenseGenotype{}
var _ Forwarder = &DenseGenotype{}
var _ Crossover[*DenseGenotype] = &DenseCrossoverUniform{}
var _ Mutation[*DenseGenotype] = &DenseMutationStd{}

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
