package goevo

type PhenotypeNode struct {
	Value               float64
	RecurrentValue      float64
	Successors          []*PhenotypeNode
	RecurrentSuccessors []*PhenotypeNode
	Weights             []float64
	RecurrentWeights    []float64
	Activation          func(float64) float64
}

type Phenotype struct {
	Nodes       []*PhenotypeNode
	InputNodes  []*PhenotypeNode
	OutputNodes []*PhenotypeNode
}

func (p *Phenotype) Calculate(inputs []float64) []float64 {
	// Recurrent pass
	for _, n := range p.Nodes {
		n.RecurrentValue = 0
	}
	for i := len(p.Nodes) - 1; i >= 0; i-- {
		n := p.Nodes[i]
		for i, n2 := range n.RecurrentSuccessors {
			n2.RecurrentValue += n.Value * n.RecurrentWeights[i]
		}
	}
	// Forward pass
	for _, n := range p.Nodes {
		n.Value = 0
	}
	for i, n := range p.Nodes {
		if i < len(inputs) {
			p.InputNodes[i].Value = inputs[i]
		}
		n.Value = n.Activation(n.Value + n.RecurrentValue)
		for i, n2 := range n.Successors {
			n2.Value += n.Value * n.Weights[i]
		}
	}
	outs := make([]float64, len(p.OutputNodes))
	for i, n := range p.OutputNodes {
		outs[i] = n.Value
	}
	return outs
}

func GrowPhenotype(g *Genotype) *Phenotype {
	nodes := make([]*PhenotypeNode, len(g.Nodes))
	inodes := make([]*PhenotypeNode, len(g.Nodes))
	onodes := make([]*PhenotypeNode, len(g.Nodes))
	ic := 0
	oc := 0
	for i := range nodes {
		p := &PhenotypeNode{}
		nodes[i] = p
	}
	for i := range nodes {
		connections := make([]*PhenotypeNode, 0)
		rConnections := make([]*PhenotypeNode, 0)
		weights := make([]float64, 0)
		rWeights := make([]float64, 0)
		for _, c := range g.Connections {
			if c.In == g.Layers[i].ID && c.Enabled {
				oi := g.GetNode(c.Out)
				if !c.Recurrent {
					connections = append(connections, nodes[oi.Layer])
					weights = append(weights, c.Weight)
				} else {
					rConnections = append(rConnections, nodes[oi.Layer])
					rWeights = append(rWeights, c.Weight)
				}
			}
		}
		if g.Layers[i].Function == InputNode {
			nodes[i].Activation = LinearActivation
			inodes[ic] = nodes[i]
			ic++
		} else if g.Layers[i].Function == OutputNode {
			nodes[i].Activation = LinearActivation
			onodes[oc] = nodes[i]
			oc++
		} else {
			nodes[i].Activation = LinearActivation
		}
		nodes[i].Weights = weights
		nodes[i].Successors = connections
		nodes[i].RecurrentWeights = rWeights
		nodes[i].RecurrentSuccessors = rConnections
	}

	return &Phenotype{
		Nodes:       nodes,
		InputNodes:  inodes[:ic],
		OutputNodes: onodes[:oc],
	}
}
