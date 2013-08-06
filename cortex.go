package neurgo

import (
	"log"
)

type Cortex struct {
	Sensors   []*Sensor
	Neurons   []*Neuron
	Actuators []*Actuator
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

	sensor := cortex.Sensors[0]

	// install function to sensor which will stream training samples
	sensorFunc := func(syncCounter int) []float64 {
		sampleX := samples[syncCounter]
		return sampleX.SampleInputs[0]
	}

	// install function to actuator which will collect outputs
	sensor.SensorFunction = sensorFunc

	/*for _, sample := range samples {
		cortex.SyncSensors()
		cortex.SyncActuators()
	}*/

	// make sure collected outputs match expected outputs

	return 0
}
