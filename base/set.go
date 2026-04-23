package base

type Set map[int]struct{}

func (set *Set) Insert(i int) {
	(*set)[i] = struct{}{}
}

func (set *Set) Contains(i int) bool {
	_, contains := (*set)[i]
	return contains
}

func (set *Set) Remove(i int) {
	if set.Contains(i) {
		delete(*set, i)
	}
}
