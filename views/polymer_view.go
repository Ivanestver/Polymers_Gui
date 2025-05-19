package views

import (
	"encoding/json"
	"polymers/datatypes"
)

type PolymerView struct {
	name    string
	polymer *datatypes.Polymer
}

func NewPolymerView(polymer *datatypes.Polymer) *PolymerView {
	newPolymerView := new(PolymerView)
	newPolymerView.name = polymer.Name()
	newPolymerView.polymer = polymer

	prev := newPolymerView.polymer.GetMonomerByIdx(0)
	curr := newPolymerView.polymer.GetMonomerByIdx(1)
	for curr != nil {
		datatypes.MakeConnection(prev, curr, datatypes.CONNECTION_TYPE_ONE)
		prev = curr
		curr = curr.NextMonomer
	}

	return newPolymerView
}

func (this *PolymerView) Name() string {
	return this.name
}

func (this *PolymerView) Len() int {
	return this.polymer.Len()
}

func (this *PolymerView) GetStartEndMonomers() (*datatypes.Monomer, *datatypes.Monomer) {
	if this.polymer.Len() == 0 {
		panic("The PolymerView cannot contain an empty polymer")
	}

	return this.polymer.GetMonomerByIdx(0), this.polymer.GetMonomerByIdx(this.polymer.Len() - 1)
}

func ForEachMonomer(polymer *PolymerView, pred func(*datatypes.Monomer)) {
	for i := 0; i < polymer.Len(); i++ {
		pred(polymer.polymer.GetMonomerByIdx(i))
	}
}

func (this *PolymerView) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name    string
		Polymer *datatypes.Polymer
	}{
		Name:    this.name,
		Polymer: this.polymer,
	})
}
