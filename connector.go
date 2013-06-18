package neurgo

import (

)

type Connector interface { 

	ConnectBidirectional(target *Node)
	ConnectBidirectionalWeighted(target *Node, weights []float64)

	connectOutboundWithChannel(target *Node, channel VectorChannel) 
	connectInboundWithChannel(source *Node, channel VectorChannel, weights []float64)

	DisconnectBidirectional(target *Node)
	disconnectOutbound(target *Node)
	disconnectInbound(source *Node)

	outboundConnections() []*connection
	inboundConnections() []*connection

	appendOutboundConnection(target *connection)
	appendInboundConnection(source *connection)

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
