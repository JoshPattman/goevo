package goevo

type Reproduction interface {
	Reproduce(a, b *Genotype) *Genotype
}
