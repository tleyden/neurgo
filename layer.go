package neurgo

import (
	"log"
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
	keys := l.Keys()
	log.Printf("keys: %v", keys)
	sort.Float64s(keys)
	log.Printf("keys sorted: %v", keys)
	eligibleKeys := make([]float64, 0)
	for _, layerIndexKey := range keys {
		if layerIndexKey < layerIndex {
			eligibleKeys = append(eligibleKeys, layerIndexKey)
		}
	}
	log.Printf("eligible: %v", eligibleKeys)
	chosenKeyIndex := RandomIntInRange(0, len(eligibleKeys))
	chosenLayerIndex := keys[chosenKeyIndex]
	neuronsChosenLayer := l[chosenLayerIndex]
	chosenNeuronIndex := RandomIntInRange(0, len(neuronsChosenLayer))
	return neuronsChosenLayer[chosenNeuronIndex]
}

func (l LayerToNeuronMap) ChooseNeuronFollowingLayer(layerIndex float64) *Neuron {
	return nil
}
