package util

import (
	"math"
)

func AdjustSizes(numerator, denominator, minValue int) (int, int) {
	var value = ceil(numerator, denominator)
	if value < minValue {
		value = minValue

		denominator = ceil(numerator, value)
	}

	return denominator, value
}

func ceil(numerator, denominator int) int {
	return int(math.Ceil(float64(numerator) / float64(denominator)))
}
