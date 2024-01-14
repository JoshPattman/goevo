package goevo

import "math"

// Activation is a string representing an activation function.
type Activation string

const (
	// ActivationFunction representing [y = x]
	ActivationLinear Activation = "linear"
	// ActivationFunction representing [y = x {x > 0} | y = 0 {x <= 0}]
	ActivationReLU Activation = "relu"
	// ActivationFunction representing [y = tanh(x)]
	ActivationTanh Activation = "tanh"
	// ActivationFunction representing [y = ln(x) {x > 0} | y = 0 {x <= 0}]. I have found this can benefit recurrent networks by not allowing the values to explode
	ActivationReLn Activation = "reln"
	// ActivationFunction representing [y = sigmoid(x)]
	ActivationSigmoid Activation = "sigmoid"
	// ActivationFunction representing [y = x {1 > x > 0} | y = 0 {x <= 0} | y = 1 {x >= 1}]
	ActivationReLUMax Activation = "relumax"
	// ActivationFunction representing [y = 0 {x < 0.5} | y = 1 {x >= 0.5}]
	ActivationStep Activation = "step"
	// ActivationFunction representing [y = sin(x)]
	ActivationSin Activation = "sin"
	// ActivationFunction representing [y = cos(x)]
	ActivationCos Activation = "cos"
)

const (
	// Alias for ActivationLinear
	AcLin = ActivationLinear
	// Alias for ActivationReLU
	AcReLU = ActivationReLU
	// Alias for ActivationTanh
	AcTanh = ActivationTanh
	// Alias for ActivationReLn
	AcReLn = ActivationReLn
	// Alias for ActivationSigmoid
	AcSig = ActivationSigmoid
	// Alias for ActivationReLUMax
	AcReLUM = ActivationReLUMax
	// Alias for ActivationStep
	AcStep = ActivationStep
	// Alias for ActivationSin
	AcSin = ActivationSin
	// Alias for ActivationCos
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
