package savers

import (
	"math"
	"polymers/base"
	dt "polymers/datatypes"
	"polymers/globaldata"
	"polymers/views"
	"slices"

	serializer "github.com/Ivanestver/lammps-file-parser/serialize"
	lammps_structs "github.com/Ivanestver/lammps-file-parser/structs"
)

func addString(source *string, str string) {
	*source += str + "\n"
}

func addNewLine(source *string) {
	*source += "\n"
}

func bondTypeMass(connType dt.ConnectionType) float64 {
	switch connType {
	case dt.ConnectionTypeOne:
		return 1
	case dt.ConnectionTypeCrosslinks:
		return 1
	case dt.ConnectionTypeCrossLinear:
		return 1
	case dt.ConnectionTypeCrossSurface:
		return math.Sqrt(2)
	case dt.ConnectionTypeCrossSpacial:
		return math.Sqrt(3)
	default:
		return -1
	}
}

type pair[TFirst, TSecond comparable] struct {
	First  TFirst
	Second TSecond
}

type intPair = pair[int64, int64]

func makeOrderedIntPair(first, second int64) intPair {
	if first <= second {
		return intPair{First: first, Second: second}
	} else {
		return intPair{First: second, Second: first}
	}
}

func SaveToLammps(globula *views.GlobulaView) (string, error) {
	lammpsStruct, err := turnGlobulaIntoLammpsStruct(globula)
	if err != nil {
		return "", err
	}
	return serializer.Serialize(lammpsStruct)
}

func turnGlobulaIntoLammpsStruct(globula *views.GlobulaView) (*lammps_structs.LammpsStruct, error) {
	lammpsStruct := lammps_structs.NewLammpsStruct(0, 0, 0, 0)
	writeSpaceDimention(lammpsStruct)
	writeAtoms(globula, lammpsStruct)
	writeBonds(globula, lammpsStruct)
	return lammpsStruct, nil
}

func writeSpaceDimention(lammpsStruct *lammps_structs.LammpsStruct) {
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	lammpsStruct.SpaceDimention[lammps_structs.DimentionTypeX] = [2]float64{spaceDimention[base.AxisX].Lower, spaceDimention[base.AxisX].Higher}
	lammpsStruct.SpaceDimention[lammps_structs.DimentionTypeY] = [2]float64{spaceDimention[base.AxisY].Lower, spaceDimention[base.AxisY].Higher}
	lammpsStruct.SpaceDimention[lammps_structs.DimentionTypeZ] = [2]float64{spaceDimention[base.AxisZ].Lower, spaceDimention[base.AxisZ].Higher}
}

func writeAtoms(globula *views.GlobulaView, lammpsStruct *lammps_structs.LammpsStruct) {
	atomsTypes := make(map[base.MendeleevTableElement]lammps_structs.Pair[int, string])
	updateAtomsInfo := createUpdateAtomsInfo(lammpsStruct, atomsTypes, globula)
	// Write atoms and gather atom types info
	for polNumber := 0; polNumber < globula.Len(); polNumber++ {
		polymer := globula.GetPolymerByIdx(polNumber)
		for monNumber := 0; monNumber < polymer.Len(); monNumber++ {
			monomer := polymer.GetMonomerByIdx(monNumber)
			if monomer.MonomerType == base.MendeleevTableElementUndefined {
				panic("Monomer cannot be undefined inside a polymer")
			}
			updateAtomsInfo(monomer, polNumber+1)
		}
	}

	// Write the atom types info gathered
	for _, p := range atomsTypes {
		lammpsStruct.AtomTypes = append(lammpsStruct.AtomTypes, lammps_structs.AtomType{
			AtomType:  p.Item1,
			AtomMass:  1.0,
			AtomLabel: p.Item2,
		})
	}
	slices.SortFunc(lammpsStruct.AtomTypes, func(atomType1, atomType2 lammps_structs.AtomType) int {
		if atomType1.AtomType < atomType2.AtomType {
			return -1
		} else if atomType1.AtomType == atomType2.AtomType {
			return 0
		} else {
			return 1
		}
	})
}

func createUpdateAtomsInfo(lammpsStruct *lammps_structs.LammpsStruct, atomsTypes map[base.MendeleevTableElement]lammps_structs.Pair[int, string], globula *views.GlobulaView) func(*dt.Monomer, int) {
	return func(monomer *dt.Monomer, polymerID int) {
		p, ok := atomsTypes[monomer.MonomerType]
		if !ok {
			atomsTypes[monomer.MonomerType] = lammps_structs.Pair[int, string]{
				Item1: len(atomsTypes) + 1,
				Item2: string(monomer.MonomerType),
			}
			p = atomsTypes[monomer.MonomerType]
		}
		lammpsStruct.Atoms = append(lammpsStruct.Atoms, lammps_structs.Atom{
			Label:      p.Item2,
			AtomID:     int(monomer.Number),
			MoleculeID: polymerID,
			AtomType:   p.Item1,
			Q:          0.0,
			AtomCoords: lammps_structs.AtomCoords{
				X: monomer.Coords()[base.AxisX],
				Y: monomer.Coords()[base.AxisY],
				Z: monomer.Coords()[base.AxisZ],
			},
		})
	}
}

func writeBonds(globula *views.GlobulaView, lammpsStruct *lammps_structs.LammpsStruct) {
	bondID := 1
	bondTypes := make(map[dt.ConnectionType]lammps_structs.BondType)
	updateBondInfo := createUpdateBondsInfo(lammpsStruct, bondTypes, globula, &bondID)
	usedPairs := make(map[[2]int]bool)
	for polNumber := 0; polNumber < globula.Len(); polNumber++ {
		polymer := globula.GetPolymerByIdx(polNumber)
		for monNumber := 0; monNumber < polymer.Len(); monNumber++ {
			mon := polymer.GetMonomerByIdx(monNumber)
			siblings := mon.GetSiblings()
			for _, sibling := range siblings {
				pair := [2]int{}
				if mon.Number < sibling.Number {
					pair[0] = int(mon.Number)
					pair[1] = int(sibling.Number)
				} else {
					pair[0] = int(sibling.Number)
					pair[1] = int(mon.Number)
				}
				if _, ok := usedPairs[pair]; ok {
					continue
				}
				usedPairs[pair] = true
				updateBondInfo(mon, sibling)
			}
		}
	}

	// Write the atom types info gathered
	for _, p := range bondTypes {
		lammpsStruct.BondTypes = append(lammpsStruct.BondTypes, p)
	}
}

func createUpdateBondsInfo(lammpsStruct *lammps_structs.LammpsStruct, bondTypes map[dt.ConnectionType]lammps_structs.BondType, globula *views.GlobulaView, bondID *int) func(*dt.Monomer, *dt.Monomer) {
	return func(mon1, mon2 *dt.Monomer) {
		side := mon1.GetSideOfSibling(mon2)
		if side == dt.SideUndefined {
			return
		}
		connectionType := mon1.GetTypeOfConnectionWithSide(side)
		if connectionType == dt.ConnectionTypeUndefined { // Just in case
			return
		}

		value, ok := bondTypes[connectionType]
		if !ok {
			value = lammps_structs.BondType{
				BondID: len(bondTypes) + 1,
				Sth1:   1.0,
				Sth2:   100.0,
			}
			bondTypes[connectionType] = value
		}

		lammpsStruct.Bonds = append(lammpsStruct.Bonds, lammps_structs.Bond{
			BondID:         *bondID,
			ConnectionType: value.BondID,
			Ends:           [2]int{int(mon1.Number), int(mon2.Number)},
		})
	}
}
