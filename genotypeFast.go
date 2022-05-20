package goevo

import "errors"

type GenotypeFast struct {
	Layers      []NodeID
	Nodes       map[NodeID]*FastNodeGene
	Connections map[ConnectionID]*FastConnectionGene
	NumInputs   int
	NumOutputs  int
}

type FastNodeGene struct {
	ID NodeID
}

type FastConnectionGene struct {
	ID          ConnectionID
	FromID      NodeID
	ToID        NodeID
	Weight      float64
	IsRecurrent bool
}

func NewGenotypeFast(numIn, numOut int, counter InnovationCounter) *GenotypeFast {
	layers := make([]NodeID, numIn+numOut)
	nodes := make(map[NodeID]*FastNodeGene)
	for i := 0; i < numIn+numOut; i++ {
		id := NodeID(counter.Next())
		n := &FastNodeGene{id}
		nodes[id] = n
		layers[i] = id
	}
	conns := make(map[ConnectionID]*FastConnectionGene)
	return &GenotypeFast{
		Layers:      layers,
		Nodes:       nodes,
		Connections: conns,
		NumInputs:   numIn,
		NumOutputs:  numOut,
	}
}

func (g *GenotypeFast) GetConnectionWeight(cid ConnectionID) (float64, error) {
	v, found := g.Connections[cid]
	if !found {
		return 0, errors.New("Connection not found")
	}
	return v.Weight, nil
}
func (g *GenotypeFast) MutateConnectionWeight(cid ConnectionID, w float64) error {
	v, found := g.Connections[cid]
	if !found {
		return errors.New("Connection not found")
	}
	g.Connections[cid].Weight = v.Weight + w
	return nil
}
func (g *GenotypeFast) SetConnectionWeight(cid ConnectionID, w float64) error {
	_, found := g.Connections[cid]
	if !found {
		return errors.New("Connection not found")
	}
	g.Connections[cid].Weight = w
	return nil
}
func (g *GenotypeFast) GetNumNodes() (int, int, int) {
	n := len(g.Nodes)
	return g.NumInputs, n - g.NumOutputs - g.NumInputs, g.NumOutputs
}
func (g *GenotypeFast) GetNodeIDAtLayer(l int) (NodeID, error) {
	if l >= len(g.Layers) || l < 0 {
		return 0, errors.New("Layer out of range")
	}
	return g.Layers[l], nil
}
func (g *GenotypeFast) GetLayerOfNode(nid NodeID) (int, error) {
	n, found := g.Nodes[nid]
	if !found {
		return 0, errors.New("Could not find node")
	}
	for l, n2id := range g.Layers {
		if n2id == n.ID {
			return l, nil
		}
	}
	return 0, errors.New("Malformed layer list")
}
func (g *GenotypeFast) GetConnectionsFrom(nid NodeID) []ConnectionID {
	cons := make([]ConnectionID, 0)
	for _, c := range g.Connections {
		if c.FromID == nid {
			cons = append(cons, c.ID)
		}
	}
	return cons
}
func (g *GenotypeFast) GetConnectionsTo(nid NodeID) []ConnectionID {
	cons := make([]ConnectionID, 0)
	for _, c := range g.Connections {
		if c.ToID == nid {
			cons = append(cons, c.ID)
		}
	}
	return cons
}
func (g *GenotypeFast) GetConnectionBetween(nida, nidb NodeID) (ConnectionID, bool) {
	for _, c := range g.Connections {
		if c.FromID == nida && c.ToID == nidb {
			return c.ID, true
		}
	}
	return 0, false
}

