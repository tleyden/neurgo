package neurgo

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

func SafeScalarInverse(x float64) float64 {
	if x == 0 {
		x += 0.000000001
	}
	return 1.0 / x
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

func IntModuloProper(x, y int) bool {
	if x > 0 && math.Mod(float64(x), float64(y)) == 0 {
		return true
	}
	return false
}

func RandomInRange(min, max float64) float64 {

	return rand.Float64()*(max-min) + min
}

// return a random number between min and max - 1
// eg, if you call it with 0,1 it will always return 0
// if you call it between 0,2 it will return 0 or 1
func RandomIntInRange(min, max int) int {
	if min == max {
		log.Printf("min==max")
		return min
	}
	return rand.Intn(max-min) + min
}

func SeedRandom() {
	rand.Seed(time.Now().UTC().UnixNano())
}
