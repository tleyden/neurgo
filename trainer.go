package neurgo

import (
	"fmt"
)

type TrainingSample struct {

	// for each sensor in the network, provide a sample input vector
	SampleInputs [][]float64

	// for each actuator in the network, provide an expected output vector
	ExpectedOutputs [][]float64
}

func (t *TrainingSample) String() string {
	return fmt.Sprintf("Inputs: %v, Expected: %v",
		t.SampleInputs,
		t.ExpectedOutputs)
}

type Trainer interface {
	Train(cortex *Cortex, examples []*TrainingSample) *Cortex
}
