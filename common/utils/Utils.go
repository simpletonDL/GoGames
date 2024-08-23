package utils

import "math/rand"

// RandInRange generates random integer from [min, max)
func RandInRange(min, max int) int {
	return min + rand.Intn(max-min)
}

func Filter[E any](s []E, f func(E) bool) []E {
	s2 := make([]E, 0, len(s))
	for _, e := range s {
		if f(e) {
			s2 = append(s2, e)
		}
	}
	return s2
}
