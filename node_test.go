package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
	"time"
)

func TestConnectedTo(t *testing.T) {

	// create nodes
	neuron := &Node{Name: "neuron", processor: &Neuron{}}
	sensor := &Node{Name: "sensor", processor: &Sensor{}}

	// assert not connected
	assert.False(t, neuron.hasOutboundConnectionTo(sensor))

	// make connection
	sensor.ConnectBidirectionalWeighted(neuron, []float64{0})

	// assert connected
	assert.True(t, sensor.hasOutboundConnectionTo(neuron))

}

func TestConnectBidirectional(t *testing.T) {

	// create nodes
	neuron := &Node{Name: "neuron", processor: &Neuron{}}
	sensor := &Node{Name: "sensor", processor: &Sensor{}}

	// make connection
	weights := []float64{20, 20, 20, 20, 20}
	sensor.ConnectBidirectionalWeighted(neuron, weights)

	// make sure the reverse connection points back to correct node type
	_, isSensor := neuron.inbound[0].other.processor.(*Sensor)
	assert.True(t, isSensor)

	// assert that it worked
	assert.Equals(t, len(sensor.outbound), 1)
	assert.Equals(t, len(neuron.inbound), 1)
	assert.True(t, neuron.inbound[0].channel != nil)
	assert.True(t, sensor.outbound[0].channel != nil)
	assert.Equals(t, len(neuron.inbound[0].weights), len(weights))
	assert.Equals(t, neuron.inbound[0].weights[0], weights[0])

	// make a new node and connect it
	actuator := &Node{processor: &Actuator{}}
	neuron.ConnectBidirectional(actuator)

	// make sure the reverse connection points back to correct node type
	_, isNeuron := actuator.inbound[0].other.processor.(*Neuron)
	assert.True(t, isNeuron)

	// assert that it worked
	assert.Equals(t, len(neuron.outbound), 1)
	assert.Equals(t, len(actuator.inbound), 1)
	assert.Equals(t, len(actuator.inbound[0].weights), 0)

}

func TestRemoveConnection(t *testing.T) {

	// create network nodes
	neuronProcessor1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuronProcessor2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuron1 := &Node{Name: "neuron1", processor: neuronProcessor1}
	neuron2 := &Node{Name: "neuron2", processor: neuronProcessor2}
	sensor := &Node{Name: "sensor", processor: &Sensor{}}

	// connect nodes together
	weights := []float64{20, 20, 20, 20, 20}
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

func TestNodeShutdown(t *testing.T) {

	neuronProcessor1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuron1 := &Node{Name: "neuron1", processor: neuronProcessor1}

	injector := &Node{}
	injector.Name = "injector"
	injector.ConnectBidirectionalWeighted(neuron1, []float64{0})

	log.Printf("call neuron1.Run()")
	neuron1.Run()
	log.Printf("called neuron1.Run()")

	log.Printf("shutting down neuron")
	neuron1.Shutdown()
	log.Printf("shut down neuron")

	timeoutChannel := time.After(time.Second / 100)
	doneChannel := make(chan bool)

	go func() {
		log.Printf("injecting value")
		injector.outbound[0].channel <- []float64{0}
		log.Printf("injected value")
		// time.Sleep(time.Second * 10)
		doneChannel <- true
	}()

	log.Printf("select()")
	select {
	case <-doneChannel:
		// neuron is shutdown, not expecting to propagate value
		assert.True(t, false)
	case <-timeoutChannel:
		// neuron is shutdown, expecting timeout
		assert.True(t, true)
	}
	log.Printf(".")

}

func identity_activation(x float64) float64 {
	return x
}
