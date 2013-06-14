package neurgo

type Connector interface { 

	ConnectBidirectional(target Connector)
	ConnectBidirectionalWeighted(target Connector, weights []float64)

	connectOutboundWithChannel(target Connector, channel VectorChannel) 
	connectInboundWithChannel(source Connector, channel VectorChannel, weights []float64)

	DisconnectBidirectional(target Connector)
	disconnectOutbound(target Connector)
	disconnectInbound(source Connector)

}


func removeConnection(connections []*connection, index int) []*connection {

	newConnections := make([]*connection, len(connections)-1)
	newConnectionsIndex := 0

	for i, connection := range connections {
		if i != index {
			newConnections[newConnectionsIndex] = connection
			newConnectionsIndex += 1
		}
	}
	return newConnections

}
