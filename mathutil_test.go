package neurgo

import (
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	"math"
	"testing"
)

func TestSigmoid(t *testing.T) {
	assert.Equals(t, Sigmoid(0), 0.5)
	assert.True(t, Sigmoid(10) > 0.9)
	assert.True(t, Sigmoid(-10) < 0.1)
}

func TestEqualsWithMaxDelta(t *testing.T) {
	assert.True(t, EqualsWithMaxDelta(0.99999, 1.00000, .01))
	assert.False(t, EqualsWithMaxDelta(0.95, 1.00000, .01))
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
	nearlyEqualsPoint25 := EqualsWithMaxDelta(error, 0.25, .01)
	assert.True(t, nearlyEqualsPoint25)
}

func TestSafeScalarInverse(t *testing.T) {
	value := SafeScalarInverse(0)
	assert.True(t, value > 1000000)
	assert.True(t, EqualsWithMaxDelta(SafeScalarInverse(1), 1.0, .0001))
}

func TestRandomInRange(t *testing.T) {
	SeedRandom()
	assert.True(t, RandomInRange(-1*math.Pi, math.Pi) <= math.Pi)
	assert.True(t, RandomInRange(-1*math.Pi, math.Pi) >= -1*math.Pi)

	foundLessThanZero := false
	foundGreaterThanZero := false
	for i := 0; i < 1000; i++ {
		randFloat := RandomInRange(-1*math.Pi, math.Pi)
		if randFloat > 0 {
			foundGreaterThanZero = true
		} else if randFloat < 0 {
			foundLessThanZero = true
		}
	}
	assert.True(t, foundLessThanZero)
	assert.True(t, foundGreaterThanZero)

}

func TestRandomIntInRange(t *testing.T) {
	SeedRandom()
	result := RandomIntInRange(1, 4)
	assert.True(t, result >= 1)
	assert.True(t, result <= 4)

	found0 := false
	found1 := false
	found2 := false
	for i := 0; i < 1000; i++ {
		randInt := RandomIntInRange(0, 2)
		if randInt == 0 {
			found0 = true
		}
		if randInt == 1 {
			found1 = true
		}
		if randInt == 2 {
			found2 = true
		}

	}
	assert.True(t, found0)
	assert.True(t, found1)
	assert.False(t, found2)

}

func TestIntModuleProper(t *testing.T) {
	assert.False(t, IntModuloProper(0, 100))
	assert.True(t, IntModuloProper(500, 100))
}

func TestNormalizeInRange(t *testing.T) {

	logg.LogKeys["TEST"] = true
	params := NormalizeParams{
		SourceRangeStart: 0.0,
		SourceRangeEnd:   31.0,
		TargetRangeStart: -1.0,
		TargetRangeEnd:   1.0,
	}

	testPairs := [][]float64{
		[]float64{0.0, -1.0},
		[]float64{31.0, 1.0},
		[]float64{15.5, 0.0},
	}

	for _, testPair := range testPairs {
		origValue := testPair[0]
		expectedValue := testPair[1]
		actualValue := NormalizeInRange(params, origValue)
		delta := 0.01
		logg.LogTo("TEST", "actual: %v, expected: %v", actualValue, expectedValue)
		assert.True(t, EqualsWithMaxDelta(actualValue, expectedValue, delta))
	}

}
