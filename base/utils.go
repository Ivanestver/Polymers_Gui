package base

import "math"

func EcludianDistance(c1, c2 Vector3D) float64 {
	return math.Sqrt(float64((c2.X-c1.X)*(c2.X-c1.X) + (c2.Y-c1.Y)*(c2.Y-c1.Y) + (c2.Z-c1.Z)*(c2.Z-c1.Z)))
}

func Contains[T comparable](container []T, value T) bool {
	for _, v := range container {
		if v == value {
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

func Min_int[T int](container []T) *T {
	if len(container) == 0 {
		return nil
	}

	var minValue *T = &container[0]
	for i := 1; i < len(container); i++ {
		if *minValue > container[i] {
			minValue = &container[i]
		}
	}
	return minValue
}

func Min_float[T float64](container []T) *T {
	if len(container) == 0 {
		return nil
	}

	var minValue *T = &container[0]
	for i := 1; i < len(container); i++ {
		if *minValue > container[i] {
			minValue = &container[i]
		}
	}
	return minValue
}

func Max_int[T int](container []T) *T {
	if len(container) == 0 {
		return nil
	}

	var maxValue *T = &container[0]
	for i := 1; i < len(container); i++ {
		if *maxValue < container[i] {
			maxValue = &container[i]
		}
	}
	return maxValue
}

func Max_float[T float64](container []T) T {
	if len(container) == 0 {
		return T(math.NaN())
	}

	var maxValue T = container[0]
	for i := 1; i < len(container); i++ {
		if maxValue < container[i] {
			maxValue = container[i]
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

func Index_if[T any](container []T, pred func(t T) bool) int {
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
