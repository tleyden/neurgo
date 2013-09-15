package neurgo

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
)

func TestActuatorJsonMarshal(t *testing.T) {

	actuatorNodeId := NewActuatorId("actuator", 0.5)

	actuator := &Actuator{
		NodeId:       actuatorNodeId,
		VectorLength: 0,
	}

	json, err := json.Marshal(actuator)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)
	jsonString := fmt.Sprintf("%s", json)
	log.Printf("jsonString: %v", jsonString)
}

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
		VectorLength:     0,
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

	actuator.Shutdown()

	// make sure our actuator function was called
	assert.Equals(t, collectedActuatorIndex, 1)
	collectedActuatorVal := collectedActuatorVals[0]
	assert.True(t, vectorEqualsWithMaxDelta(collectedActuatorVal, fakeInput, 0.1))

}
