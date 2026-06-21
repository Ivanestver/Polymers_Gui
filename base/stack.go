package base

type Stack []any

func (s *Stack) Push(item any) {
	*s = append(*s, item)
}

func (s *Stack) Pop() (any, bool) {
	if s.IsEmpty() {
		return nil, false // Or handle error appropriately
	}
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element, true
}

func (s *Stack) Peek() (any, bool) {
	if s.IsEmpty() {
		return nil, false // Or handle error appropriately
	}
	index := len(*s) - 1
	return (*s)[index], true
}

func (s *Stack) PeekNotSafe() any {
	return (*s)[len(*s)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}
