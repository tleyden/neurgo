package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
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
	neuron2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	sensor := &Sensor{}
	actuator := &Actuator{}
	wiretap := &Wiretap{}
	injector := &Injector{}

	// connect nodes together 
	injector.ConnectBidirectional(sensor)
	weights := []float64{20,20,20,20,20}
	sensor.ConnectBidirectionalWeighted(neuron1, weights)
	sensor.ConnectBidirectionalWeighted(neuron2, weights)
	neuron1.ConnectBidirectional(actuator)
	neuron2.ConnectBidirectional(actuator)
	actuator.ConnectBidirectional(wiretap)

	// spinup node goroutines
	signallers := []Signaller{neuron1, neuron2, sensor, actuator}
	for _, signaller := range signallers {
		go Run(signaller)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)

	// inject a value into sensor
	go func() {
		testValue := []float64{1,1,1,1,1}
		injector.outbound[0].channel <- testValue
		wg.Done()
	}()

	// read the value from wiretap (which taps into actuator)
	go func() {
		value := <- wiretap.inbound[0].channel
		assert.Equals(t, len(value), 2)  
		assert.Equals(t, value[0], float64(110)) 
		assert.Equals(t, value[1], float64(110))
		wg.Done() 
	}()

	wg.Wait()

}



func identity_activation(x float64) float64 {
	return x
}
