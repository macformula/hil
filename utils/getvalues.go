package utils

//type GetValuesFunc[K comparable, V any] func(map[K]V) []V

// Returns all the values of a map
func GetValues[K comparable, V any](myMap map[K]V) []V {
	values := make([]V, 0, len(myMap))
	for _, v := range myMap {
		values = append(values, v)
	}
	return values
}
