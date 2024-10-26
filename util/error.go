package util

func UnwrapError[T any](t T, _ error) T {
	return t
}
