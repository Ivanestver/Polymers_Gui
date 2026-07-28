package base

type Queue[T any] []T

func (s *Queue[T]) Push(item T) {
	*s = append(*s, item)
}

func (s *Queue[T]) Pop() (any, bool) {
	if s.IsEmpty() {
		return nil, false // Or handle error appropriately
	}
	element := (*s)[0]
	*s = (*s)[1:]
	return element, true
}

func (s Queue[T]) Peek() (any, bool) {
	if s.IsEmpty() {
		return nil, false // Or handle error appropriately
	}
	return s[0], true
}

func (s Queue[T]) PeekNotSafe() any {
	return s[0]
}

func (s Queue[T]) IsEmpty() bool {
	return len(s) == 0
}
