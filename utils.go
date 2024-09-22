package msg

func GetPtr[T any](x T) *T {
	return &x
}
