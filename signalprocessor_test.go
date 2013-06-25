package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
	"time"
)

func TestCanPropagateSignalShutdown(t *testing.T) {

	// create nodes
	neuronProcessor1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuronProcessor2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuron1 := &Node{Name: "neuron1", processor: neuronProcessor1}
	neuron2 := &Node{Name: "neuron2", processor: neuronProcessor2}
	actuatorProcessor := &Actuator{}
	actuator := &Node{Name: "actuator", processor: actuatorProcessor}

	actuator.Run()

	// connect nodes together
	neuron1.ConnectBidirectional(actuator)
	neuron2.ConnectBidirectional(actuator)

	isShutdown := actuator.processor.waitCanPropagate(actuator)
	assert.False(t, isShutdown)

	go func() {
		time.Sleep(time.Second / 100)
		actuator.Shutdown()
	}()

	neuron1.DisconnectBidirectional(actuator)
	neuron2.DisconnectBidirectional(actuator)

	isShutdown = actuator.processor.waitCanPropagate(actuator)
	assert.True(t, isShutdown)

}

func TestCanPropagateSignalReAddConnection(t *testing.T) {

	// create nodes
	neuronProcessor1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuronProcessor2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuron1 := &Node{Name: "neuron1", processor: neuronProcessor1}
	neuron2 := &Node{Name: "neuron2", processor: neuronProcessor2}
	actuatorProcessor := &Actuator{}
	actuator := &Node{Name: "actuator", processor: actuatorProcessor}

	actuator.Run()

	// connect nodes together
	neuron1.ConnectBidirectional(actuator)
	neuron2.ConnectBidirectional(actuator)

	isShutdown := actuator.processor.waitCanPropagate(actuator)
	assert.False(t, isShutdown)

	go func() {
		time.Sleep(time.Second / 100)
		neuron1.ConnectBidirectional(actuator)
	}()

	neuron1.DisconnectBidirectional(actuator)
	neuron2.DisconnectBidirectional(actuator)

	isShutdown = actuator.processor.waitCanPropagate(actuator)
	assert.False(t, isShutdown)

	log.Printf("done")

}
