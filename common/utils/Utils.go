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

func All[E any](xs []E, f func(E) bool) bool {
	for _, x := range xs {
		if !f(x) {
			return false
		}
	}
	return true
}

func AllEntries[K comparable, V any](xs map[K]V, f func(key K, value V) bool) bool {
	for k, v := range xs {
		if !f(k, v) {
			return false
		}
	}
	return true
}
