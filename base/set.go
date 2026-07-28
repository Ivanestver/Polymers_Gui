package base

type UnorderedSet[T comparable] map[T]struct{}

func (set *UnorderedSet[T]) Insert(i T) {
	(*set)[i] = struct{}{}
}

func (set *UnorderedSet[T]) Contains(i T) bool {
	_, contains := (*set)[i]
	return contains
}

func (set *UnorderedSet[T]) Remove(i T) {
	if set.Contains(i) {
		delete(*set, i)
	}
}

func (set UnorderedSet[T]) IsEmpty() bool {
	return len(set) == 0
}

func (set UnorderedSet[T]) Size() int {
	return len(set)
}
