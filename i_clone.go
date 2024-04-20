package goevo

// Cloneable is an interface that must be implemented by any object that wants to be cloned.
// The clone method must return a new object that is a deep copy of the original object.
// This new method is typed as any. To clone while including the type, use [Clone], which is generic so will perform the cast.
type Cloneable interface {
	Clone() any
}

// Clone clones an object that implements the [Cloneable] interface.
// It also casts the child object to the type of the parent object.
func Clone[T Cloneable](obj T) T {
	return obj.Clone().(T)
}
