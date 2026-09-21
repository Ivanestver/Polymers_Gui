package buildglobula

import (
	"math"
	"polymers/base"
	"polymers/datatypes"
	"polymers/views"
)

type Crystall2InputDataBuilder struct {
}

func (builder Crystall2InputDataBuilder) CreateInputData(defaultParams []string) (ICalcAlgInputData, error) {
	return Crystall2InputData{}, nil
}

type Crystall2InputData struct {
}

func (inputData Crystall2InputData) GetGlobulaType() views.GlobulaProperty {
	return views.GlobulaCrystal2Type
}

type Crystall2CalcAlg AbstractAlg[Crystall2InputData]

type Plate []base.Vector3DF

func (alg Crystall2CalcAlg) Calc() []datatypes.IPolymer {
	return alg.createV3()
}

func (alg *Crystall2CalcAlg) createV1() []datatypes.IPolymer {
	plates := alg.createPlates()
	polymer := alg.turnPlatesIntoPolymer(plates)
	return []datatypes.IPolymer{polymer}
}

func (alg *Crystall2CalcAlg) createPlates() []Plate {
	startingPlate := alg.createStartingPlate()
	return alg.multiplyStartingPlate(startingPlate)
}

func (alg *Crystall2CalcAlg) createStartingPlate() Plate {
	startPoint := base.IdentityVectorF()
	plate := Plate{startPoint}
	for i := range alg.getCirclesCount() {
		alg.createCircle(&plate, i+1)
	}
	return plate
}

func (alg *Crystall2CalcAlg) createCircle(plate *Plate, currCircleNumber int) {
	if currCircleNumber < 1 {
		return
	}

	newCircleLen := 6 * currCircleNumber
	prevCircleLen := max(6*(currCircleNumber-1), 1)
	currPoint := (*plate)[len(*plate)-prevCircleLen]
	var direction base.Vector3DF
	direction = base.Vector3DF{
		math.Cos(math.Pi / 3),
		math.Sin(math.Pi / 3),
	}
	stepLength := 1.0
	currPoint = base.AddVecF(currPoint, base.MultiplyByConstantF(direction, stepLength))
	*plate = append(*plate, currPoint)
	direction = base.RotateVector(direction, -2*math.Pi/3, base.AxisZVec)
	for i := 1; i < newCircleLen; i++ {
		currPoint = base.AddVecF(currPoint, base.MultiplyByConstantF(direction, stepLength))
		*plate = append(*plate, currPoint)
		if i%currCircleNumber == 0 {
			direction = base.RotateVector(direction, -math.Pi/3, base.AxisZVec)
		}
	}
}

func (alg *Crystall2CalcAlg) multiplyStartingPlate(startingPlate Plate) []Plate {
	plates := []Plate{startingPlate}
	direction := base.AxisZVec
	stepLength := 1.0
	for range 23 {
		prevPlate := plates[len(plates)-1]
		plate := make(Plate, len(startingPlate))
		for i, point := range prevPlate {
			plate[i] = base.AddVecF(point, base.MultiplyByConstantF(direction, stepLength))
		}
	}
	return plates
}

func (alg *Crystall2CalcAlg) turnPlatesIntoPolymer(plates []Plate) datatypes.IPolymer {
	field := datatypes.NewRealField(
		[3][2]float64{
			{-100.0, 100.0},
			{-100.0, 100.0},
			{-100.0, 100.0},
		})
	polymer := datatypes.NewRealPolymer(field, 0)
	for _, plate := range plates {
		for _, point := range plate {
			monomer := field.GetMonomerByCoords(point)
			polymer.AddMonomer(monomer)
		}
	}
	return polymer
}

func (alg *Crystall2CalcAlg) getCirclesCount() int {
	return 3
}

func (alg *Crystall2CalcAlg) createV2() []datatypes.IPolymer {
	field := datatypes.NewRealField(
		[3][2]float64{
			{-100.0, 100.0},
			{-100.0, 100.0},
			{-100.0, 100.0},
		})
	polymer := datatypes.NewRealPolymer(field, 0)
	currPoint := base.IdentityVectorF()
	direction := base.Vector3DF{
		math.Cos(math.Pi / 3),
		math.Sin(math.Pi / 3),
	}
	currPoint.AddF(base.MultiplyByConstantF(direction, float64(alg.getCirclesCount())))
	currPoint.AddF(base.AxisZVecReversed)
	direction = base.Vector3DF{
		math.Cos(math.Pi / 3),
		-math.Sin(math.Pi / 3),
	}
	for range alg.getCirclesCount() {
		for range alg.getCirclesCount() - 1 {
			for range alg.getPlatesCount() {
				currPoint.AddF(base.AxisZVec)
				polymer.AddMonomer(field.GetMonomerByCoords(currPoint))
			}
		}
	}
	return nil
}

