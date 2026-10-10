package buildglobula

import (
	"errors"
	"math"
	"polymers/base"
	"polymers/datatypes"
	"polymers/globaldata"
	"polymers/outputformat"
	"polymers/views"
	"strconv"
)

var _ConnectionTypeRhombusShorter = datatypes.ConnectionType{Number: datatypes.ConnectionTypeCount.Number, Length: 1.56798}
var _ConnectionTypeRhombusLonger = datatypes.ConnectionType{Number: _ConnectionTypeRhombusShorter.Number + 1, Length: 2.35402}
var _ConnectionTypeRhombusEdge = datatypes.ConnectionType{Number: _ConnectionTypeRhombusLonger.Number + 1, Length: 1.41421}
var _ConnectionTypeRhombusSurface = datatypes.ConnectionType{Number: _ConnectionTypeRhombusEdge.Number + 1, Length: 1.73205}
var _ConnectionTypeAmorphous = datatypes.ConnectionType{Number: _ConnectionTypeRhombusSurface.Number + 1, Length: 1.41421}

type Crystall2InputDataBuilder struct {
}

func (builder Crystall2InputDataBuilder) CreateInputData(defaultParams []string) (ICalcAlgInputData, error) {
	if len(defaultParams) != 6 {
		return nil, errors.New("входное количество параметров должно быть 6")
	}
	inputData := Crystall2InputData{}
	if stepLength, err := strconv.ParseFloat(defaultParams[0], 64); err == nil {
		inputData.StepLength = stepLength
	} else {
		return nil, err
	}

	if monNumberInRow, err := strconv.Atoi(defaultParams[1]); err == nil {
		inputData.MonomersNumberInRow = monNumberInRow
	} else {
		return nil, err
	}

	if lamelaeCount, err := strconv.Atoi(defaultParams[2]); err == nil {
		inputData.LamelaeCount = lamelaeCount
	} else {
		return nil, err
	}

	if horCount, err := strconv.Atoi(defaultParams[3]); err == nil {
		inputData.HorCount = horCount
	} else {
		return nil, err
	}

	if verCount, err := strconv.Atoi(defaultParams[4]); err == nil {
		inputData.VerCount = verCount
	} else {
		return nil, err
	}

	switch defaultParams[5] {
	case "true":
		inputData.Harden = true
	case "false":
		inputData.Harden = false
	default:
		return nil, errors.New("параметр для придания жёсткости должен быть true или false")
	}

	return inputData, nil
}

type Crystall2InputData struct {
	StepLength          float64
	MonomersNumberInRow int
	LamelaeCount        int
	HorCount            int
	VerCount            int
	Harden              bool
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
	stepLength := alg.inputData.StepLength
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
	stepLength := alg.inputData.StepLength
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
			polymer.AddMonomerWithConnection(monomer, datatypes.ConnectionType{
				Number: datatypes.ConnectionTypeOne.Number,
				Length: base.EcludianDistanceF(monomer.Coords(), polymer.LastMonomer().Coords()),
			})
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
	return alg.inputData.LamelaeCount
}

func (alg *Crystall2CalcAlg) getWholeNumberOfPoints() int {
	n := alg.getCirclesCount()
	return 1 + 3*n*(n-1)
}

func (alg *Crystall2CalcAlg) getMonomersInAmorphousLoop() (horCount, verCount int) {
	horCount = alg.inputData.HorCount
	verCount = alg.inputData.VerCount
	return
}

func (alg *Crystall2CalcAlg) getMonomersNumberInRow() int {
	return alg.inputData.MonomersNumberInRow
}

func (alg *Crystall2CalcAlg) addMonomer(point base.Vector3DF, polymer datatypes.IPolymer, monomerType base.MendeleevTableElement) {
	field := polymer.GetField()
	monomer := field.GetMonomerByCoords(point)
	if monomer == nil {
		outputformat.GetPrint().PrintflnError("точка находится вне заданного пространства: %v", point)
		return
	}
	monomer.MonomerType = monomerType
	if polymer.Len() > 0 {
		polymer.AddMonomerWithConnection(monomer, datatypes.ConnectionType{
			Number: datatypes.ConnectionTypeOne.Number,
			Length: base.EcludianDistanceF(polymer.LastMonomer().Coords(), monomer.Coords()),
		})
	} else {
		polymer.AddMonomer(monomer)
	}
}

func (alg *Crystall2CalcAlg) addCrystallMonomer(point base.Vector3DF, polymer datatypes.IPolymer) {
	alg.addMonomer(point, polymer, base.O)
}

func (alg *Crystall2CalcAlg) addAmorphousMonomer(point base.Vector3DF, polymer datatypes.IPolymer) {
	alg.addMonomer(point, polymer, base.N)
}

