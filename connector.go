package neurgo

type Connector interface { 

	ConnectBidirectional(target Connector)
	ConnectBidirectionalWeighted(target Connector, weights []float64)

	connectOutboundWithChannel(target Connector, channel VectorChannel) 
	connectInboundWithChannel(source Connector, channel VectorChannel, weights []float64)

}
