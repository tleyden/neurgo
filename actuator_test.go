package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"testing"
)

func TestActuatorRun(t *testing.T) {

	fakeNodeId := NewNeuronId("fake-node", 0.25)
	actuatorNodeId := NewActuatorId("actuator", 0.5)

	collectedActuatorVals := make([][]float64, 1)
	collectedActuatorIndex := 0
	actuatorFunc := func(outputs []float64) {
		collectedActuatorVals[collectedActuatorIndex] = outputs
		collectedActuatorIndex += 1
	}

	actuator := &Actuator{
		NodeId:           actuatorNodeId,
		VectorLength:     2,
		ActuatorFunction: actuatorFunc,
	}
	actuator.Init()
	go actuator.Run()

	// send it a message
	fakeInput := []float64{1}
	dataMessage := &DataMessage{
		SenderId: fakeNodeId,
		Inputs:   fakeInput,
	}

	actuator.DataChan <- dataMessage

	// make sure our actuator function was called
	assert.Equals(t, collectedActuatorIndex, 1)
	collectedActuatorVal := collectedActuatorVals[0]
	assert.True(t, vectorEqualsWithMaxDelta(collectedActuatorVal, fakeInput, 0.1))

}
