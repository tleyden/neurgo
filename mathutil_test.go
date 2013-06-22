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

func TestSumOfSquaresError(t *testing.T) {
	expected := []float64{.5}
	actual := []float64{1}
	error := SumOfSquaresError(expected, actual)
	nearlyEqualsPoint25 := equalsWithMaxDelta(error, 0.25, .01)
	assert.True(t, nearlyEqualsPoint25)
}

func TestSafeScalarInverse(t *testing.T) {
	value := SafeScalarInverse(0)
	assert.True(t, value > 1000000)
	assert.True(t, equalsWithMaxDelta(SafeScalarInverse(1), 1.0, .0001))
}
