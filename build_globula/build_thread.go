package build_globula

import (
	"fmt"
	"polymers/base"
	"polymers/datatypes"
	"polymers/global_data"
	"strconv"
)

type BuildThreadAlgInputData struct {
	Cell struct {
		Lx int
		Ly int
		Lz int
	}
	ThreadRadius float64
	ThreadLength float64
}

func (alg BuildThreadAlgInputData) GetName() string {
	return "Thread"
}

type BuildThreadAlgInputDataBuilder struct {
}

func (builder BuildThreadAlgInputDataBuilder) CreateInputData(algType AlgType) (ICalcAlgInputData, error) {
	inputData := BuildThreadAlgInputData{}

	fmt.Print("Input the cell's sizes: ")
	//fmt.Scanln(&inputData.Cell.Lx, &inputData.Cell.Ly, &inputData.Cell.Lx)
	inputData.Cell.Lx = 1
	inputData.Cell.Ly = 1
	inputData.Cell.Lz = 1

	fmt.Print("Input the thread diameter: ")
	//fmt.Scanln(&inputData.ThreadRadius)
	inputData.ThreadRadius = 4.0
	inputData.ThreadRadius /= 2

	fmt.Print("Input the thread length: ")
	//fmt.Scanln(&inputData.ThreadLength)
	inputData.ThreadLength = 12

	return inputData, nil
}

type BuildThreadAlg struct {
	inputData BuildThreadAlgInputData
}

func (alg *BuildThreadAlg) Calc() []*datatypes.Polymer {
	// define start monomers
	fmt.Println("[INFO] Define new field")
	field := datatypes.NewField(uint64(alg.inputData.ThreadRadius))
	fmt.Println("[INFO] Define start monomers")
	startPositions := alg.defineStartMonomers()
	// create threads
	var polymers []*datatypes.Polymer = make([]*datatypes.Polymer, len(startPositions))
	// build the polymers
	for i, startPosition := range startPositions {
		fmt.Println("[INFO] The start position is (" + strconv.FormatInt(startPosition.X, 10) + ", " + strconv.FormatInt(startPosition.Y, 10) + ", " + strconv.FormatInt(startPosition.Z, 10) + ")")
		polymers[i] = datatypes.NewPolymer(field, int64(i))
		polymer := polymers[i]
		// add a start monomer
		mon := field.GetMonomerByCoords(*startPosition)
		polymer.AddMonomer(mon)
		field.MakeFilled(mon)
		// move forward until the distance between a current monomer and the start monomers are more than inputData.ThreadLength
		forwardVector := base.Vector3D{X: 0, Y: 0, Z: 1}
		currPosition := &base.Vector3D{X: startPosition.X, Y: startPosition.Y, Z: startPosition.Z}
		for {
			nextPosition := base.AddVec(currPosition, &forwardVector)
			if base.EcludianDistance(*startPosition, *nextPosition) > alg.inputData.ThreadLength {
				break
			}
			mon = field.GetMonomerByCoords(*nextPosition)
			polymer.AddMonomer(mon)
			field.MakeFilled(mon)
			currPosition = nextPosition
		}
	}
	return polymers
}

func (alg *BuildThreadAlg) defineStartMonomers() []*base.Vector3D {
	globalData := global_data.GetGlobalData()
	center := &base.Vector3D{
		X: globalData.SpaceDimention / 2,
		Y: globalData.SpaceDimention / 2,
		Z: globalData.SpaceDimention / 2,
	}
	center.Z = int64(float64(center.Z) - alg.inputData.ThreadLength)
	startPositions := make([]*base.Vector3D, 0)
	toVisit := make([]*base.Vector3D, 0)
	toVisit = append(toVisit, center)
	for len(toVisit) != 0 {
		currPoint := toVisit[0]
		toVisit = toVisit[1:]
		if base.Contains_if(startPositions, currPoint, func(it *base.Vector3D, value *base.Vector3D) bool { return base.VectorsAreEqual(it, value) }) || base.EcludianDistance(*center, *currPoint) > alg.inputData.ThreadRadius {
			continue
		} else {
			startPositions = append(startPositions, currPoint)
			directions := []base.Vector3D{
				{X: 1, Y: 0, Z: 0},
				{X: 0, Y: 1, Z: 0},
				{X: -1, Y: 0, Z: 0},
				{X: 0, Y: -1, Z: 0},
			}

			for _, direction := range directions {
				toVisit = append(toVisit, base.AddVec(currPoint, &direction))
			}
		}
	}
	return startPositions
}

func (alg *BuildThreadAlg) GetLiteralsTable() map[datatypes.MonomerType]string {
	m := make(map[datatypes.MonomerType]string)
	m[datatypes.MONOMER_TYPE_UNDEFINED] = ""
	m[datatypes.MONOMER_TYPE_USUAL] = "O"
	m[datatypes.MONOMER_TYPE_VYNIL] = "C"
	m[datatypes.MONOMER_TYPE_O_CONTAINING] = "N"
	m[datatypes.MONOMER_TYPE_FWISE] = "F"
	m[datatypes.MONOMER_TYPE_CLWISE] = "Cl"
	m[datatypes.MONOMER_TYPE_CROSSLINKED] = "H"
	return m
}
