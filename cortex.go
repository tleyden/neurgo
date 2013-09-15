package neurgo

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/logg"
	"log"
	"time"
)

const FITNESS_THRESHOLD = 1e8

type Cortex struct {
	NodeId    *NodeId
	Sensors   []*Sensor
	Neurons   []*Neuron
	Actuators []*Actuator
	SyncChan  chan *NodeId
}

type ActuatorBarrier map[*NodeId]bool // TODO: fixme!! totally broken
type UUIDToNeuronMap map[string]*Neuron

func (cortex *Cortex) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			NodeId    *NodeId
			Sensors   []*Sensor
			Neurons   []*Neuron
			Actuators []*Actuator
		}{
			NodeId:    cortex.NodeId,
			Sensors:   cortex.Sensors,
			Neurons:   cortex.Neurons,
			Actuators: cortex.Actuators,
		})
}

func (cortex *Cortex) MarshalJSONToFile(filename string) error {
	json, err := json.Marshal(cortex)
	if err != nil {
		return err
	}
	jsonString := fmt.Sprintf("%s", json)
	logg.Log("%v", jsonString)
	WriteStringToFile(jsonString, filename)
	return nil
}

func (cortex *Cortex) String() string {
	return JsonString(cortex)
}

func (cortex *Cortex) Copy() *Cortex {

	// serialize to json
	jsonBytes, err := json.Marshal(cortex)
	if err != nil {
		log.Fatal(err)
	}

	// new cortex
	cortexCopy := &Cortex{}

	// deserialize json into new cortex
	err = json.Unmarshal(jsonBytes, cortexCopy)
	if err != nil {
		log.Fatal(err)
	}

	cortexCopy.Init(false)
	cortexCopy.initOutboundConnections()

	return cortexCopy

}

func (cortex *Cortex) Run() {

	cortex.checkRunnable()

	// TODO: merge slices, create Runnable() interface
	// and make into single loop

	for _, sensor := range cortex.Sensors {
		go sensor.Run()
	}
	for _, neuron := range cortex.Neurons {
		go neuron.Run()
	}
	for _, actuator := range cortex.Actuators {
		go actuator.Run()
	}
}

func (cortex *Cortex) Shutdown() {
	for _, sensor := range cortex.Sensors {
		sensor.Shutdown()
	}
	for _, neuron := range cortex.Neurons {
		neuron.Shutdown()
	}
	for _, actuator := range cortex.Actuators {
		actuator.Shutdown()
	}
	cortex.SyncChan = nil
}

// Initialize/re-initialize the cortex.
// reInit: basically this is a messy hack to solve the issue:
// - neuron.Init() function is called and DataChan buffer len = X
// - new recurrent connections are added
// - since the DataChan buffer len is X, and needs to be X+1, network is wedged
// So by doing a "destructive reInit" it will rebuild all DataChan's
// and all outbound connections which contain DataChan's, thus solving
// the problem.
// TODO: fix this hack
func (cortex *Cortex) Init(reInit bool) {

	if cortex.SyncChan == nil {
		cortex.SyncChan = make(chan *NodeId, 1)
	}

	for _, sensor := range cortex.Sensors {
		sensor.Init(false)
	}
	for _, neuron := range cortex.Neurons {
		neuron.Init(false)
	}
	for _, actuator := range cortex.Actuators {
		actuator.Init(false)
	}

	// cortex.initOutboundConnections()

}

func (cortex *Cortex) SetSensors(sensors []*Sensor) {
	cortex.Sensors = sensors
	for _, sensor := range cortex.Sensors {
		sensor.Cortex = cortex
	}
}

func (cortex *Cortex) SetNeurons(neurons []*Neuron) {
	cortex.Neurons = neurons
	for _, neuron := range cortex.Neurons {
		neuron.Cortex = cortex
	}
}

func (cortex *Cortex) SetActuators(actuators []*Actuator) {
	cortex.Actuators = actuators
	for _, actuator := range cortex.Actuators {
		actuator.Cortex = cortex
	}
}

func (cortex *Cortex) NeuronUUIDMap() UUIDToNeuronMap {
	neuronUUIDMap := make(UUIDToNeuronMap)
	for _, neuron := range cortex.Neurons {
		neuronUUIDMap[neuron.NodeId.UUID] = neuron
	}
	return neuronUUIDMap
}

func (cortex *Cortex) CreateNeuronInLayer(layerIndex float64) *Neuron {
	uuid := NewUuid()
	neuron := &Neuron{
		ActivationFunction: RandomEncodableActivation(),
		NodeId:             NewNeuronId(uuid, layerIndex),
		Bias:               RandomBias(),
	}
	neuron.Cortex = cortex

	reInit := false
	neuron.Init(reInit)

	cortex.Neurons = append(cortex.Neurons, neuron)

	return neuron
}

