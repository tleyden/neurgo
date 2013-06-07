package neurgo

import (
	"math"
)

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Pow(math.E,-1*x))
}

