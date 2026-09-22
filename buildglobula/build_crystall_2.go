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
	return 5
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
	return 10
}

func (alg *Crystall2CalcAlg) addMonomer(point base.Vector3DF, polymer datatypes.IPolymer, monomerType base.MendeleevTableElement) {
	field := polymer.GetField()
	monomer := field.GetMonomerByCoords(point)
	monomer.MonomerType = monomerType
	polymer.AddMonomer(monomer)
}

func (alg *Crystall2CalcAlg) addCrystallMonomer(point base.Vector3DF, polymer datatypes.IPolymer) {
	alg.addMonomer(point, polymer, base.O)
}

func (alg *Crystall2CalcAlg) addAmorphousMonomer(point base.Vector3DF, polymer datatypes.IPolymer) {
	alg.addMonomer(point, polymer, base.N)
}

func getNextPoint(point, direction base.Vector3DF) base.Vector3DF {
	stepLength := 1.0
	return base.AddVecF(point, base.MultiplyByConstantF(direction, stepLength))
}

func moveForward(currPoint *base.Vector3DF, direction base.Vector3DF) {
	*currPoint = getNextPoint(*currPoint, direction)
}

func (alg *Crystall2CalcAlg) createStick(currPoint *base.Vector3DF, polymer datatypes.IPolymer) {
	direction := alg.getDirectionOfZ()
	for range alg.getPlatesCount() {
		alg.addCrystallMonomer(*currPoint, polymer)
		moveForward(currPoint, direction)
	}
}

func (alg *Crystall2CalcAlg) createAmorphousPart(currPoint *base.Vector3DF, monomersCount int, polymer datatypes.IPolymer, direction base.Vector3DF) {
	for range monomersCount - 1 {
		alg.addAmorphousMonomer(*currPoint, polymer)
		moveForward(currPoint, direction)
	}
	alg.addAmorphousMonomer(*currPoint, polymer)
}

func (alg *Crystall2CalcAlg) createAmorphousLoop(currPoint *base.Vector3DF, polymer datatypes.IPolymer, baseDirection base.Vector3DF) {
	// Create the first vertical part
	horCount, verCount := alg.getMonomersInAmorphousLoop()
	alg.createAmorphousPart(currPoint, verCount, polymer, alg.getDirectionOfZ())

	// Make the turn to the horizontal part
	direction := alg.getDirectionForLoops()
	angle := math.Asin(baseDirection[base.AxisY])
	if !alongX && !base.CompareFloat(angle, 0.0) {
		angle += math.Pi
	}
	direction = base.RotateVector(direction, angle, base.AxisZVec)
	moveForward(currPoint, direction)

	// Make the horizontal part
	alg.createAmorphousPart(currPoint, horCount, polymer, baseDirection)
	alongZ = !alongZ
	// Make the turn to the vertical part
	direction = alg.getDirectionForLoops()
	angle = math.Asin(baseDirection[base.AxisY])
	if !alongX && !base.CompareFloat(angle, 0.0) {
		angle += math.Pi
	}
	direction = base.RotateVector(direction, angle, base.AxisZVec)
	moveForward(currPoint, direction)

	// Make the second vertical part
	alg.createAmorphousPart(currPoint, verCount, polymer, alg.getDirectionOfZ())
	moveForward(currPoint, alg.getDirectionOfZ())
}

func (alg *Crystall2CalcAlg) getAmorphousLoopLength() float64 {
	prevAlongX := alongX
	prevAlongZ := alongZ
	alongX = true
	alongZ = true
	direction := alg.getDirectionForLoops()
	alongX = prevAlongX
	alongZ = prevAlongZ
	angle := base.GetAngleInRad(direction, base.AxisXVec)
	length := 1.0 * math.Cos(angle)
	length *= 2
	horCount, _ := alg.getMonomersInAmorphousLoop()
	length += float64(horCount - 1)
	return length
}

type AlongX = bool
type AlongZ = bool

var alongX AlongX = AlongX(true)
var alongZ AlongZ = AlongZ(true)

