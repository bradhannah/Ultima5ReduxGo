package helpers

import "golang.org/x/exp/constraints"

// Generic Max function
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Generic Min function
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
