package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"testing"
)

func TestCreateEmptyWeightedInputs(t *testing.T) {

	nodeId_1 := &NodeId{UUID: "node-1", NodeType: NEURON}
	nodeId_2 := &NodeId{UUID: "node-2", NodeType: NEURON}

	weights_1 := []float64{1, 1, 1, 1, 1}
	weights_2 := []float64{1}

	inboundConnection1 := &InboundConnection{
		NodeId:  nodeId_1,
		Weights: weights_1,
	}
	inboundConnection2 := &InboundConnection{
		NodeId:  nodeId_2,
		Weights: weights_2,
	}

	inbound := []*InboundConnection{
		inboundConnection1,
		inboundConnection2,
	}

	weightedInputs := createEmptyWeightedInputs(inbound)
	assert.Equals(t, len(inbound), len(weightedInputs))
	assert.Equals(t, inbound[0].NodeId, weightedInputs[0].senderNodeId)

}
