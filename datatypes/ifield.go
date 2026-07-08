package datatypes

import (
	"polymers/base"
	"polymers/globaldata"
)

type FieldType = int

const (
	FieldTypeReal FieldType = iota
	FieldTypeLattice
)

type IField interface {
	GetType() FieldType
	MakeFilled(monomer *Monomer)
	MakeFree(monomer *Monomer)
	IsFree(coords base.Vector3DF) bool
	GetMonomerByCoords(coords base.Vector3DF) *Monomer
	MoveMonomer(monomer *Monomer, to base.Vector3DF) error
	GetMonomersWithin(lower, higher base.Vector3DF) []*Monomer
	GetMinMonomersByAxis(axis base.Axis) []*Monomer
	GetMaxMonomersByAxis(axis base.Axis) []*Monomer
	DeepCopy() IField
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
	globalData := globaldata.GetGlobalData()
	return uint64(max(globalData.SpaceDimention[base.AxisX].Higher, globalData.SpaceDimention[base.AxisY].Higher, globalData.SpaceDimention[base.AxisZ].Higher))
}
