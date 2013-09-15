package neurgo

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
	"time"
)

func identityActivationFunction() ActivationFunction {
	return func(x float64) float64 { return x }
}

func encodableIdentityActivationFunction() *EncodableActivation {
	return &EncodableActivation{
		Name:               "identityActivation",
		ActivationFunction: identityActivationFunction(),
	}
}

func TestNeuronJsonMarshal(t *testing.T) {
	neuronNodeId := NewNeuronId("neuron", 0.25)
	neuron := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             neuronNodeId,
		Bias:               -30,
	}
	json, err := json.Marshal(neuron)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)
	jsonString := fmt.Sprintf("%s", json)
	log.Printf("jsonString: %v", jsonString)
}

func TestRecurrentNeuron(t *testing.T) {

	// injector -> n1 -> n2 -> wiretap where n2 has recurrent connection
	// back to n1.  send a two passes of inputs, and make sure the
	// wiretap gets a signal

	activation := encodableIdentityActivationFunction()

	shouldReInit := false

	injectorNodeId_1 := &NodeId{
		UUID:       "injector-1",
		NodeType:   SENSOR,
		LayerIndex: 0.0,
	}

	neuron1NodeId := &NodeId{
		UUID:       "neuron1",
		NodeType:   NEURON,
		LayerIndex: 0.125,
	}

	neuron2NodeId := &NodeId{
		UUID:       "neuron2",
		NodeType:   NEURON,
		LayerIndex: 0.25,
	}

	inboundConnectionToN2 := &InboundConnection{
		NodeId:  neuron1NodeId,
		Weights: []float64{1},
	}

	inboundN2 := []*InboundConnection{inboundConnectionToN2}

	closingN2 := make(chan chan bool)
	dataN2 := make(chan *DataMessage, 1)

	inboundConnectionToN1 := &InboundConnection{
		NodeId:  injectorNodeId_1,
		Weights: []float64{1, 1, 1, 1, 1},
	}

	recurrentInboundConnectionToN1 := &InboundConnection{
		NodeId:  neuron2NodeId,
		Weights: []float64{1},
	}

	outboundConnectionToN2 := &OutboundConnection{
		NodeId:   neuron2NodeId,
		DataChan: dataN2,
	}

	outbound := []*OutboundConnection{
		outboundConnectionToN2,
	}

	inboundN1 := []*InboundConnection{
		inboundConnectionToN1,
		recurrentInboundConnectionToN1,
	}

	closingN1 := make(chan chan bool)
	dataN1 := make(chan *DataMessage, 1)

	neuronN1 := &Neuron{
		ActivationFunction: activation,
		NodeId:             neuron1NodeId,
		Bias:               20,
		Inbound:            inboundN1,
		Outbound:           outbound,
		Closing:            closingN1,
		DataChan:           dataN1,
	}

	wiretapNodeId := &NodeId{
		UUID:       "wireteap-node",
		NodeType:   ACTUATOR,
		LayerIndex: 0.5,
	}
	wiretapDataChan := make(chan *DataMessage, 1)
	wiretapConnection := &OutboundConnection{
		NodeId:   wiretapNodeId,
		DataChan: wiretapDataChan,
	}

	n2ToN1Connection := &OutboundConnection{
		NodeId:   neuron1NodeId,
		DataChan: dataN1,
	}

	outboundN2 := []*OutboundConnection{
		wiretapConnection,
		n2ToN1Connection,
	}

	neuronN2 := &Neuron{
		ActivationFunction: activation,
		NodeId:             neuron2NodeId,
		Bias:               20,
		Inbound:            inboundN2,
		Outbound:           outboundN2,
		Closing:            closingN2,
		DataChan:           dataN2,
	}

	recurrentConnections := neuronN2.RecurrentOutboundConnections()
	assert.Equals(t, len(recurrentConnections), 1)

	neuronN1.Init(shouldReInit)
	neuronN2.Init(shouldReInit)

	go neuronN1.Run()
	go neuronN2.Run()

	// send one input
	inputs_1 := []float64{20, 20, 20, 20, 20}
	dataMessage := &DataMessage{
		SenderId: injectorNodeId_1,
		Inputs:   inputs_1,
	}
	neuronN1.DataChan <- dataMessage

	// wait for output
	log.Printf("wiretapDataChan: %v", wiretapDataChan)
	select {
	case outputDataMessage := <-wiretapDataChan:
		outputVector := outputDataMessage.Inputs
		log.Printf("outputVector: %v", outputVector)
		outputValue := outputVector[0]
		expectedOut := 100 + 20 + 20 // inputs plus two biases
		assert.Equals(t, int(outputValue), expectedOut)
	case <-time.After(time.Second):
		assert.Errorf(t, "Did not get result at wiretap")
	}

	// send another input
	inputs_2 := []float64{20, 20, 20, 20, 20}
	dataMessage2 := &DataMessage{
		SenderId: injectorNodeId_1,
		Inputs:   inputs_2,
	}
	neuronN1.DataChan <- dataMessage2

	// wait for output
	select {
	case outputDataMessage2 := <-wiretapDataChan:
		outputVector2 := outputDataMessage2.Inputs
		log.Printf("outputVector #2: %v", outputVector2)
		outputValue2 := outputVector2[0]
		expectedOut2 := (100 + 20 + 20) * 2
		assert.Equals(t, int(outputValue2), expectedOut2)
	case <-time.After(time.Second):
		assert.Errorf(t, "Did not get result #2 at wiretap")
	}

	neuronN1.Shutdown()
	neuronN2.Shutdown()

}

