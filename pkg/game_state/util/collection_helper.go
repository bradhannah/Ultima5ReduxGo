package util

func FindIndexFromSliceT[T comparable](arr []T, target T) int {
	for i, v := range arr {
		if v == target {
			return i
		}
	}
	return -1
}

func FindKeyByValueT[K comparable, V comparable](data map[K]V, value V) (K, bool) {
	for key, v := range data {
		if v == value {
			return key, true
		}
	}
	var empty K
	return empty, false // Value not found
}
