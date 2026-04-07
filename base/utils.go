package base

import (
	"cmp"
	"math"
)

func EcludianDistance(c1, c2 Vector3D) float64 {
	return math.Sqrt(float64((c2.X-c1.X)*(c2.X-c1.X) + (c2.Y-c1.Y)*(c2.Y-c1.Y) + (c2.Z-c1.Z)*(c2.Z-c1.Z)))
}

func EcludianDistanceF(c1, c2 Vector3DF) float64 {
	return math.Sqrt((c2.X-c1.X)*(c2.X-c1.X) + (c2.Y-c1.Y)*(c2.Y-c1.Y) + (c2.Z-c1.Z)*(c2.Z-c1.Z))
}

func Contains[T comparable](container []T, value T) bool {
	for _, v := range container {
		if v == value {
			return true
		}
	}

	return false
}

func ContainsIf[T comparable](container []T, value T, pred func(it T, value T) bool) bool {
	for _, v := range container {
		if pred(v, value) {
			return true
		}
	}
	return false
}

func All[T comparable](container []T, pred func(T) bool) bool {
	for _, value := range container {
		if !pred(value) {
			return false
		}
	}

	return true
}

func Any[T comparable](container []T, pred func(T) bool) bool {
	for _, value := range container {
		if pred(value) {
			return true
		}
	}

	return false
}

func MinInt[T int](container []T) *T {
	if len(container) == 0 {
		return nil
	}

	var minValue = &container[0]
	for i := 1; i < len(container); i++ {
		if *minValue > container[i] {
			minValue = &container[i]
		}
	}
	return minValue
}

func MinFloat[T float64](container []T) *T {
	if len(container) == 0 {
		return nil
	}

	var minValue = &container[0]
	for i := 1; i < len(container); i++ {
		if *minValue > container[i] {
			minValue = &container[i]
		}
	}
	return minValue
}

func Min[T cmp.Ordered](container []T) *T {
	if len(container) == 0 {
		return nil
	}

	var minValue = &container[0]
	for i := 1; i < len(container); i++ {
		m := min(*minValue, container[i])
		if m != *minValue {
			minValue = &container[i]
		}
	}
	return minValue
}

func MaxInt[T int](container []T) *T {
	if len(container) == 0 {
		return nil
	}

	var maxValue = &container[0]
	for i := 1; i < len(container); i++ {
		if *maxValue < container[i] {
			maxValue = &container[i]
		}
	}
	return maxValue
}

func MaxFloat[T float64](container []T) T {
	if len(container) == 0 {
		return T(math.NaN())
	}

	var maxValue = container[0]
	for i := 1; i < len(container); i++ {
		if maxValue < container[i] {
			maxValue = container[i]
		}
	}
	return maxValue
}

func Max[T cmp.Ordered](container []T) *T {
	if len(container) == 0 {
		return nil
	}

	var maxValue = &container[0]
	for i := 1; i < len(container); i++ {
		m := max(*maxValue, container[i])
		if m != *maxValue {
			maxValue = &container[i]
		}
	}
	return maxValue
}

func Index[T comparable](container []T, value T) int {
	for idx, item := range container {
		if item == value {
			return idx
		}
	}

	return -1
}

func IndexIf[T any](container []T, pred func(t T) bool) int {
	for idx, item := range container {
		if pred(item) {
			return idx
		}
	}

	return -1
}

type Summable interface {
	int | int64 | float64
}

func Sum[T Summable](container []T) T {
	var sum T
	for _, item := range container {
		sum += item
	}
	return sum
}

func CompareFloat(left, right float64) bool {
	return math.Abs(left-right) <= 0.001
}

func PointInSpace(coords, lower, higher *Vector3DF) bool {
	return (lower.X < coords.X || CompareFloat(lower.X, coords.X)) &&
		(coords.X < higher.X || CompareFloat(higher.X, coords.X)) &&
		(lower.Y < coords.Y || CompareFloat(lower.Y, coords.Y)) &&
		(coords.Y < higher.Y || CompareFloat(coords.Y, higher.Y)) &&
		(lower.Z < coords.Z || CompareFloat(lower.Z, coords.Z)) &&
		(coords.Z < higher.Z || CompareFloat(coords.Z, higher.Z))
}
