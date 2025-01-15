package goevo

// Factory is an interface that defines somthing that can create new things of type T.
// This is primarily used for genotypes.
type Factory[T any] interface {
	New() T
}
