package neurgo

type InboundConnection struct {
	NodeId  *NodeId
	Weights []float64
}

type OutboundConnection struct {
	NodeId *NodeId
	Data   chan DataMessage
}
