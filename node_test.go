package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
	"log"
	"sync"
)

func TestConnectBidirectional(t *testing.T) {

	neuron := &Neuron{}
	sensor := &Sensor{}

	weights := []float32{20,20,20,20,20}
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
	neuron := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	sensor := &Sensor{}
	actuator := &Actuator{}
	
	// connect nodes together 
	weights := []float32{20,20,20,20,20}
	sensor.ConnectBidirectionalWeighted(neuron, weights)
	neuron.ConnectBidirectional(actuator)

	// spinup node goroutines
	go neuron.Run()
	go sensor.Run()
	go actuator.Run()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)

	// inject a value from sensor -> neuron
	testValue := []float32{1,1,1,1,1}
	log.Printf("Injecting value: %v", testValue)
	go func() {
		sensorChannel := sensor.outbound[0].channel
		sensorChannel <- testValue
		wg.Done()
	}()

	// read the value from actuator
	go func() {
		actuatorChannel := actuator.inbound[0].channel
		value := <- actuatorChannel
		log.Printf("Received value: %v", value)
		assert.Equals(t, len(value), len(testValue))
		assert.Equals(t, value[0], testValue[0])  // TODO: is there a better way to check slice equality?
		wg.Done()
	}()

	wg.Wait()
	log.Println("done")

}

func identity_activation(x float32) float32 {
	return x
}
