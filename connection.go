package neurgo

type InboundConnection struct {
	SourceNodeId NodeId
	Weights      []float32
}

type OutboundConnection struct {
	TargetNodeId NodeId
	Data         chan DataMessage
}
