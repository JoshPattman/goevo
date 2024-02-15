package goevo

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type Activation int

const (
	Relu Activation = iota
	Linear
	Sigmoid
	Tanh
	Sin
	Cos
)

func (a Activation) String() string {
	switch a {
	case Relu:
		return "relu"
	case Linear:
		return "linear"
	case Sigmoid:
		return "sigmoid"
	case Tanh:
		return "tanh"
	case Sin:
		return "sin"
	case Cos:
		return "cos"
	}
	panic("unknown activation")
}

func activate(x float64, a Activation) float64 {
	switch a {
	case Relu:
		if x < 0 {
			return 0
		}
		return x
	case Linear:
		return x
	case Sigmoid:
		return 1 / (1 + math.Exp(-x))
	case Tanh:
		return math.Tanh(x)
	case Sin:
		return math.Sin(x)
	case Cos:
		return math.Cos(x)
	}
	panic("unknown activation")
}

var _ json.Marshaler = Relu
var dummy = Relu
var _ json.Unmarshaler = &dummy

// UnmarshalJSON implements json.Unmarshaler.
func (a *Activation) UnmarshalJSON(bs []byte) error {
	s := strings.TrimPrefix(strings.TrimSuffix(string(bs), "\""), "\"")
	switch s {
	case Relu.String():
		*a = Relu
	case Linear.String():
		*a = Linear
	case Sigmoid.String():
		*a = Sigmoid
	case Tanh.String():
		*a = Tanh
	case Sin.String():
		*a = Sin
	case Cos.String():
		*a = Cos
	}
	return fmt.Errorf("invalid activation: %s", s)
}

// MarshalJSON implements json.Marshaler.
func (a Activation) MarshalJSON() ([]byte, error) {
	return []byte("\"" + a.String() + "\""), nil
}