func (alg *Crystall2CalcAlg) getPlatesCount() int {
	return 3
}

func (alg *Crystall2CalcAlg) getWholeNumberOfPoints() int {
	n := alg.getCirclesCount()
	return 1 + 3*n*(n-1)
}

func (alg *Crystall2CalcAlg) getMonomersInAmorphousLoop() (horCount, verCount int) {
	horCount = 1
	verCount = 1
	return
}

func (alg *Crystall2CalcAlg) getMonomersNumberInRow() int {
	return 3
}

func (alg *Crystall2CalcAlg) createV3() []datatypes.IPolymer {
	field := datatypes.NewRealField(
		[3][2]float64{
			{-100.0, 100.0},
			{-100.0, 100.0},
			{-100.0, 100.0},
		})
	polymer := datatypes.NewIPolymer(datatypes.FieldTypeReal, field, int64(0))
	addMonomer := func(point base.Vector3DF) { polymer.AddMonomer(field.GetMonomerByCoords(point)) }
	stepLength := 1.0
	getNextPoint := func(point, direction base.Vector3DF) base.Vector3DF {
		return base.AddVecF(point, base.MultiplyByConstantF(direction, stepLength))
	}
	currPoint := base.IdentityVectorF()
	monomerNumberInRow := alg.getMonomersNumberInRow()
	type AlongX bool
	type AlongZ bool
	alongX := AlongX(true)
	alongZ := AlongZ(true)
	directionsForLoops := make(map[AlongX]map[AlongZ]base.Vector3DF)
	directionsForLoops[true] = make(map[AlongZ]base.Vector3DF)
	directionsForLoops[true][true] = base.AddVecF(base.AxisXVec, base.AxisZVec)
	directionsForLoops[true][false] = base.AddVecF(base.AxisXVec, base.AxisZVecReversed)
	directionsForLoops[false] = make(map[AlongZ]base.Vector3DF)
	directionsForLoops[false][true] = base.AddVecF(base.AxisXVecReversed, base.AxisZVec)
	directionsForLoops[false][false] = base.AddVecF(base.AxisXVec, base.AxisZVecReversed)
	directionsOfZ := make(map[AlongZ]base.Vector3DF)
	directionsOfZ[true] = base.AxisZVec
	directionsOfZ[false] = base.AxisZVecReversed
	directionsOfX := make(map[AlongX]base.Vector3DF)
	directionsOfX[true] = base.AxisXVec
	directionsOfX[false] = base.AxisXVecReversed
	direction := base.AxisZVec
	for range monomerNumberInRow {
		// Create a row
		for range monomerNumberInRow {
			// Create a stick
			for range alg.getPlatesCount() {
				addMonomer(currPoint)
				currPoint = getNextPoint(currPoint, direction)
			}

			// Create an amorphous part
			// Create the first vertical part
			horCount, verCount := alg.getMonomersInAmorphousLoop()
			for range verCount - 1 {
				addMonomer(currPoint)
				currPoint = getNextPoint(currPoint, direction)
			}
			addMonomer(currPoint)
			// Make the turn to the horizontal part
			direction = directionsForLoops[alongX][alongZ]
			currPoint = getNextPoint(currPoint, direction)

			// Make the horizontal part
			for range horCount - 1 {
				addMonomer(currPoint)
				currPoint = getNextPoint(currPoint, direction)
			}
			addMonomer(currPoint)
			alongZ = !alongZ
			// Make the turn to the vertical part
			direction = directionsForLoops[alongX][alongZ]
			currPoint = getNextPoint(currPoint, direction)

			// Make the second vertical part
			direction = directionsOfZ[alongZ]
			for range verCount - 1 {
				addMonomer(currPoint)
				currPoint = getNextPoint(currPoint, direction)
			}
			addMonomer(currPoint)
		}
		// Now move to the next row
		direction = base.Vector3DF{
			math.Cos(1.1752),
			math.Sin(1.1752),
			0.0,
		}
		currPoint = getNextPoint(currPoint, direction)
		alongX = !alongX
	}
	return []datatypes.IPolymer{polymer}
}