func (cortex *Cortex) SensorNodeIds() []*NodeId {
	nodeIds := make([]*NodeId, 0)
	for _, sensor := range cortex.Sensors {
		nodeIds = append(nodeIds, sensor.NodeId)
	}
	return nodeIds
}

func (cortex *Cortex) NeuronNodeIds() []*NodeId {
	nodeIds := make([]*NodeId, 0)
	for _, neuron := range cortex.Neurons {
		nodeIds = append(nodeIds, neuron.NodeId)
	}
	return nodeIds
}

func (cortex *Cortex) ActuatorNodeIds() []*NodeId {
	nodeIds := make([]*NodeId, 0)
	for _, actuator := range cortex.Actuators {
		nodeIds = append(nodeIds, actuator.NodeId)
	}
	return nodeIds

}

func (cortex *Cortex) AllNodeIds() []*NodeId {
	neuronNodeIds := cortex.NeuronNodeIds()
	sensorNodeIds := cortex.SensorNodeIds()
	actuatorNodeIds := cortex.ActuatorNodeIds()
	availableNodeIds := append(neuronNodeIds, sensorNodeIds...)
	availableNodeIds = append(availableNodeIds, actuatorNodeIds...)
	return availableNodeIds
}

func (cortex *Cortex) NeuronLayerMap() LayerToNeuronMap {
	layerToNeuronMap := make(LayerToNeuronMap)
	for _, neuron := range cortex.Neurons {
		if _, ok := layerToNeuronMap[neuron.NodeId.LayerIndex]; !ok {
			neurons := make([]*Neuron, 0)
			neurons = append(neurons, neuron)
			layerToNeuronMap[neuron.NodeId.LayerIndex] = neurons
		} else {
			neurons := layerToNeuronMap[neuron.NodeId.LayerIndex]
			neurons = append(neurons, neuron)
			layerToNeuronMap[neuron.NodeId.LayerIndex] = neurons
		}

	}
	return layerToNeuronMap
}

func (cortex *Cortex) NodeIdLayerMap() LayerToNodeIdMap {
	layerToNodeIdMap := make(LayerToNodeIdMap)
	for _, nodeId := range cortex.AllNodeIds() {
		if _, ok := layerToNodeIdMap[nodeId.LayerIndex]; !ok {
			nodeIds := make([]*NodeId, 0)
			nodeIds = append(nodeIds, nodeId)
			layerToNodeIdMap[nodeId.LayerIndex] = nodeIds
		} else {
			nodeIds := layerToNodeIdMap[nodeId.LayerIndex]
			nodeIds = append(nodeIds, nodeId)
			layerToNodeIdMap[nodeId.LayerIndex] = nodeIds
		}

	}
	return layerToNodeIdMap
}

// We may be in a state where the outbound connections
// do not have data channels associated with them, even
// though the data channels exist.  (eg, when deserializing
// from json).  Fix this by seeking out those outbound
// connections and setting the data channels.
func (cortex *Cortex) initOutboundConnections() {

	// build a nodeId -> dataChan map
	nodeIdToDataMsg := cortex.nodeIdToDataMsg()

	// walk all sensors and neurons and fix up their outbound connections
	for _, sensor := range cortex.Sensors {
		sensor.initOutboundConnections(nodeIdToDataMsg)
	}
	for _, neuron := range cortex.Neurons {
		neuron.initOutboundConnections(nodeIdToDataMsg)
	}

}

func (cortex *Cortex) shutdownOutboundConnections() {

	// walk all sensors and neurons and shutdown their outbound connections
	for _, sensor := range cortex.Sensors {
		sensor.shutdownOutboundConnections()
	}
	for _, neuron := range cortex.Neurons {
		neuron.shutdownOutboundConnections()
	}

}

func (cortex *Cortex) nodeIdToDataMsg() nodeIdToDataMsgMap {
	nodeIdToDataMsg := make(nodeIdToDataMsgMap)
	for _, neuron := range cortex.Neurons {
		nodeIdToDataMsg[neuron.NodeId.UUID] = neuron.DataChan
	}
	for _, actuator := range cortex.Actuators {
		nodeIdToDataMsg[actuator.NodeId.UUID] = actuator.DataChan
	}
	return nodeIdToDataMsg

}

func (cortex *Cortex) checkRunnable() {
	if cortex.SyncChan == nil {
		log.Panicf("cortex.SyncChan is nil")
	}
}

func (cortex *Cortex) Verify(samples []*TrainingSample) bool {
	fitness := cortex.Fitness(samples)
	return fitness >= FITNESS_THRESHOLD
}

