package server

import (
	"math/rand"
)

func Random(low, high float64) float64 {
	return low + rand.Float64()*(high-low)
}

func RandomVector3(low, high float64) Vector3 {
	return Vector3{
		Random(low, high),
		Random(low, high),
		Random(low, high),
	}
}
