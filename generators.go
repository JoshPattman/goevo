package goevo

import (
	"fmt"
	"math/rand/v2"
)

// Generator is an interface the can create objects of type T.
type Generator[T any] interface {
	Next() T
}

// NormalGenerator generates numbers in a normal distribution then casts them to the type.
type NormalGenerator[T floatType] struct {
	Mean float64
	Std  float64
}

// Next implements Generator.
func (s *NormalGenerator[T]) Next() T {
	MustValidate(s)
	v := rand.NormFloat64()*s.Std + s.Mean
	return T(v)
}

func (s *NormalGenerator[T]) Validate() error {
	if s.Mean < 0 {
		return fmt.Errorf("cannot have mean less than 0 (%.3f)", s.Std)
	}
	return nil
}

// ChoiceGenerator is a generator that randomly chooses one of the values
type ChoiceGenerator[T any] struct {
	Choices []T
}

// Next implements Generator.
func (c *ChoiceGenerator[T]) Next() T {
	MustValidate(c)
	return c.Choices[rand.IntN(len(c.Choices))]
}

func (c *ChoiceGenerator[T]) Validate() error {
	if len(c.Choices) == 0 {
		return fmt.Errorf("must have at least one choice for choices generator")
	}
	return nil
}
