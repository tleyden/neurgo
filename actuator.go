package neurgo

import (
	"fmt"
)

type Actuator struct {
}

func (actuator *Actuator) canPropagateSignal(node *Node) bool {
	return len(node.inbound) > 0
}

func (actuator *Actuator) propagateSignal(node *Node) {

	// this implemenation is just a stub which makes it easy to test.
	// at some point, actuators will act as proxies to real virtual actuators
	// and probably be pushing their outputs to sockets.

	gatheredInputs := actuator.gatherInputs(node)
	node.scatterOutput(gatheredInputs)

}

func (actuator *Actuator) gatherInputs(node *Node) []float64 {

	outputVector := make([]float64, 0)

	for _, inboundConnection := range node.inbound {
		if inputVector, ok := <-inboundConnection.channel; ok {
			actuator.validateInputs(inputVector)
			inputValue := inputVector[0]
			outputVector = append(outputVector, inputValue)
		}
	}

	return outputVector

}

func (actuator *Actuator) validateInputs(inputs []float64) {
	if len(inputs) != 1 {
		t := "%T got invalid input vector: %v"
		message := fmt.Sprintf(t, actuator, inputs)
		panic(message)
	}
}
