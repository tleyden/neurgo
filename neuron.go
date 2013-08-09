package neurgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/proxypoke/vector"
	"log"
	"sync"
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
	wg                 *sync.WaitGroup
}

func (neuron *Neuron) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			NodeId   *NodeId
			Bias     float64
			Inbound  []*InboundConnection
			Outbound []*OutboundConnection
		}{
			NodeId:   neuron.NodeId,
			Bias:     neuron.Bias,
			Inbound:  neuron.Inbound,
			Outbound: neuron.Outbound,
		})
}

func (neuron *Neuron) Run() {

	log.Printf("%v Run() started", neuron.NodeId.UUID)

	defer neuron.wg.Done()

	neuron.checkRunnable()

	neuron.sendEmptySignalRecurrentOutbound()

	weightedInputs := createEmptyWeightedInputs(neuron.Inbound)

	closed := false

	for {

		log.Printf("Neuron %v select().  datachan: %v", neuron.NodeId.UUID, neuron.DataChan)

		select {
		case responseChan := <-neuron.Closing:
			closed = true
			responseChan <- true
			break // TODO: do we need this for anything??
		case dataMessage := <-neuron.DataChan:
			log.Printf("Neuron %v recording input: %v", neuron.NodeId.UUID, dataMessage)
			recordInput(weightedInputs, dataMessage)
			log.Printf("Neuron %v new weightedInputs: %v", neuron.NodeId.UUID, weightedInputs)
		}

		if closed {
			neuron.Closing = nil
			neuron.DataChan = nil
			break
		}

		if receiveBarrierSatisfied(weightedInputs) {

			log.Printf("Neuron %v barrier satisfied via inputs: %v", neuron.NodeId.UUID, weightedInputs)
			scalarOutput := neuron.computeScalarOutput(weightedInputs)

			dataMessage := &DataMessage{
				SenderId: neuron.NodeId,
				Inputs:   []float64{scalarOutput},
			}

			neuron.scatterOutput(dataMessage)

			weightedInputs = createEmptyWeightedInputs(neuron.Inbound)

		} else {
			log.Printf("Neuron %v receive barrier not satisfied.  weightedInputs: %v", neuron.NodeId.UUID, weightedInputs)
		}

	}

	log.Printf("%v Run() finished", neuron.NodeId.UUID)

}

func (neuron *Neuron) String() string {
	return JsonString(neuron)
}

func (neuron *Neuron) ConnectOutbound(connectable OutboundConnectable) {
	if neuron.Outbound == nil {
		neuron.Outbound = make([]*OutboundConnection, 0)
	}
	connection := &OutboundConnection{
		NodeId:   connectable.nodeId(),
		DataChan: connectable.dataChan(),
	}
	neuron.Outbound = append(neuron.Outbound, connection)
}

func (neuron *Neuron) ConnectInboundWeighted(connectable InboundConnectable, weights []float64) {
	if neuron.Inbound == nil {
		neuron.Inbound = make([]*InboundConnection, 0)
	}
	connection := &InboundConnection{
		NodeId:  connectable.nodeId(),
		Weights: weights,
	}
	neuron.Inbound = append(neuron.Inbound, connection)

}

// In order to prevent deadlock, any neurons we have recurrent outbound
// connections to must be "primed" by sending an empty signal.  A recurrent
// outbound connection simply means that it's a connection to ourself or
// to a neuron in a previous (eg, to the left) layer.  If we didn't do this,
// that previous neuron would be waiting forever for a signal that will
// never come, because this neuron wouldn't fire until it got a signal.
func (neuron *Neuron) sendEmptySignalRecurrentOutbound() {

	recurrentConnections := neuron.recurrentOutboundConnections()
	for _, recurrentConnection := range recurrentConnections {

		inputs := []float64{0}
		dataMessage := &DataMessage{
			SenderId: neuron.NodeId,
			Inputs:   inputs,
		}
		recurrentConnection.DataChan <- dataMessage
	}

}

// Find the subset of outbound connections which are "recurrent" - meaning
// that the connection is to this neuron itself, or to a neuron in a previous
// (eg, to the left) layer.
func (neuron *Neuron) recurrentOutboundConnections() []*OutboundConnection {
	result := make([]*OutboundConnection, 0)
	for _, outboundConnection := range neuron.Outbound {
		if neuron.isConnectionRecurrent(outboundConnection) {
			result = append(result, outboundConnection)
		}
	}
	return result
}

// a connection is considered recurrent if it has a connection
// to itself or to a node in a previous layer.  Previous meaning
// if you look at a feedforward from left to right, with the input
// layer being on the far left, and output layer on the far right,
// then any layer to the left is considered previous.
func (neuron *Neuron) isConnectionRecurrent(connection *OutboundConnection) bool {
	if connection.NodeId.LayerIndex <= neuron.NodeId.LayerIndex {
		return true
	}
	return false
}

func (neuron *Neuron) scatterOutput(dataMessage *DataMessage) {
	for _, outboundConnection := range neuron.Outbound {
		dataChan := outboundConnection.DataChan
		log.Printf("Neuron %v scatter %v to: %v", neuron.NodeId.UUID, dataMessage, outboundConnection)
		dataChan <- dataMessage
	}
}

func (neuron *Neuron) Init() {
	if neuron.Closing == nil {
		neuron.Closing = make(chan chan bool)
	}

	if neuron.DataChan == nil {
		neuron.DataChan = make(chan *DataMessage, len(neuron.Inbound))
	}

	if neuron.ActivationFunction == nil {

		// TODO: fix this .. we need to serialize the name of
		// the function, and when we deserialize, resolve to
		// actual function
		neuron.ActivationFunction = Sigmoid
	}

	if neuron.wg == nil {
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
		log.Printf("inputVector: %v", inputVector)
		log.Printf("weightVector: %v", weightVector)
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
