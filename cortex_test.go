package neurgo

import (
	"encoding/json"
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	"log"
	"testing"
)

func TestCortexCopy(t *testing.T) {

	logg.LogKeys["DEBUG"] = true

	xnorCortex := XnorCortex()
	xnorCortexCopy := xnorCortex.Copy()

	for _, neuron := range xnorCortex.Neurons {
		assert.False(t, neuron.Cortex == nil)
	}

	// inputs + expected outputs
	examples := XnorTrainingSamples()

	// get the fitness
	fitness := xnorCortexCopy.Fitness(examples)

	assert.True(t, fitness >= FITNESS_THRESHOLD)

}

func TestCortexJsonMarshal(t *testing.T) {
	xnorCortex := XnorCortex()
	xnorCortex.MarshalJSONToFile("/tmp/output.json")
	// TODO: add some assertions about this  (try reading file back in)
}

func TestCortexJsonUnmarshal(t *testing.T) {

	jsonString := `{"NodeId":{"UUID":"cortex","NodeType":"CORTEX","LayerIndex":0},"Sensors":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"VectorLength":2,"Outbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25}}]}],"Neurons":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Bias":-30,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"Weights":[20,20]}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}],"ActivationFunction":{"Name":"sigmoid"}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Bias":10,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"Weights":[-20,-20]}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}],"ActivationFunction":{"Name":"sigmoid"}},{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Bias":-10,"Inbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Weights":[20]},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Weights":[20]}],"Outbound":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5}}],"ActivationFunction":{"Name":"sigmoid"}}],"Actuators":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5},"VectorLength":1,"Inbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Weights":null}]}]}`

	jsonBytes := []byte(jsonString)

	cortex := &Cortex{}
	err := json.Unmarshal(jsonBytes, cortex)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)
	log.Printf("cortex: %v", cortex)

}

func TestCortexInit(t *testing.T) {

	jsonString := `{"NodeId":{"UUID":"cortex","NodeType":"CORTEX","LayerIndex":0},"Sensors":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"VectorLength":2,"Outbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25}}]}],"Neurons":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Bias":-30,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"Weights":[20,20]}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}],"ActivationFunction":{"Name":"sigmoid"}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Bias":10,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"Weights":[-20,-20]}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}],"ActivationFunction":{"Name":"sigmoid"}},{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Bias":-10,"Inbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Weights":[20]},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Weights":[20]}],"Outbound":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5}}],"ActivationFunction":{"Name":"sigmoid"}}],"Actuators":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5},"VectorLength":1,"Inbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Weights":null}]}]}`

	jsonBytes := []byte(jsonString)

	cortex := &Cortex{}
	err := json.Unmarshal(jsonBytes, cortex)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)

	cortex.Init()

	neuron := cortex.Neurons[0]
	assert.True(t, neuron.DataChan != nil)

	cortex.Run()
	cortex.Shutdown()

}

func TestSyncActuators(t *testing.T) {

	actuatorNodeId := NewActuatorId("actuator", 0.5)
	actuator := &Actuator{
		NodeId:       actuatorNodeId,
		VectorLength: 1,
	}

	syncChan := make(chan *NodeId, 1)

	cortexNodeId := NewCortexId("cortex")
	cortex := &Cortex{
		NodeId:    cortexNodeId,
		Actuators: []*Actuator{actuator},
		SyncChan:  syncChan,
	}

	cortex.Init()

	syncChan <- actuatorNodeId

	cortex.SyncActuators()

}

func TestRecurrentCortex(t *testing.T) {

	jsonString := `{"NodeId":{"UUID":"cortex","NodeType":"CORTEX","LayerIndex":0},"Sensors":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"VectorLength":2,"Outbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25}}]}],"Neurons":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Bias":-30,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"Weights":[20,20]}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}],"ActivationFunction":{"Name":"sigmoid"}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Bias":10,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"Weights":[-20,-20]}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}],"ActivationFunction":{"Name":"sigmoid"}},{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Bias":-10,"Inbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Weights":[20]},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Weights":[20]},{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Weights":[0.0955837638877588]}],"Outbound":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5}},{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}],"ActivationFunction":{"Name":"sigmoid"}}],"Actuators":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5},"VectorLength":1,"Inbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Weights":null}]}]}`

	logg.LogKeys["NODE_SEND"] = true
	logg.LogKeys["NODE_RECV"] = true
	logg.LogKeys["MISC"] = true

	jsonBytes := []byte(jsonString)

	cortex := &Cortex{}
	err := json.Unmarshal(jsonBytes, cortex)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)

	examples := XnorTrainingSamples()
	logg.LogTo("MISC", "training samples: %v", examples)

	fitness := cortex.Fitness(examples)
	assert.True(t, fitness >= 0)

}

func TestCortexFitness(t *testing.T) {

	xnorCortex := XnorCortex()
	assert.True(t, xnorCortex != nil)

	// inputs + expected outputs
	examples := XnorTrainingSamples()
	log.Printf("training samples: %v", examples)

	// get the fitness
	fitness := xnorCortex.Fitness(examples)
	log.Printf("cortex fitness: %v", fitness)

	assert.True(t, fitness >= 1e8)

}

func TestNeuronLayerMap(t *testing.T) {
	xnorCortex := XnorCortex()
	layerToNeuronMap := xnorCortex.NeuronLayerMap()
	log.Printf("layerToNeuron: %v", layerToNeuronMap)
	hiddenNeurons := layerToNeuronMap[0.25]
	assert.Equals(t, len(hiddenNeurons), 2)
	outputNeurons := layerToNeuronMap[0.35]
	assert.Equals(t, len(outputNeurons), 1)
}
