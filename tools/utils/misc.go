// License: GPLv3 Copyright: 2022, Kovid Goyal, <kovid at kovidgoyal.net>

package utils

import (
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

var _ = fmt.Print

func Reverse[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func Reversed[T any](s []T) []T {
	ans := make([]T, len(s))
	for i, x := range s {
		ans[len(s)-1-i] = x
	}
	return ans
}

func Remove[T comparable](s []T, q T) []T {
	idx := slices.Index(s, q)
	if idx > -1 {
		return slices.Delete(s, idx, idx+1)
	}
	return s
}

func RemoveAll[T comparable](s []T, q T) []T {
	ans := s
	for {
		idx := slices.Index(ans, q)
		if idx < 0 {
			break
		}
		ans = slices.Delete(ans, idx, idx+1)
	}
	return ans
}

func Filter[T any](s []T, f func(x T) bool) []T {
	ans := make([]T, 0, len(s))
	for _, x := range s {
		if f(x) {
			ans = append(ans, x)
		}
	}
	return ans
}

func Map[T any, O any](f func(x T) O, s []T) []O {
	ans := make([]O, 0, len(s))
	for _, x := range s {
		ans = append(ans, f(x))
	}
	return ans
}

func Sort[T any](s []T, less func(a, b T) bool) []T {
	sort.Slice(s, func(i, j int) bool { return less(s[i], s[j]) })
	return s
}

func StableSort[T any](s []T, less func(a, b T) bool) []T {
	sort.SliceStable(s, func(i, j int) bool { return less(s[i], s[j]) })
	return s
}

func sort_with_key[T any, C constraints.Ordered](impl func(any, func(int, int) bool), s []T, key func(a T) C) []T {
	temp := make([]struct {
		key C
		val T
	}, len(s))
	for i, x := range s {
		temp[i].val, temp[i].key = x, key(x)
	}
	impl(temp, func(i, j int) bool {
		return temp[i].key < temp[j].key
	})
	for i, x := range temp {
		s[i] = x.val
	}
	return s
}

func SortWithKey[T any, C constraints.Ordered](s []T, key func(a T) C) []T {
	return sort_with_key(sort.Slice, s, key)
}

func StableSortWithKey[T any, C constraints.Ordered](s []T, key func(a T) C) []T {
	return sort_with_key(sort.SliceStable, s, key)
}

func Max[T constraints.Ordered](a T, items ...T) (ans T) {
	ans = a
	for _, q := range items {
		if q > ans {
			ans = q
		}
	}
	return ans
}

func Min[T constraints.Ordered](a T, items ...T) (ans T) {
	ans = a
	for _, q := range items {
		if q < ans {
			ans = q
		}
	}
	return ans
}

func Memset[T any](dest []T, pattern ...T) []T {
	if len(pattern) == 0 {
		var zero T
		switch any(zero).(type) {
		case byte: // special case this as the compiler can generate efficient code for memset of a byte slice to zero
			bd := any(dest).([]byte)
			for i := range bd {
				bd[i] = 0
			}
		default:
			for i := range dest {
				dest[i] = zero
			}
		}
		return dest
	}
	bp := copy(dest, pattern)
	for bp < len(dest) {
		copy(dest[bp:], dest[:bp])
		bp *= 2
	}
	return dest
}
