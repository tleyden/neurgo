
package neurgo

type Node struct {
	inbound  []*connection
	outbound []*connection
}

type connection struct {
	channel     VectorChannel
	weights     []float32
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


