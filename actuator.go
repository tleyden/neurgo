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

func (actuator *Actuator) HasBias() bool {
	return false
}

func (actuator *Actuator) BiasValue() float64 {
	panic("Actuators don't have bias parameter")
}

func (actuator *Actuator) SetBias(newBias float64) {
	panic("Actuators don't have bias parameter")
}

func (actuator *Actuator) copy() SignalProcessor {
	actuatorCopy := &Actuator{}
	return actuatorCopy
}

func (actuator *Actuator) canPropagate(node *Node) bool {

	return len(node.inbound) > 0

}

func (actuator *Actuator) propagateSignal(node *Node) bool {

	// this implemenation is just a stub which makes it easy to test.
	// at some point, actuators will act as proxies to real virtual actuators
	// and probably be pushing their outputs to sockets.

	gatheredInputs, isShutdown := actuator.gatherInputs(node)
	if isShutdown {
		return true
	}

	node.scatterOutput(gatheredInputs)
	return false

}

func (actuator *Actuator) gatherInputs(node *Node) (outputVector []float64, isShutdown bool) {

	outputVector = make([]float64, 0)

	for _, connection := range node.inbound {

		var inputs []float64
		var ok bool

		select {
		case inputs = <-connection.channel:
			ok = true
		case <-connection.closing: // skip this connection since its closed
		case <-node.closing:
			isShutdown = true
			return
		}

		if ok {
			actuator.validateInputs(inputs)
			inputValue := inputs[0]
			outputVector = append(outputVector, inputValue)
		}
	}

	return

}

func (actuator *Actuator) validateInputs(inputs []float64) {
	if len(inputs) != 1 {
		t := "%T got invalid input vector: %v"
		message := fmt.Sprintf(t, actuator, inputs)
		panic(message)
	}
}
