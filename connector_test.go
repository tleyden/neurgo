package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
)


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


func TestRemoveConnection(t *testing.T) {

	// create network nodes
	neuron1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}  
	neuron2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	sensor := &Sensor{}

	// connect nodes together 
	weights := []float64{20,20,20,20,20}
	sensor.ConnectBidirectionalWeighted(neuron1, weights)
	sensor.ConnectBidirectionalWeighted(neuron2, weights)

	neuron1.inbound = removeConnection(neuron1.inbound, 0) 
	sensor.outbound = removeConnection(sensor.outbound, 0) 

	assert.Equals(t, len(neuron1.inbound), 0)
	assert.Equals(t, len(neuron2.inbound), 1)
	assert.Equals(t, len(sensor.outbound), 1)
	assert.Equals(t, sensor.outbound[0].channel, neuron2.inbound[0].channel)

}


func identity_activation(x float64) float64 {
	return x
}
