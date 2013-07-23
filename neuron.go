package neurgo

import (
	"fmt"
	"github.com/proxypoke/vector"
	"log"
)

type ActivationFunction func(float64) float64

type Neuron struct {
	LayerIndex         float64
	Bias               float64
	Inbound            []*InboundConnection
	Outbound           []*OutboundConnection
	Closing            chan chan bool
	Data               chan *DataMessage
	ActivationFunction ActivationFunction
}

func (neuron *Neuron) Run() {

	log.Printf("%v Run() called.", neuron)

	neuron.checkRunnable()

	weightedInputs := createEmptyWeightedInputs(neuron.Inbound)
	log.Printf("weightedInputs: %v", weightedInputs)

	closed := false

	for {
		select {
		case responseChan := <-neuron.Closing:
			log.Printf("%v got value on closing channel", neuron)
			closed = true
			responseChan <- true
			break
		case dataMessage := <-neuron.Data:
			log.Printf("%v got data value %v", neuron, dataMessage)
		}

		if closed {
			neuron.Closing = nil // TODO: move to defer()?
			neuron.Data = nil
			break
		}

	}

	log.Printf("%v Run() finishing", neuron)

}

func (neuron *Neuron) checkRunnable() {

	if neuron.Inbound == nil {
		msg := fmt.Sprintf("not expecting neuron.Inbound to be nil")
		panic(msg)
	}

	if neuron.Closing == nil {
		msg := fmt.Sprintf("not expecting neuron.Closing to be nil")
		panic(msg)
	}

	if neuron.Data == nil {
		msg := fmt.Sprintf("not expecting neuron.Data to be nil")
		panic(msg)
	}

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
