package goevo

const (
	InputNode NodeFunction = iota
	HiddenNode
	OutputNode
)

type NodeFunction int
type NodeID int
type ConnectionID int

type Genotype interface {
	// Creating synapses
	ConnectNodes(NodeID, NodeID, float64, InnovationCounter) (ConnectionID, error)
	GetConnectionWeight(ConnectionID) (float64, error)
	SetConnectionWeight(ConnectionID, float64) error
	MutateConnectionWeight(ConnectionID, float64) error
	// Creating nodes
	CreateNodeOn(ConnectionID, InnovationCounter) (NodeID, ConnectionID, error)
	// Getting info
	GetNumNodes() (int, int, int)
	GetConnectionEndpoints(ConnectionID) (NodeID, NodeID, error)
	IsConnectionRecurrent(ConnectionID) (bool, error)
	GetNodeIDAtLayer(int) (NodeID, error)
	GetLayerOfNode(NodeID) (int, error)
	GetConnectionsFrom(NodeID) []ConnectionID
	GetConnectionsTo(NodeID) []ConnectionID
	GetConnectionBetween(NodeID, NodeID) (ConnectionID, bool)
	GetAllConnectionIDs() []ConnectionID
	Copy() Genotype
}
