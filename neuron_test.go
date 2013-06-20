package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
)

func TestComputeOutput(t *testing.T) {

	activation := func(x float64) float64 { return x }
	neuron := &Neuron{Bias: 0, ActivationFunction: activation}

	weights := []float64{1, 1, 1, 1, 1}
	inputs := []float64{20, 20, 20, 20, 20}

	weightedInput1 := &weightedInput{weights: weights, inputs: inputs}
	weightedInput2 := &weightedInput{weights: weights, inputs: inputs}
	weightedInputs := []*weightedInput{weightedInput1, weightedInput2}

	result := neuron.computeScalarOutput(weightedInputs)

	assert.Equals(t, result, float64(200))

}

func TestCopyNeuron(t *testing.T) {

	activation := func(x float64) float64 { return x }
	neuron := &Neuron{Bias: 0, ActivationFunction: activation}
	neuronCopyProcessor := neuron.copy()
	neuronCopy := neuronCopyProcessor.(*Neuron)
	assert.Equals(t, neuron.Bias, neuronCopy.Bias)
	assert.Equals(t, neuron.ActivationFunction(1.0), neuronCopy.ActivationFunction(1.0))
	log.Printf("neuronCopy: %v", neuronCopy)

}
