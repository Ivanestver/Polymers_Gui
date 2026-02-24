package views

import (
	"encoding/json"
	"polymers/datatypes"
)

type PolymerView struct {
	name    string
	polymer datatypes.IPolymer
}

func NewPolymerView(polymer datatypes.IPolymer) *PolymerView {
	newPolymerView := new(PolymerView)
	newPolymerView.name = polymer.Name()
	newPolymerView.polymer = polymer

	prev := newPolymerView.polymer.GetMonomerByIdx(0)
	curr := newPolymerView.polymer.GetMonomerByIdx(1)
	for curr != nil && polymer.GetPolymerType() != datatypes.POLYMER_TYPE_REAL {
		datatypes.MakeConnection(prev, curr, datatypes.CONNECTION_TYPE_ONE)
		prev = curr
		curr = curr.NextMonomer
	}

	return newPolymerView
}

func (polymerView *PolymerView) Name() string {
	return polymerView.name
}

func (polymerView *PolymerView) Len() int {
	return polymerView.polymer.Len()
}

func (polymerView *PolymerView) GetStartEndMonomers() (*datatypes.Monomer, *datatypes.Monomer) {
	if polymerView.polymer.Len() == 0 {
		panic("The PolymerView cannot contain an empty polymer")
	}

	return polymerView.polymer.GetMonomerByIdx(0), polymerView.polymer.GetMonomerByIdx(polymerView.polymer.Len() - 1)
}

func ForEachMonomer(polymer *PolymerView, pred func(*datatypes.Monomer) bool) {
	for i := 0; i < polymer.Len(); i++ {
		if !pred(polymer.polymer.GetMonomerByIdx(i)) {
			return
		}
	}
}

func (polymerView *PolymerView) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name    string
		Polymer datatypes.IPolymer
	}{
		Name:    polymerView.name,
		Polymer: polymerView.polymer,
	})
}

func (polymerView *PolymerView) DeepCopy(field *datatypes.Field) *PolymerView {
	newPolymerView := new(PolymerView)
	newPolymerView.name = polymerView.name
	newPolymerView.polymer = polymerView.polymer.DeepCopy(field).(*datatypes.Polymer)
	return newPolymerView
}

func (polymerView *PolymerView) GetUnderlinedField() *datatypes.Field {
	if polymerView.polymer.GetPolymerType() == datatypes.POLYMER_TYPE_LATTICE {
		return polymerView.polymer.(*datatypes.Polymer).Field()
	} else {
		return nil
	}
}

func (polymerView *PolymerView) TrunkTo(newSize int) {
	if newSize >= polymerView.Len() {
		return
	}

	for polymerView.Len() > newSize {
		polymerView.polymer.(*datatypes.Polymer).MakeStepBack()
	}
}
