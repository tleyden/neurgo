package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"testing"
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

func TestVectorEqualsWithMaxDelta(t *testing.T) {

	xValues := []float64{0.99999, 0.00000}
	yValues := []float64{1.00000, 0.00001}

	assert.True(t, vectorEqualsWithMaxDelta(xValues, yValues, .01))

	xValues = []float64{0.95, 0.00000}
	yValues = []float64{1.00000, 1.00000}

	assert.False(t, vectorEqualsWithMaxDelta(xValues, yValues, .01))
}
