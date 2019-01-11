package ddsort

import (
	"sort"
)

// int64 排序
type int64Sort []int64

func (a int64Sort) Len() int {
	return len(a)
}
func (a int64Sort) Less(i, j int) bool {
	return a[i] < a[j]
}
func (a int64Sort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// 顺序排序
func SortInt64(b []int64) []int64 {
	if len(b) <= 0 {
		return nil
	}
	nb := make([]int64, len(b))
	copy(nb, b)
	sort.Sort(int64Sort(nb))
	return nb
}

// 逆序排序
func SortReverseInt64(b []int64) []int64 {
	if len(b) <= 0 {
		return nil
	}
	nb := make([]int64, len(b))
	copy(nb, b)
	sort.Sort(sort.Reverse(int64Sort(nb)))
	return nb
}
