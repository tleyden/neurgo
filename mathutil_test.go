package neurgo

import (
	"testing"
	"github.com/couchbaselabs/go.assert"
)

func TestSigmoid(t *testing.T) {

	assert.Equals(t, sigmoid(0), 0.5)
	assert.True(t, sigmoid(10) > 0.9)
	assert.True(t, sigmoid(-10) < 0.1)

}
