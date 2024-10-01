package utils

func SetDefault[T comparable](v *T, defaultValue T) {
	if v == nil {
		return
	}

	var zero T
	if *v == zero {
		*v = defaultValue
	}
}

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}

	return falseVal
}
