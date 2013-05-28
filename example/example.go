
package main

import "fmt"
import "github.com/tleyden/neurgo"

func main() {

	// create network nodes
	neuron := &neurgo.Neuron{Bias: 10, ActivationFunction: identity_activation}
	sensor := &neurgo.Sensor{InputChannel: make(chan []float32)}
	actuator := &neurgo.Actuator{OutputChannel: make(chan float32)}

	// connect nodes together
	sensor.Connect_with_weights(neuron, []float32{20,20,20,20,20})
	neuron.Connect(actuator)

	// spinup node goroutines
	go neuron.Run()
	go sensor.Run()
	go actuator.Run()

	// push test value into sensor input channel
	// sensor.SendInput([]float32{1,1,1,1,1})

	// read value from actuator output channel
	result := actuator.ReadOutput()

	fmt.Printf("result: %f\n", result)

	// make sure it's the expected value

	// debug crap ..
	fmt.Printf("neuron bias: %f\n", neuron.Bias)
	fmt.Printf("sensor: %v\n", sensor)
	fmt.Printf("actuator: %v\n", actuator)

}

func identity_activation(x float32) float32 {
	return x
}

func simulated_sensor() []float32 {
	return []float32 { 1.0, 1.0, 1.0, 1.0, 1.0 }
}
