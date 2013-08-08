package neurgo

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
)

func TestCortexCopy(t *testing.T) {

	xnorCortex := XnorCortex(t)
	xnorCortexCopy := xnorCortex.Copy()

	xnorCortexCopy.Init()
	sensor := xnorCortexCopy.Sensors[0]
	assert.True(t, sensor.Outbound[0].DataChan != nil)

	// inputs + expected outputs
	examples := xnorTrainingSamples()

	// get the fitness
	fitness := xnorCortexCopy.Fitness(examples)

	assert.True(t, fitness >= 1e8)

}

func TestCortexJsonMarshal(t *testing.T) {
	xnorCortex := XnorCortex(t)
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

	xnorCortex := XnorCortex(t)
	assert.True(t, xnorCortex != nil)

	// inputs + expected outputs
	examples := xnorTrainingSamples()
	log.Printf("training samples: %v", examples)

	// get the fitness
	fitness := xnorCortex.Fitness(examples)
	log.Printf("cortex fitness: %v", fitness)

	assert.True(t, fitness >= 1e8)

}

func XnorCortex(t *testing.T) *Cortex {

	// create network nodes

	sensorNodeId := NewSensorId("sensor", 0.0)
	hiddenNeuron1NodeId := NewNeuronId("hidden-neuron1", 0.25)
	hiddenNeuron2NodeId := NewNeuronId("hidden-neuron2", 0.25)
	outputNeuronNodeIde := NewNeuronId("output-neuron", 0.35)

	actuatorNodeId := NewActuatorId("actuator", 0.5)

	hiddenNeuron1 := &Neuron{
		ActivationFunction: Sigmoid,
		NodeId:             hiddenNeuron1NodeId,
		Bias:               -30,
	}
	hiddenNeuron1.Init()

	hiddenNeuron2 := &Neuron{
		ActivationFunction: Sigmoid,
		NodeId:             hiddenNeuron2NodeId,
		Bias:               10,
	}
	hiddenNeuron2.Init()

	outputNeuron := &Neuron{
		ActivationFunction: Sigmoid,
		NodeId:             outputNeuronNodeIde,
		Bias:               -10,
	}
	outputNeuron.Init()

	sensor := &Sensor{
		NodeId:       sensorNodeId,
		VectorLength: 2,
	}
	sensor.Init()

	actuator := &Actuator{
		NodeId:       actuatorNodeId,
		VectorLength: 1,
	}
	actuator.Init()

	sensor.ConnectOutbound(hiddenNeuron1)
	hiddenNeuron1.ConnectInboundWeighted(sensor, []float64{20, 20})

	sensor.ConnectOutbound(hiddenNeuron2)
	hiddenNeuron2.ConnectInboundWeighted(sensor, []float64{-20, -20})

	assert.Equals(t, len(sensor.Outbound), 2)
	assert.Equals(t, len(hiddenNeuron1.Inbound), 1)
	assert.Equals(t, len(hiddenNeuron2.Inbound), 1)

	hiddenNeuron1.ConnectOutbound(outputNeuron)
	outputNeuron.ConnectInboundWeighted(hiddenNeuron1, []float64{20})

	hiddenNeuron2.ConnectOutbound(outputNeuron)
	outputNeuron.ConnectInboundWeighted(hiddenNeuron2, []float64{20})

	assert.Equals(t, len(hiddenNeuron1.Outbound), 1)
	assert.Equals(t, len(hiddenNeuron2.Outbound), 1)
	assert.Equals(t, len(outputNeuron.Inbound), 2)

	outputNeuron.ConnectOutbound(actuator)
	actuator.ConnectInbound(outputNeuron)
	assert.Equals(t, len(outputNeuron.Outbound), 1)
	assert.Equals(t, len(actuator.Inbound), 1)

	nodeId := NewCortexId("cortex")

	cortex := &Cortex{
		NodeId:    nodeId,
		Sensors:   []*Sensor{sensor},
		Neurons:   []*Neuron{hiddenNeuron1, hiddenNeuron2, outputNeuron},
		Actuators: []*Actuator{actuator},
	}

	return cortex

}

func xnorTrainingSamples() []*TrainingSample {

	// inputs + expected outputs
	examples := []*TrainingSample{

		// TODO: how to wrap this?
		{SampleInputs: [][]float64{[]float64{0, 1}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{1, 1}}, ExpectedOutputs: [][]float64{[]float64{1}}},
		{SampleInputs: [][]float64{[]float64{1, 0}}, ExpectedOutputs: [][]float64{[]float64{0}}},
		{SampleInputs: [][]float64{[]float64{0, 0}}, ExpectedOutputs: [][]float64{[]float64{1}}}}

	return examples

}
