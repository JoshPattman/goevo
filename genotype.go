package goevo

import "math"

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
}

type ConnectionID int

type ConnectionGene struct {
	ID      ConnectionID
	In      NodeID
	Out     NodeID
	Weight  float64
	Enabled bool
}

type Genotype struct {
	Nodes       []NodeGene
	Connections []ConnectionGene
	numInput    int
	numOutput   int
}

func CreateGenotype(numIn, numOut int, counter InnovationCounter) *Genotype {
	nodes := make([]NodeGene, numIn+numOut)
	for i := 0; i < numIn; i++ {
		nodes[i] = NodeGene{
			ID:       NodeID(counter.Next()),
			Function: InputNode,
		}
	}
	for i := numIn; i < numIn+numOut; i++ {
		nodes[i] = NodeGene{
			ID:       NodeID(counter.Next()),
			Function: OutputNode,
		}
	}
	g := &Genotype{
		Nodes:       nodes,
		Connections: make([]ConnectionGene, 0),
		numInput:    numIn,
		numOutput:   numOut,
	}
	return g
}

func (g *Genotype) GetConnectionByEndpoints(a, b NodeID) (*ConnectionGene, int) {
	for ci := range g.Connections {
		c := &g.Connections[ci]
		if c.In == a && c.Out == b {
			return c, ci
		}
	}
	return nil, -1
}

func (g *Genotype) GetNodeTypeCounts() (int, int, int) {
	return g.numInput, len(g.Nodes) - g.numInput - g.numOutput, g.numOutput
}

func (g *Genotype) IsConnected(a, b NodeID) bool {
	c, _ := g.GetConnectionByEndpoints(a, b)
	return c != nil
}

func (g *Genotype) GetConnection(cid ConnectionID) (*ConnectionGene, int) {
	for ci := range g.Connections {
		c := &g.Connections[ci]
		if c.ID == cid {
			return c, ci
		}
	}
	return nil, -1
}

func (g *Genotype) GetNode(nid NodeID) (*NodeGene, int) {
	for ci := range g.Nodes {
		c := &g.Nodes[ci]
		if c.ID == nid {
			return c, ci
		}
	}
	return nil, -1
}

func (g *Genotype) CreateConnection(in, out NodeID, weight float64, counter InnovationCounter) bool {
	if g.IsConnected(in, out) {
		return false
	}
	inNode, inNodeI := g.GetNode(in)
	outNode, outNodeI := g.GetNode(out)
	if inNode == nil || outNode == nil {
		return false
	}
	if inNode.Function == OutputNode || outNode.Function == InputNode {
		return false
	}
	if inNodeI >= outNodeI {
		return false
	}
	newID := ConnectionID(counter.Next())
	con := ConnectionGene{
		ID:      newID,
		In:      in,
		Out:     out,
		Weight:  weight,
		Enabled: true,
	}
	g.Connections = append(g.Connections, con)
	return true
}

func (g *Genotype) CreateNode(conID ConnectionID, counter InnovationCounter) bool {
	c, _ := g.GetConnection(conID)
	if c == nil {
		return false
	}

	_, na := g.GetNode(c.In)
	_, nb := g.GetNode(c.Out)

	n := NodeGene{
		NodeID(counter.Next()),
		HiddenNode,
	}
	inpC, _, _ := g.GetNodeTypeCounts()
	insertionPoint := integerMidpoint(na, nb)
	if insertionPoint < inpC {
		insertionPoint = inpC
	}
	g.Nodes = append(g.Nodes, n)
	copy(g.Nodes[insertionPoint+1:], g.Nodes[insertionPoint:])
	g.Nodes[insertionPoint] = n
	c.Enabled = false
	g.CreateConnection(c.In, n.ID, c.Weight, counter)
	g.CreateConnection(n.ID, c.Out, 1, counter)
	return true
}

func (g *Genotype) MutateConnectionBy(cid ConnectionID, v float64) bool {
	c, _ := g.GetConnection(cid)
	if c == nil {
		return false
	}
	c.Weight += v
	return true
}

func CopyGenotype(g *Genotype) *Genotype {
	nodes := make([]NodeGene, len(g.Nodes))
	cons := make([]ConnectionGene, len(g.Connections))
	copy(nodes, g.Nodes)
	copy(cons, g.Connections)
	g1 := &Genotype{
		Nodes:       nodes,
		Connections: cons,
		numInput:    g.numInput,
		numOutput:   g.numOutput,
	}
	return g1
}

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
