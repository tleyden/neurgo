package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
)

func TestComputeOutput(t *testing.T) {

	activation := func(x float64) float64 { return x }
	neuron := &Neuron{Bias: 0, ActivationFunction: activation} 
	
	weights := []float64{1,1,1,1,1}
	inputs := []float64{20,20,20,20,20}

	weightedInput1 := &weightedInput{weights: weights, inputs: inputs}
	weightedInput2 := &weightedInput{weights: weights, inputs: inputs}
	weightedInputs := []*weightedInput{weightedInput1, weightedInput2}

	result := neuron.computeScalarOutput(weightedInputs)

	assert.Equals(t, result, float64(200))

}

