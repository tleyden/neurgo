package neurgo

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
)

func TestCortexCopy(t *testing.T) {

	xnorCortex := XnorCortex()
	xnorCortexCopy := xnorCortex.Copy()

	xnorCortexCopy.Init()
	sensor := xnorCortexCopy.Sensors[0]
	assert.True(t, sensor.Outbound[0].DataChan != nil)

	// inputs + expected outputs
	examples := XnorTrainingSamples()

	// get the fitness
	fitness := xnorCortexCopy.Fitness(examples)

	assert.True(t, fitness >= 1e8)

}

func TestCortexJsonMarshal(t *testing.T) {
	xnorCortex := XnorCortex()
	json, err := json.Marshal(xnorCortex)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)
	jsonString := fmt.Sprintf("%s", json)
	log.Printf("jsonString: %v", jsonString)
	writeStringToFile(jsonString, "/tmp/output.json")

}

func TestCortexJsonUnmarshal(t *testing.T) {

	jsonString := `{"NodeId":{"UUID":"cortex","NodeType":"CORTEX","LayerIndex":0},"Sensors":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"VectorLength":2,"Outbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25}}]}],"Neurons":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Bias":-30,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0}}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}]},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Bias":10,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0}}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}]},{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Bias":-10,"Inbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25}}],"Outbound":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5}}]}],"Actuators":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5},"VectorLength":1,"Inbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}]}]}`
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

	jsonString := `{"NodeId":{"UUID":"cortex","NodeType":"CORTEX","LayerIndex":0},"Sensors":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"VectorLength":2,"Outbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25}}]}],"Neurons":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Bias":-30,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0}}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}]},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Bias":10,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0}}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}]},{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Bias":-10,"Inbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25}}],"Outbound":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5}}]}],"Actuators":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5},"VectorLength":1,"Inbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}]}]}`
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
