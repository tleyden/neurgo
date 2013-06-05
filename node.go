
package neurgo

import (
	"log"
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

type weightedInput struct {
	weights     []float64
	inputs      []float64
}

func (node *Node) weightedInputs() []*weightedInput {
	weightedInputs := make([]*weightedInput, len(node.inbound))
	for i, inboundConnection := range node.inbound {
		log.Printf("%v reading from channel: %v", node.Name, inboundConnection.channel)
		inputs := <- inboundConnection.channel
		log.Printf("%v got data from channel: %v", node.Name, inboundConnection.channel)
		weightedInputs[i] = &weightedInput{weights: inboundConnection.weights, inputs: inputs}
	}
	return weightedInputs
}

func (node *Node) propagateSignal() {

	log.Printf("%s: Run()", node.Name) // TODO: how do I print the type of this struct w/o using Name field?
	
	//weightedInputs := node.weightedInputs()
	//scalarOutput := node.computeOutput(weightedInputs)
	//log.Printf("scalarOutput: %v", scalarOutput)

	// OLD

	if len(node.inbound) > 0 {  

		// TODO: for a neuron, output vector dimension should be [1]
		outputVectorDimension := len(node.inbound)
		outputVector := make([]float64,outputVectorDimension) 

		// TODO: sum up the dot products and then add the bias?

		/*
                neurlang snippet:
                    weighted_inputs = get_weighted_inputs(neuron)
                    scalar_output = compute_output(weighted_inputs, bias, activation_function)
                    [ scalar_output ]
                */

		for i, inboundConnection := range node.inbound {
			log.Printf("%v reading from channel: %v", node.Name, inboundConnection.channel)
			inputVector := <- inboundConnection.channel
			log.Printf("%v got data from channel: %v", node.Name, inboundConnection.channel)
			// TODO multiply by weights, run through activation function  (in its own method, with a test)
			fakeOutputValue := float64(len(inputVector))
			outputVector[i] = fakeOutputValue
		}

		if len(node.outbound) > 0 {
			node.outbound[0].channel <- outputVector  // TODO: loop over output channels
		}

	}

	// stub: read from input and send to output
	//if len(node.inbound) > 0 && len(node.outbound) > 0 {
	//	val := <- node.inbound[0].channel   
	//	node.outbound[0].channel <- val
	//}

	// loop through all the inbound_connections 

	    // get channel 

	    // read value from channel

	// Via "Activatable" interface..
	// calculate output value: dot product + bias of all values read from SignalEmitters

	// loop over all outbound_connections 

	    // get channel 

	    // send output value to channel

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


