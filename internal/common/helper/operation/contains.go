package operation

func Contains[T comparable](target T, list []T) bool {
	for _, item := range list {
		if target == item {
			return true
		}
	}

	return false
}
