package goevo

import "math/rand"

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
