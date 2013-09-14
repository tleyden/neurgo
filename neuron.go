package neurgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/couchbaselabs/logg"
	"github.com/proxypoke/vector"
	"log"
	"sync"
	"time"
)

type Neuron struct {
	NodeId             *NodeId
	Bias               float64
	Inbound            []*InboundConnection
	Outbound           []*OutboundConnection
	Closing            chan chan bool
	DataChan           chan *DataMessage
	ActivationFunction *EncodableActivation
	wg                 *sync.WaitGroup
	Cortex             *Cortex
}

func (neuron *Neuron) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			NodeId             *NodeId
			Bias               float64
			Inbound            []*InboundConnection
			Outbound           []*OutboundConnection
			ActivationFunction *EncodableActivation
		}{
			NodeId:             neuron.NodeId,
			Bias:               neuron.Bias,
			Inbound:            neuron.Inbound,
			Outbound:           neuron.Outbound,
			ActivationFunction: neuron.ActivationFunction,
		})
}

func (neuron *Neuron) Run() {

	defer neuron.wg.Done()

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
			neuron.logDataReceive(dataMessage)
			recordInput(weightedInputs, dataMessage)
		}

		if closed {
			neuron.Closing = nil
			neuron.DataChan = nil
			break
		}

		if receiveBarrierSatisfied(weightedInputs) {

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
	return JsonString(neuron)
}

func (neuron *Neuron) ConnectOutbound(connectable OutboundConnectable) *OutboundConnection {
	return ConnectOutbound(neuron, connectable)
}

func (neuron *Neuron) ConnectInboundWeighted(connectable InboundConnectable, weights []float64) *InboundConnection {
	return ConnectInboundWeighted(neuron, connectable, weights)
}

func (neuron *Neuron) outbound() []*OutboundConnection {
	return neuron.Outbound
}

func (neuron *Neuron) setOutbound(newOutbound []*OutboundConnection) {
	neuron.Outbound = newOutbound
}

func (neuron *Neuron) inbound() []*InboundConnection {
	return neuron.Inbound
}

func (neuron *Neuron) setInbound(newInbound []*InboundConnection) {
	neuron.Inbound = newInbound
}

// In order to prevent deadlock, any neurons we have recurrent outbound
// connections to must be "primed" by sending an empty signal.  A recurrent
// outbound connection simply means that it's a connection to ourself or
// to a neuron in a previous (eg, to the left) layer.  If we didn't do this,
// that previous neuron would be waiting forever for a signal that will
// never come, because this neuron wouldn't fire until it got a signal.
func (neuron *Neuron) sendEmptySignalRecurrentOutbound() {

	recurrentConnections := neuron.RecurrentOutboundConnections()

	for _, recurrentConnection := range recurrentConnections {

		inputs := []float64{0}
		dataMessage := &DataMessage{
			SenderId: neuron.NodeId,
			Inputs:   inputs,
		}
		if recurrentConnection.DataChan == nil {
			log.Panicf("Can't sendEmptySignalRecurrentOutbound to %v, DataChan is nil", recurrentConnection)
		}

		select {
		case recurrentConnection.DataChan <- dataMessage:
		case <-time.After(time.Second):
			log.Panicf("Timeout sendEmptySignalRecurrentOutbound to %v, DataChan is nil", recurrentConnection)
		}

	}

}

// Find the subset of outbound connections which are "recurrent" - meaning
// that the connection is to this neuron itself, or to a neuron in a previous
// (eg, to the left) layer.
func (neuron *Neuron) RecurrentOutboundConnections() []*OutboundConnection {
	result := make([]*OutboundConnection, 0)
	for _, outboundConnection := range neuron.Outbound {
		if neuron.IsConnectionRecurrent(outboundConnection) {
			result = append(result, outboundConnection)
		}
	}
	return result
}

func (neuron *Neuron) RecurrentInboundConnections() []*InboundConnection {
	result := make([]*InboundConnection, 0)
	for _, inboundConnection := range neuron.Inbound {
		if neuron.IsInboundConnectionRecurrent(inboundConnection) {
			result = append(result, inboundConnection)
		}
	}
	return result
}

// a connection is considered recurrent if it has a connection
// to itself or to a node in a previous layer.  Previous meaning
// if you look at a feedforward from left to right, with the input
// layer being on the far left, and output layer on the far right,
// then any layer to the left is considered previous.
func (neuron *Neuron) IsConnectionRecurrent(connection *OutboundConnection) bool {
	if connection.NodeId.LayerIndex <= neuron.NodeId.LayerIndex {
		return true
	}
	return false
}

// same as isConnectionRecurrent, but for inbound connections
// TODO: use interfaces to eliminate code duplication
func (neuron *Neuron) IsInboundConnectionRecurrent(connection *InboundConnection) bool {
	if neuron.NodeId.LayerIndex <= connection.NodeId.LayerIndex {
		return true
	}
	return false
}

