package neurgo

import (
	"testing"
	"fmt"
	"github.com/couchbaselabs/go.assert"
)

func TestConnectBidirectionalWeighted(t *testing.T) {

	fmt.Println("test is running!")

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


}
