
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

func (node *Node) connectOutbound(target Connectable, channel VectorChannel) {
	connection := &connection{channel: channel}
	node.outbound = append(node.outbound, connection)
}

func (node *Node) connectInboundWeighted(source Connectable, channel VectorChannel, weights []float32) {
	connection := &connection{channel: channel, weights: weights}
	node.inbound = append(node.inbound, connection)
}


// Create a bi-directional connection between node <-> target with the given weights
// used for the inbound connection node of target <- node (inbound from target's perspective).
// In the process, a channel will be created and both nodes will have a reference to it.
func (node *Node) ConnectBidirectional(target Connectable, weights []float32) {

	fmt.Println("neural node connect w/ weights")
	channel := make(VectorChannel)		
	node.connectOutbound(target, channel)
	target.connectInboundWeighted(node, channel, weights)

	// dest2src := &connection{channel: channel, source: destination, destination: node}
	// destination.inbound = append(destination.inbound, dest2src)

	// node.add_outbound_connection(target, connection_channel)
	// target.add_inbound_connection_with_weights(neuralnode, weights, connection_channel)

}

// Same as Connect_with_weights, except neither connection will have any weights associated with it
//func (node *Node) Connect(destination Connectable) {
//	fmt.Println("neural node connect")
//}

/*
func (node Connectable) addOutboundConnection(destination Connectable, connectionChannel VectorChannel) {
}


func (node Connectable) addInboundConnectionWithWeights(source Connectable, weights []float32, connectionChannel VectorChannel) {
}
*/
