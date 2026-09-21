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
	AddMonomerNoConnection(monomer *Monomer)
	AddMonomerAtStart(monomer *Monomer)
	GetMonomerByIdx(idx int) *Monomer
	GetFieldType() FieldType
	GetField() IField
}

func NewIPolymer(fieldType FieldType, field IField, args ...any) IPolymer {
	return map[FieldType]func(args ...any) IPolymer{
		FieldTypeLattice: func(args ...any) IPolymer {
			return NewPolymer(field.(*Field), args[0].(int64))
		},
		FieldTypeReal: func(args ...any) IPolymer {
			return NewRealPolymer(field.(*RealField), args[0].(int64))
		},
	}[fieldType](args...)
}
