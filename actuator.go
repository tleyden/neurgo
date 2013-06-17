
package neurgo

import (
	"fmt"
)

type Actuator struct {
	Node
}

// implementation needed here because when it was a method on *Node, it was calling 
// connectInboundWithChannel with a *Node instance and losing the fact it was a actuator
func (actuator *Actuator) ConnectBidirectional(target Connector) {
	actuator.ConnectBidirectionalWeighted(target, nil)
}

// implementation needed here because when it was a method on *Node, it was calling 
// connectInboundWithChannel with a *Node instance and losing the fact it was an actuator
func (actuator *Actuator) ConnectBidirectionalWeighted(target Connector, weights []float64) {
	channel := make(VectorChannel)		
	actuator.connectOutboundWithChannel(target, channel)
	target.connectInboundWithChannel(actuator, channel, weights)
}


func (actuator *Actuator) validateInputs(inputs []float64) {
	if len(inputs) != 1 {
		t := "%T got invalid input vector: %v"
		message := fmt.Sprintf(t, actuator, inputs)
		panic(message)
	}
}

func (actuator *Actuator) gatherInputs() []float64 {

	outputVector := make([]float64,0) 

	for _, inboundConnection := range actuator.inbound {
		if inputVector, ok := <- inboundConnection.channel; ok {
			actuator.validateInputs(inputVector)
			inputValue := inputVector[0]
			outputVector = append(outputVector, inputValue)
		}
	}

	return outputVector

}

func (actuator *Actuator) propagateSignal() {

	// this implemenation is just a stub which makes it easy to test.
	// at some point, actuators will act as proxies to real virtual actuators
	// and probably be pushing their outputs to sockets.

	gatheredInputs := actuator.gatherInputs()
	actuator.scatterOutput(gatheredInputs)

}
