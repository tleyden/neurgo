package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
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


func identity_activation(x float64) float64 {
	return x
}
