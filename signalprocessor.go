package neurgo

import (
	"time"
)

type SignalProcessor interface {

	// read inputs from inbound connections, calculate output, then
	// propagate the output to outbound connections
	propagateSignal(node *Node)

	// is this signaller actually able to propagate a signal?
	canPropagateSignal(node *Node) bool
}

// continually propagate incoming signals -> outgoing signals
func Run(processor SignalProcessor, node *Node) {

	for {
		if !processor.canPropagateSignal(node) {
			time.Sleep(1 * 1e9)
		} else {
			processor.propagateSignal(node)
		}

	}

}
