package arrayx

import (
	"golang.org/x/exp/constraints"
	"sort"
)

type ArrayX[T constraints.Ordered] []T

func (a ArrayX[T]) Sort() {
	sort.Sort(&a)
}

func (a ArrayX[T]) Len() int {
	return len(a)
}

func (a ArrayX[T]) Less(i, j int) bool {
	return a[i] < a[j]
}

func (a ArrayX[T]) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ArrayX[T]) GetMax() T {
	a.Sort()
	return a[len(a)-1]
}

func (a ArrayX[T]) GetMin() T {
	a.Sort()
	return a[0]
}

func (a ArrayX[T]) Mid() (T, int) {
	a.Sort()
	ll := len(a)
	if ll%2 == 0 {
		return a[ll/2] + a[ll/2+1], 2
	}
	return a[ll/2+1], 1
}

func (a ArrayX[T]) Sum() (sum T) {
	for i := range a {
		sum += a[i]
	}
	return
}

func (a ArrayX[T]) IsContain(v T) bool {
	for i := range a {
		if a[i] == v {
			return true
		}
	}
	return false
}

func (a ArrayX[T]) Append(v ...T) ArrayX[T] {
	c := make(ArrayX[T], a.Len()+len(v))
	copy(c, a)
	pre := a.Len()
	for i := range v {
		c[pre+i] = v[i]
	}
	return c
}

func Delete[T comparable](src []T, v T) []T {
	delId := 0
	if len(src) == 0 {
		return src
	}
	for i := range src {
		if src[i] == v {
			delId = i
			break
		}
		if i == len(src)-1 {
			return src
		}
	}
	return Append[T]([]T{}, Append(src[:delId], src[delId+1:]...)...)
}

func Append[T comparable](src []T, v ...T) []T {
	c := make([]T, len(src)+len(v))
	copy(c, src)
	for i := range v {
		c[len(src)+i] = v[i]
	}
	return c
}

// Equal Common
func Equal[T comparable](p1, p2 T) bool {
	if p1 == p2 {
		return true
	}
	return false
}

//func EqualCommon[T any](p1, p2 T) {
//	Equal[T](p1, p2)
//}
