package goevo

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"gonum.org/v1/gonum/mat"
)

type floatType interface {
	float32 | float64
}

type numberType interface {
	floatType | int | int16 | int32 | int64
}

func stdN(std float64) int {
	v := math.Abs(rand.NormFloat64() * std)
	if v > std*10 {
		v = std * 10 // Lets just cap this at 10 std to prevent any sillyness
	}
	return int(math.Round(v))
}

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

type mutMat interface {
	mat.Matrix
	Set(r, c int, v float64)
}

type mutVecWrapper struct {
	*mat.VecDense
}

func (m *mutVecWrapper) Set(r, c int, v float64) {
	if c != 0 {
		panic("cannot set vector at anything other than column 0")
	}
	m.SetVec(r, v)
}

// make sure to check the shapes first!!
func randomChoiceMatrix(into mutMat, ms []mutMat) {
	rs, cs := into.Dims()
	for ri := range rs {
		for ci := range cs {
			idx := rand.Intn(len(ms))
			into.Set(ri, ci, ms[idx].At(ri, ci))
		}
	}
}

// Utility struct to build up a simple graphviz graph one statement at a time
type simpleGraphvizWriter struct {
	lines []string
}

func newSimpleGraphvizWriter() *simpleGraphvizWriter {
	return &simpleGraphvizWriter{[]string{}}
}

func (s *simpleGraphvizWriter) writeGraphParam(name, value string) {
	s.lines = append(s.lines, fmt.Sprintf("%s=\"%s\";", name, value))
}

func (s *simpleGraphvizWriter) writeMinRank(nodes []string) {
	nodesList := strings.Join(nodes, "; ") + ";"
	s.lines = append(s.lines, fmt.Sprintf("{rank=min; %s}", nodesList))
}

func (s *simpleGraphvizWriter) writeMaxRank(nodes []string) {
	nodesList := strings.Join(nodes, "; ") + ";"
	s.lines = append(s.lines, fmt.Sprintf("{rank=max; %s}", nodesList))
}

func (s *simpleGraphvizWriter) writeNode(name, label, color string) {
	s.lines = append(s.lines, fmt.Sprintf("%s [label=\"%s\", color=\"%s\", shape=rect];", name, label, color))
}

func (s *simpleGraphvizWriter) writeEdge(from, to, label, color string) {
	s.lines = append(s.lines, fmt.Sprintf("%s -> %s [label=\"%s\", color=\"%s\"];", from, to, label, color))
}

func (s *simpleGraphvizWriter) dot() string {
	body := strings.Join(s.lines, "\n\t")
	return fmt.Sprintf("digraph G {\n\t%s\n}", body)
}
