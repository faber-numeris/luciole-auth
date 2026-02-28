package tools

// Must is a helper function that panics if err is not nil.
// It returns the value if err is nil.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
