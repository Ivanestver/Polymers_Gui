package build_globula

import (
	"polymers/base"
	"polymers/datatypes"
	"polymers/global_data"
	"polymers/output_format"
	"polymers/views"
	"strconv"
)

type BuildThreadAlgInputData struct {
	Cell struct {
		Lx int
		Ly int
		Lz int
	}
	ThreadRadius     float64
	ThreadLength     float64
	MaxPolymersCount int
	particleName     string
}

func (alg BuildThreadAlgInputData) GetName() string {
	return alg.particleName
}

func (alg BuildThreadAlgInputData) GetGlobulaType() views.GlobulaProperty {
	return views.GLOBULA_THREAD_TYPE
}

type BuildThreadAlgInputDataBuilder struct {
}

func (builder BuildThreadAlgInputDataBuilder) CreateInputData(algType AlgType, predefinedParams []string, particleName string) (ICalcAlgInputData, error) {
	inputData := BuildThreadAlgInputData{}

	predefinedParamsCount := len(predefinedParams)
	if predefinedParamsCount < 3 {
		output_format.Print("Input the cell's sizes: ")
		output_format.Readln(&inputData.Cell.Lx, &inputData.Cell.Ly, &inputData.Cell.Lx)
	} else {
		x, err := strconv.Atoi(predefinedParams[0])
		if err != nil {
			return inputData, err
		}
		y, err := strconv.Atoi(predefinedParams[1])
		if err != nil {
			return inputData, err
		}
		z, err := strconv.Atoi(predefinedParams[2])
		if err != nil {
			return inputData, err
		}
		inputData.Cell.Lx = x
		inputData.Cell.Ly = y
		inputData.Cell.Lz = z
	}

	if predefinedParamsCount < 4 {
		output_format.Print("Input the thread diameter: ")
		output_format.Readln(&inputData.ThreadRadius)
	} else {
		radius, err := strconv.ParseFloat(predefinedParams[3], 64)
		if err != nil {
			return inputData, err
		}
		inputData.ThreadRadius = radius
		inputData.ThreadRadius /= 2
	}

	if predefinedParamsCount < 5 {
		output_format.Print("Input the thread length: ")
		output_format.Readln(&inputData.ThreadLength)
	} else {
		inputData.ThreadLength = 12
		threadLength, err := strconv.Atoi(predefinedParams[4])
		if err != nil {
			return inputData, err
		}
		inputData.ThreadLength = float64(threadLength)
	}

	if predefinedParamsCount < 6 {
		output_format.Print("Input the max polymers count: ")
		output_format.Readln(&inputData.MaxPolymersCount)
	} else {
		inputData.MaxPolymersCount = -1
		maxPolymersCount, err := strconv.Atoi(predefinedParams[5])
		if err != nil {
			return inputData, err
		}
		inputData.MaxPolymersCount = maxPolymersCount
	}

	inputData.particleName = particleName

	return inputData, nil
}

type BuildThreadAlg struct {
	inputData BuildThreadAlgInputData
}

func (alg *BuildThreadAlg) Calc() []*datatypes.Polymer {
	// define start monomers
	output_format.Println("[INFO] Define new field")
	field := datatypes.NewField(uint64(alg.inputData.ThreadRadius))
	output_format.Println("[INFO] Define start monomers")
	startPositions := alg.defineStartMonomers()
	if startPositions == nil {
		return nil
	}
	// create threads
	var polymers []*datatypes.Polymer = make([]*datatypes.Polymer, len(startPositions))
	// build the polymers
	for i, startPosition := range startPositions {
		output_format.Println("[INFO] The start position is (" + strconv.FormatInt(startPosition.X, 10) + ", " + strconv.FormatInt(startPosition.Y, 10) + ", " + strconv.FormatInt(startPosition.Z, 10) + ")")
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
			if base.EcludianDistance(*startPosition, *nextPosition) >= alg.inputData.ThreadLength {
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
	if globalData.SpaceDimention.Z < int64(alg.inputData.ThreadLength) {
		return nil
	}
	center := &base.Vector3D{
		X: globalData.SpaceDimention.X / 2,
		Y: globalData.SpaceDimention.Y / 2,
		Z: 0,
	}
	startPositions := make([]*base.Vector3D, 0)
	toVisit := make([]*base.Vector3D, 0)
	toVisit = append(toVisit, center)
	for len(toVisit) != 0 && (alg.inputData.MaxPolymersCount <= 0 || len(startPositions) < alg.inputData.MaxPolymersCount) {
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