func (alg *Crystall2CalcAlg) getDirectionForLoops() base.Vector3DF {
	directionsForLoops := make(map[AlongX]map[AlongZ]base.Vector3DF)
	directionsForLoops[true] = make(map[AlongZ]base.Vector3DF)
	directionsForLoops[true][true] = base.AddVecF(base.AxisXVec, base.AxisZVec).Normalized()
	directionsForLoops[true][false] = base.AddVecF(base.AxisXVec, base.AxisZVecReversed).Normalized()
	directionsForLoops[false] = make(map[AlongZ]base.Vector3DF)
	directionsForLoops[false][true] = base.AddVecF(base.AxisXVecReversed, base.AxisZVec).Normalized()
	directionsForLoops[false][false] = base.AddVecF(base.AxisXVecReversed, base.AxisZVecReversed).Normalized()
	return directionsForLoops[alongX][alongZ]
}

func (alg *Crystall2CalcAlg) getDirectionOfZ() base.Vector3DF {
	directionsOfZ := make(map[AlongZ]base.Vector3DF)
	directionsOfZ[true] = base.AxisZVec
	directionsOfZ[false] = base.AxisZVecReversed
	return directionsOfZ[alongZ]
}

func (alg *Crystall2CalcAlg) getDirectionOfX() base.Vector3DF {
	directionsOfX := make(map[AlongX]base.Vector3DF)
	directionsOfX[true] = base.AxisXVec
	directionsOfX[false] = base.AxisXVecReversed
	return directionsOfX[alongX]
}

func (alg *Crystall2CalcAlg) createRow(currPoint *base.Vector3DF, polymer datatypes.IPolymer) {
	monomerNumberInRow := alg.getMonomersNumberInRow()
	for range monomerNumberInRow - 1 {
		// Create a stick
		alg.createStick(currPoint, polymer)

		// Create an amorphous part
		alg.createAmorphousLoop(currPoint, polymer, base.IdentityVectorF())
	}
	alg.createStick(currPoint, polymer)
}

func (alg *Crystall2CalcAlg) getLateralDirection() base.Vector3DF {
	return base.Vector3DF{
		math.Cos(1.1752),
		math.Sin(1.1752),
		0.0,
	}
}

func (alg *Crystall2CalcAlg) createV3() []datatypes.IPolymer {
	field := datatypes.NewRealField(
		[3][2]float64{
			{-100.0, 100.0},
			{-100.0, 100.0},
			{-100.0, 100.0},
		})
	polymer := datatypes.NewIPolymer(datatypes.FieldTypeReal, field, int64(0))
	currPoint := base.IdentityVectorF()
	monomerNumberInRow := alg.getMonomersNumberInRow()
	for range monomerNumberInRow - 1 {
		// Create a row
		alg.createRow(&currPoint, polymer)
		// Now move to the next row
		direction := alg.getLateralDirection()
		alg.createAmorphousLoop(&currPoint, polymer, direction)
		alongX = !alongX
	}
	alg.createRow(&currPoint, polymer)

	alg.hardenPolymer(polymer)
	return []datatypes.IPolymer{polymer}
}

type Cell struct {
	LeftLowerCloser   *datatypes.Monomer
	LeftUpperCloser   *datatypes.Monomer
	RightLowerCloser  *datatypes.Monomer
	RightUpperCloser  *datatypes.Monomer
	LeftLowerFurther  *datatypes.Monomer
	LeftUpperFurther  *datatypes.Monomer
	RightLowerFurther *datatypes.Monomer
	RightUpperFurther *datatypes.Monomer
}

