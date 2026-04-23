package base

type Queue []interface{}

func (s *Queue) Push(item interface{}) {
	*s = append(*s, item)
}

func (s *Queue) Pop() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false // Or handle error appropriately
	}
	element := (*s)[0]
	*s = (*s)[1:]
	return element, true
}

func (s *Queue) Peek() (interface{}, bool) {
	if s.IsEmpty() {
		return nil, false // Or handle error appropriately
	}
	return (*s)[0], true
}

func (s *Queue) PeekNotSafe() interface{} {
	return (*s)[0]
}

func (s *Queue) IsEmpty() bool {
	return len(*s) == 0
}
