package neurgo

import ()

type SignalProcessor interface {
	CalculateOutput(weightedInputs []*weightedInput) []float64

	// create a copy of this signalprocessor
	copy() SignalProcessor

	// does this signal processor have a "bias"
	HasBias() bool

	BiasValue() float64

	SetBias(newBias float64)
}
