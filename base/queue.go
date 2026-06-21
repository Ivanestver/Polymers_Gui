package base

type Queue []any

func (s *Queue) Push(item any) {
	*s = append(*s, item)
}

func (s *Queue) Pop() (any, bool) {
	if s.IsEmpty() {
		return nil, false // Or handle error appropriately
	}
	element := (*s)[0]
	*s = (*s)[1:]
	return element, true
}

func (s *Queue) Peek() (any, bool) {
	if s.IsEmpty() {
		return nil, false // Or handle error appropriately
	}
	return (*s)[0], true
}

func (s *Queue) PeekNotSafe() any {
	return (*s)[0]
}

func (s *Queue) IsEmpty() bool {
	return len(*s) == 0
}
