package neurgo

import (
	"log"
	"math"
	"sort"
)

type LayerToNodeIdMap map[float64][]*NodeId

func (layerToNodeIdMap LayerToNodeIdMap) Keys() []float64 {
	// TODO: better/easier way to get list of keys?
	keys := make([]float64, len(layerToNodeIdMap))
	i := 0
	for key, _ := range layerToNodeIdMap {
		keys[i] = key
		i += 1
	}
	sort.Float64s(keys)
	return keys
}

func (l LayerToNodeIdMap) IntegerIndexOf(layer float64) int {
	keys := l.Keys()
	for i, key := range keys {
		if key == layer {
			return i
		}
	}
	log.Panicf("Unable to find integer index of layer: %v", layer)
	return -1
}

func (l LayerToNodeIdMap) LayerOfIntegerIndex(layerInteger int) float64 {
	keys := l.Keys()
	return keys[layerInteger]
}

func (l LayerToNodeIdMap) NewLayerBetween(initial, final float64) float64 {
	return (initial + final) / 2.0
}

func (l LayerToNodeIdMap) LayersAdjacent(initial, final float64) bool {
	if initial == final {
		return true
	}

	initialIntegerIndex := l.IntegerIndexOf(initial)
	finalIntegerIndex := l.IntegerIndexOf(final)

	delta := finalIntegerIndex - initialIntegerIndex
	return math.Abs(float64(delta)) == 1.0

}

func (l LayerToNodeIdMap) LayerBetweenOrNew(initial, final float64) float64 {

	if initial == final {
		return initial
	}

	initialIntegerIndex := l.IntegerIndexOf(initial)
	finalIntegerIndex := l.IntegerIndexOf(final)

	delta := finalIntegerIndex - initialIntegerIndex
	if math.Abs(float64(delta)) == 1.0 {
		// adjacent layer, make a new layer
		return l.NewLayerBetween(initial, final)
	} else if delta > 1 {
		// there is an existing layer between?  use it
		nextLayerIntegerIndex := initialIntegerIndex + 1
		nextLayerFractalIndex := l.LayerOfIntegerIndex(nextLayerIntegerIndex)
		return nextLayerFractalIndex
	} else if delta < -1 {
		// there is an existing layer between?  use it
		nextLayerIntegerIndex := initialIntegerIndex - 1
		nextLayerFractalIndex := l.LayerOfIntegerIndex(nextLayerIntegerIndex)
		return nextLayerFractalIndex
	} else {
		log.Panicf("Unexpected delta: %v", delta)
		return -1.0
	}

}

func (layerToNodeIdMap LayerToNodeIdMap) ChooseRandomLayer() float64 {
	keys := layerToNodeIdMap.Keys()
	randomKeyIndex := RandomIntInRange(0, len(keys))
	return keys[randomKeyIndex]
}

func (l LayerToNodeIdMap) ChooseNodeIdPrecedingLayer(layerIndex float64) *NodeId {
	chooser := func(layerIndexKey float64) bool {
		return layerIndexKey < layerIndex
	}
	return l.chooseNodeIdFromLayer(chooser)
}

func (l LayerToNodeIdMap) ChooseNodeIdFollowingLayer(layerIndex float64) *NodeId {
	chooser := func(layerIndexKey float64) bool {
		return layerIndexKey > layerIndex
	}
	return l.chooseNodeIdFromLayer(chooser)
}

func (l LayerToNodeIdMap) chooseNodeIdFromLayer(chooser func(float64) bool) *NodeId {
	keys := l.Keys()
	eligibleKeys := make([]float64, 0)
	for _, layerIndexKey := range keys {
		if chooser(layerIndexKey) == true {
			eligibleKeys = append(eligibleKeys, layerIndexKey)
		}
	}
	if len(eligibleKeys) == 0 {
		return nil
	}
	chosenKeyIndex := RandomIntInRange(0, len(eligibleKeys))
	chosenLayerIndex := eligibleKeys[chosenKeyIndex]
	nodeIdsChosenLayer := l[chosenLayerIndex]
	chosenNodeIdIndex := RandomIntInRange(0, len(nodeIdsChosenLayer))
	chosenNodeId := nodeIdsChosenLayer[chosenNodeIdIndex]
	return chosenNodeId

}

// LayerToNeuronMap ..
type LayerToNeuronMap map[float64][]*Neuron

func (layerToNeuronMap LayerToNeuronMap) Keys() []float64 {
	// TODO: better/easier way to get list of keys?
	keys := make([]float64, len(layerToNeuronMap))
	i := 0
	for key, _ := range layerToNeuronMap {
		keys[i] = key
		i += 1
	}
	return keys
}

func (layerToNeuronMap LayerToNeuronMap) ChooseRandomLayer() float64 {
	keys := layerToNeuronMap.Keys()
	randomKeyIndex := RandomIntInRange(0, len(keys))
	return keys[randomKeyIndex]
}
