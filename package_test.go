package goevo

import (
	"testing"
)

func assertEq[T comparable](t *testing.T, a T, b T, name string) {
	if a != b {
		t.Fatalf("error in check '%s' (not equal): '%v' and '%v'", name, a, b)
	}
}

func TestNewGenotype(t *testing.T) {
	c := NewCounter()
	g := NewGenotype(c, 10, 5, Tanh)
	assertEq(t, g.numInputs, 10, "inputs")
	assertEq(t, g.numOutputs, 5, "outputs")
	p := g.Build()
	outs := p.Forward(make([]float64, 10))
	assertEq(t, len(outs), 5, "output length")
}
