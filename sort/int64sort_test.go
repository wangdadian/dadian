package ddsort

import (
	"testing"
)

func TestSortInt64(t *testing.T) {
	a := []int64{123, 456, 1, 8, 9, 0, 123451, 99, 77, 55, 0, 1, 99, 789, 456, 99, 0, 8, 99, 10, 0, 123456}
	t.Logf("orig: %v\n", a)
	as := SortInt64(a)
	a[0] = 10001
	t.Logf("sort: %v\n", as)
	ars := SortReverseInt64(a)
	t.Logf("reverse sort: %v\n", ars)
}
