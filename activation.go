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
	Binary
	Relum
	Reln
)

var AllActivations = []Activation{Relu, Linear, Sigmoid, Tanh, Sin, Cos, Binary, Reln, Relum}

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
	case Binary:
		return "binary"
	case Reln:
		return "reln"
	case Relum:
		return "relum"
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
	case Binary:
		if x > 0 {
			return 1
		} else {
			return 0
		}
	case Relum:
		if x < 0 {
			return 0
		} else if x > 1 {
			return 1
		}
		return x
	case Reln:
		if x < 0 {
			return 0
		}
		return math.Log(x + 1)
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
	case Binary.String():
		*a = Binary
	case Reln.String():
		*a = Reln
	case Relum.String():
		*a = Relum
	default:
		return fmt.Errorf("invalid activation: '%s'", s)
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (a Activation) MarshalJSON() ([]byte, error) {
	return []byte("\"" + a.String() + "\""), nil
}
