package savers

import (
	"math"
	"polymers/base"
	dt "polymers/datatypes"
	"polymers/global_data"
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
	case dt.CONNECTION_TYPE_ONE:
		return 1
	case dt.CONNECTION_TYPE_CROSSLINKS:
		return 1
	case dt.CONNECTION_TYPE_CROSS_LINEAR:
		return 1
	case dt.CONNECTION_TYPE_CROSS_SURFACE:
		return math.Sqrt(2)
	case dt.CONNECTION_TYPE_CROSS_SPACIAL:
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
	spaceDimention := global_data.GetGlobalData().SpaceDimention
	lammpsStruct.SpaceDimention[lammps_structs.DIMENTION_TYPE_X] = [2]float64{spaceDimention[base.X_AXIS].Lower, spaceDimention[base.X_AXIS].Higher}
	lammpsStruct.SpaceDimention[lammps_structs.DIMENTION_TYPE_Y] = [2]float64{spaceDimention[base.Y_AXIS].Lower, spaceDimention[base.Y_AXIS].Higher}
	lammpsStruct.SpaceDimention[lammps_structs.DIMENTION_TYPE_Z] = [2]float64{spaceDimention[base.Z_AXIS].Lower, spaceDimention[base.Z_AXIS].Higher}
}

func writeAtoms(globula *views.GlobulaView, lammpsStruct *lammps_structs.LammpsStruct) {
	polymerID := 1
	atomsTypes := make(map[dt.MonomerType]lammps_structs.Pair[int, string])
	updateAtomsInfo := createUpdateAtomsInfo(lammpsStruct, atomsTypes, globula)
	// Write atoms and gather atom types info
	views.ForEachPolymer(globula, func(polymer *views.PolymerView) {
		views.ForEachMonomer(polymer, func(monomer *dt.Monomer) bool {
			if monomer.MonomerType == dt.MONOMER_TYPE_UNDEFINED {
				panic("Monomer cannot be undefined inside a polymer")
			}
			updateAtomsInfo(monomer, polymerID)
			return true
		})
		polymerID++
	})

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

func createUpdateAtomsInfo(lammpsStruct *lammps_structs.LammpsStruct, atomsTypes map[dt.MonomerType]lammps_structs.Pair[int, string], globula *views.GlobulaView) func(*dt.Monomer, int) {
	return func(monomer *dt.Monomer, polymerID int) {
		p, ok := atomsTypes[monomer.MonomerType]
		if !ok {
			atomsTypes[monomer.MonomerType] = lammps_structs.Pair[int, string]{
				Item1: len(atomsTypes) + 1,
				Item2: globula.GetLiteral(monomer.MonomerType),
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
				X: monomer.Coords().X,
				Y: monomer.Coords().Y,
				Z: monomer.Coords().Z,
			},
		})
	}
}

func writeBonds(globula *views.GlobulaView, lammpsStruct *lammps_structs.LammpsStruct) {
	bondID := 1
	bondTypes := make(map[dt.ConnectionType]lammps_structs.BondType)
	updateBondInfo := createUpdateBondsInfo(lammpsStruct, bondTypes, globula, &bondID)
	allSides := dt.GetAllSides()
	usedPairs := make(map[[2]int]bool)
	views.ForEachPolymer(globula, func(polymer *views.PolymerView) {
		views.ForEachMonomer(polymer, func(mon *dt.Monomer) bool {
			for _, side := range allSides {
				otherMon, err := mon.GetSibling(side)
				if err != nil {
					continue
				}
				pair := [2]int{}
				if mon.Number < otherMon.Number {
					pair[0] = int(mon.Number)
					pair[1] = int(otherMon.Number)
				} else {
					pair[0] = int(otherMon.Number)
					pair[1] = int(mon.Number)
				}
				if _, ok := usedPairs[pair]; ok {
					continue
				}
				usedPairs[pair] = true
				updateBondInfo(mon, otherMon)
			}
			return true
		})
	})

	// Write the atom types info gathered
	for _, p := range bondTypes {
		lammpsStruct.BondTypes = append(lammpsStruct.BondTypes, p)
	}
}

func createUpdateBondsInfo(lammpsStruct *lammps_structs.LammpsStruct, bondTypes map[dt.ConnectionType]lammps_structs.BondType, globula *views.GlobulaView, bondID *int) func(*dt.Monomer, *dt.Monomer) {
	return func(mon1, mon2 *dt.Monomer) {
		side := mon1.GetSideOfSibling(mon2)
		if side == dt.SIDE_Undefined {
			return
		}
		connectionType := mon1.GetTypeOfConnectionWithSide(side)
		if connectionType == dt.CONNECTION_TYPE_UNDEFINED { // Just in case
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