func (cell *Cell) CreateConnections() {
	datatypes.MakeConnection(cell.LeftLowerCloser, cell.RightUpperCloser, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftUpperCloser, cell.RightLowerCloser, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftLowerCloser, cell.LeftUpperFurther, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftLowerFurther, cell.LeftUpperCloser, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftLowerFurther, cell.RightUpperFurther, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftUpperFurther, cell.RightLowerFurther, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.RightLowerCloser, cell.RightUpperFurther, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.RightUpperCloser, cell.RightLowerFurther, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftLowerCloser, cell.RightLowerFurther, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftLowerFurther, cell.RightLowerCloser, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftUpperCloser, cell.RightUpperFurther, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftUpperFurther, cell.RightUpperCloser, datatypes.ConnectionTypeCrossSurface)
	datatypes.MakeConnection(cell.LeftLowerCloser, cell.RightUpperFurther, datatypes.ConnectionTypeCrossSpacial)
	datatypes.MakeConnection(cell.LeftUpperCloser, cell.RightLowerFurther, datatypes.ConnectionTypeCrossSpacial)
	datatypes.MakeConnection(cell.LeftLowerFurther, cell.RightUpperCloser, datatypes.ConnectionTypeCrossSpacial)
	datatypes.MakeConnection(cell.LeftUpperFurther, cell.RightLowerCloser, datatypes.ConnectionTypeCrossSpacial)
}

func (cell *Cell) MakeCell(leftLowerCloser *datatypes.Monomer, field datatypes.IField, lateralDirection base.Vector3DF, stepLengthInStick, stepLengthInAmorphousPart float64) {
	cell.LeftLowerCloser = leftLowerCloser
	getMonomer := func(mon *datatypes.Monomer, direction base.Vector3DF, stepLength float64) *datatypes.Monomer {
		return field.GetMonomerByCoords(base.AddVecF(mon.Coords(), base.MultiplyByConstantF(direction, stepLength)))
	}
	cell.LeftUpperCloser = getMonomer(cell.LeftLowerCloser, base.AxisZVec, stepLengthInStick)
	cell.RightLowerCloser = getMonomer(cell.LeftLowerCloser, base.AxisXVec, stepLengthInAmorphousPart)
	cell.RightUpperCloser = getMonomer(cell.LeftUpperCloser, base.AxisXVec, stepLengthInAmorphousPart)
	cell.LeftLowerFurther = getMonomer(cell.LeftLowerCloser, lateralDirection, stepLengthInAmorphousPart)
	cell.LeftUpperFurther = getMonomer(cell.LeftUpperCloser, lateralDirection, stepLengthInAmorphousPart)
	cell.RightLowerFurther = getMonomer(cell.RightLowerCloser, lateralDirection, stepLengthInAmorphousPart)
	cell.RightUpperFurther = getMonomer(cell.RightUpperCloser, lateralDirection, stepLengthInAmorphousPart)
}

func (alg *Crystall2CalcAlg) hardenPolymer(polymer datatypes.IPolymer) {
	field := polymer.GetField()
	currMonomer := polymer.GetMonomerByIdx(0)
	alongX = true
	alongZ = true

	stepLengthInAmorphousPart := alg.getAmorphousLoopLength()
	stepLengthInStick := 1.0

	for range alg.getMonomersNumberInRow() - 1 { // The width
		for j := range alg.getMonomersNumberInRow() - 1 { // The row
			for k := range alg.getPlatesCount() - 1 { // The height
				cell := Cell{}
				cell.MakeCell(currMonomer, field, alg.getLateralDirection(), stepLengthInStick, stepLengthInAmorphousPart)
				cell.CreateConnections()
				if k != alg.getPlatesCount()-2 {
					currMonomer = field.GetMonomerByCoords(base.AddVecF(currMonomer.Coords(), alg.getDirectionOfZ()))
				}
			}
			alongZ = !alongZ
			if j != alg.getMonomersNumberInRow()-2 {
				currMonomer = field.GetMonomerByCoords(base.AddVecF(currMonomer.Coords(), base.MultiplyByConstantF(alg.getDirectionOfX(), stepLengthInAmorphousPart)))
			}
		}
		alongX = !alongX
		currMonomer = field.GetMonomerByCoords(base.AddVecF(currMonomer.Coords(), base.MultiplyByConstantF(alg.getLateralDirection(), stepLengthInAmorphousPart)))
	}
}
