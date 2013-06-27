package neurgo

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/go.assert"
	"testing"
)

func TestPerturbParameters(t *testing.T) {

	neuralNet := xnorCondensedNetwork()

	nnJson, _ := json.Marshal(neuralNet)
	nnJsonString := fmt.Sprintf("%s", nnJson)

	shc := new(StochasticHillClimber)

	shc.perturbParameters(neuralNet)

	nnJsonAfter, _ := json.Marshal(neuralNet)
	nnJsonStringAfter := fmt.Sprintf("%s", nnJsonAfter)

	// the json should be different after we perturb it
	assert.NotEquals(t, nnJsonString, nnJsonStringAfter)

}
