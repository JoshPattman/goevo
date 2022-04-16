package goevo

import "math"

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
