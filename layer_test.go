package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"testing"
)

func TestChooseRandomLayer(t *testing.T) {
	layerToNeuronMap := make(LayerToNeuronMap)
	neurons := make([]*Neuron, 0)
	layerToNeuronMap[0.0] = neurons
	layerToNeuronMap[0.25] = neurons

	foundFirstLayer := false
	foundSecondLayer := false

	for i := 0; i < 20; i++ {
		layerIndex := layerToNeuronMap.ChooseRandomLayer()
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

func TestChooseNeuronPrecedingLayer(t *testing.T) {

	layerToNeuronMap := make(LayerToNeuronMap)

	neuron1 := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("neuron1", 0.25),
		Bias:               -30,
	}
	neurons25 := []*Neuron{neuron1}
	layerToNeuronMap[0.25] = neurons25

	neuron2 := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("neuron2", 0.35),
		Bias:               -30,
	}
	neurons35 := []*Neuron{neuron2}
	layerToNeuronMap[0.35] = neurons35

	neuron3 := &Neuron{
		ActivationFunction: EncodableSigmoid(),
		NodeId:             NewNeuronId("neuron3", 0.45),
		Bias:               -30,
	}
	neurons45 := []*Neuron{neuron3}
	layerToNeuronMap[0.45] = neurons45

	foundNeuron1 := false
	foundNeuron2 := false
	foundNeuron3 := false

	for i := 0; i < 20; i++ {
		chosenNeuron := layerToNeuronMap.ChooseNeuronPrecedingLayer(0.45)
		switch chosenNeuron.NodeId.UUID {
		case neuron1.NodeId.UUID:
			foundNeuron1 = true
		case neuron2.NodeId.UUID:
			foundNeuron2 = true
		case neuron3.NodeId.UUID:
			foundNeuron3 = true
		}

	}

	assert.True(t, foundNeuron1)
	assert.True(t, foundNeuron2)
	assert.False(t, foundNeuron3)

}
