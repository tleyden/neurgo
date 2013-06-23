package neurgo

import (
	"encoding/json"
	"fmt"
)

type Actuator struct {
}

func (actuator *Actuator) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Type string `json:"type"`
		}{
			Type: "Actuator",
		})
}

func (actuator *Actuator) hasBias() bool {
	return false
}

func (actuator *Actuator) bias() float64 {
	panic("Actuators don't have bias parameter")
}

func (actuator *Actuator) setBias(newBias float64) {
	panic("Actuators don't have bias parameter")
}

func (actuator *Actuator) copy() SignalProcessor {
	actuatorCopy := &Actuator{}
	return actuatorCopy
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

	for _, connection := range node.inbound {

		var inputs []float64
		var ok bool

		select {
		case inputs = <-connection.channel:
			ok = true
		case <-connection.closing:
		}

		if ok {
			actuator.validateInputs(inputs)
			inputValue := inputs[0]
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
