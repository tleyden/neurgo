package neurgo

import (
	"encoding/json"
	"fmt"
)

type Sensor struct {
}

func (sensor *Sensor) HasBias() bool {
	return false
}

func (sensor *Sensor) BiasValue() float64 {
	panic("Sensors don't have bias parameter")
}

func (sensor *Sensor) SetBias(newBias float64) {
	panic("Sensors don't have bias parameter")
}

func (sensor *Sensor) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Type string `json:"type"`
		}{
			Type: "Sensor",
		})
}

func (sensor *Sensor) copy() SignalProcessor {
	sensorCopy := &Sensor{}
	return sensorCopy
}

func (sensor *Sensor) CalculateOutput(weightedInputs []*weightedInput) []float64 {

	// sensors will eventually be reading their input from sockets
	// in the meantime, the should only have one input
	if len(weightedInputs) != 1 {
		msg := fmt.Sprintf("sensors should only have one input channel")
		panic(msg)
	}

	outputs := weightedInputs[0].inputs
	return outputs

}
