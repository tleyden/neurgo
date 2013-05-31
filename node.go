
package neurgo

import (
	"log"
)

type connection struct {
	channel     VectorChannel
	weights     []float32
}

type Node struct {
	inbound  []*connection
	outbound []*connection
}

func (node *Node) Run() {

	// stub: read from input and send to output
	if len(node.inbound) > 0 && len(node.outbound) > 0 {
		val := <- node.inbound[0].channel   
		node.outbound[0].channel <- val
	}

	// loop through all the inbound_connections 

	    // get channel 

	    // read value from channel

	// Via "Activatable" interface..
	// calculate output value: dot product + bias of all values read from SignalEmitters

	// loop over all outbound_connections 

	    // get channel 

	    // send output value to channel

}


// Create a bi-directional connection between node <-> target with the given weights.
// In the process, a channel will be created and both nodes will have a reference to it.
func (node *Node) ConnectBidirectionalWeighted(target Connectable, weights []float32) {
	channel := make(VectorChannel)		
	node.connectOutbound(target, channel)
	target.connectInbound(node, channel, weights)
}

// Create a bi-directional connection between node <-> target with no weights associated
// with the connection
func (node *Node) ConnectBidirectional(target Connectable) {
	node.ConnectBidirectionalWeighted(target, nil)
}

// Create outbound connection from node -> target
func (node *Node) connectOutbound(target Connectable, channel VectorChannel) {
	connection := &connection{channel: channel}
	node.outbound = append(node.outbound, connection)
}

// Create inbound connection from source -> node
func (node *Node) connectInbound(source Connectable, channel VectorChannel, weights []float32) {
	connection := &connection{channel: channel, weights: weights}
	node.inbound = append(node.inbound, connection)
}


