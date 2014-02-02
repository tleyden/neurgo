package neurgo

import (
	"encoding/json"
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	logg.LogKeys["NODE_PRE_SEND"] = true
	logg.LogKeys["NODE_POST_SEND"] = true
	logg.LogKeys["NODE_POST_RECV"] = true
	logg.LogKeys["MISC"] = true
	logg.LogKeys["MAIN"] = true
	logg.LogKeys["DEBUG"] = true
	logg.LogKeys["SENSOR_SYNC"] = true
	logg.LogKeys["ACTUATOR_SYNC"] = true

}

func TestCortexCopy(t *testing.T) {

	// create original cortex
	xnorCortex := XnorCortex()

	// add a sensor function to the sensor
	sensorFunc := func(syncCounter int) []float64 {
		return []float64{1.0}
	}
	xnorCortex.Sensors[0].SensorFunction = sensorFunc

	// get the fitness
	examples := XnorTrainingSamples()
	fitness := xnorCortex.Fitness(examples)
	assert.True(t, fitness >= FITNESS_THRESHOLD)
	logg.LogTo("DEBUG", "Original cortex has fitness: %v", fitness)

	// copy the cortex
	xnorCortexCopy := xnorCortex.Copy()

	// validate the copy
	assert.True(t, xnorCortexCopy.Validate())

	// make sure actuator has reference to cortex in both orig and copy
	assert.True(t, xnorCortex.Actuators[0].Cortex != nil)
	assert.True(t, xnorCortexCopy.Actuators[0].Cortex != nil)

	// make sure the sensor function got copied over
	assert.True(t, xnorCortexCopy.Sensors[0].SensorFunction != nil)

	// make sure each neuron is associated with a cortex
	for _, neuron := range xnorCortex.Neurons {
		assert.False(t, neuron.Cortex == nil)
	}

	// get the fitness of the copied cortex
	fitness = xnorCortexCopy.Fitness(examples)
	assert.True(t, fitness >= FITNESS_THRESHOLD)

}

func TestCortexJsonMarshal(t *testing.T) {
	xnorCortex := XnorCortex()
	xnorCortex.MarshalJSONToFile("/tmp/output.json")
	// TODO: add some assertions about this  (try reading file back in)
}

func exampleCortexJson() string {
	jsonString := `{"NodeId":{"UUID":"cortex","NodeType":"CORTEX","LayerIndex":0},"Sensors":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"VectorLength":2,"Outbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25}}]}],"Neurons":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Bias":-30,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"Weights":[20,20]}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}],"ActivationFunction":{"Name":"sigmoid"}},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Bias":10,"Inbound":[{"NodeId":{"UUID":"sensor","NodeType":"SENSOR","LayerIndex":0},"Weights":[-20,-20]}],"Outbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35}}],"ActivationFunction":{"Name":"sigmoid"}},{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Bias":-10,"Inbound":[{"NodeId":{"UUID":"hidden-neuron1","NodeType":"NEURON","LayerIndex":0.25},"Weights":[20]},{"NodeId":{"UUID":"hidden-neuron2","NodeType":"NEURON","LayerIndex":0.25},"Weights":[20]}],"Outbound":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5}}],"ActivationFunction":{"Name":"sigmoid"}}],"Actuators":[{"NodeId":{"UUID":"actuator","NodeType":"ACTUATOR","LayerIndex":0.5},"VectorLength":1,"Inbound":[{"NodeId":{"UUID":"output-neuron","NodeType":"NEURON","LayerIndex":0.35},"Weights":null}]}]}`
	return jsonString

}

func TestCortexJsonUnmarshal(t *testing.T) {

	jsonBytes := []byte(exampleCortexJson())

	cortex := &Cortex{}
	err := json.Unmarshal(jsonBytes, cortex)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)
	log.Printf("cortex: %v", cortex)

}

func TestCortexInit(t *testing.T) {

	jsonBytes := []byte(exampleCortexJson())

	cortex := &Cortex{}
	err := json.Unmarshal(jsonBytes, cortex)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)

	cortex.Init()
	cortex.LinkNodesToCortex()

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

	jsonBytes := []byte(jsonString)

	cortex := &Cortex{}
	err := json.Unmarshal(jsonBytes, cortex)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)

	examples := XnorTrainingSamples()

	cortex.LinkNodesToCortex()

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

func TestMarshalJsonToFile(t *testing.T) {

	filename := "xnor.json"
	xnorCortex := XnorCortex()
	xnorCortex.MarshalJSONToFile(filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	contentStr := string(content)
	assert.True(t, len(contentStr) > 0)

}

func TestNewCortexFromJSONString(t *testing.T) {
	cortex, err := NewCortexFromJSONString(exampleCortexJson())
	assert.True(t, err == nil)
	assert.True(t, cortex != nil)
}
