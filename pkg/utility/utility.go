package utitlity

func MapGetValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	i := 0
	for _, value := range m {
		values[i] = value
		i++
	}
	return values
}

func Ternary[T any](condition bool, valTrue, valFalse T) T {
	if condition {
		return valTrue
	}
	return valFalse
}
