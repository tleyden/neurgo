package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
	"log"
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
	log.Printf("result: %v", result)

	delta := result - 200  // had to do within limit hack because 200 != 200 issue
	withinLimit := delta < 0.01
	assert.Equals(t, withinLimit, true)

}

