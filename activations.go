package goevo

import "math"

var DefaultActivationInfo = ActivationInfo{
	InputActivation:  LinearActivation,
	HiddenActivation: RelnActivation,
	OutputActivation: TanhActivation,
}

func LinearActivation(x float64) float64 {
	return x
}

func ReluActivation(x float64) float64 {
	if x < 0 {
		return 0
	}
	return x
}

func SigmoidActivation(x float64) float64 {
	return 1 / (1 + math.Pow(math.E, -x))
}

func TanhActivation(x float64) float64 {
	return math.Tanh(x)
}

func RelnActivation(x float64) float64 {
	if x < 0 {
		return 0
	}
	return math.Log(x + 1)
}

type ActivationInfo struct {
	InputActivation  func(float64) float64
	HiddenActivation func(float64) float64
	OutputActivation func(float64) float64
}
