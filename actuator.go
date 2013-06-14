
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

	// TODO!! deal with closed channels  (and write test to exercise this)

	outputVectorDimension := len(actuator.inbound)
	outputVector := make([]float64,outputVectorDimension) 

	for i, inboundConnection := range actuator.inbound {
		inputVector := <- inboundConnection.channel
		actuator.validateInputs(inputVector)
		inputValue := inputVector[0]
		outputVector[i] = inputValue 
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
