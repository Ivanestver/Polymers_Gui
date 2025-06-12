package build_globula

import (
	"fmt"
	"math"
	"math/rand"
	"polymers/base"
	"polymers/datatypes"
)

type CalcAlgInputData struct {
	GlobulaCount     int
	PolymersCount    int
	AcceptThreshold  float64
	MaxMonomersCount int
	SphereRadius     int
}

type CalcAlg struct {
	inputData CalcAlgInputData
}

func CreateCalcAlg(inputData CalcAlgInputData) CalcAlg {
	return CalcAlg{
		inputData: inputData,
	}
}

func (this *CalcAlg) Calc() []*datatypes.Polymer {
	field := datatypes.NewField(uint64(this.inputData.SphereRadius))
	polymers := make([]*datatypes.Polymer, this.inputData.PolymersCount)
	for i := 0; i < this.inputData.PolymersCount; i++ {
		polymers[i] = datatypes.NewPolymer(field, int64(i))
	}

	return this.calc_impl(polymers, field)
}

func (this *CalcAlg) calc_impl(polymers []*datatypes.Polymer, field *datatypes.Field) []*datatypes.Polymer {
	finishedPolymers := make([]*datatypes.Polymer, 0)
	for _, p := range polymers {
		startMonomer := field.DefineStartMonomer()
		p.AddMonomer(startMonomer)
		field.MakeFilled(startMonomer)
	}

	blacklist := make([]int, 0)
	for len(finishedPolymers) != len(polymers) {
		for i, polymer := range polymers {
			if base.Contains[int](blacklist, i) {
				continue
			}

			if polymer.Len() == this.inputData.MaxMonomersCount {
				continue
			}

			currentMonomer := polymer.LastMonomer()
			availableCells := field.GetAvailableCells(currentMonomer.Coords())
			if len(availableCells) == 0 {
				finishedPolymers = append(finishedPolymers, polymer)
				blacklist = append(blacklist, i)
				continue
			}

			continuations := this.getContinuations(len(availableCells), availableCells)
			potentialConfigs := make([]*datatypes.Polymer, len(continuations))
			for i := 0; i < len(potentialConfigs); i++ {
				potentialConfigs[i] = this.getNextConfig(polymer, continuations[i].Copy())
			}

			currentPosition := this.getNextCurrentPosition(potentialConfigs, polymer.CalcEnergy())
			if currentPosition.IsInvalid() {
				continue
			}
			currentMonomer = field.GetMonomerByCoords(currentPosition)
			polymer.AddMonomer(currentMonomer)
			field.MakeFilled(currentMonomer)

			persentage := float64(polymer.Len()) / float64(this.inputData.MaxMonomersCount) * 100
			intPersentage := int(persentage)
			if persentage-float64(intPersentage) < 0.1 {
				fmt.Printf("\t\t%s. Done: %d percent out of 100.\n", polymer.Name(), intPersentage)
			}

			//fmt.Printf("%s's monomersc count: %d\n", polymer.Name(), polymer.Len())
			if polymer.Len() == this.inputData.MaxMonomersCount {
				finishedPolymers = append(finishedPolymers, polymer)
				blacklist = append(blacklist, i)
			}
		}
	}

	return finishedPolymers
}

func (this *CalcAlg) getContinuations(kFree int, availableCells []*datatypes.Monomer) []*datatypes.Monomer {
	chosen_continuations_idxs := rand.Perm(kFree)
	continuations := make([]*datatypes.Monomer, len(chosen_continuations_idxs))
	for i := 0; i < len(chosen_continuations_idxs); i++ {
		continuations[i] = availableCells[chosen_continuations_idxs[i]]
	}
	return continuations
}

func (this *CalcAlg) getNextConfig(currConfig *datatypes.Polymer, continuation *datatypes.Monomer) *datatypes.Polymer {
	configCopy := datatypes.NewPolymer(currConfig.Field(), -1)
	for i := 0; i < currConfig.Len(); i++ {
		mon := currConfig.GetMonomerByIdx(i)
		newMon := mon.Copy()
		configCopy.AddMonomer(newMon)
	}
	configCopy.AddMonomer(continuation)
	return configCopy
}

func (this *CalcAlg) getNextCurrentPosition(potentialConfigs []*datatypes.Polymer, U_current float64) base.Vector3D {
	deltasOfPotentialConfigs := make([]float64, len(potentialConfigs))
	for i := 0; i < len(potentialConfigs); i++ {
		deltasOfPotentialConfigs[i] = U_current - potentialConfigs[i].CalcEnergy()
	}
	maxDelta := base.Max_float(deltasOfPotentialConfigs)
	if math.IsNaN(maxDelta) {
		return base.InvalidVector()
	}

	if !base.All[float64](deltasOfPotentialConfigs, func(delta float64) bool { return delta == maxDelta }) {
		var maxDeltasIndices []int
		for i := 0; i < len(deltasOfPotentialConfigs); i++ {
			if deltasOfPotentialConfigs[i] == maxDelta {
				maxDeltasIndices = append(maxDeltasIndices, i)
			}
		}
		choise := rand.Intn(len(maxDeltasIndices))
		return potentialConfigs[maxDeltasIndices[choise]].LastMonomer().Coords()
	} else {
		for true {
			r := rand.Float64()
			if r < this.inputData.AcceptThreshold {
				return potentialConfigs[rand.Intn(len(potentialConfigs))].LastMonomer().Coords()
			}
		}
	}

	return base.InvalidVector()
}
