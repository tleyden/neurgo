
package main

import (
	"github.com/tleyden/neurgo"
	"log"
)

func main() {

	// create network nodes
	neuron := &neurgo.Neuron{Bias: 10, ActivationFunction: identity_activation}
	sensor := &neurgo.Sensor{InputChannel: make(neurgo.VectorChannel)}
	actuator := &neurgo.Actuator{OutputChannel: make(neurgo.VectorChannel)}
	
	// connect nodes together 
	weights = []float32{20,20,20,20,20}
	sensor.ConnectBidirectional(neuron, weights)

	// neuron.Connect(actuator)

	// spinup node goroutines
	go neuron.Run()
	go sensor.Run()
	go actuator.Run()

	// push test value into sensor input channel
	// sensor.SendInput([]float32{1,1,1,1,1})

	// read value from actuator output channel
	result := actuator.ReadOutput()

	log.Printf("result: %f", result)

	// make sure it's the expected value

	// debug crap ..
	log.Printf("neuron bias: %f", neuron.Bias)
	log.Printf("sensor: %v", sensor)
	log.Printf("actuator: %v", actuator)

}

func identity_activation(x float32) float32 {
	return x
}

func simulated_sensor() []float32 {
	return []float32 { 1.0, 1.0, 1.0, 1.0, 1.0 }
}
