package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
)

func TestPruneEmptyElements(t *testing.T) {

	weightedInputs := make([]*weightedInput, 3)
	weightedInputs[0] = &weightedInput{}
	weightedInputs[2] = &weightedInput{}
	
	pruned := pruneEmptyElements(weightedInputs)
	assert.Equals(t, len(pruned), 2)
	
}
