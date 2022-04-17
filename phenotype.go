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

func (p *Phenotype) ResetRecurrent() {
	for _, n := range p.Nodes {
		n.RecurrentValue = 0
	}
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
func GrowPhenotype(g Genotype) *Phenotype {
	numNodesIn, numNodesHid, numNodesOut := g.GetNumNodes()
	numNodes := numNodesIn + numNodesHid + numNodesOut
	nodes := make([]*PhenotypeNode, numNodes)
	inodes := make([]*PhenotypeNode, numNodes)
	onodes := make([]*PhenotypeNode, numNodes)
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
		thisGNode, _ := g.GetNodeIDAtLayer(i)
		for _, cid := range g.GetConnectionsFrom(thisGNode) {
			_, outID, _ := g.GetConnectionEndpoints(cid)
			outLayer, _ := g.GetLayerOfNode(outID)
			w, _ := g.GetConnectionWeight(cid)
			if r, _ := g.IsConnectionRecurrent(cid); !r {
				connections = append(connections, nodes[outLayer])
				weights = append(weights, w)
			} else {
				rConnections = append(rConnections, nodes[outLayer])
				rWeights = append(rWeights, w)
			}
		}
		if i < numNodesIn {
			nodes[i].Activation = LinearActivation
			inodes[ic] = nodes[i]
			ic++
		} else if i >= numNodesIn+numNodesHid {
			nodes[i].Activation = LinearActivation
			onodes[oc] = nodes[i]
			oc++
		} else {
			nodes[i].Activation = SigmoidActivation
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
