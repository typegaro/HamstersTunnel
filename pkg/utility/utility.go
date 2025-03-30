package utility

func MapGetValues[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return []V{}
	}
	values := make([]V, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

func Ternary[T any](condition bool, valTrue, valFalse T) T {
	if condition {
		return valTrue
	}
	return valFalse
}
