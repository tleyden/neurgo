package neurgo

import ()

type SignalProcessor interface {

	// read inputs from inbound connections, calculate output, then
	// propagate the output to outbound connections.
	// returns true if the node was detected to be shutdown
	propagateSignal(node *Node) bool

	// is this signaller actually able to propagate a signal?
	canPropagate(node *Node) bool

	// create a copy of this signalprocessor
	copy() SignalProcessor

	// does this signal processor have a "bias"
	hasBias() bool

	bias() float64

	setBias(newBias float64)
}
