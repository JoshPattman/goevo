package goevo

// Buildable is an interface that defines a method to build a [Forwarder].
// If the genotype is already a [Forwarder], it should return itself.
type Buildable interface {
	// Build converts this genotype into a [Forwarder].
	// This may be an expensive operation, so it should be called sparingly.
	Build() Forwarder
}