func (g *GenotypeFast) ConnectNodes(nida, nidb NodeID, weight float64, counter InnovationCounter) (ConnectionID, error) {
	if nida == nidb {
		return 0, errors.New("Those are the same node")
	}
	if _, isConnected := g.GetConnectionBetween(nida, nidb); isConnected {
		return 0, errors.New("Those nodes were already connected")
	}
	la, err1 := g.GetLayerOfNode(nida)
	lb, err2 := g.GetLayerOfNode(nidb)
	if err1 != nil || err2 != nil {
		return 0, errors.New("Node did not exist")
	}
	if lb > la {
		//not recurrent
		if lb < g.NumInputs {
			return 0, errors.New("Cannot connect to an input")
		} else if la >= len(g.Nodes)-g.NumOutputs {
			return 0, errors.New("Cannot connect from an output")
		}
		id := ConnectionID(counter.Next())
		g.Connections[id] = &FastConnectionGene{
			ID:          id,
			FromID:      nida,
			ToID:        nidb,
			Weight:      weight,
			IsRecurrent: false,
		}
		return id, nil
	} else {
		//recurrent
		if la < g.NumInputs {
			return 0, errors.New("Cannot connect from an input on recurrent")
		} else if lb >= len(g.Nodes)-g.NumOutputs {
			return 0, errors.New("Cannot connect to an output on recurrent")
		}
		id := ConnectionID(counter.Next())
		g.Connections[id] = &FastConnectionGene{
			ID:          id,
			FromID:      nida,
			ToID:        nidb,
			Weight:      weight,
			IsRecurrent: true,
		}
		return id, nil
	}
}
func (g *GenotypeFast) CreateNodeOn(cid ConnectionID, counter InnovationCounter) (NodeID, ConnectionID, error) {
	c, exists := g.Connections[cid]
	if !exists {
		return 0, 0, errors.New("Connection does not exist")
	}
	if c.IsRecurrent {
		return 0, 0, errors.New("Cannot create nodes on recurrent connection at the moment")
	}
	// Create IDs
	newNodeID := NodeID(counter.Next())
	newConID := ConnectionID(counter.Next())
	// Create node
	n := &FastNodeGene{newNodeID}
	g.Nodes[newNodeID] = n
	// Create extra connection
	c2 := &FastConnectionGene{
		ID:          newConID,
		FromID:      newNodeID,
		ToID:        c.ToID,
		Weight:      1,
		IsRecurrent: false,
	}
	c.ToID = newNodeID
	// Find new node insertion layer and insert node
	g.Connections[newConID] = c2
	la, _ := g.GetLayerOfNode(c.FromID)
	lb, _ := g.GetLayerOfNode(c2.ToID)
	insertionPoint := integerMidpoint(la, lb)
	if insertionPoint < g.NumInputs {
		insertionPoint = g.NumInputs
	}
	g.Layers = append(g.Layers, newNodeID)
	copy(g.Layers[insertionPoint+1:], g.Layers[insertionPoint:])
	g.Layers[insertionPoint] = newNodeID
	return newNodeID, newConID, nil
}

func (g *GenotypeFast) GetConnectionEndpoints(cid ConnectionID) (NodeID, NodeID, error) {
	c, isIn := g.Connections[cid]
	if !isIn {
		return 0, 0, errors.New("Connection does not exist")
	}
	return c.FromID, c.ToID, nil
}

func (g *GenotypeFast) GetAllConnectionIDs() []ConnectionID {
	keys := make([]ConnectionID, len(g.Connections))
	i := 0
	for k := range g.Connections {
		keys[i] = k
		i++
	}
	return keys
}

func (g *GenotypeFast) IsConnectionRecurrent(cid ConnectionID) (bool, error) {
	c, isIn := g.Connections[cid]
	if !isIn {
		return false, errors.New("Connection does not exist")
	}
	return c.IsRecurrent, nil
}

func (g *GenotypeFast) Copy() Genotype {
	layers := make([]NodeID, len(g.Layers))
	copy(layers, g.Layers)
	nodes := make(map[NodeID]*FastNodeGene)
	for nid, n := range g.Nodes {
		copy := *n
		nodes[nid] = &copy
	}
	cons := make(map[ConnectionID]*FastConnectionGene)
	for cid, n := range g.Connections {
		copy := *n
		cons[cid] = &copy
	}
	return &GenotypeFast{
		Layers:      layers,
		Nodes:       nodes,
		Connections: cons,
		NumInputs:   g.NumInputs,
		NumOutputs:  g.NumOutputs,
	}
}
