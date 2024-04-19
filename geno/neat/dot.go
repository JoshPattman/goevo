package neat

import (
	"fmt"
	"image"

	"github.com/goccy/go-graphviz"
)

// RenderDot returns a string in the DOT language that represents the genotype.
// This DOT code cannot be use to recreate the genotype, but can be used to visualise it using Graphviz.
func (g *Genotype) RenderDot(width, height float64) string {
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

// RenderImage returns an image of the genotype using graphviz.
func (g *Genotype) RenderImage(width, height float64) image.Image {
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
