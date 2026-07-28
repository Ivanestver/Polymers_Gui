package base

import "container/list"

type SmartSet[v comparable] struct {
	list     *list.List
	registry map[v]*list.Element
}

func NewSmartSet[v comparable]() *SmartSet[v] {
	return &SmartSet[v]{
		list:     list.New(),
		registry: make(map[v]*list.Element),
	}
}

func (s SmartSet[v]) Size() int {
	return s.list.Len()
}

func (s SmartSet[v]) IsEmpty() bool {
	return s.Size() == 0
}

func (s SmartSet[v]) Contains(val v) bool {
	_, ok := s.registry[val]
	return ok
}

func (s *SmartSet[v]) Add(val v) {
	if !s.Contains(val) {
		elem := s.list.PushBack(val)
		s.registry[val] = elem
	}
}

func (s *SmartSet[v]) Remove(val v) {
	if s.Contains(val) {
		elem := s.registry[val]
		delete(s.registry, val)
		s.list.Remove(elem)
	}
}

func (s SmartSet[v]) Get(idx int) v {
	if 0 <= idx && idx < s.Size() {
		elem := s.list.Front()
		for range idx {
			elem = elem.Next()
		}
		return elem.Value.(v)
	}
	var ret v
	return ret
}
