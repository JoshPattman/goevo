package goevo

type Agent[T any] struct {
	Genotype T
	Fitness  float64
}

func NewAgent[T any](gt T) *Agent[T] {
	return &Agent[T]{Genotype: gt}
}