func (cortex *Cortex) Fitness(samples []*TrainingSample) float64 {

	shouldReInit := true
	cortex.Init(shouldReInit)

	errorAccumulated := float64(0)

	// assumes there is only one sensor and one actuator
	// (to support more, this method will require more coding)
	if len(cortex.Sensors) != 1 {
		log.Panicf("Must have exactly one sensor")
	}
	if len(cortex.Actuators) != 1 {
		log.Panicf("Must have exactly one actuator")
	}

	// install function to sensor which will stream training samples
	sensor := cortex.Sensors[0]
	sensorFunc := func(syncCounter int) []float64 {
		sampleX := samples[syncCounter]
		return sampleX.SampleInputs[0]
	}
	sensor.SensorFunction = sensorFunc

	// install function to actuator which will collect outputs
	actuator := cortex.Actuators[0]
	numTimesFuncCalled := 0
	actuatorFunc := func(outputs []float64) {
		logg.LogTo("MISC", "numTimesFuncCalled: %v", numTimesFuncCalled)
		logg.LogTo("MISC", "len(samples): %v", len(samples))
		logg.LogTo("MISC", "samples: %v", samples)
		expected := samples[numTimesFuncCalled].ExpectedOutputs[0]
		error := SumOfSquaresError(expected, outputs)
		errorAccumulated += error
		numTimesFuncCalled += 1
		cortex.SyncChan <- actuator.NodeId
	}
	actuator.ActuatorFunction = actuatorFunc

	go cortex.Run()

	for _ = range samples {
		cortex.SyncSensors()
		cortex.SyncActuators()
	}

	cortex.Shutdown()

	// calculate fitness
	fitness := float64(1) / errorAccumulated

	return fitness

}

func (cortex *Cortex) FindSensor(nodeId *NodeId) *Sensor {
	for _, sensor := range cortex.Sensors {
		if sensor.NodeId.UUID == nodeId.UUID {
			return sensor
		}
	}
	return nil
}

func (cortex *Cortex) FindNeuron(nodeId *NodeId) *Neuron {
	for _, neuron := range cortex.Neurons {
		if neuron.NodeId.UUID == nodeId.UUID {
			return neuron
		}
	}
	return nil
}

func (cortex *Cortex) FindActuator(nodeId *NodeId) *Actuator {
	for _, actuator := range cortex.Actuators {
		if actuator.NodeId.UUID == nodeId.UUID {
			return actuator
		}
	}
	return nil
}

// TODO: rename to FindOutboundConnector
func (cortex *Cortex) FindConnector(nodeId *NodeId) OutboundConnector {
	for _, sensor := range cortex.Sensors {
		if sensor.NodeId.UUID == nodeId.UUID {
			return sensor
		}
	}
	for _, neuron := range cortex.Neurons {
		if neuron.NodeId.UUID == nodeId.UUID {
			return neuron
		}
	}
	return nil
}

func (cortex *Cortex) FindInboundConnector(nodeId *NodeId) InboundConnector {
	for _, neuron := range cortex.Neurons {
		if neuron.NodeId.UUID == nodeId.UUID {
			return neuron
		}
	}
	for _, actuator := range cortex.Actuators {
		if actuator.NodeId.UUID == nodeId.UUID {
			return actuator
		}
	}

	return nil
}

func (cortex *Cortex) SyncSensors() {
	for _, sensor := range cortex.Sensors {
		select {
		case sensor.SyncChan <- true:
		case <-time.After(time.Second):
			log.Panicf("Cortex unable to send Sync message to sensor %v", sensor)
		}
	}

}

func (cortex *Cortex) SyncActuators() {
	actuatorBarrier := cortex.createActuatorBarrier()
	for {

		select {
		case senderNodeId := <-cortex.SyncChan:
			actuatorBarrier[senderNodeId] = true
		case <-time.After(time.Second):
			log.Panicf("Timeout waiting for actuator sync message")
		}

		if cortex.isBarrierSatisfied(actuatorBarrier) {
			break
		}

	}
}

func (cortex *Cortex) Validate() bool {

	for _, neuron := range cortex.Neurons {
		if neuron.Cortex == nil {
			logg.LogWarn("Neuron: %v has no cortex", neuron)
			return false
		}
	}
	return true
}

func (cortex *Cortex) Repair() {

	for _, neuron := range cortex.Neurons {
		if neuron.Cortex == nil {
			neuron.Cortex = cortex
		}
	}

}

func (cortex *Cortex) createActuatorBarrier() ActuatorBarrier {
	actuatorBarrier := make(ActuatorBarrier)
	for _, actuator := range cortex.Actuators {
		actuatorBarrier[actuator.NodeId] = false
	}
	return actuatorBarrier
}

func (cortex *Cortex) isBarrierSatisfied(barrier ActuatorBarrier) bool {
	for _, value := range barrier {
		if value == false {
			return false
		}
	}
	return true
}
