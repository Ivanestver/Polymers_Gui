package datatypes

type PolymerType int

const (
	POLYMER_TYPE_LATTICE PolymerType = iota
	POLYMER_TYPE_REAL
)

type IPolymer interface {
	Len() int
	LastMonomer() *Monomer
	Name() string
	Number() int64
	Copy() IPolymer
	DeepCopy(args ...any) IPolymer
	AddMonomer(monomer *Monomer)
	GetMonomerByIdx(idx int) *Monomer
	GetPolymerType() PolymerType
}

func NewIPolymer(polymerType PolymerType) IPolymer {
	return map[PolymerType]func(args ...any) IPolymer{
		POLYMER_TYPE_LATTICE: func(args ...any) IPolymer {
			return NewPolymer(args[0].(*Field), args[1].(int64))
		},
		POLYMER_TYPE_REAL: func(args ...any) IPolymer {
			return NewRealPolymer(args[0].(int64))
		},
	}[polymerType]()
}
