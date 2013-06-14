package neurgo

import (
	"time"
)

type Signaller interface {

	// read inputs from inbound connections, calculate output, then
	// propagate the output to outbound connections
	propagateSignal()

	// is this signaller actually able to propagate a signal?
	canPropagateSignal() bool

}

// continually propagate incoming signals -> outgoing signals
func Run(signaller Signaller) {

	for {
		if !signaller.canPropagateSignal() {
			time.Sleep(1 * 1e9)
		} else {
			signaller.propagateSignal()	
		}

		
	}

}
