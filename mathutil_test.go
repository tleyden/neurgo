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

func TestEqualsWithMaxDelta(t *testing.T) {
	assert.True(t, equalsWithMaxDelta(0.99999, 1.00000, .01))
	assert.False(t, equalsWithMaxDelta(0.95, 1.00000, .01))
}
