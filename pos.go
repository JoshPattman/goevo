package goevo

// Pos is a type representing an n dimensional position
type Pos []float64

func P(ds ...float64) Pos {
	return Pos(ds)
}

// PosGrid creates a grid of positions in a hypercube, with min and max values.
// The nth dimension will have dims[n] points.
func PosGrid(dims []int, min, max float64) []Pos {
	panic("Not implemented yet")
}

// PosArray creates an array of 'dims' points with minimum and maximum values.
// It returns 1D points.
func PosArray(dims int, min, max float64) []Pos {
	return PosGrid([]int{dims}, min, max)
}
