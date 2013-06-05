package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
	"log"
	"sync"
)

type Wiretap struct {
	Node
}

type Injector struct {
	Node
}


func TestConnectBidirectional(t *testing.T) {

	neuron := &Neuron{}
	sensor := &Sensor{}

	weights := []float64{20,20,20,20,20}
	sensor.ConnectBidirectionalWeighted(neuron, weights)

	assert.Equals(t, len(sensor.outbound), 1)
	assert.Equals(t, len(neuron.inbound), 1)
	assert.True(t, neuron.inbound[0].channel != nil)
	assert.True(t, sensor.outbound[0].channel != nil)
	assert.Equals(t, len(neuron.inbound[0].weights), len(weights))
	assert.Equals(t, neuron.inbound[0].weights[0], weights[0])

	actuator := &Actuator{}
	neuron.ConnectBidirectional(actuator)
	assert.Equals(t, len(neuron.outbound), 1)
	assert.Equals(t, len(actuator.inbound), 1)
	assert.Equals(t, len(actuator.inbound[0].weights), 0)

}

func TestNetwork(t *testing.T) {

	// create network nodes
	neuron1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}  
	neuron1.Name = "neuron1" // TODO: why doesn't this work in literal above?
	neuron2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuron2.Name = "neuron2"
	sensor := &Sensor{}
	sensor.Name = "sensor"
	actuator := &Actuator{}
	actuator.Name = "actuator"
	wiretap := &Wiretap{}
	wiretap.Name = "wiretap"
	injector := &Injector{}
	injector.Name = "injector"

	// connect nodes together 
	injector.ConnectBidirectional(sensor)
	weights := []float64{20,20,20,20,20}
	sensor.ConnectBidirectionalWeighted(neuron1, weights)
	sensor.ConnectBidirectionalWeighted(neuron2, weights)
	neuron1.ConnectBidirectional(actuator)
	neuron2.ConnectBidirectional(actuator)
	actuator.ConnectBidirectional(wiretap)

	// spinup node goroutines
	go Run(neuron1)
	go Run(neuron2)
	go Run(sensor)
	go Run(actuator)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)

	// inject a value from sensor -> neuron
	testValue := []float64{1,1,1,1,1}

	go func() {

		log.Printf("%v Injecting value2: %v via outbound[0].  channel: %v", injector.Name, testValue, injector.outbound[0].channel)
		injector.outbound[0].channel <- testValue

		log.Println("injector goroutine done")
		wg.Done()

	}()

	// read the value from actuator
	go func() {

		log.Printf("%v Getting value on chan: %v", wiretap.Name, wiretap.inbound[0].channel)
		value := <- wiretap.inbound[0].channel
		log.Printf("%v Received value on chan: %v: %v", wiretap.Name, wiretap.inbound[0].channel, value)
		
		log.Println("wiretap goroutine done")
		wg.Done()
		
		assert.Equals(t, len(value), 2)  // accumulate values from 2 neurons, therefore 2 elt vector    
		//assert.Equals(t, value[0], 110) // n1 output: 110 
		//assert.Equals(t, value[1], 110) // n2 output: 110 

	}()

	wg.Wait()
	log.Println("test done")

}



func identity_activation(x float64) float64 {
	return x
}
