package base

type CmpFunc[T any] func(left, right T) bool

type _KeyValuePair[TKey, TValue any] struct {
	Key   TKey
	Value TValue
}

type Dictionary[TKey, TValue any] struct {
	table []_KeyValuePair[TKey, TValue]
	cmp   CmpFunc[TKey]
}

func NewDictionary[TKey, TValue any](cmpFunc CmpFunc[TKey]) *Dictionary[TKey, TValue] {
	return &Dictionary[TKey, TValue]{
		table: make([]_KeyValuePair[TKey, TValue], 0),
		cmp:   cmpFunc,
	}
}

func (m *Dictionary[TKey, TValue]) SetValue(key TKey, value TValue) {
	if idx := m.index(key); idx != -1 {
		m.table[idx].Value = value
	} else {
		m.table = append(m.table, _KeyValuePair[TKey, TValue]{key, value})
	}
}

func (m *Dictionary[TKey, TValue]) AddValue(key TKey, value TValue) {
	if idx := m.index(key); idx == -1 {
		m.table = append(m.table, _KeyValuePair[TKey, TValue]{key, value})
	}
}

func (m *Dictionary[TKey, TValue]) GetValue(key TKey) TValue {
	if idx := m.index(key); idx != -1 {
		return m.table[idx].Value
	} else {
		panic("no such a key")
	}
}

func (m *Dictionary[TKey, TValue]) ContainsKey(key TKey) bool {
	idx := m.index(key)
	return idx > 0
}

func (m *Dictionary[TKey, TValue]) ForEach(callback func(key TKey, value TValue)) {
	for _, p := range m.table {
		callback(p.Key, p.Value)
	}
}

func (m *Dictionary[TKey, TValue]) Len() int {
	return len(m.table)
}

func (m *Dictionary[TKey, TValue]) index(key TKey) int {
	for i, p := range m.table {
		if m.cmp(p.Key, key) {
			return i
		}
	}
	return -1
}
