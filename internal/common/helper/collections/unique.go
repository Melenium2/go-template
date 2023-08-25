package collections

import (
	"cmp"
)

func Unique[T cmp.Ordered](input []T) []T {
	u := make([]T, 0)
	m := make(map[T]struct{})

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = struct{}{}

			u = append(u, val)
		}
	}

	return u
}
