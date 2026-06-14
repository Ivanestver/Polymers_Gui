package datatypes

import (
	"encoding/json"
	"math"
	"polymers/base"
	"strconv"
)

type Polymer struct {
	field         *Field
	polymer       []*Monomer
	polymerNumber int64
}

func NewPolymer(field *Field, polymerNumber int64) *Polymer {
	newPolymer := new(Polymer)
	newPolymer.field = field
	newPolymer.polymerNumber = polymerNumber
	return newPolymer
}

func (polymer *Polymer) AddMonomer(monomer *Monomer) {
	if len(polymer.polymer) != 0 {
		lastMonomer := polymer.polymer[len(polymer.polymer)-1]
		lastMonomer.NextMonomer = monomer
		monomer.PrevMonomer = lastMonomer
		MakeConnection(polymer.polymer[len(polymer.polymer)-1], monomer, ConnectionTypeOne)
	}
	polymer.polymer = append(polymer.polymer, monomer)
	polymer.field.MakeFilled(monomer)
}

func (polymer *Polymer) AddMonomerNoConnection(monomer *Monomer) {
	if len(polymer.polymer) != 0 {
		lastMonomer := polymer.polymer[len(polymer.polymer)-1]
		lastMonomer.NextMonomer = monomer
		monomer.PrevMonomer = lastMonomer
	}
	polymer.polymer = append(polymer.polymer, monomer)
	polymer.field.MakeFilled(monomer)
}

func (polymer *Polymer) Len() int {
	return len(polymer.polymer)
}

func (polymer *Polymer) LastMonomer() *Monomer {
	return polymer.polymer[polymer.Len()-1]
}

func (polymer *Polymer) Name() string {
	return string("Polymer ") + strconv.FormatInt(polymer.polymerNumber, 10)
}

func (polymer *Polymer) Number() int64 {
	return polymer.polymerNumber
}

func (polymer *Polymer) GetField() IField {
	return polymer.field
}

func (polymer *Polymer) Copy() IPolymer {
	newPolymer := NewPolymer(polymer.field, polymer.polymerNumber)
	newPolymer.polymer = polymer.polymer
	return newPolymer
}

func (polymer *Polymer) DeepCopy(args ...any) IPolymer {
	field := args[0].(*Field)
	newPolymer := new(Polymer)
	newPolymer.polymerNumber = polymer.polymerNumber
	if field == nil {
		newPolymer.field = polymer.field.DeepCopy().(*Field)
	} else {
		newPolymer.field = field
	}

	newPolymer.polymer = make([]*Monomer, len(polymer.polymer))
	for i, mon := range polymer.polymer {
		newPolymer.polymer[i] = newPolymer.field.GetMonomerByCoords(mon.coords)
	}

	return newPolymer
}

func (polymer *Polymer) GetMonomerByIdx(idx int) *Monomer {
	if idx < 0 || idx >= polymer.Len() {
		return nil
	}

	return polymer.polymer[idx]
}

func CalcEnergy(polymer IPolymer) float64 {
	u := 0.0
	lastPoint := polymer.LastMonomer()
	var prelastPoint *Monomer
	if polymer.Len() > 1 {
		prelastPoint = polymer.GetMonomerByIdx(polymer.Len() - 2)
	} else {
		prelastPoint = lastPoint
	}
	for _, side := range GetMovementSides() {
		sibling, err := lastPoint.GetSibling(side)
		if err == nil && sibling.IsNotTypeOf(base.MendeleevTableElementUndefined) && !MonomersAreEqual(sibling, prelastPoint) {
			u += -1.0
		}
	}

	return u
}

func (polymer *Polymer) CalcLagevenEnergy() float64 {
	u := 0.0
	for i := 0; i < polymer.Len(); i++ {
		for j := 0; j < polymer.Len(); j++ {
			if i == j {
				continue
			}

			rij := 1 / distanceOfMonomers(polymer.polymer[i], polymer.polymer[j])
			u += 4 * 0.01 * (math.Pow(rij, 12) - math.Pow(rij, 6))
		}
	}
	return u
}

func (polymer *Polymer) MakeStepBack() bool {
	if polymer.Len() <= 1 {
		return false
	}

	polymer.field.MakeFree(polymer.polymer[polymer.Len()-1])
	mon := polymer.polymer[polymer.Len()-1]
	allSides := GetAllSides()
	for _, side := range allSides {
		other, _ := mon.GetSibling(side)
		if other != nil {
			BreakConnection(mon, other, side)
		}
	}
	if mon.PrevMonomer != nil {
		mon.PrevMonomer.NextMonomer = nil
		mon.PrevMonomer = nil
	}
	polymer.polymer = polymer.polymer[:polymer.Len()-1]
	return true
}

func (polymer *Polymer) MarshalJSON() ([]byte, error) {
	pol := make([]base.Vector3DF, polymer.Len())
	for i, item := range polymer.polymer {
		pol[i] = item.coords
	}
	return json.Marshal(&struct {
		Field         *Field
		Polymer       []base.Vector3DF
		PolymerNumber int64
	}{
		Field:         polymer.field,
		Polymer:       pol,
		PolymerNumber: polymer.polymerNumber,
	})
}

func (polymer *Polymer) GetFieldType() FieldType {
	return FieldTypeLattice
}
