package savers

import (
	"math"
	"polymers/base"
	"polymers/datatypes"
	dt "polymers/datatypes"
	"polymers/global_data"
	"polymers/views"
	"slices"
	"strconv"
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
	var content string
	addString(&content, "LAMMPS data file via write_data, version 24 Dec 2020, timestep = 40000000")
	addNewLine(&content)

	monomerTypes := getMonomerTypes(globula)
	bondTypes, bondCount := getBondTypes(globula)
	lowerestNumberBond := base.Min(bondTypes)
	highestNumberBond := base.Max(bondTypes)

	addString(&content, strconv.Itoa(getAtomsCount(globula))+" atoms")
	addString(&content, strconv.Itoa(len(monomerTypes))+" atom types")
	addString(&content, strconv.Itoa(bondCount)+" bonds")
	addString(&content, strconv.Itoa(int(*highestNumberBond))+" bond types")
	addNewLine(&content)

	spaceDim := global_data.GetGlobalData().SpaceDimention
	addString(&content, "0 "+strconv.FormatInt(spaceDim.X, 10)+" xlo xhi")
	addString(&content, "0 "+strconv.FormatInt(spaceDim.Y, 10)+" ylo yhi")
	addString(&content, "0 "+strconv.FormatInt(spaceDim.Z, 10)+" zlo zhi")
	addNewLine(&content)

	addString(&content, "Masses")
	addNewLine(&content)
	mapMonomerTypeNumber := make(map[dt.MonomerType]int)
	for i, t := range monomerTypes {
		if t == dt.MONOMER_TYPE_UNDEFINED {
			continue
		}
		mapMonomerTypeNumber[t] = i + 1
		addString(&content, strconv.Itoa(mapMonomerTypeNumber[t])+" 1 # "+globula.GetLiteral(t))
	}

	addNewLine(&content)

	addString(&content, "Bond Coeffs # harmonic")
	addNewLine(&content)

	for i := (int)(*lowerestNumberBond); i <= (int)(*highestNumberBond); i++ {
		addString(&content, strconv.Itoa(i)+" 100 "+strconv.FormatFloat(bondTypeMass((dt.ConnectionType)(i)), 'f', 3, 64))
	}
	addNewLine(&content)

	addString(&content, "Atoms # full")
	addNewLine(&content)

	maxNumber := -1
	polymerNumber := 0
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		polymerNumber++
		views.ForEachMonomer(pol, func(mon *dt.Monomer) bool {
			if mon.MonomerType == dt.MONOMER_TYPE_UNDEFINED {
				return true
			}
			monCoords := mon.Coords()
			maxNumber = *base.Max_int([]int{maxNumber, int(mon.Number)})
			addString(&content,
				strconv.Itoa(int(mon.Number))+
					" "+strconv.Itoa(polymerNumber)+" "+
					strconv.Itoa(mapMonomerTypeNumber[mon.MonomerType])+
					" 0.00000 "+
					strconv.FormatInt(monCoords.X, 10)+" "+
					strconv.FormatInt(monCoords.Y, 10)+" "+
					strconv.FormatInt(monCoords.Z, 10)+" "+
					"0 0 0")
			return true
		})
	})
	if globula.Is(views.GLOBULA_WATERIZED) {
		views.ForEachPolymer_If(globula, func(pv *views.PolymerView) bool {
			field := pv.GetUnderlinedField()
			var globalData *global_data.GlobalData = global_data.GetGlobalData()
			shape := [...]int64{globalData.SpaceDimention.X, globalData.SpaceDimention.Y, globalData.SpaceDimention.Z}
			var i int64
			var j int64
			var k int64
			for i = 0; i < shape[0]; i++ {
				for j = 0; j < shape[1]; j++ {
					for k = 0; k < shape[2]; k++ {
						mon := field.GetMonomerByCoords(base.Vector3D{X: i, Y: j, Z: k})
						if mon == nil || mon.MonomerType != datatypes.MONOMER_TYPE_WATER {
							continue
						}
						maxNumber++
						addString(&content, strconv.Itoa(maxNumber)+
							" 1 "+
							strconv.Itoa(mapMonomerTypeNumber[mon.MonomerType])+
							" 0.00000 "+
							strconv.FormatInt(i, 10)+" "+
							strconv.FormatInt(j, 10)+" "+
							strconv.FormatInt(k, 10)+" "+
							"0 0 0")
					}
				}
			}
			return false
		})
	}
	addNewLine(&content)

	addString(&content, "Bonds")
	addNewLine(&content)

	bondNumber := 1
	used_pairs := make(map[intPair]bool)
	allSides := dt.GetAllSides()
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		views.ForEachMonomer(pol, func(mon *dt.Monomer) bool {
			for _, side := range allSides {
				otherMon, err := mon.GetSibling(side)
				connType := mon.GetTypeOfConnectionWithSide(side)
				if err == nil && connType != dt.CONNECTION_TYPE_UNDEFINED && otherMon.Number != -1 {
					newPair := makeOrderedIntPair(mon.Number, otherMon.Number)
					if _, ok := used_pairs[newPair]; !ok {
						addString(&content, strconv.Itoa(bondNumber)+" "+strconv.Itoa(int(connType))+" "+strconv.Itoa(int(mon.Number))+" "+strconv.Itoa(int(otherMon.Number)))
						bondNumber += 1
						used_pairs[newPair] = true
					}
				}
			}
			return true
		})
	})

	/*
		addNewLine(&content)
		addString(&content, "Angles")
		addNewLine(&content)

		angleNumber := 1
		views.ForEachPolymer(globula, func(pol *views.PolymerView) {
			for i := 1; i < pol.Len(); i++ {
				addString(&content, strconv.Itoa(angleNumber)+" 1 "+strconv.Itoa(i)+" "+strconv.Itoa(i+1)+" "+strconv.Itoa(i+1))
				angleNumber++
			}
		})
	*/

	return content, nil
}

func getAtomsCount(globula *views.GlobulaView) int {
	atomsCount := 0
	if !globula.Is(views.GLOBULA_WATERIZED) {
		views.ForEachPolymer(globula, func(pol *views.PolymerView) {
			views.ForEachMonomer(pol, func(mon *dt.Monomer) bool {
				if mon.MonomerType != dt.MONOMER_TYPE_UNDEFINED {
					atomsCount++
				}
				return true
			})
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
