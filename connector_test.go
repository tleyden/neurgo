package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
	"log"
	"sync"
	"time"
)


func TestConnectBidirectional(t *testing.T) {

	// create nodes
	neuron := &Neuron{}
	sensor := &Sensor{}

	// give names
	neuron.Name = "neuron"
	sensor.Name = "sensor"

	// make connection
	weights := []float64{20,20,20,20,20}
	sensor.ConnectBidirectionalWeighted(neuron, weights)

	// assert that it worked
	assert.Equals(t, len(sensor.outbound), 1)
	assert.Equals(t, len(neuron.inbound), 1)
	assert.True(t, neuron.inbound[0].channel != nil)
	assert.True(t, sensor.outbound[0].channel != nil)
	assert.Equals(t, len(neuron.inbound[0].weights), len(weights))
	assert.Equals(t, neuron.inbound[0].weights[0], weights[0])

	// make a new node and connect it
	actuator := &Actuator{}
	neuron.ConnectBidirectional(actuator)

	// assert that it worked
	assert.Equals(t, len(neuron.outbound), 1)
	assert.Equals(t, len(actuator.inbound), 1)
	assert.Equals(t, len(actuator.inbound[0].weights), 0)

}


func TestRemoveConnection(t *testing.T) {

	// create network nodes
	neuron1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}  
	neuron2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	sensor := &Sensor{}

	// give names to nodes
	neuron1.Name = "neuron1"
	neuron2.Name = "neuron2"
	sensor.Name = "sensor"

	// connect nodes together 
	weights := []float64{20,20,20,20,20}
	sensor.ConnectBidirectionalWeighted(neuron1, weights)
	sensor.ConnectBidirectionalWeighted(neuron2, weights)

	// remove connections
	neuron1.inbound = removeConnection(neuron1.inbound, 0) 
	sensor.outbound = removeConnection(sensor.outbound, 0) 

	// assert that it worked
	assert.Equals(t, len(neuron1.inbound), 0)
	assert.Equals(t, len(neuron2.inbound), 1)
	assert.Equals(t, len(sensor.outbound), 1)
	assert.Equals(t, sensor.outbound[0].channel, neuron2.inbound[0].channel)

}

func TestRemoveConnectionFromRunningNode(t *testing.T) {

	// create nodes
	sensor1 := &Sensor{}
	sensor2 := &Sensor{}
	neuron := &Neuron{Bias: 10, ActivationFunction: identity_activation}

	// give names to nodes
	neuron.Name = "neuron"
	sensor1.Name = "sensor1"
	sensor2.Name = "sensor2"

	// connect nodes together
	weights := []float64{20}
	sensor1.ConnectBidirectionalWeighted(neuron, weights)
	sensor2.ConnectBidirectionalWeighted(neuron, weights)

	// TODO
	// close other channel
	// call weightedInputs and get zero results 


	// basic sanity check, send two inputs to neuron inbound channels
	// and verify that weightedInputs() returns both inputs
	go func() {
		sensor1.outbound[0].channel <- []float64{0}
	}()
	go func() {
		sensor2.outbound[0].channel <- []float64{0}
	}()
	weightedInputs := neuron.weightedInputs()
	assert.Equals(t, len(weightedInputs), 2)
	
	// close one channel while a neuron is reading from
	// both inbound connections, make sure it returns one value
	
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)
	
	go func() {
		weightedInputs := neuron.weightedInputs()
		log.Printf("len(weightedInputs): %v", len(weightedInputs))
		assert.Equals(t, len(weightedInputs), 1)
		wg.Done() 
	}()

	go func() {
		time.Sleep(0.1 * 1e9)
		sensor1.DisconnectBidirectional(neuron)
		sensor2.outbound[0].channel <- []float64{0}
		wg.Done() 
	}()

	wg.Wait()

	log.Printf("done")

}

func identity_activation(x float64) float64 {
	return x
}
