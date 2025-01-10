package utils

func Keys[K comparable, V any](dict map[K]V) []K {
	keys := make([]K, len(dict))
	index := 0
	for key := range dict {
		keys[index] = key
		index++
	}
	return keys
}
