package datatypes

import (
	"polymers/base"
	"polymers/global_data"
)

type FieldType = int

const (
	FieldTypeReal FieldType = iota
	FieldTypeLattice
)

type IField interface {
	MakeFilled(monomer *Monomer)
	MakeFree(monomer *Monomer)
	IsFree(coords base.Vector3DF) bool
	GetMonomerByCoords(coords base.Vector3DF) *Monomer
	GetMonomersWithin(lower, higher base.Vector3DF) []*Monomer
}

func CreateField(fieldType FieldType, args ...any) IField {
	switch fieldType {
	case FieldTypeReal:
		return NewRealField(args[0].([3][2]float64))
	case FieldTypeLattice:
		return NewField(getMaxDimention())
	}
	return nil
}

func getMaxDimention() uint64 {
	globalData := global_data.GetGlobalData()
	return uint64(max(globalData.SpaceDimention[base.AxisX].Higher, globalData.SpaceDimention[base.AxisY].Higher, globalData.SpaceDimention[base.AxisZ].Higher))
}
