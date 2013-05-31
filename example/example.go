
package main

import (
	"github.com/tleyden/neurgo"
	"log"
	"sync"
)

func main() {

	// create network nodes
	neuron := &neurgo.Neuron{Bias: 10, ActivationFunction: identity_activation}
	sensor := &neurgo.Sensor{}
	actuator := &neurgo.Actuator{}
	
	// connect nodes together 
	weights := []float32{20,20,20,20,20}
	sensor.ConnectBidirectionalWeighted(neuron, weights)
	neuron.ConnectBidirectional(actuator)

	// spinup node goroutines
	go neuron.Run()
	go sensor.Run()
	go actuator.Run()

	inbound_channels := sensor.inbound

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)

	// inject a value from sensor -> neuron
	go func() {
		sensorChannel := sensor.outbound[0].channel
		sensorChannel <- []float32{1,1,1,1,1}
		wg.Done()
	}()

	// read the value from actuator
	go func() {
		actuatorChannel := actuator.inbound[0].channel
		value := <- actuatorChannel
		log.Printf("Value: %v", value)
		wg.Done()
	}()

	wg.Wait()
	log.Println("done")


}




func identity_activation(x float32) float32 {
	return x
}

func simulated_sensor() []float32 {
	return []float32 { 1.0, 1.0, 1.0, 1.0, 1.0 }
}
