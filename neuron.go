package neurgo

import (
	"fmt"
	"github.com/proxypoke/vector"
	"log"
)

type ActivationFunction func(float64) float64

type Neuron struct {
	NodeId             *NodeId
	Bias               float64
	Inbound            []*InboundConnection
	Outbound           []*OutboundConnection
	Closing            chan chan bool
	DataChan           chan *DataMessage
	ActivationFunction ActivationFunction
}

func (neuron *Neuron) Run() {

	log.Printf("")

	neuron.checkRunnable()

	neuron.sendEmptySignalRecurrentOutbound()

	weightedInputs := createEmptyWeightedInputs(neuron.Inbound)

	closed := false

	for {
		select {
		case responseChan := <-neuron.Closing:
			closed = true
			responseChan <- true
			break // TODO: do we need this for anything??
		case dataMessage := <-neuron.DataChan:
			neuron.recordInput(weightedInputs, dataMessage)
		}

		if closed {
			neuron.Closing = nil
			neuron.DataChan = nil
			break
		}

		if neuron.receiveBarrierSatisfied(weightedInputs) {

			scalarOutput := neuron.computeScalarOutput(weightedInputs)

			dataMessage := &DataMessage{
				SenderId: neuron.NodeId,
				Inputs:   []float64{scalarOutput},
			}

			neuron.scatterOutput(dataMessage)

			weightedInputs = createEmptyWeightedInputs(neuron.Inbound)

		}

	}

}

func (neuron *Neuron) String() string {
	return fmt.Sprintf("%v", neuron.NodeId)
}

func (neuron *Neuron) receiveBarrierSatisfied(weightedInputs []*weightedInput) bool {
	satisfied := true
	for _, weightedInput := range weightedInputs {
		if weightedInput.inputs == nil {
			satisfied = false
			break
		}

	}
	return satisfied
}

// In order to prevent deadlock, any neurons we have recurrent outbound
// connections to must be "primed" by sending an empty signal.  A recurrent
// outbound connection simply means that it's a connection to ourself or
// to a neuron in a previous (eg, to the left) layer.  If we didn't do this,
// that previous neuron would be waiting forever for a signal that will
// never come, because this neuron wouldn't fire until it got a signal.
func (neuron *Neuron) sendEmptySignalRecurrentOutbound() {

	recurrentConnections := neuron.recurrentOutboundConnections()
	for recurrentConnection := range recurrentConnections {
		inputs := []float64{0}
		dataMessage := &DataMessage{
			SenderId: neuron.NodeId,
			Inputs:   inputs,
		}
		neuron.DataChan <- dataMessage
	}

}

// Find the subset of outbound connections which are "recurrent" - meaning
// that the connection is to this neuron itself, or to a neuron in a previous
// (eg, to the left) layer.
func (neuron *Neuron) recurrentOutboundConnections() []*OutboundConnection {
	result := make([]*OutboundConnection, 0)
	for outboundConnection := range neuron.Outbound {
		if outboundConnection.isRecurrent(neuron.NodeId) {
			result = append(result, outboundConnection)
		}
	}
	return result
}

func (neuron *Neuron) recordInput(weightedInputs []*weightedInput, dataMessage *DataMessage) {
	for _, weightedInput := range weightedInputs {
		if weightedInput.senderNodeId == dataMessage.SenderId {
			weightedInput.inputs = dataMessage.Inputs
		}
	}

}

func (neuron *Neuron) scatterOutput(dataMessage *DataMessage) {
	for _, outboundConnection := range neuron.Outbound {
		dataChan := outboundConnection.DataChan
		dataChan <- dataMessage
	}
}

func (neuron *Neuron) checkRunnable() {

	if neuron.NodeId == nil {
		msg := fmt.Sprintf("not expecting neuron.NodeId to be nil")
		panic(msg)
	}

	if neuron.Inbound == nil {
		msg := fmt.Sprintf("not expecting neuron.Inbound to be nil")
		panic(msg)
	}

	if neuron.Closing == nil {
		msg := fmt.Sprintf("not expecting neuron.Closing to be nil")
		panic(msg)
	}

	if neuron.DataChan == nil {
		msg := fmt.Sprintf("not expecting neuron.DataChan to be nil")
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
