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
	case dt.CONNECTION_TYPE_TWO:
		return math.Sqrt(2)
	case dt.CONNECTION_TYPE_THREE:
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
	addString(&content, strconv.Itoa(getAtomsCount(globula))+" atoms")
	addString(&content, strconv.Itoa(len(monomerTypes))+" atom types")
	addString(&content, strconv.Itoa(bondCount)+" bonds")
	addString(&content, strconv.Itoa(len(bondTypes))+" bond types")
	addNewLine(&content)

	spaceDim_Str := strconv.Itoa(int(global_data.GetGlobalData().SpaceDimention))
	addString(&content, "0 "+spaceDim_Str+" xlo xhi")
	addString(&content, "0 "+spaceDim_Str+" ylo yhi")
	addString(&content, "0 "+spaceDim_Str+" zlo zhi")
	addNewLine(&content)

	mapMonomerTypeNumber := make(map[dt.MonomerType]int)
	for i, monomerType := range monomerTypes {
		mapMonomerTypeNumber[monomerType] = i + 1
	}
	addString(&content, "Masses")
	addNewLine(&content)

	for monType, number := range mapMonomerTypeNumber {
		if monType == dt.MONOMER_TYPE_UNDEFINED {
			continue
		}
		literal := globula.GetLiteral(monType)
		addString(&content, strconv.Itoa(number)+" 1 # "+literal)
	}

	addNewLine(&content)

	addString(&content, "Bond Coeffs # harmonic")
	addNewLine(&content)

	for i, bondType := range bondTypes {
		addString(&content, strconv.Itoa(i+1)+" 100 "+strconv.FormatFloat(bondTypeMass(bondType), 'f', 3, 64))
	}
	addNewLine(&content)

	addString(&content, "Atoms # full")
	addNewLine(&content)

	maxNumber := -1
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		views.ForEachMonomer(pol, func(mon *dt.Monomer) {
			if mon.MonomerType == dt.MONOMER_TYPE_UNDEFINED {
				return
			}
			monCoords := mon.Coords()
			maxNumber = *base.Max_int([]int{maxNumber, int(mon.Number)})
			addString(&content, strconv.Itoa(int(mon.Number))+
				" 1 "+
				strconv.Itoa(mapMonomerTypeNumber[mon.MonomerType])+
				" 0.00000 "+
				strconv.FormatInt(monCoords.X, 10)+" "+
				strconv.FormatInt(monCoords.Y, 10)+" "+
				strconv.FormatInt(monCoords.Z, 10)+" "+
				"0 0 0")
		})
	})
	if globula.Is(views.GLOBULA_WATERIZED) {
		views.ForEachPolymer_If(globula, func(pv *views.PolymerView) bool {
			field := pv.GetUnderlinedField()
			var globalData *global_data.GlobalData = global_data.GetGlobalData()
			shape := [...]int64{globalData.SpaceDimention, globalData.SpaceDimention, globalData.SpaceDimention}
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
		views.ForEachMonomer(pol, func(mon *dt.Monomer) {
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
			views.ForEachMonomer(pol, func(mon *dt.Monomer) {
				if mon.MonomerType != dt.MONOMER_TYPE_UNDEFINED {
					atomsCount++
				}
			})
		})
	} else {
		globalData := global_data.GetGlobalData()
		atomsCount = int(globalData.SpaceDimention) * int(globalData.SpaceDimention) * int(globalData.SpaceDimention)
	}
	return atomsCount
}

func getMonomerTypes(globula *views.GlobulaView) []dt.MonomerType {
	monomersTypes_map := make(map[dt.MonomerType]bool)
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		views.ForEachMonomer(pol, func(mon *datatypes.Monomer) {
			monomersTypes_map[mon.MonomerType] = true
		})
	})

	monomersTypes := make([]dt.MonomerType, 0)
	for key := range monomersTypes_map {
		monomersTypes = append(monomersTypes, key)
	}
	if globula.Is(views.GLOBULA_WATERIZED) {
		monomersTypes = append(monomersTypes, dt.MONOMER_TYPE_WATER)
	}
	slices.Sort(monomersTypes)
	return monomersTypes
}

func getBondTypes(globula *views.GlobulaView) ([]dt.ConnectionType, int) {
	connTypes_map := make(map[dt.ConnectionType]bool)
	allSides := dt.GetAllSides()
	bondsCount := 0
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		views.ForEachMonomer(pol, func(mon *dt.Monomer) {
			for _, side := range allSides {
				connType := mon.GetTypeOfConnectionWithSide(side)
				sibling, err := mon.GetSibling(side)
				if err == nil && connType != dt.CONNECTION_TYPE_UNDEFINED && sibling.Number != -1 {
					connTypes_map[connType] = true
					bondsCount += 1
				}
			}
		})
	})
	types := make([]dt.ConnectionType, 0)
	for key := range connTypes_map {
		types = append(types, key)
	}
	slices.Sort(types)
	return types, int(bondsCount / 2)
}