func TestRunningNeuron(t *testing.T) {

	log.Printf("")

	shouldReInit := false

	activation := encodableIdentityActivationFunction()

	neuronNodeId := &NodeId{
		UUID:       "neuron",
		NodeType:   NEURON,
		LayerIndex: 0.25,
	}
	nodeId_1 := &NodeId{UUID: "node-1", NodeType: SENSOR, LayerIndex: 0.0}
	nodeId_2 := &NodeId{UUID: "node-2", NodeType: SENSOR, LayerIndex: 0.0}
	nodeId_3 := &NodeId{UUID: "node-3", NodeType: SENSOR, LayerIndex: 0.0}

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

	wiretapNodeId := &NodeId{
		UUID:       "wireteap-node",
		NodeType:   ACTUATOR,
		LayerIndex: 0.5,
	}
	wiretapDataChan := make(chan *DataMessage, 1)
	wiretapConnection := &OutboundConnection{
		NodeId:   wiretapNodeId,
		DataChan: wiretapDataChan,
	}
	outbound := []*OutboundConnection{
		wiretapConnection,
	}

	neuron := &Neuron{
		ActivationFunction: activation,
		NodeId:             neuronNodeId,
		Bias:               20,
		Inbound:            inbound,
		Outbound:           outbound,
		Closing:            closing,
		DataChan:           data,
	}

	neuron.Init(shouldReInit)
	go neuron.Run()

	// send one input
	inputs_1 := []float64{20, 20, 20, 20, 20}
	dataMessage := &DataMessage{
		SenderId: nodeId_1,
		Inputs:   inputs_1,
	}
	data <- dataMessage

	// wait for output - should timeout
	select {
	case output := <-wiretapDataChan:
		assert.Errorf(t, "Got unexpected output: %v", output)
	case <-time.After(time.Second / 100):
	}

	// send rest of inputs
	inputs_2 := []float64{20}
	dataMessage2 := &DataMessage{
		SenderId: nodeId_2,
		Inputs:   inputs_2,
	}
	data <- dataMessage2

	inputs_3 := []float64{20}
	dataMessage3 := &DataMessage{
		SenderId: nodeId_3,
		Inputs:   inputs_3,
	}
	data <- dataMessage3

	// get output - should receive something
	select {
	case outputDataMessage := <-wiretapDataChan:
		outputVector := outputDataMessage.Inputs
		outputValue := outputVector[0]
		assert.Equals(t, int(outputValue), int(160))
	case <-time.After(time.Second):
		assert.Errorf(t, "Timed out waiting for output")
	}

	neuron.Shutdown()

}

func TestComputeScalarOutput(t *testing.T) {

	activation := encodableIdentityActivationFunction()

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

func TestNeuronShutdown(t *testing.T) {

	shouldReInit := false

	sensor := &Sensor{
		NodeId:       NewSensorId("sensor", 0.0),
		VectorLength: 2,
	}
	sensor.Init(shouldReInit)

	neuron := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("neuron", 0.35),
		Bias:               -10,
	}
	neuron.Init(shouldReInit)

	sensor.ConnectOutbound(neuron)
	neuron.ConnectInboundWeighted(sensor, []float64{20, 20})

	go sensor.Run()
	go neuron.Run()

	sensor.Shutdown()
	neuron.Shutdown()

	outboundConnection := sensor.Outbound[0]
	assert.True(t, outboundConnection.DataChan == nil)

	assert.True(t, neuron.wg == nil)
	assert.True(t, sensor.wg == nil)

}

func TestRecurrentOutboundConnections(t *testing.T) {

	// make a recurrent connection
	neuron1NodeId := &NodeId{
		UUID:       "neuron1",
		NodeType:   NEURON,
		LayerIndex: 0.0,
	}

	neuron2NodeId := &NodeId{
		UUID:       "neuron2",
		NodeType:   NEURON,
		LayerIndex: 0.5,
	}

	outboundConnectionN2ToN1 := &OutboundConnection{
		NodeId:   neuron1NodeId,
		DataChan: make(chan *DataMessage, 1),
	}

	outboundN2 := []*OutboundConnection{
		outboundConnectionN2ToN1,
	}

	neuronN2 := &Neuron{
		ActivationFunction: nil,
		NodeId:             neuron2NodeId,
		Bias:               20,
		Inbound:            nil,
		Outbound:           outboundN2,
		Closing:            nil,
		DataChan:           nil,
	}

	recurrentConnections := neuronN2.RecurrentOutboundConnections()

	assert.Equals(t, len(recurrentConnections), 1)

}
