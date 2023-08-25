package operation

import (
	"cmp"
)

func Max[T cmp.Ordered](a, b T) T {
	if a < b {
		return b
	}

	return a
}
