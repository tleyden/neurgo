
package neurgo

import "fmt"

type Node struct {
	
}

// Connectable interface implementations

func (node *Node) Connect_with_weights(target Connectable, weights []float32) {
	fmt.Println("neural node connect_with_weights")

	// TODO: unit test!
	
	// connection_channel = (make a channel)
	// neuralnode.add_outbound_connection(target, connection_channel)
	// target.add_inbound_connection_with_weights(neuralnode, weights, connection_channel)

}

func (node *Node) Connect(target Connectable) {
	fmt.Println("neural node connect")
}

