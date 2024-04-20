package goevo

// GeneticDistance is an interface for calculating the genetic distance between two genotypes.
type GeneticDistance[T any] interface {
	DistanceBetween(a, b T) float64
}
