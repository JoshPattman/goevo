package goevo

import "math"

// A string representing an activation function. It must be one of the consts that start with 'Activation...'
type Activation string

const (
	// y = x
	ActivationLinear Activation = "linear"
	// y = x {x > 0} | y = 0 {x <= 0}
	ActivationReLU Activation = "relu"
	// y = tanh(x)
	ActivationTanh Activation = "tanh"
	// y = ln(x) {x > 0} | y = 0 {x <= 0}. I have found this can benefit recurrent networks by not allowing the values to explode
	ActivationReLn Activation = "reln"
	// y = sigmoid(x)
	ActivationSigmoid Activation = "sigmoid"
	// y = x {1 > x > 0} | y = 0 {x <= 0} | y = 1 {x >= 1}
	ActivationReLUMax Activation = "relumax"
	// y = 0 {x < 0.5} | y = 1 {x >= 0.5}
	ActivationStep Activation = "step"
	// y = sin(x)
	ActivationSin Activation = "sin"
	// y = cos(x)
	ActivationCos Activation = "cos"
)

const (
	// These are just aliases for the above consts
	// ActivationLinear
	AcLin = ActivationLinear
	// ActivationReLU
	AcReLU = ActivationReLU
	// ActivationTanh
	AcTanh = ActivationTanh
	// ActivationReLn
	AcReLn = ActivationReLn
	// ActivationSigmoid
	AcSig = ActivationSigmoid
	// ActivationReLUMax
	AcReLUM = ActivationReLUMax
	// ActivationStep
	AcStep = ActivationStep
	// ActivationSin
	AcSin = ActivationSin
	// ActivationCos
	AcCos = ActivationCos
)

var activationMap = map[Activation](func(float64) float64){
	ActivationLinear:  linearActivation,
	ActivationReLU:    reluActivation,
	ActivationTanh:    tanhActivation,
	ActivationReLn:    relnActivation,
	ActivationSigmoid: sigmoidActivation,
	ActivationReLUMax: relumaxActivation,
	ActivationStep:    stepActivation,
	ActivationSin:     sinActivation,
	ActivationCos:     cosActivation,
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
func stepActivation(x float64) float64 {
	if x < 0.5 {
		return 0
	}
	return 1
}

func relumaxActivation(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
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

func sinActivation(x float64) float64 {
	return math.Sin(x)
}

func cosActivation(x float64) float64 {
	return math.Cos(x)
}
