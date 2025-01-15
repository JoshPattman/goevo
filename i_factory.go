package goevo

// Factory is an interface that defines somthing that can create new things of type T.
// This is primarily used for genotypes.
type Factory[T any] interface {
	New() T
}

// ValidateableFactory is an interface that incudes [Factory] and [Validateable].
// Most genotype factories will implement this.
type ValidateableFactory[T any] interface {
	Validateable
	Factory[T]
}
