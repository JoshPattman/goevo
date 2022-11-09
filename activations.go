package goevo

import "math"

type Activation string

const (
	ActivationLinear  Activation = "linear"
	ActivationReLU    Activation = "relu"
	ActivationTanh    Activation = "tanh"
	ActivationReLn    Activation = "reln"
	ActivationSigmoid Activation = "sigmoid"
)

var activationMap = map[Activation](func(float64) float64){
	ActivationLinear:  linearActivation,
	ActivationReLU:    reluActivation,
	ActivationTanh:    tanhActivation,
	ActivationReLn:    relnActivation,
	ActivationSigmoid: sigmoidActivation,
}

func linearActivation(x float64) float64 {
	return x
}

func reluActivation(x float64) float64 {
	if x < 0 {
		return 0
	}
	return x
}

func sigmoidActivation(x float64) float64 {
	return 1 / (1 + math.Pow(math.E, -x))
}

func tanhActivation(x float64) float64 {
	return math.Tanh(x)
}

func relnActivation(x float64) float64 {
	if x < 0 {
		return 0
	}
	return math.Log(x + 1)
}
