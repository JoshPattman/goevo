package goevo

import "math"

var DefaultActivationsMap = map[string]func(float64) float64{
	"linear":  LinearActivation,
	"relu":    ReluActivation,
	"sigmoid": SigmoidActivation,
}
var DefaultActivations = []string{
	"linear",
	"relu",
	"sigmoid",
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
