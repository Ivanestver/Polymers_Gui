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

func (data CalcAlgInputData) GetName() string {
	return "Globula"
}

type CalcAlgInputDataBuilder struct {
}

func (creator CalcAlgInputDataBuilder) CreateInputData(algType AlgType) (ICalcAlgInputData, error) {
	inputData := CalcAlgInputData{}

	fmt.Print("Enter the globula count: ")
	fmt.Scanln(&inputData.GlobulaCount)

	fmt.Print("Enter the polymers count: ")
	fmt.Scanln(&inputData.PolymersCount)

	fmt.Print("Enter the accept threshold count: ")
	fmt.Scanln(&inputData.AcceptThreshold)

	fmt.Print("Enter the max monomers count: ")
	fmt.Scanln(&inputData.MaxMonomersCount)

	fmt.Print("Enter the sphere radius: ")
	fmt.Scanln(&inputData.SphereRadius)
	return inputData, nil
}

type CalcAlg struct {
	inputData CalcAlgInputData
}

func (alg *CalcAlg) Calc() []*datatypes.Polymer {
	field := datatypes.NewField(uint64(alg.inputData.SphereRadius))
	polymers := make([]*datatypes.Polymer, alg.inputData.PolymersCount)
	for i := 0; i < alg.inputData.PolymersCount; i++ {
		polymers[i] = datatypes.NewPolymer(field, int64(i))
	}

	return alg.calc_impl(polymers, field)
}

func (alg *CalcAlg) calc_impl(polymers []*datatypes.Polymer, field *datatypes.Field) []*datatypes.Polymer {
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

			if polymer.Len() == alg.inputData.MaxMonomersCount {
				continue
			}

			currentMonomer := polymer.LastMonomer()
			availableCells := field.GetAvailableCells(currentMonomer.Coords())
			if len(availableCells) == 0 {
				finishedPolymers = append(finishedPolymers, polymer)
				blacklist = append(blacklist, i)
				continue
			}

			continuations := alg.getContinuations(len(availableCells), availableCells)
			potentialConfigs := make([]*datatypes.Polymer, len(continuations))
			for i := 0; i < len(potentialConfigs); i++ {
				potentialConfigs[i] = alg.getNextConfig(polymer, continuations[i].Copy())
			}

			currentPosition := alg.getNextCurrentPosition(potentialConfigs, polymer.CalcEnergy())
			if currentPosition.IsInvalid() {
				continue
			}
			currentMonomer = field.GetMonomerByCoords(currentPosition)
			polymer.AddMonomer(currentMonomer)
			field.MakeFilled(currentMonomer)

			persentage := float64(polymer.Len()) / float64(alg.inputData.MaxMonomersCount) * 100
			intPersentage := int(persentage)
			if persentage-float64(intPersentage) < 0.1 {
				fmt.Printf("\t\t%s. Done: %d percent out of 100.\n", polymer.Name(), intPersentage)
			}

			//fmt.Printf("%s's monomersc count: %d\n", polymer.Name(), polymer.Len())
			if polymer.Len() == alg.inputData.MaxMonomersCount {
				finishedPolymers = append(finishedPolymers, polymer)
				blacklist = append(blacklist, i)
			}
		}
	}

	return finishedPolymers
}

func (alg *CalcAlg) getContinuations(kFree int, availableCells []*datatypes.Monomer) []*datatypes.Monomer {
	chosen_continuations_idxs := rand.Perm(kFree)
	continuations := make([]*datatypes.Monomer, len(chosen_continuations_idxs))
	for i := 0; i < len(chosen_continuations_idxs); i++ {
		continuations[i] = availableCells[chosen_continuations_idxs[i]]
	}
	return continuations
}

func (alg *CalcAlg) getNextConfig(currConfig *datatypes.Polymer, continuation *datatypes.Monomer) *datatypes.Polymer {
	configCopy := datatypes.NewPolymer(currConfig.Field(), -1)
	for i := 0; i < currConfig.Len(); i++ {
		mon := currConfig.GetMonomerByIdx(i)
		newMon := mon.Copy()
		configCopy.AddMonomer(newMon)
	}
	configCopy.AddMonomer(continuation)
	return configCopy
}

func (alg *CalcAlg) getNextCurrentPosition(potentialConfigs []*datatypes.Polymer, U_current float64) base.Vector3D {
	deltasOfPotentialConfigs := make([]float64, len(potentialConfigs))
	for i := 0; i < len(potentialConfigs); i++ {
		deltasOfPotentialConfigs[i] = U_current - potentialConfigs[i].CalcEnergy()
	}
	maxDelta := base.Max_float(deltasOfPotentialConfigs)
	if math.IsNaN(maxDelta) {
		return base.InvalidVector()
	}

	if !base.All(deltasOfPotentialConfigs, func(delta float64) bool { return delta == maxDelta }) {
		var maxDeltasIndices []int
		for i := 0; i < len(deltasOfPotentialConfigs); i++ {
			if deltasOfPotentialConfigs[i] == maxDelta {
				maxDeltasIndices = append(maxDeltasIndices, i)
			}
		}
		choise := rand.Intn(len(maxDeltasIndices))
		return potentialConfigs[maxDeltasIndices[choise]].LastMonomer().Coords()
	} else {
		for {
			r := rand.Float64()
			if r < alg.inputData.AcceptThreshold {
				return potentialConfigs[rand.Intn(len(potentialConfigs))].LastMonomer().Coords()
			}
		}
	}
}

func (alg *CalcAlg) GetLiteralsTable() map[datatypes.MonomerType]string {
	m := make(map[datatypes.MonomerType]string)
	m[datatypes.MONOMER_TYPE_UNDEFINED] = ""
	m[datatypes.MONOMER_TYPE_USUAL] = "C"
	m[datatypes.MONOMER_TYPE_VYNIL] = "O"
	m[datatypes.MONOMER_TYPE_O_CONTAINING] = "N"
	m[datatypes.MONOMER_TYPE_FWISE] = "F"
	m[datatypes.MONOMER_TYPE_CLWISE] = "Cl"
	m[datatypes.MONOMER_TYPE_CROSSLINKED] = "H"
	m[datatypes.MONOMER_TYPE_WATER] = "I"
	return m
}
