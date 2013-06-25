package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
	"time"
)

func TestCanPropagateSignal(t *testing.T) {

	// create nodes
	neuronProcessor1 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuronProcessor2 := &Neuron{Bias: 10, ActivationFunction: identity_activation}
	neuron1 := &Node{Name: "neuron1", processor: neuronProcessor1}
	neuron2 := &Node{Name: "neuron2", processor: neuronProcessor2}
	actuatorProcessor := &Actuator{}
	actuator := &Node{Name: "actuator", processor: actuatorProcessor}

	go actuator.Run()

	// this hack is needed ... I think because the closing channels
	// aren't setup yet and it breaks things.  TODO: modify
	// node.Run() to be synchronous and kick off go routine internally
	// after closing channels have been setup
	time.Sleep(time.Second / 100)

	// connect nodes together
	neuron1.ConnectBidirectional(actuator)
	neuron2.ConnectBidirectional(actuator)

	isShutdown := actuator.processor.waitCanPropagate(actuator)
	assert.False(t, isShutdown)

	go func() {
		log.Printf("goroutine sleeping")
		time.Sleep(time.Second / 100)
		log.Printf("goroutine calling actuator.Shutdown()")
		actuator.Shutdown()
		log.Printf("Shutdown() finished")
	}()

	log.Printf("-- disconnect nodes from actuator")
	neuron1.DisconnectBidirectional(actuator)
	neuron2.DisconnectBidirectional(actuator)

	log.Printf("-- call waitCanPropagate()")
	isShutdown = actuator.processor.waitCanPropagate(actuator)
	log.Printf("called waitCanPropagate() and got isShutdown: %v", isShutdown)
	assert.True(t, isShutdown)

	log.Printf("done")

}
