package goevo

import (
	"fmt"
	"strings"
)

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
