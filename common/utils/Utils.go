package utils

import (
	"math/rand"
	"strings"
)

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

func Map[K comparable, V any, R any](m map[K]V, f func(K, V) R) []R {
	result := []R{}
	for k, v := range m {
		result = append(result, f(k, v))
	}
	return result
}

func First[T1, T2 any](x T1, y T2) T1 {
	return x
}

func Second[T1, T2 any](x T1, y T2) T2 {
	return y
}

func Keys[K comparable, V any](m map[K]V) []K {
	return Map(m, First)
}

func Values[K comparable, V any](m map[K]V) []V {
	return Map(m, Second)
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

func AdjustString(s string, size int) string {
	s = s[:min(len(s), size)]
	return s + strings.Repeat(" ", size-len(s))
}
