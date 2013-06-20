package neurgo

import (
	"math"
)

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Pow(math.E, -1*x))
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
