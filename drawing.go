package goevo

import (
	"fmt"
	"image"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// Render this genotype to an image.Image, with a width and height in inches
func (g *Genotype) Draw(width, height float64) image.Image {
	gv := graphviz.New()

	graph, err := gv.Graph()
	if err != nil {
		panic(fmt.Sprintf("error when creating a graph, this should not have happened: %v", err))
	}

	defer func() {
		graph.Close()
		gv.Close()
	}()

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
		} else if no >= len(g.neuronOrder)-g.numOutputs {
			nodes[nid].SetColor("red")
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
}
