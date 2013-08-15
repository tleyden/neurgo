package neurgo

import (
	"encoding/json"
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

type ActuatorBarrier map[*NodeId]bool
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
func (cortex *Cortex) Init(reInit bool) {

	if reInit == true {
		cortex.shutdownOutboundConnections()
	}

	if reInit == true {
		cortex.SyncChan = make(chan *NodeId, 1)
	} else if cortex.SyncChan == nil {
		cortex.SyncChan = make(chan *NodeId, 1)
	}

	for _, sensor := range cortex.Sensors {
		sensor.Init(reInit)
	}
	for _, neuron := range cortex.Neurons {
		neuron.Init(reInit)
	}
	for _, actuator := range cortex.Actuators {
		actuator.Init(reInit)
	}

	cortex.initOutboundConnections()

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
	log.Printf("Fitness() started")

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
	log.Printf("Fitness() finished: %v", fitness)

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

func (cortex *Cortex) SyncSensors() {
	for _, sensor := range cortex.Sensors {
		select {
		case sensor.SyncChan <- true:
			log.Printf("Cortex send Sync message to %v", sensor)
		case <-time.After(time.Second):
			log.Panicf("Unable to send Sync message to sensor %v", sensor)
		}
	}

}

func (cortex *Cortex) SyncActuators() {
	actuatorBarrier := cortex.createActuatorBarrier()
	for {

		select {
		case senderNodeId := <-cortex.SyncChan:
			log.Printf("Cortex received Sync from -> %v", senderNodeId)
			actuatorBarrier[senderNodeId] = true
		case <-time.After(time.Second):
			log.Panicf("Timeout waiting for actuator sync message")
		}

		if cortex.isBarrierSatisfied(actuatorBarrier) {
			break
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
