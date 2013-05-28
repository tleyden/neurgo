
package main

import "fmt"
import "github.com/tleyden/neurgo"

func main() {

	neuron := &neurgo.Neuron{Bias: 10, ActivationFunction: identity_activation}
	sensor := &neurgo.Sensor{SyncFunction: simulated_sensor}
	actuator := &neurgo.Actuator{}

	sensor.Connect_with_weights(neuron, []float32{1,1,1,1,1})
	neuron.Connect(actuator)



	fmt.Printf("neuron bias: %f\n", neuron.Bias)
	fmt.Printf("neuron activation result: %f\n", neuron.Activate(1.0))
	fmt.Printf("sensor: %v\n", sensor)
	fmt.Printf("actuator: %v\n", actuator)

}

func identity_activation(x float32) float32 {
	return x
}

func simulated_sensor() []float32 {
	return []float32 { 1.0, 1.0, 1.0, 1.0, 1.0 }
}
