
package neurgo

import (

)

type connection struct {
	channel     VectorChannel
	weights     []float64
}

type Node struct {
	Name     string
	inbound  []*connection
	outbound []*connection
}

func (node *Node) canPropagateSignal() bool {
	return len(node.inbound) > 0 
}

func (node *Node) scatterOutput(outputs []float64) {
	for _, outboundConnection := range node.outbound {
		outboundConnection.channel <- outputs
	}
}

// Create a bi-directional connection between node <-> target with no weights associated
// with the connection
func (node *Node) ConnectBidirectional(target Connectable) {
	node.ConnectBidirectionalWeighted(target, nil)
}

// Create a bi-directional connection between node <-> target with the given weights.
func (node *Node) ConnectBidirectionalWeighted(target Connectable, weights []float64) {
	channel := make(VectorChannel)		
	node.connectOutboundWithChannel(target, channel)
	target.connectInboundWithChannel(node, channel, weights)
}

// Create outbound connection from node -> target
func (node *Node) connectOutboundWithChannel(target Connectable, channel VectorChannel) {
	connection := &connection{channel: channel}
	node.outbound = append(node.outbound, connection)
}

// Create inbound connection from source -> node
func (node *Node) connectInboundWithChannel(source Connectable, channel VectorChannel, weights []float64) {
	connection := &connection{channel: channel, weights: weights}
	node.inbound = append(node.inbound, connection)
}


