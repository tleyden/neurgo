package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"testing"
)

func TestCanPropagateSignal(t *testing.T) {

	// create nodes
	neuronProcessor1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuronProcessor2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuron1 := &Node{Name: "neuron1", processor: neuronProcessor1}
	neuron2 := &Node{Name: "neuron2", processor: neuronProcessor2}
	actuatorProcessor := &Actuator{}
	actuator := &Node{Name: "actuator", processor: actuatorProcessor}

	// connect nodes together
	neuron1.ConnectBidirectional(actuator)
	neuron2.ConnectBidirectional(actuator)

	canPropagate := actuator.processor.canPropagate(actuator)
	assert.True(t, canPropagate)

	neuron1.DisconnectBidirectional(actuator)
	neuron2.DisconnectBidirectional(actuator)

	canPropagateAfter := actuator.processor.canPropagate(actuator)
	assert.False(t, canPropagateAfter)

}
