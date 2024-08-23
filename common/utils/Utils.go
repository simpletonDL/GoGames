package utils

import "math/rand"

// RandInRange generates random integer from [min, max)
func RandInRange(min, max int) int {
	return min + rand.Intn(max-min)
}
