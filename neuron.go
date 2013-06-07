
package neurgo

import (
	"fmt"
	"github.com/proxypoke/vector"
)

type activationFunction func(float64) float64

type Neuron struct {
	Bias               float64
	ActivationFunction activationFunction
	Node
}

type weightedInput struct {
	weights     []float64
	inputs      []float64
}

func (neuron *Neuron) propagateSignal() {
	weightedInputs := neuron.weightedInputs()
	scalarOutput := neuron.computeScalarOutput(weightedInputs)
	outputs := []float64{scalarOutput}  
	neuron.scatterOutput(outputs)
}

// read each inbound channel and get the inputs, and pair this vector
// with the weight vector for that inbound channel, then return the
// list of those weight/input pairings.
func (neuron *Neuron) weightedInputs() []*weightedInput {
	weightedInputs := make([]*weightedInput, len(neuron.inbound))
	for i, connection := range neuron.inbound {
		inputs := <- connection.channel
		weights := connection.weights
		weightedInputs[i] = &weightedInput{weights: weights, inputs: inputs}
	}
	return weightedInputs
}
	
// compute the scalar output for the neuron
func (neuron *Neuron) computeScalarOutput(weightedInputs []*weightedInput) float64 {
	output := neuron.weightedInputDotProductSum(weightedInputs)
	output += neuron.Bias
	output = neuron.ActivationFunction(output)
	return output
}

// for each weighted input vector, calculate the (inputs * weights) dot product
// and sum all of these dot products together to produce a sum  
func (neuron *Neuron) weightedInputDotProductSum(weightedInputs []*weightedInput) float64 {

	var dotProductSummation float64
	dotProductSummation = 0

	for _, weightedInput := range weightedInputs {
		inputs := weightedInput.inputs
		weights := weightedInput.weights
		inputVector := vector.NewFrom(inputs) 
		weightVector := vector.NewFrom(weights)
		dotProduct, error := vector.DotProduct(inputVector, weightVector)
		if error != nil {
			t := "%v error performing dot product between %v and %v"
			message := fmt.Sprintf(t, neuron.Name, inputVector, weightVector) 
			panic(message)
		}
		dotProductSummation += dotProduct
	}

	return dotProductSummation
	

}
