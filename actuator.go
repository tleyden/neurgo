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

func (actuator *Actuator) CalculateOutput(weightedInputs []*weightedInput) []float64 {

	outputs := make([]float64, 0)
	for _, weightedInput := range weightedInputs {
		inputs := weightedInput.inputs
		actuator.validateInputs(inputs)
		inputValue := inputs[0]
		outputs = append(outputs, inputValue)
	}

	return outputs

}

func (actuator *Actuator) validateInputs(inputs []float64) {
	if len(inputs) != 1 {
		t := "%T got invalid input vector: %v"
		message := fmt.Sprintf(t, actuator, inputs)
		panic(message)
	}
}
