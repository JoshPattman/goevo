package goevo

import (
	"math"
	"math/rand"
)

func randomMapKey[T comparable, U any](m map[T]U) T {
	n := rand.Intn(len(m))
	i := 0
	for k := range m {
		if i == n {
			return k
		}
		i++
	}
	panic("cannot occur")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func stdN(std float64) int {
	v := math.Abs(rand.NormFloat64() * std)
	if v > std*10 {
		v = std * 10 // Lets just cap this at 10 std to prevent any sillyness
	}
	return int(math.Round(v))
}
