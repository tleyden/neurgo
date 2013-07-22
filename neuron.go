package neurgo

import (
	"fmt"
	"github.com/proxypoke/vector"
)

type ActivationFunction func(float64) float64

type Neuron struct {
	Bias               float64
	Inbound            []*InboundConnection
	Outbound           []*OutboundConnection
	Closing            chan bool
	Data               chan DataMessage
	ActivationFunction ActivationFunction
}

type weightedInput struct {
	weights []float64
	inputs  []float64
}

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
			t := "%T error performing dot product between %v and %v"
			message := fmt.Sprintf(t, neuron, inputVector, weightVector)
			panic(message)
		}
		dotProductSummation += dotProduct
	}

	return dotProductSummation

}
