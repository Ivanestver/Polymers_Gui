package datatypes

import "polymers/base"

type FieldType = int

const (
	FIELD_TYPE_REAL FieldType = iota
	FIELD_TYPE_LATTICE
)

type IField interface {
	MakeFilled(monomer *Monomer)
	MakeFree(monomer *Monomer)
	IsFree(coords base.Vector3DF) bool
	GetMonomerByCoords(coords base.Vector3DF) *Monomer
}

func CreateField(fieldType FieldType, args ...any) IField {
	switch fieldType {
	case FIELD_TYPE_REAL:
		return NewRealField()
	case FIELD_TYPE_LATTICE:
		return NewField(args[0].(uint64))
	}
	return nil
}
