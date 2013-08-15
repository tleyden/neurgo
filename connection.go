package neurgo

import (
	"encoding/json"
	"fmt"
	"log"
)

type InboundConnection struct {
	NodeId  *NodeId
	Weights []float64
}

type OutboundConnection struct {
	NodeId   *NodeId
	DataChan chan *DataMessage
}

type OutboundConnectable interface {
	nodeId() *NodeId
	dataChan() chan *DataMessage
}

type InboundConnectable interface {
	nodeId() *NodeId
}

type OutboundConnector interface {
	outbound() []*OutboundConnection
	setOutbound([]*OutboundConnection)
}

type InboundConnector interface {
	inbound() []*InboundConnection
	setInbound([]*InboundConnection)
}

type weightedInput struct {
	senderNodeUUID string
	weights        []float64
	inputs         []float64
}

type UUIDToInboundConnection map[string]*InboundConnection

func (connection *InboundConnection) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			NodeId  *NodeId
			Weights []float64
		}{
			NodeId:  connection.NodeId,
			Weights: connection.Weights,
		})
}

func (connection *OutboundConnection) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			NodeId *NodeId
		}{
			NodeId: connection.NodeId,
		})
}

func createEmptyWeightedInputs(inbound []*InboundConnection) []*weightedInput {

	weightedInputs := make([]*weightedInput, len(inbound))
	for i, inboundConnection := range inbound {
		weightedInput := &weightedInput{
			senderNodeUUID: inboundConnection.NodeId.UUID,
			weights:        inboundConnection.Weights,
			inputs:         nil,
		}
		weightedInputs[i] = weightedInput
	}
	return weightedInputs

}

func recordInput(weightedInputs []*weightedInput, dataMessage *DataMessage) {
	for _, weightedInput := range weightedInputs {
		if weightedInput.senderNodeUUID == dataMessage.SenderId.UUID {
			weightedInput.inputs = dataMessage.Inputs
		}
	}
}

func receiveBarrierSatisfied(weightedInputs []*weightedInput) bool {
	satisfied := true
	for _, weightedInput := range weightedInputs {
		if weightedInput.inputs == nil {
			satisfied = false
			break
		}

	}
	return satisfied
}

func (connection *OutboundConnection) String() string {
	return JsonString(connection)
}

func (connection *InboundConnection) String() string {
	return JsonString(connection)
}

func (weightedInput *weightedInput) String() string {
	return fmt.Sprintf("sender: %v, weights: %v, inputs: %v",
		weightedInput.senderNodeUUID,
		weightedInput.weights,
		weightedInput.inputs,
	)
}

func ConnectOutbound(connector OutboundConnector, connectable OutboundConnectable) *OutboundConnection {
	if connector.outbound() == nil {
		connector.setOutbound(make([]*OutboundConnection, 0))
	}

	if connectable.dataChan() == nil {
		log.Panicf("Cannot make outbound connection, dataChan == nil")
	}

	connection := &OutboundConnection{
		NodeId:   connectable.nodeId(),
		DataChan: connectable.dataChan(),
	}

	connector.setOutbound(append(connector.outbound(), connection))
	return connection

}

func ConnectInbound(connector InboundConnector, connectable InboundConnectable) *InboundConnection {
	return ConnectInboundWeighted(connector, connectable, nil)
}

func ConnectInboundWeighted(connector InboundConnector, connectable InboundConnectable, weights []float64) *InboundConnection {
	if connector.inbound() == nil {
		connector.setInbound(make([]*InboundConnection, 0))
	}

	connection := &InboundConnection{
		NodeId:  connectable.nodeId(),
		Weights: weights,
	}

	connector.setInbound(append(connector.inbound(), connection))
	return connection
}
