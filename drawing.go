package goevo

import (
	"fmt"
	"image"
	"strings"

	"github.com/goccy/go-graphviz"
)

// Render this genotype to an image.Image, with a width and height in inches
/*func (g *Genotype) Draw(width, height float64) image.Image {
	gv := graphviz.New()

	graph, err := gv.Graph()
	if err != nil {
		panic(fmt.Sprintf("error when creating a graph, this should not have happened: %v", err))
	}

	defer func() {
		graph.Close()
		gv.Close()
	}()

	//graph.SetLayout("neato")
	graph.SetRankDir(cgraph.LRRank)
	graph.SetRatio(cgraph.FillRatio)
	graph.SetSize(width, height)

	nodes := make(map[NeuronID]*cgraph.Node)
	for no, nid := range g.neuronOrder {
		nodes[nid], err = graph.CreateNode(fmt.Sprintf("N%v [%v]\n%v", nid, no, g.activations[nid]))
		if err != nil {
			panic(fmt.Sprintf("error when creating node on a graph, this should not have happened: %v", err))
		}
		if no < g.numInputs {
			nodes[nid].SetColor("green")
			//nodes[nid].SetPos(0, float64(no))
		} else if no >= len(g.neuronOrder)-g.numOutputs {
			nodes[nid].SetColor("red")
			//nodes[nid].SetPos(0, float64(no))
		}
		nodes[nid].SetShape(cgraph.RectShape)
	}

	for wid, w := range g.weights {
		ep := g.synapseEndpointLookup[wid]
		of, ot := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]
		edge, _ := graph.CreateEdge(fmt.Sprintf("%v->%v", ep.From, ep.To), nodes[ep.From], nodes[ep.To])
		edge.SetLabel(fmt.Sprintf("C%v\n%.3f", wid, w))
		if of >= ot {
			edge.SetColor("red")
		}
	}
	img, err := gv.RenderImage(graph)
	if err != nil {
		panic(fmt.Sprintf("error when creating an image, this should not have happened: %v", err))
	}
	return img
}*/

// This is almost certainly a worse way to implement this, but I cannot find a way to force the input and output nodes to align correctly.

func (g *NEATGenotype) RenderDot(width, height float64) string {
	graphDrawer := newSimpleGraphvizWriter()
	graphDrawer.writeGraphParam("rankdir", "LR")
	graphDrawer.writeGraphParam("ratio", "fill")
	graphDrawer.writeGraphParam("size", fmt.Sprintf("%v,%v", width, height))
	graphDrawer.writeGraphParam("layout", "dot")

	inputRanks := []string{}
	outputRanks := []string{}

	for no, nid := range g.neuronOrder {
		name := fmt.Sprintf("N%v", nid)
		label := fmt.Sprintf("N%v [%v]\n%v", nid, no, g.activations[nid])
		color := "black"
		if no < g.numInputs {
			color = "green"
			inputRanks = append(inputRanks, name)
		} else if no >= len(g.neuronOrder)-g.numOutputs {
			color = "red"
			outputRanks = append(outputRanks, name)
		}
		graphDrawer.writeNode(name, label, color)
	}

	graphDrawer.writeMinRank(inputRanks)
	graphDrawer.writeMaxRank(outputRanks)

	for wid, w := range g.weights {
		ep := g.synapseEndpointLookup[wid]
		of, ot := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]
		fromName := fmt.Sprintf("N%v", ep.From)
		toName := fmt.Sprintf("N%v", ep.To)
		label := fmt.Sprintf("C%v\n%.3f", wid, w)
		color := "black"
		if of >= ot {
			color = "red"
		}
		graphDrawer.writeEdge(fromName, toName, label, color)
	}
	return graphDrawer.dot()
}

// RenderImage returns an image of the genotype
func (g *NEATGenotype) RenderImage(width, height float64) image.Image {
	graph, err := graphviz.ParseBytes([]byte(g.RenderDot(width, height)))
	if err != nil {
		panic(fmt.Sprintf("error when creating a dot graph, this should not have happened (please report bug): %v", err))
	}
	gv := graphviz.New()
	img, err := gv.RenderImage(graph)
	if err != nil {
		panic(fmt.Sprintf("error when creating an image from dot, this should not have happened (please report bug): %v", err))
	}
	return img
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
