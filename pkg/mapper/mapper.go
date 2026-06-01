package mapper

func Map[T any, U any](value []T, fun func(T) U) []U {
	output := make([]U, 0, len(value))

	for _, t := range value {
		o := fun(t)
		output = append(output, o)
	}

	return output
}
