
package neurgo

import (
	"fmt"
)

type Actuator struct {
	Node
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
