package neurgo

type Signaller interface {

	// read inputs from inbound connections, calculate output, then
	// propagate the output to outbound connections
	propagateSignal()

}

// continually propagate incoming signals -> outgoing signals
func Run(signaller Signaller) {

	for {
		signaller.propagateSignal()
	}

}
