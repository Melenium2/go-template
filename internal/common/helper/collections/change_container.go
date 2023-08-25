package collections

import (
	"cmp"
)

type ExtendedConstraint interface {
	cmp.Ordered | ~[16]byte // uuid.UUID ([16]byte) type
}

func ChangeContainer[K ExtendedConstraint, V any](m map[K]V) []V {
	list := make([]V, 0, len(m))

	for _, v := range m {
		list = append(list, v)
	}

	return list
}
