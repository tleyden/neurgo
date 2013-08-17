package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"testing"
)

func fakeLayerToNodeIdMap() (LayerToNodeIdMap, *NodeId, *NodeId, *NodeId) {
	layerToNodeIdMap := make(LayerToNodeIdMap)

	nodeId1 := NewNeuronId("nodeId1", 0.25)
	nodeIds25 := []*NodeId{nodeId1}
	layerToNodeIdMap[0.25] = nodeIds25

	nodeId2 := NewNeuronId("nodeId2", 0.35)
	nodeIds35 := []*NodeId{nodeId2}
	layerToNodeIdMap[0.35] = nodeIds35

	nodeId3 := NewNeuronId("nodeId3", 0.45)
	nodeIds45 := []*NodeId{nodeId3}
	layerToNodeIdMap[0.45] = nodeIds45
	return layerToNodeIdMap, nodeId1, nodeId2, nodeId3

}

func TestLayerBetweenOrNew(t *testing.T) {

	layerToNodeIdMap, nodeId1, nodeId2, nodeId3 := fakeLayerToNodeIdMap()

	initialLayer := nodeId1.LayerIndex
	finalLayer := nodeId2.LayerIndex
	layer := layerToNodeIdMap.LayerBetweenOrNew(initialLayer, finalLayer)
	assert.True(t, layer > initialLayer)
	assert.True(t, layer < finalLayer)

	initialLayer = nodeId1.LayerIndex
	finalLayer = nodeId3.LayerIndex
	layer = layerToNodeIdMap.LayerBetweenOrNew(initialLayer, finalLayer)
	assert.True(t, layer == nodeId2.LayerIndex)

}

func TestChooseRandomLayer(t *testing.T) {
	layerToNodeIdMap := make(LayerToNodeIdMap)
	nodeIds := make([]*NodeId, 0)
	layerToNodeIdMap[0.0] = nodeIds
	layerToNodeIdMap[0.25] = nodeIds

	foundFirstLayer := false
	foundSecondLayer := false

	for i := 0; i < 20; i++ {
		layerIndex := layerToNodeIdMap.ChooseRandomLayer()
		if layerIndex == 0.0 {
			foundFirstLayer = true
		}
		if layerIndex == 0.25 {
			foundSecondLayer = true
		}
	}

	assert.True(t, foundFirstLayer)
	assert.True(t, foundSecondLayer)

}

func TestChooseNodeIdPrecedingLayer(t *testing.T) {

	layerToNodeIdMap, nodeId1, nodeId2, nodeId3 := fakeLayerToNodeIdMap()

	foundNodeId1 := false
	foundNodeId2 := false
	foundNodeId3 := false

	for i := 0; i < 20; i++ {
		chosenNodeId := layerToNodeIdMap.ChooseNodeIdPrecedingLayer(0.45)
		switch chosenNodeId.UUID {
		case nodeId1.UUID:
			foundNodeId1 = true
		case nodeId2.UUID:
			foundNodeId2 = true
		case nodeId3.UUID:
			foundNodeId3 = true
		}

	}

	assert.True(t, foundNodeId1)
	assert.True(t, foundNodeId2)
	assert.False(t, foundNodeId3)

}

func TestChooseNodeIdFollowingLayer(t *testing.T) {

	layerToNodeIdMap, nodeId1, nodeId2, nodeId3 := fakeLayerToNodeIdMap()

	foundNodeId1 := false
	foundNodeId2 := false
	foundNodeId3 := false

	for i := 0; i < 20; i++ {
		chosenNodeId := layerToNodeIdMap.ChooseNodeIdFollowingLayer(0.25)
		switch chosenNodeId.UUID {
		case nodeId1.UUID:
			foundNodeId1 = true
		case nodeId2.UUID:
			foundNodeId2 = true
		case nodeId3.UUID:
			foundNodeId3 = true
		}

	}

	assert.False(t, foundNodeId1)
	assert.True(t, foundNodeId2)
	assert.True(t, foundNodeId3)

}
