package neurgo

type Neuron struct {
	Inbound  []InboundConnection
	Outbound []OutboundConnection
	Closing  chan bool
	Data     chan DataMessage
}

type InboundConnection struct {
	SourceNodeId NodeId
	Weights      []float32
}

type OutboundConnection struct {
	TargetNodeId NodeId
	Data         chan DataMessage
}

type DataMessage struct {
	SenderId NodeId
	Inputs   []float32
}
