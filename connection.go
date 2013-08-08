package neurgo

import (
	"encoding/json"
	"fmt"
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

type weightedInput struct {
	senderNodeId *NodeId
	weights      []float64
	inputs       []float64
}

func (connection *InboundConnection) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			NodeId *NodeId
		}{
			NodeId: connection.NodeId,
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
			senderNodeId: inboundConnection.NodeId,
			weights:      inboundConnection.Weights,
			inputs:       nil,
		}
		weightedInputs[i] = weightedInput
	}
	return weightedInputs

}

func recordInput(weightedInputs []*weightedInput, dataMessage *DataMessage) {
	for _, weightedInput := range weightedInputs {
		if weightedInput.senderNodeId == dataMessage.SenderId {
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
	return fmt.Sprintf("node: %v, datachan: %v",
		connection.NodeId,
		connection.DataChan,
	)
}

func (weightedInput *weightedInput) String() string {
	return fmt.Sprintf("sender: %v, weights: %v, inputs: %v",
		weightedInput.senderNodeId,
		weightedInput.weights,
		weightedInput.inputs,
	)
}
