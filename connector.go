package neurgo

import (
	"time"
)

type Connector interface { 

	ConnectBidirectional(target Connector)
	ConnectBidirectionalWeighted(target Connector, weights []float64)

	connectOutboundWithChannel(target Connector, channel VectorChannel) 
	connectInboundWithChannel(source Connector, channel VectorChannel, weights []float64)

	DisconnectBidirectional(target Connector)
	disconnectOutbound(target Connector)
	disconnectInbound(source Connector)

	outboundConnections() []*connection
	inboundConnections() []*connection

	appendOutboundConnection(target *connection)
	appendInboundConnection(source *connection)

	// read inputs from inbound connections, calculate output, then
	// propagate the output to outbound connections
	propagateSignal()

	// is this signaller actually able to propagate a signal?
	canPropagateSignal() bool


}

// continually propagate incoming signals -> outgoing signals
func Run(connector Connector) {

	for {
		if !connector.canPropagateSignal() {
			time.Sleep(1 * 1e9)
		} else {
			connector.propagateSignal()	
		}

		
	}

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
