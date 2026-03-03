package datatypes

import "strconv"

type RealPolymer struct {
	number   int64
	monomers []*Monomer
	field    *RealField
}

func NewRealPolymer(field *RealField, polymerNumber int64) *RealPolymer {
	return &RealPolymer{
		number: polymerNumber,
		field:  field,
	}
}

func (realPolymer *RealPolymer) Len() int {
	return len(realPolymer.monomers)
}

func (realPolymer *RealPolymer) LastMonomer() *Monomer {
	if realPolymer.Len() > 0 {
		return realPolymer.monomers[realPolymer.Len()-1]
	}
	return nil
}

func (realPolymer *RealPolymer) Name() string {
	return string("Polymer ") + strconv.FormatInt(realPolymer.number, 10)
}

func (realPolymer *RealPolymer) Number() int64 {
	return realPolymer.number
}

func (realPolymer *RealPolymer) Copy() IPolymer {
	newPolymer := NewRealPolymer(realPolymer.field, realPolymer.number)
	newPolymer.monomers = realPolymer.monomers
	return newPolymer
}

func (realPolymer *RealPolymer) DeepCopy(args ...any) IPolymer {
	return nil
}

func (realPolymer *RealPolymer) AddMonomer(monomer *Monomer) {
	if monomer == nil {
		return
	}

	if realPolymer.Len() == 0 {
		realPolymer.monomers = append(realPolymer.monomers, monomer)
		return
	}
	lastMonomer := realPolymer.LastMonomer()
	lastMonomer.NextMonomer = monomer
	monomer.PrevMonomer = lastMonomer
	realPolymer.monomers = append(realPolymer.monomers, monomer)
}

func (realPolymer *RealPolymer) GetMonomerByIdx(idx int) *Monomer {
	if idx < 0 || idx >= realPolymer.Len() {
		return nil
	}

	return realPolymer.monomers[idx]
}

func (realPolymer *RealPolymer) GetFieldType() FieldType {
	return FIELD_TYPE_REAL
}
