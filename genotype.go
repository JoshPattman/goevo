package goevo

import (
	"fmt"
	"strconv"
)

const (
	InputNode NodeFunction = iota
	HiddenNode
	OutputNode
)

type NodeFunction int

type NodeID int

type NodeGene struct {
	ID       NodeID
	Function NodeFunction
	Layer    int
}

type ConnectionID int

type ConnectionGene struct {
	ID        ConnectionID
	In        NodeID
	Out       NodeID
	Weight    float64
	Recurrent bool
	Enabled   bool
}

type Genotype struct {
	Nodes       map[NodeID]*NodeGene
	Connections map[ConnectionID]*ConnectionGene
	Layers      []*NodeGene
	numInput    int
	numOutput   int
}

func NewGenotype(numIn, numOut int, counter InnovationCounter) *Genotype {
	nodes := make([]*NodeGene, numIn+numOut)
	for i := 0; i < numIn; i++ {
		nodes[i] = &NodeGene{
			ID:       NodeID(counter.Next()),
			Function: InputNode,
			Layer:    i,
		}
	}
	for i := numIn; i < numIn+numOut; i++ {
		nodes[i] = &NodeGene{
			ID:       NodeID(counter.Next()),
			Function: OutputNode,
			Layer:    i,
		}
	}
	nodesMap := make(map[NodeID]*NodeGene)
	for _, n := range nodes {
		nodesMap[n.ID] = n
	}
	g := &Genotype{
		Layers:      nodes,
		Connections: make(map[ConnectionID]*ConnectionGene, 0),
		Nodes:       nodesMap,
		numInput:    numIn,
		numOutput:   numOut,
	}
	return g
}

// TODO: Speed this up
func (g *Genotype) GetConnectionByEndpoints(a, b NodeID) *ConnectionGene {
	for _, c := range g.Connections {
		if c.In == a && c.Out == b {
			return c
		}
	}
	return nil
}

func (g *Genotype) GetNodeTypeCounts() (int, int, int) {
	return g.numInput, len(g.Nodes) - g.numInput - g.numOutput, g.numOutput
}

func (g *Genotype) IsConnected(a, b NodeID) bool {
	c := g.GetConnectionByEndpoints(a, b)
	return c != nil
}

func (g *Genotype) GetConnection(cid ConnectionID) *ConnectionGene {
	if v, ok := g.Connections[cid]; ok {
		return v
	}
	return nil
}

func (g *Genotype) GetNode(nid NodeID) *NodeGene {
	if v, ok := g.Nodes[nid]; ok {
		return v
	}
	return nil
}

func (g *Genotype) CreateConnection(in, out NodeID, weight float64, counter InnovationCounter) bool {
	if g.IsConnected(in, out) {
		return false
	}
	inNode := g.GetNode(in)
	outNode := g.GetNode(out)
	if inNode == nil || outNode == nil {
		return false
	}
	if inNode.Function == OutputNode || outNode.Function == InputNode {
		return false
	}
	if inNode.Layer >= outNode.Layer {
		return false
	}
	newID := ConnectionID(counter.Next())
	con := &ConnectionGene{
		ID:      newID,
		In:      in,
		Out:     out,
		Weight:  weight,
		Enabled: true,
	}
	g.Connections[con.ID] = con
	return true
}
func (g *Genotype) CreateRecurrentConnection(in, out NodeID, weight float64, counter InnovationCounter) bool {
	if g.IsConnected(in, out) {
		return false
	}
	inNode := g.GetNode(in)
	outNode := g.GetNode(out)
	if inNode == nil || outNode == nil {
		return false
	}
	if inNode.Function == InputNode || outNode.Function == OutputNode {
		return false
	}
	if inNode.Layer < outNode.Layer {
		return false
	}
	newID := ConnectionID(counter.Next())
	con := &ConnectionGene{
		ID:        newID,
		In:        in,
		Out:       out,
		Weight:    weight,
		Enabled:   true,
		Recurrent: true,
	}
	g.Connections[con.ID] = con
	return true
}

func (g *Genotype) CreateNode(conID ConnectionID, counter InnovationCounter) bool {
	c := g.GetConnection(conID)
	if c == nil {
		return false
	}

	na := g.GetNode(c.In)
	nb := g.GetNode(c.Out)

	insertionPoint := integerMidpoint(na.Layer, nb.Layer)
	inpC, _, _ := g.GetNodeTypeCounts()
	if insertionPoint < inpC {
		insertionPoint = inpC
	}
	n := &NodeGene{
		NodeID(counter.Next()),
		HiddenNode,
		insertionPoint,
	}
	g.Nodes[n.ID] = n
	// insertion
	g.Layers = append(g.Layers, n)
	copy(g.Layers[insertionPoint+1:], g.Layers[insertionPoint:])
	g.Layers[insertionPoint] = n
	for l := insertionPoint + 1; l < len(g.Layers); l++ {
		g.Layers[l].Layer++
	}
	c.Enabled = false
	g.CreateConnection(c.In, n.ID, c.Weight, counter)
	g.CreateConnection(n.ID, c.Out, 1, counter)
	return true
}

func (g *Genotype) MutateConnectionBy(cid ConnectionID, v float64) bool {
	c := g.GetConnection(cid)
	if c == nil {
		return false
	}
	c.Weight += v
	return true
}

func CopyGenotype(g *Genotype) *Genotype {
	nodes := make(map[NodeID]*NodeGene)
	layers := make([]*NodeGene, len(g.Layers))
	cons := make(map[ConnectionID]*ConnectionGene)

	for l, n := range g.Layers {
		newNode := *n
		layers[l] = &newNode
		nodes[newNode.ID] = &newNode
	}

	for _, c := range g.Connections {
		newCon := *c
		cons[c.ID] = &newCon
	}

	g1 := &Genotype{
		Nodes:       nodes,
		Connections: cons,
		Layers:      layers,
		numInput:    g.numInput,
		numOutput:   g.numOutput,
	}
	return g1
}

/*
// ApproximateGeneticDistance : This is not the correct genetic difference, but rather a heuristic
func (g *Genotype) ApproximateGeneticDistance(g1 *Genotype) float64 {
	weightDiff := 1.0
	connDiff := 1.0
	d := 0.0
	found := make(map[ConnectionID]*ConnectionGene)
	for c := range g.Connections {
		found[g.Connections[c].ID] = &g.Connections[c]
		// assume this is the only gene. We will reverse this if the other genome has this gene
		d += connDiff
	}
	for c := range g1.Connections {
		if v, ok := found[g1.Connections[c].ID]; ok {
			// Both genomes have this connection
			d -= connDiff
			d += math.Abs(v.Weight-g1.Connections[c].Weight) * weightDiff
		} else {
			// only this genome has this connection
			d += connDiff
		}
	}
	return d
}
*/

func (g *Genotype) String() string {
	s := "(["
	for k, v := range g.Nodes {
		s += strconv.Itoa(int(k)) + ":"
		s += fmt.Sprint(*v) + ","
	}
	s += "]["
	for k, v := range g.Connections {
		s += strconv.Itoa(int(k)) + ":"
		s += fmt.Sprint(*v) + ","
	}
	s += "]["
	for k, v := range g.Layers {
		s += strconv.Itoa(int(k)) + ":"
		s += fmt.Sprint(v.ID) + ","
	}
	s += "])"
	return s
}
