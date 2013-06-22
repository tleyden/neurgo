package neurgo

import (
	"fmt"
	"math"
)

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Pow(math.E, -1*x))
}

// http://en.wikipedia.org/wiki/Residual_sum_of_squares
func SumOfSquaresError(expected []float64, actual []float64) float64 {

	result := float64(0)
	if len(expected) != len(actual) {
		msg := fmt.Sprintf("vector lengths dont match (%d != %d)", len(expected), len(actual))
		panic(msg)
	}

	for i, expectedVal := range expected {
		actualVal := actual[i]
		delta := actualVal - expectedVal
		deltaSquared := math.Pow(delta, 2)
		result += deltaSquared
	}

	return result
}

func equalsWithMaxDelta(x, y, maxDelta float64) bool {
	delta := math.Abs(x - y)
	return delta <= maxDelta
}

func vectorEqualsWithMaxDelta(xValues, yValues []float64, maxDelta float64) bool {
	equals := true
	for i, x := range xValues {
		y := yValues[i]
		if !equalsWithMaxDelta(x, y, maxDelta) {
			equals = false
		}
	}
	return equals
}
