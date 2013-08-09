package neurgo

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
)

func TestActivationFunctionMarshal(t *testing.T) {

	expectedJsonString := `{"Name":"sigmoid"}`

	encodableActivation := &EncodableActivation{
		Name:               "sigmoid",
		ActivationFunction: Sigmoid,
	}

	json, err := json.Marshal(encodableActivation)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)
	jsonString := fmt.Sprintf("%s", json)
	log.Printf("jsonString: %v", jsonString)
	assert.Equals(t, jsonString, expectedJsonString)

}

func TestActivationFunctionUnmarshal(t *testing.T) {

	jsonString := `{"Name":"sigmoid"}`
	jsonBytes := []byte(jsonString)

	encodableActivation := &EncodableActivation{}
	err := json.Unmarshal(jsonBytes, encodableActivation)
	if err != nil {
		log.Fatal(err)
	}
	assert.True(t, err == nil)
	assert.Equals(t, encodableActivation.Name, "sigmoid")
	assert.True(t, encodableActivation.ActivationFunction != nil)

}
