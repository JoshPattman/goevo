package goevo

// Buildable is an interface that defines a method to build a Forwarder.
// If the genotype is already a forwarder, it should return itself.
type Buildable interface {
	Build() Forwarder
}
