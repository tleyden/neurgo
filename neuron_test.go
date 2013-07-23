package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
	"time"
)

func identityActivationFunction() ActivationFunction {
	return func(x float64) float64 { return x }
}

func TestRunningNeuron(t *testing.T) {

	log.Printf("TestRunningNeuron")

	activation := identityActivationFunction()

	neuronNodeId := &NodeId{UUID: "neuron", NodeType: "test-neuron"}
	nodeId_1 := &NodeId{UUID: "node-1", NodeType: "test-node"}
	nodeId_2 := &NodeId{UUID: "node-2", NodeType: "test-node"}
	nodeId_3 := &NodeId{UUID: "node-3", NodeType: "test-node"}

	weights_1 := []float64{1, 1, 1, 1, 1}
	weights_2 := []float64{1}
	weights_3 := []float64{1}

	inboundConnection1 := &InboundConnection{
		NodeId:  nodeId_1,
		Weights: weights_1,
	}
	inboundConnection2 := &InboundConnection{
		NodeId:  nodeId_2,
		Weights: weights_2,
	}
	inboundConnection3 := &InboundConnection{
		NodeId:  nodeId_3,
		Weights: weights_3,
	}

	inbound := []*InboundConnection{
		inboundConnection1,
		inboundConnection2,
		inboundConnection3,
	}

	closing := make(chan chan bool)
	data := make(chan *DataMessage, len(inbound))

	neuron := &Neuron{
		ActivationFunction: activation,
		NodeId:             neuronNodeId,
		Bias:               0,
		LayerIndex:         0.5,
		Inbound:            inbound,
		Closing:            closing,
		Data:               data,
	}

	go neuron.Run()

	// send one input
	inputs_1 := []float64{20, 20, 20, 20, 20}
	dataMessage := &DataMessage{
		SenderId: nodeId_1,
		Inputs:   inputs_1,
	}
	data <- dataMessage

	// wait for output - should timeout
	time.Sleep(time.Second)

	// send rest of inputs

	// get output - should send something

	// make sure value is expected

	// send val to closing channel

	// make sure it's closed (need chan chan bool?)
	closingResponse := make(chan bool)
	closing <- closingResponse
	response := <-closingResponse
	assert.True(t, response)

}

func TestComputeScalarOutput(t *testing.T) {

	activation := identityActivationFunction()

	weights_1 := []float64{1, 1, 1, 1, 1}
	weights_2 := []float64{1}
	weights_3 := []float64{1}

	neuron := &Neuron{
		ActivationFunction: activation,
		Bias:               0,
	}

	inputs_1 := []float64{20, 20, 20, 20, 20}
	inputs_2 := []float64{10}
	inputs_3 := []float64{10}

	weightedInput1 := &weightedInput{weights: weights_1, inputs: inputs_1}
	weightedInput2 := &weightedInput{weights: weights_2, inputs: inputs_2}
	weightedInput3 := &weightedInput{weights: weights_3, inputs: inputs_3}

	weightedInputs := []*weightedInput{
		weightedInput1,
		weightedInput2,
		weightedInput3,
	}

	result := neuron.computeScalarOutput(weightedInputs)

	assert.Equals(t, result, float64(120))

}
