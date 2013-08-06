package neurgo

type NodeId struct {
	UUID       string
	NodeType   NodeType
	LayerIndex float64
}

func NewSensorId(UUID string, LayerIndex float64) *NodeId {
	return &NodeId{
		UUID:       UUID,
		NodeType:   SENSOR,
		LayerIndex: LayerIndex,
	}

}

func NewNeuronId(UUID string, LayerIndex float64) *NodeId {
	return &NodeId{
		UUID:       UUID,
		NodeType:   NEURON,
		LayerIndex: LayerIndex,
	}

}

func NewActuatorId(UUID string, LayerIndex float64) *NodeId {
	return &NodeId{
		UUID:       UUID,
		NodeType:   ACTUATOR,
		LayerIndex: LayerIndex,
	}

}
