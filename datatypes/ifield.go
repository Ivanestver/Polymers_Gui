package datatypes

import (
	"polymers/base"
	"polymers/global_data"
)

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
		return NewRealField(args[0].([3][2]float64))
	case FIELD_TYPE_LATTICE:
		return NewField(getMaxDimention())
	}
	return nil
}

func getMaxDimention() uint64 {
	globalData := global_data.GetGlobalData()
	return uint64(max(globalData.SpaceDimention[base.X_AXIS].Higher, globalData.SpaceDimention[base.Y_AXIS].Higher, globalData.SpaceDimention[base.Z_AXIS].Higher))
}
