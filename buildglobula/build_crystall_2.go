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
	for i := range 2 {
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
	return []Plate{startingPlate}
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
