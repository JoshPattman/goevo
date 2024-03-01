package goevo

type Forwarder interface {
	Forward([]float64) []float64
}
