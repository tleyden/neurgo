
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


func (neuron *Neuron) propagateSignal() {
	weightedInputs := neuron.weightedInputs()
	scalarOutput := neuron.computeScalarOutput(weightedInputs)
	outputs := []float64{scalarOutput}  
	neuron.scatterOutput(outputs)
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
			message := fmt.Sprintf("%v error performing dot product between %v and %v", neuron.Name, inputVector, weightVector) 
			panic(message)
		}
		dotProductSummation += dotProduct
	}

	return dotProductSummation
	

}
