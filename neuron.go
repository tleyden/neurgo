package neurgo

import (
	"fmt"
	"github.com/proxypoke/vector"
)

type activationFunction func(float64) float64

type Neuron struct {
	Bias               float64
	ActivationFunction activationFunction
}

type weightedInput struct {
	weights []float64
	inputs  []float64
}

func (neuron *Neuron) copy() SignalProcessor {
	neuronCopy := &Neuron{}
	neuronCopy.Bias = neuron.Bias
	neuronCopy.ActivationFunction = neuron.ActivationFunction
	return neuronCopy
}

func (neuron *Neuron) canPropagateSignal(node *Node) bool {
	return len(node.inbound) > 0
}

func (neuron *Neuron) propagateSignal(node *Node) {
	weightedInputs := neuron.weightedInputs(node)
	scalarOutput := neuron.computeScalarOutput(weightedInputs)
	outputs := []float64{scalarOutput}
	node.scatterOutput(outputs)
}

// read each inbound channel and get the inputs, and pair this vector
// with the weight vector for that inbound channel, then return the
// list of those weight/input pairings.
func (neuron *Neuron) weightedInputs(node *Node) []*weightedInput {

	weightedInputs := make([]*weightedInput, 0)
	for _, connection := range node.inbound {

		// change this to a select
		var inputs []float64
		var ok bool

		select {
		case inputs = <-connection.channel:
			ok = true
		case <-connection.closing:
			return weightedInputs // <-- todo!! think about this later, won't need ok
		}

		if ok {
			if len(inputs) == 0 {
				msg := fmt.Sprintf("%v got empty inputs", neuron)
				panic(msg)
			}
			weights := connection.weights
			weightedInputs = append(weightedInputs, &weightedInput{weights: weights, inputs: inputs})
		}
	}

	return weightedInputs
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
