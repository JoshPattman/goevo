package goevo

import "math/rand/v2"

type Generator[T any] interface {
	Next() T
}

type generatorNormal[T floatType] struct {
	mean float64
	std  float64
}

// NewGeneratorNormal creates a new [Generator] that generates floating point numbers
// within a normal distribution.
func NewGeneratorNormal[T floatType](mean, std T) Generator[T] {
	if std < 0 {
		panic("cannot have std < 0")
	}
	return &generatorNormal[T]{
		mean: float64(mean),
		std:  float64(std),
	}
}

func (s *generatorNormal[T]) Next() T {
	v := rand.NormFloat64()*s.std + s.mean
	return T(v)
}

type generatorChoice[T any] struct {
	choices []T
}

// NewGeneratorChoices creates a new [Generator] that chooses values from
// the given choices slice.
func NewGeneratorChoices[T any](choices []T) Generator[T] {
	if len(choices) == 0 {
		panic("cannot have no choices")
	}
	return &generatorChoice[T]{
		choices: choices,
	}
}

func (c *generatorChoice[T]) Next() T {
	return c.choices[rand.IntN(len(c.choices))]
}
