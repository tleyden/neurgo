
package neurgo

import "fmt"

type activationFunction func(float32) float32

type Neuron struct {
	Bias               float32
	ActivationFunction activationFunction
	Node
}


// Methods

func (neuron *Neuron) Run() {
	fmt.Println("neuron.Run()")

	// loop through all the inbound_connections (array of SignalEmitters) where SignalEmitter is interface

	    // get channel from SignalEmitter by calling signalEmitter.Channel()

	    // read value from channel

	// calculate output value: dot product + bias of all values read from SignalEmitters

	// loop over all outbound_connections (array of SignalAcceptors) where SignalAcceptor is interface

	    // get channel from SignalAcceptor by calling signalAcceptor.Channel()

	    // send output value to channel

}
