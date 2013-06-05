package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
	"log"
)

func TestComputeOutput(t *testing.T) {

	activation := func(x float32) float32 { return x }
	neuron := &Neuron{Bias: 0, ActivationFunction: activation} 
	
	weights := []float32{1,1,1,1,1}
	inputs := []float32{20,20,20,20,20}

	weightedInput1 := &weightedInput{weights: weights, inputs: inputs}
	weightedInput2 := &weightedInput{weights: weights, inputs: inputs}
	weightedInputs := []*weightedInput{weightedInput1, weightedInput2}

	result := neuron.computeOutput(weightedInputs)
	log.Printf("result: %v", result)

	assert.Equals(t, result, 200)

}

