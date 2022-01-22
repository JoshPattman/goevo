package goevo

func LinearActivation(x float64) float64 {
	return x
}

func ReluActivation(x float64) float64 {
	if x < 0 {
		return 0
	}
	return x
}
