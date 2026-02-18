package savers

import (
	"math"
	"polymers/datatypes"
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
	lammpsStruct := lammps_structs.NewEmptyLammpsStruct()
	writeSpaceDimention(lammpsStruct)
	writeAtoms(globula, lammpsStruct)
	writeBonds(globula, lammpsStruct)
	return lammpsStruct, nil
}

func writeSpaceDimention(lammpsStruct *lammps_structs.LammpsStruct) {
	spaceDimention := global_data.GetGlobalData().SpaceDimention
	lammpsStruct.SpaceDimention[lammps_structs.DIMENTION_TYPE_X] = [2]float64{0, float64(spaceDimention.X)}
	lammpsStruct.SpaceDimention[lammps_structs.DIMENTION_TYPE_Y] = [2]float64{0, float64(spaceDimention.Y)}
	lammpsStruct.SpaceDimention[lammps_structs.DIMENTION_TYPE_Z] = [2]float64{0, float64(spaceDimention.Z)}
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
			Item1: p.Item1,
			Item2: 1.0,
		})
	}
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
				Item1: len(bondTypes) + 1,
				Item2: 1.0,
				Item3: 100.0,
			}
			bondTypes[connectionType] = value
		}

		lammpsStruct.Bonds = append(lammpsStruct.Bonds, lammps_structs.Bond{
			BondID:         *bondID,
			ConnectionType: value.Item1,
			Ends:           [2]int{int(mon1.Number), int(mon2.Number)},
		})
	}
}

func getAtomsCount(globula *views.GlobulaView) int {
	atomsCount := 0
	if !globula.Is(views.GLOBULA_WATERIZED) {
		views.ForEachPolymer(globula, func(pol *views.PolymerView) {
			atomsCount += pol.Len()
		})
	} else {
		globalData := global_data.GetGlobalData()
		atomsCount = int(globalData.SpaceDimention.X) * int(globalData.SpaceDimention.Y) * int(globalData.SpaceDimention.Z)
	}
	return atomsCount
}

func getMonomerTypes(globula *views.GlobulaView) []dt.MonomerType {
	monomersTypes_map := make(map[dt.MonomerType]bool)
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		views.ForEachMonomer(pol, func(mon *datatypes.Monomer) bool {
			monomersTypes_map[mon.MonomerType] = true
			return true
		})
	})

	monomersTypes := make([]dt.MonomerType, 0)
	for key := range monomersTypes_map {
		monomersTypes = append(monomersTypes, key)
	}
	if globula.Is(views.GLOBULA_WATERIZED) {
		monomersTypes = append(monomersTypes, dt.MONOMER_TYPE_WATER)
	}
	slices.SortFunc(monomersTypes, func(a, b dt.MonomerType) int {
		if int(a) < int(b) {
			return -1
		} else if int(a) == int(b) {
			return 0
		} else {
			return 1
		}
	})
	return monomersTypes

	/*
		monomerTypes := make([]dt.MonomerType, 0)
		mTypes := globula.GetLiterals()
		for k := range *mTypes {
			monomerTypes = append(monomerTypes, k)
		}
		slices.SortFunc(monomerTypes, func(a, b dt.MonomerType) int {
			if int(a) < int(b) {
				return -1
			} else if int(a) == int(b) {
				return 0
			} else {
				return 1
			}
		})
		return monomerTypes
	*/
}

func getBondTypes(globula *views.GlobulaView) ([]dt.ConnectionType, int) {
	connTypes_map := make(map[dt.ConnectionType]bool)
	allSides := dt.GetAllSides()
	bondsCount := 0
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		views.ForEachMonomer(pol, func(mon *dt.Monomer) bool {
			for _, side := range allSides {
				connType := mon.GetTypeOfConnectionWithSide(side)
				sibling, err := mon.GetSibling(side)
				if err == nil && connType != dt.CONNECTION_TYPE_UNDEFINED && sibling.Number != -1 {
					connTypes_map[connType] = true
					bondsCount += 1
				}
			}
			return true
		})
	})
	types := make([]dt.ConnectionType, 0)
	for key := range connTypes_map {
		types = append(types, key)
	}
	slices.Sort(types)
	return types, int(bondsCount / 2)
}
