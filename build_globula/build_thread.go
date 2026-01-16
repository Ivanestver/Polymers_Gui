package build_globula

import (
	"polymers/base"
	"polymers/datatypes"
	"polymers/global_data"
	"polymers/output_format"
	"polymers/views"
	"strconv"
	"strings"
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

func (alg BuildThreadAlgInputData) GetLiterals() map[datatypes.MonomerType]string {
	m := make(map[datatypes.MonomerType]string)
	m[datatypes.MONOMER_TYPE_USUAL] = "O"
	m[datatypes.MONOMER_TYPE_O_CONTAINING] = "N"
	m[datatypes.MONOMER_TYPE_VYNIL] = "C"
	m[datatypes.MONOMER_TYPE_CROSSLINKED] = "H"
	m[datatypes.MONOMER_TYPE_S] = "S"
	return m
}

type BuildThreadAlgInputDataBuilder struct {
}

func (builder BuildThreadAlgInputDataBuilder) CreateInputData(algType AlgType, predefinedParams []string, particleName string) (ICalcAlgInputData, error) {
	inputData := BuildThreadAlgInputData{}

	predefinedParamsCount := len(predefinedParams)
	if predefinedParamsCount < 3 {
		output_format.GetPrint().Print("Input the cell's sizes: ")
		output_format.GetPrint().Readln(&inputData.Cell.Lx, &inputData.Cell.Ly, &inputData.Cell.Lx)
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
		output_format.GetPrint().Print("Input the thread diameter: ")
		output_format.GetPrint().Readln(&inputData.ThreadRadius)
	} else {
		param := predefinedParams[3]
		if param[0] == '(' {
			split := strings.Split(param, ",")
			radius, err := strconv.ParseFloat(split[0][1:], 64)
			if err != nil {
				return inputData, err
			}
			inputData.ThreadRadius = radius
			inputData.ThreadRadius /= 2
			maxPolymersCount, err := strconv.Atoi(split[1][:len(split[1])])
			if err != nil {
				return inputData, err
			}
			inputData.MaxPolymersCount = maxPolymersCount
		} else {
			polymersCount, err := strconv.Atoi(param)
			if err != nil {
				return inputData, err
			}
			inputData.ThreadRadius = 57.0 // It's enough to place 10,000 threads
			inputData.MaxPolymersCount = polymersCount
		}
	}

	if predefinedParamsCount < 5 {
		output_format.GetPrint().Print("Input the thread length: ")
		output_format.GetPrint().Readln(&inputData.ThreadLength)
	} else {
		inputData.ThreadLength = 12
		threadLength, err := strconv.Atoi(predefinedParams[4])
		if err != nil {
			return inputData, err
		}
		inputData.ThreadLength = float64(threadLength)
	}

	inputData.particleName = particleName

	return inputData, nil
}

type BuildThreadAlg struct {
	inputData BuildThreadAlgInputData
}

func (alg *BuildThreadAlg) Calc() []*datatypes.Polymer {
	// define start monomers
	output_format.GetPrint().PrintlnInfo("Define new field")
	field := datatypes.NewField(uint64(alg.inputData.ThreadRadius))
	output_format.GetPrint().PrintlnInfo("Define start monomers")
	startPositions := alg.defineStartMonomers()
	if startPositions == nil {
		return nil
	}
	// create threads
	var polymers []*datatypes.Polymer = make([]*datatypes.Polymer, len(startPositions))
	// build the polymers
	for i, startPosition := range startPositions {
		output_format.GetPrint().PrintlnInfo("The start position is (" + strconv.FormatInt(startPosition.X, 10) + ", " + strconv.FormatInt(startPosition.Y, 10) + ", " + strconv.FormatInt(startPosition.Z, 10) + ")")
		polymers[i] = datatypes.NewPolymer(field, int64(i))
		polymer := polymers[i]
		// add a start monomer
		mon := field.GetMonomerByCoords(*startPosition)
		polymer.AddMonomer(mon)
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