func moveForward(currPoint *base.Vector3DF, direction base.Vector3DF) {
	*currPoint = base.AddVecF(*currPoint, direction)
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
	turnLength := alg.getTurnLength()
	angle := math.Asin(baseDirection[base.AxisY])
	if !alongX && !base.CompareFloat(angle, 0.0) {
		angle += math.Pi
	}
	direction = base.RotateVector(direction, angle, base.AxisZVec)
	moveForward(currPoint, base.MultiplyByConstantF(direction, turnLength))
	// nMon1 := polymer.Len() - 1
	// nMon2 := nMon1 + 1

	// Make the horizontal part
	alg.createAmorphousPart(currPoint, horCount, polymer, baseDirection)
	// datatypes.SetConnectionType(
	// 	polymer.GetMonomerByIdx(nMon1),
	// 	polymer.GetMonomerByIdx(nMon2),
	// 	_ConnectionTypeAmorphous,
	// )
	alongZ = !alongZ
	// Make the turn to the vertical part
	direction = alg.getDirectionForLoops()
	angle = math.Asin(baseDirection[base.AxisY])
	if !alongX && !base.CompareFloat(angle, 0.0) {
		angle += math.Pi
	}
	direction = base.RotateVector(direction, angle, base.AxisZVec)
	moveForward(currPoint, base.MultiplyByConstantF(direction, turnLength))
	// nMon1 = polymer.Len() - 1
	// nMon2 = nMon1 + 1

	// Make the second vertical part
	alg.createAmorphousPart(currPoint, verCount, polymer, alg.getDirectionOfZ())
	// datatypes.SetConnectionType(
	// 	polymer.GetMonomerByIdx(nMon1),
	// 	polymer.GetMonomerByIdx(nMon2),
	// 	_ConnectionTypeAmorphous,
	// )
	moveForward(currPoint, alg.getDirectionOfZ())
}

func (alg *Crystall2CalcAlg) getTurnLength() float64 {
	return alg.inputData.StepLength
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
	length := alg.getTurnLength() * math.Cos(angle)
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
	return base.MultiplyByConstantF(directionsOfZ[alongZ], alg.inputData.StepLength)
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
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	field := datatypes.NewRealField(
		[3][2]float64{
			{spaceDimention[base.AxisX].Lower, spaceDimention[base.AxisX].Higher},
			{spaceDimention[base.AxisY].Lower, spaceDimention[base.AxisY].Higher},
			{spaceDimention[base.AxisZ].Lower, spaceDimention[base.AxisZ].Higher},
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

	if alg.inputData.Harden {
		alg.hardenPolymer(polymer)
	}
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
	// Edges
	if err := cell.makeConnection(cell.LeftLowerCloser, cell.RightLowerCloser, _ConnectionTypeRhombusEdge); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.RightLowerCloser, cell.RightLowerFurther, _ConnectionTypeRhombusEdge); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.RightLowerFurther, cell.LeftLowerFurther, _ConnectionTypeRhombusEdge); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftLowerFurther, cell.LeftLowerCloser, _ConnectionTypeRhombusEdge); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftUpperCloser, cell.RightUpperCloser, _ConnectionTypeRhombusEdge); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.RightUpperCloser, cell.RightUpperFurther, _ConnectionTypeRhombusEdge); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.RightUpperFurther, cell.LeftUpperFurther, _ConnectionTypeRhombusEdge); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftUpperFurther, cell.LeftUpperCloser, _ConnectionTypeRhombusEdge); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}

	// Surfaces
	if err := cell.makeConnection(cell.LeftLowerCloser, cell.RightUpperCloser, _ConnectionTypeRhombusSurface); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftUpperCloser, cell.RightLowerCloser, _ConnectionTypeRhombusSurface); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.RightLowerCloser, cell.RightUpperFurther, _ConnectionTypeRhombusSurface); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.RightUpperCloser, cell.RightLowerFurther, _ConnectionTypeRhombusSurface); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.RightLowerFurther, cell.LeftUpperFurther, _ConnectionTypeRhombusSurface); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.RightUpperFurther, cell.LeftLowerFurther, _ConnectionTypeRhombusSurface); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftLowerFurther, cell.LeftUpperCloser, _ConnectionTypeRhombusSurface); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftUpperFurther, cell.LeftLowerCloser, _ConnectionTypeRhombusSurface); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftLowerCloser, cell.RightLowerFurther, _ConnectionTypeRhombusLonger); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftLowerFurther, cell.RightLowerCloser, _ConnectionTypeRhombusShorter); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftUpperCloser, cell.RightUpperFurther, _ConnectionTypeRhombusLonger); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
	if err := cell.makeConnection(cell.LeftUpperFurther, cell.RightUpperCloser, _ConnectionTypeRhombusShorter); err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return
	}
}

func (cell *Cell) makeConnection(mon1, mon2 *datatypes.Monomer, connectionType datatypes.ConnectionType) error {
	decimalsAfterDot := 1e8
	connType := datatypes.ConnectionType{
		Number: connectionType.Number,
		Length: float64(int64(base.EcludianDistanceF(mon1.Coords(), mon2.Coords())*decimalsAfterDot)) / decimalsAfterDot,
	}
	return datatypes.MakeConnection(mon1, mon2, connType)
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
	stepLengthInStick := alg.inputData.StepLength

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
