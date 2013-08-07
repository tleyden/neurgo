package neurgo

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/go.assert"
	"log"
	"testing"
)

func TestJsonMarshal(t *testing.T) {
	fakeNodeId := NewNeuronId("fake-node", 0.25)
	json, err := json.Marshal(fakeNodeId)
	if err != nil {
		log.Fatal(err)
	}

	assert.True(t, err == nil)
	jsonString := fmt.Sprintf("%s", json)
	log.Printf("jsonString: %v", jsonString)
}
