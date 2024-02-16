package goevo

// Pos is a type representing an n dimensional position
type Pos []float64

func P(ds ...float64) Pos {
	return Pos(ds)
}

func LayoutPosLine(min, max Pos, numPoints int) []Pos {
	if len(min) != len(max) {
		panic("lengths must match")
	}
	if numPoints == 1 {
		return []Pos{min}
	}
	if numPoints < 1 {
		panic("must have some points")
	}
	increment := make(Pos, len(min))
	for i := range min {
		increment[i] = (max[i] - min[i]) / float64(numPoints-1)
	}
	current := make(Pos, len(min))
	copy(current, min)
	pts := []Pos{}
	for i := 0; i < numPoints; i++ {
		pt := make(Pos, len(min))
		copy(pt, current)
		pts = append(pts, pt)
		for j := range current {
			current[j] += increment[j]
		}
	}
	return pts
}
