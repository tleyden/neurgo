package neurgo

import (
	"encoding/json"
	"log"
	"time"
)

type Cortex struct {
	NodeId    *NodeId
	Sensors   []*Sensor
	Neurons   []*Neuron
	Actuators []*Actuator
	SyncChan  chan *NodeId
}

type ActuatorBarrier map[*NodeId]bool

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

func (cortex *Cortex) Init() {
	if cortex.SyncChan == nil {
		cortex.SyncChan = make(chan *NodeId, 1)
	}
}

func (cortex *Cortex) checkRunnable() {
	if cortex.SyncChan == nil {
		log.Panicf("cortex.SyncChan is nil")
	}
}

func (cortex *Cortex) Fitness(samples []*TrainingSample) float64 {

	errorAccumulated := float64(0)
	log.Printf("error: %v", errorAccumulated)

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

	cortex.Init()
	go cortex.Run()

	for _ = range samples {
		cortex.SyncSensors()
		cortex.SyncActuators()
	}

	cortex.Shutdown()

	// calculate fitness
	log.Printf("errorAccumulated: %v", errorAccumulated)

	return float64(1) / errorAccumulated

}

func (cortex *Cortex) SyncSensors() {
	for _, sensor := range cortex.Sensors {
		select {
		case sensor.SyncChan <- true:
			log.Printf("Sync -> %v", sensor)
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