func (neuron *Neuron) scatterOutput(dataMessage *DataMessage) {

	for _, outboundConnection := range neuron.Outbound {
		logmsg := fmt.Sprintf("%v -> %v: %v", neuron.NodeId.UUID,
			outboundConnection.NodeId.UUID, dataMessage)
		logg.LogTo("NODE_SEND", logmsg)
		dataChan := outboundConnection.DataChan
		dataChan <- dataMessage
	}
}

// Initialize/re-initialize the neuron.
// reInit: basically this is a messy hack to solve the issue:
// - neuron.Init() function is called and DataChan buffer len = X
// - new recurrent connections are added
// - since the DataChan buffer len is X, and needs to be X+1, network is wedged
// So by doing a "destructive reInit" it will rebuild all DataChan's
// and all outbound connections which contain DataChan's, thus solving
// the problem.
// TODO: fix this hack
func (neuron *Neuron) Init(reInit bool) {
	if reInit == true {
		neuron.Closing = make(chan chan bool)
	} else if neuron.Closing == nil {
		neuron.Closing = make(chan chan bool)
	}

	if reInit == true {
		neuron.DataChan = make(chan *DataMessage, len(neuron.Inbound))
	} else if neuron.DataChan == nil {
		neuron.DataChan = make(chan *DataMessage, len(neuron.Inbound))
	}

	if reInit == true {
		neuron.wg = &sync.WaitGroup{}
		neuron.wg.Add(1)
	} else if neuron.wg == nil {
		neuron.wg = &sync.WaitGroup{}
		neuron.wg.Add(1)
	}

}

func (neuron *Neuron) Shutdown() {

	closingResponse := make(chan bool)
	neuron.Closing <- closingResponse
	response := <-closingResponse
	if response != true {
		log.Panicf("Got unexpected response on closing channel")
	}

	neuron.shutdownOutboundConnections()

	neuron.wg.Wait()
	neuron.wg = nil
}

func (neuron *Neuron) InboundUUIDMap() UUIDToInboundConnection {
	inboundUUIDMap := make(UUIDToInboundConnection)
	for _, connection := range neuron.Inbound {
		inboundUUIDMap[connection.NodeId.UUID] = connection
	}
	return inboundUUIDMap
}

func (neuron *Neuron) Copy() *Neuron {

	// serialize to json
	jsonBytes, err := json.Marshal(neuron)
	if err != nil {
		log.Fatal(err)
	}

	// new neuron
	neuronCopy := &Neuron{}

	// deserialize json into new neuron
	err = json.Unmarshal(jsonBytes, neuronCopy)
	if err != nil {
		log.Fatal(err)
	}

	return neuronCopy

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

	if neuron.ActivationFunction == nil {
		msg := fmt.Sprintf("not expecting neuron.ActivationFunction to be nil")
		panic(msg)
	}

	if err := neuron.validateOutbound(); err != nil {
		msg := fmt.Sprintf("invalid outbound connection(s): %v", err.Error())
		panic(msg)
	}

	inboundRecurrentCxns := neuron.RecurrentInboundConnections()
	if cap(neuron.DataChan) < len(inboundRecurrentCxns) {
		msg := fmt.Sprintf("dataChan buffer capacity %v less than # of inbound recurrent connections: %v", cap(neuron.DataChan), len(inboundRecurrentCxns))
		panic(msg)
	}

}

func (neuron *Neuron) validateOutbound() error {
	for _, connection := range neuron.Outbound {
		if connection.DataChan == nil {
			msg := fmt.Sprintf("%v has empty DataChan", connection)
			return errors.New(msg)
		}
	}
	return nil
}

func (neuron *Neuron) computeScalarOutput(weightedInputs []*weightedInput) float64 {
	output := neuron.weightedInputDotProductSum(weightedInputs)
	output += neuron.Bias
	output = neuron.ActivationFunction.ActivationFunction(output)
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

func (neuron *Neuron) dataChan() chan *DataMessage {
	return neuron.DataChan
}

func (neuron *Neuron) nodeId() *NodeId {
	return neuron.NodeId
}

func (neuron *Neuron) initOutboundConnections(nodeIdToDataMsg nodeIdToDataMsgMap) {
	for _, outboundConnection := range neuron.Outbound {
		if outboundConnection.DataChan == nil {
			dataChan := nodeIdToDataMsg[outboundConnection.NodeId.UUID]
			if dataChan != nil {
				outboundConnection.DataChan = dataChan
			}
		}
	}
}

func (neuron *Neuron) shutdownOutboundConnections() {
	for _, outboundConnection := range neuron.Outbound {
		outboundConnection.DataChan = nil
	}
}

func (neuron *Neuron) logDataReceive(dataMessage *DataMessage) {
	sender := dataMessage.SenderId.UUID
	logmsg := fmt.Sprintf("%v -> %v: %v", sender,
		neuron.NodeId.UUID, dataMessage)
	logg.LogTo("NODE_RECV", logmsg)
}
