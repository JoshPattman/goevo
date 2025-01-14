package goevo

import "math/rand/v2"

type Generator[T any] interface {
	Next() T
}

var _ Generator[float64] = &NormalGenerator[float64]{}

// NormalGenerator generates numbers in a normal distribution then casts them to the type.
type NormalGenerator[T floatType] struct {
	Mean float64
	Std  float64
}

// Next implements Generator.
func (s *NormalGenerator[T]) Next() T {
	v := rand.NormFloat64()*s.Std + s.Mean
	return T(v)
}
