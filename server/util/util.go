// DANGER: if you are adding to this, you are likely doing something wrong.
// Only add simple, standalone functions that are frequently reused.
// This file should essentially be a shim providing std-lib-like functions.
package util

func RemoveElementFromSlice[T any](s []T, i int) []T {
	ret := make([]T, 0)
	ret = append(ret, s[:i]...)
	return append(ret, s[i+1:]...)
}
