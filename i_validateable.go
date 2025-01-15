package goevo

import "fmt"

// Validateable is an interface for types that can be validated.
// For example, you can check if a neat genotype actually contains all neurons that are used by each synapse.
// This is aimed to be used for either tests, or validating loaded data.
type Validateable interface {
	// Validate checks if the object is valid.
	// If the object is valid, it should return nil.
	// If the object is invalid, it should return an error.
	Validate() error
}

func MustValidate(v Validateable) {
	err := v.Validate()
	if err != nil {
		panic(fmt.Errorf("failed to validate: %v", err))
	}
}
