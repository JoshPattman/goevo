package goevo

// PointCrossoverable is an interface that must be implemented by any object that wants to be point crossovered.
// It should return a new object that is the crossover of two parents.
// The mutations should have no spatial importance, i.e. two genes next to each other will have re relation.
type PointCrossoverable interface {
	PointCrossoverWith(other PointCrossoverable) PointCrossoverable
}

// PointCrossover performs a point crossover between two parents.
// It is generic, so it will perform the cast back to the parent type.
func PointCrossover[T PointCrossoverable](parent1, parent2 T) T {
	return parent1.PointCrossoverWith(parent2).(T)
}
