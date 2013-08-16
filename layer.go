package neurgo

import (
	"sort"
)

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

func (l LayerToNeuronMap) ChooseNeuronPrecedingLayer(layerIndex float64) *Neuron {
	chooser := func(layerIndexKey float64) bool {
		return layerIndexKey < layerIndex
	}
	return l.chooseNeuronFromLayer(chooser)
}

func (l LayerToNeuronMap) ChooseNeuronFollowingLayer(layerIndex float64) *Neuron {
	chooser := func(layerIndexKey float64) bool {
		return layerIndexKey > layerIndex
	}
	return l.chooseNeuronFromLayer(chooser)
}

func (l LayerToNeuronMap) chooseNeuronFromLayer(chooser func(float64) bool) *Neuron {
	keys := l.Keys()
	sort.Float64s(keys)
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
	neuronsChosenLayer := l[chosenLayerIndex]
	chosenNeuronIndex := RandomIntInRange(0, len(neuronsChosenLayer))
	chosenNeuron := neuronsChosenLayer[chosenNeuronIndex]
	return chosenNeuron

}
