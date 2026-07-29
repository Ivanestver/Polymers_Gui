package buildglobula

import (
	"polymers/base"
	"polymers/datatypes"
	"polymers/globaldata"
	"polymers/outputformat"
	"polymers/views"
	"strconv"
	"strings"
)

type ThreadInputData struct {
	Cell struct {
		Lx int
		Ly int
		Lz int
	}
	ThreadRadius     float64
	ThreadLength     float64
	MaxPolymersCount int
}

func (alg ThreadInputData) GetGlobulaType() views.GlobulaProperty {
	return views.GlobulaThreadType
}

type ThreadInputDataBuilder struct {
}

func (builder ThreadInputDataBuilder) CreateInputData(predefinedParams []string, particleName string) (ICalcAlgInputData, error) {
	inputData := ThreadInputData{}

	predefinedParamsCount := len(predefinedParams)
	if predefinedParamsCount < 3 {
		outputformat.GetPrint().Print("Input the cell's sizes: ")
		outputformat.GetPrint().Readln(&inputData.Cell.Lx, &inputData.Cell.Ly, &inputData.Cell.Lx)
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
		outputformat.GetPrint().Print("Input the thread diameter: ")
		outputformat.GetPrint().Readln(&inputData.ThreadRadius)
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
			maxPolymersCount, err := strconv.Atoi(split[1][:len(split[1])-1])
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
		outputformat.GetPrint().Print("Input the thread length: ")
		outputformat.GetPrint().Readln(&inputData.ThreadLength)
	} else {
		inputData.ThreadLength = 12
		threadLength, err := strconv.Atoi(predefinedParams[4])
		if err != nil {
			return inputData, err
		}
		inputData.ThreadLength = float64(threadLength)
	}

	return inputData, nil
}

type ThreadCalcAlg AbstractAlg[ThreadInputData]

func (alg *ThreadCalcAlg) Calc() []*datatypes.Polymer {
	// define start monomers
	outputformat.GetPrint().PrintlnInfo("Define new field")
	field := datatypes.NewField(uint64(alg.inputData.ThreadRadius))
	outputformat.GetPrint().PrintlnInfo("Define start monomers")
	startPositions := alg.defineStartMonomers()
	if startPositions == nil {
		return nil
	}
	// create threads
	polymers := make([]*datatypes.Polymer, len(startPositions))
	// build the polymers
	for i, startPosition := range startPositions {
		outputformat.GetPrint().PrintflnInfo("The start position is (%v)", startPosition)
		polymers[i] = datatypes.NewPolymer(field, int64(i))
		polymer := polymers[i]
		// add a start monomer
		mon := field.GetMonomerByCoords(startPosition)
		polymer.AddMonomer(mon)
		// move forward until the distance between a current monomer and the start monomers are more than inputData.ThreadLength
		forwardVector := base.Vector3DF{0, 0, 1}
		currPosition := startPosition
		for {
			nextPosition := base.AddVecF(currPosition, forwardVector)
			if base.EcludianDistanceF(startPosition, nextPosition) >= alg.inputData.ThreadLength {
				break
			}
			mon = field.GetMonomerByCoords(nextPosition)
			polymer.AddMonomer(mon)
			currPosition = nextPosition
		}
	}
	return polymers
}

func (alg *ThreadCalcAlg) defineStartMonomers() []base.Vector3DF {
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	if spaceDimention[base.AxisZ].Higher-spaceDimention[base.AxisZ].Lower < alg.inputData.ThreadLength {
		return nil
	}
	center := spaceDimention.GetCenter()
	center[base.AxisZ] = 0.0
	startPositions := make([]base.Vector3DF, 0)
	toVisit := make([]base.Vector3DF, 0)
	toVisit = append(toVisit, center)
	for len(toVisit) != 0 && (alg.inputData.MaxPolymersCount <= 0 || len(startPositions) < alg.inputData.MaxPolymersCount) {
		currPoint := toVisit[0]
		toVisit = toVisit[1:]
		if base.ContainsIf(startPositions, currPoint, func(it base.Vector3DF, value base.Vector3DF) bool {
			return base.VectorsAreEqualF(it, value)
		}) ||
			base.EcludianDistanceF(center, currPoint) > alg.inputData.ThreadRadius ||
			!spaceDimention.PointInSpace(&currPoint) {
			continue
		} else {
			startPositions = append(startPositions, currPoint)
			directions := []base.Vector3DF{
				{1, 0, 0},
				{0, 1, 0},
				{-1, 0, 0},
				{0, -1, 0},
			}

			for _, direction := range directions {
				toVisit = append(toVisit, base.AddVecF(currPoint, direction))
			}
		}
	}
	return startPositions
}
