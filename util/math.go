package util

import (
	"math"
)

func AdjustPoolSize(poolSize, amount, minChunkSize int) (int, int) {
	var chunkSize = ceil(amount, poolSize)
	if chunkSize < minChunkSize {
		chunkSize = minChunkSize

		poolSize = ceil(amount, chunkSize)
	}

	return poolSize, chunkSize
}

func ceil(numerator, denominator int) int {
	return int(math.Ceil(float64(numerator) / float64(denominator)))
}
