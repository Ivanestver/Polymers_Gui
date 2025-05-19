package datatypes

import (
	"encoding/json"
	"math"
	"polymers/base"
	"polymers/global_data"
	"strconv"
)

type Polymer struct {
	field         *Field
	polymer       []*Monomer
	polymerNumber int64
}

func NewPolymer(field *Field, polymerNumber int64) *Polymer {
	var newPolymer *Polymer = new(Polymer)
	newPolymer.field = field
	newPolymer.polymerNumber = polymerNumber
	return newPolymer
}

func (this *Polymer) AddMonomer(monomer *Monomer) {
	if len(this.polymer) != 0 {
		var lastMonomer *Monomer = this.polymer[len(this.polymer)-1]
		lastMonomer.NextMonomer = monomer
		monomer.PrevMonomer = lastMonomer
	}
	this.polymer = append(this.polymer, monomer)
}

func (this *Polymer) Len() int {
	return len(this.polymer)
}

func (this *Polymer) LastMonomer() *Monomer {
	return this.polymer[this.Len()-1]
}

func (this *Polymer) Name() string {
	return string("Polymer ") + strconv.FormatInt(this.polymerNumber, 10)
}

func (this *Polymer) Number() int64 {
	return this.polymerNumber
}

func (this *Polymer) Field() *Field {
	return this.field
}

func (this *Polymer) Copy() *Polymer {
	newPolymer := NewPolymer(this.field, this.polymerNumber)
	newPolymer.polymer = this.polymer
	return newPolymer
}

func (this *Polymer) GetMonomerByIdx(idx int) *Monomer {
	if idx < 0 || idx >= this.Len() {
		return nil
	}

	return this.polymer[idx]
}

func (this *Polymer) CalcEnergy() float64 {
	u := 0.0
	last_point := this.LastMonomer()
	var prelast_point *Monomer
	if this.Len() > 1 {
		prelast_point = this.polymer[this.Len()-2]
	} else {
		prelast_point = last_point
	}
	for _, side := range GetMovementSides() {
		sibling, err := last_point.GetSibling(side)
		if err == nil && sibling.IsNotTypeOf(MONOMER_TYPE_UNDEFINED) && !MonomersAreEqual(sibling, prelast_point) {
			u += -1.0
		}
	}

	return u
}

func (this *Polymer) CalcLagevenEnergy() float64 {
	u := 0.0
	for i := 0; i < this.Len(); i++ {
		for j := 0; j < this.Len(); j++ {
			if i == j {
				continue
			}

			r_ij := 1 / distance_of_monomers(this.polymer[i], this.polymer[j])
			u += 4 * 0.01 * (math.Pow(r_ij, 12) - math.Pow(r_ij, 6))
		}
	}
	return u
}

func (this *Polymer) MakeStepBack() bool {
	if this.Len() <= 1 {
		return false
	}

	this.polymer = this.polymer[:this.Len()-1]
	return true
}

func (this *Polymer) GetMinMaxWidthHeight() (float64, float64, float64, float64) {
	globalData := global_data.GetGlobalData()
	minWidth := float64(globalData.SpaceDimention)
	maxWidth := 0.0
	minHeight := float64(globalData.SpaceDimention)
	maxHeight := 0.0

	for _, mon := range this.polymer {
		minWidth = math.Min(minWidth, float64(mon.coords.X))
		maxWidth = math.Max(maxWidth, float64(mon.coords.Y))
		minHeight = math.Min(minHeight, float64(mon.coords.X))
		maxHeight = math.Max(maxHeight, float64(mon.coords.Y))
	}

	return minWidth, maxWidth, minHeight, maxHeight
}

func (this *Polymer) MarshalJSON() ([]byte, error) {
	pol := make([]base.Vector3D, this.Len())
	for i, item := range this.polymer {
		pol[i] = item.coords
	}
	return json.Marshal(&struct {
		Field         *Field
		Polymer       []base.Vector3D
		PolymerNumber int64
	}{
		Field:         this.field,
		Polymer:       pol,
		PolymerNumber: this.polymerNumber,
	})
}
