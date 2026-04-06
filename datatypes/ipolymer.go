package datatypes

type PolymerType int

type IPolymer interface {
	Len() int
	LastMonomer() *Monomer
	Name() string
	Number() int64
	Copy() IPolymer
	DeepCopy(args ...any) IPolymer
	AddMonomer(monomer *Monomer)
	GetMonomerByIdx(idx int) *Monomer
	GetFieldType() FieldType
	GetField() IField
}

func NewIPolymer(fieldType FieldType, args ...any) IPolymer {
	return map[FieldType]func(args ...any) IPolymer{
		FIELD_TYPE_LATTICE: func(args ...any) IPolymer {
			return NewPolymer(args[0].(*Field), args[1].(int64))
		},
		FIELD_TYPE_REAL: func(args ...any) IPolymer {
			return NewRealPolymer(args[0].(*RealField), args[1].(int64))
		},
	}[fieldType](args...)
}
