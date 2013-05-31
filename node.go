
package neurgo

import "fmt"

type Node struct {
	inbound  []*connection
	outbound []*connection
}


type connection struct {
	channel     VectorChannel
	weights     []float32
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


// Create a bi-directional connection between node <-> target with the given weights
// used for the inbound connection node of target <- node (inbound from target's perspective).
// In the process, a channel will be created and both nodes will have a reference to it.
func (node *Node) ConnectBidirectionalWeighted(target Connectable, weights []float32) {

	fmt.Println("neural node connect w/ weights")
	channel := make(VectorChannel)		
	node.connectOutbound(target, channel)
	target.connectInbound(node, channel, weights)
}

func (node *Node) ConnectBidirectionalUnweighted(target Connectable) {
	node.ConnectBidirectionalWeighted(target, nil)
}

