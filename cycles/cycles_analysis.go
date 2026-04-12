package cycles

import (
	"fmt"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/globaldata"
	"polymers/outputformat"
	"polymers/views"
	"slices"
)

type CyclesAnalyzer struct {
	globula *views.GlobulaView
	axises  []base.Axis
	printer outputformat.IPrint
	file    *os.File
}

func Analyze(globula *views.GlobulaView, axises []base.Axis) {
	analyzer := CyclesAnalyzer{
		globula: globula,
		axises:  axises,
		printer: outputformat.GetPrint(),
	}
	analyzer.file, _ = os.Create("cycles.log")
	analyzer.analyzeSurfaceIntersections()
	analyzer.analyzeCounts()
}

func (analyzer *CyclesAnalyzer) analyzeCounts() {
	cyclesMap := make(map[base.Axis][][]*datatypes.Monomer)
	for _, axis := range analyzer.axises {
		startingPoints, leftBorder, rightBorder := analyzer.getStartingPointsForCycles(axis)
		if startingPoints == nil && len(startingPoints) == 0 {
			analyzer.printer.PrintlnError("no starting points have been defined")
			continue
		}
		cycles := make([][]*datatypes.Monomer, 0)
		dfs := makeDFSAlg(axis, leftBorder, rightBorder)
		for _, startingPoint := range startingPoints {
			cycles = append(cycles, dfs.FindCycles(startingPoint)...)
		}
		cyclesMap[axis] = cycles
	}

	// Print axises
	analyzer.file.WriteString("Across cutting planes X,Y,Z :\n")
	for _, axis := range []base.Axis{base.AxisX, base.AxisY, base.AxisZ} {
		analyzer.file.WriteString(fmt.Sprintf("\t%s: %d\n", axis.ToString(), len(cyclesMap[axis])))
	}
}

func (analyzer *CyclesAnalyzer) getStartingPointsForCycles(axisAlong base.Axis) ([]*datatypes.Monomer, float64, float64) {
	var field datatypes.IField
	views.ForEachPolymerIf(analyzer.globula, func(pv *views.PolymerView) bool {
		field = pv.GetUnderlinedField()
		return false
	})
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	startingPoints, leftBorder := defineStartingPoints(axisAlong, spaceDimention, field)
	finishingPoints, rightBorder := defineFinishingPoints(axisAlong, spaceDimention, field)

	return slices.DeleteFunc(startingPoints, func(p *datatypes.Monomer) bool {
		return p.IsTypeOf(datatypes.MonomerTypeUndefined) || slices.ContainsFunc(finishingPoints, func(m *datatypes.Monomer) bool {
			return datatypes.MonomersAreEqual(p, m)
		})
	}), leftBorder, rightBorder
}

func defineStartingPoints(axisAlong base.Axis, spaceDimention globaldata.SpaceDimention, field datatypes.IField) ([]*datatypes.Monomer, float64) {
	var points []*datatypes.Monomer
	var moveDirection datatypes.Side
	switch axisAlong {
	case base.AxisX:
		points = field.GetMonomersWithin(
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Lower,
				Y: spaceDimention[base.AxisY].Lower,
				Z: spaceDimention[base.AxisZ].Lower,
			},
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Lower,
				Y: spaceDimention[base.AxisY].Higher,
				Z: spaceDimention[base.AxisZ].Higher,
			},
		)
		moveDirection = datatypes.SideForward
	case base.AxisY:
		points = field.GetMonomersWithin(
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Lower,
				Y: spaceDimention[base.AxisY].Lower,
				Z: spaceDimention[base.AxisZ].Lower,
			},
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Higher,
				Y: spaceDimention[base.AxisY].Lower,
				Z: spaceDimention[base.AxisZ].Higher,
			},
		)
		moveDirection = datatypes.SideLeft
	case base.AxisZ:
		points = field.GetMonomersWithin(
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Lower,
				Y: spaceDimention[base.AxisY].Lower,
				Z: spaceDimention[base.AxisZ].Lower,
			},
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Higher,
				Y: spaceDimention[base.AxisY].Higher,
				Z: spaceDimention[base.AxisZ].Lower,
			},
		)
		moveDirection = datatypes.SideUp
	default:
		return nil, 0.0
	}

	return moveSurfaceAlong(&points, moveDirection, axisAlong)
}

func defineFinishingPoints(axisAlong base.Axis, spaceDimention globaldata.SpaceDimention, field datatypes.IField) ([]*datatypes.Monomer, float64) {
	var points []*datatypes.Monomer
	var moveDirection datatypes.Side
	switch axisAlong {
	case base.AxisX:
		points = field.GetMonomersWithin(
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Higher,
				Y: spaceDimention[base.AxisY].Lower,
				Z: spaceDimention[base.AxisZ].Lower,
			},
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Higher,
				Y: spaceDimention[base.AxisY].Higher,
				Z: spaceDimention[base.AxisZ].Higher,
			},
		)
		moveDirection = datatypes.SideBackward
	case base.AxisY:
		points = field.GetMonomersWithin(
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Lower,
				Y: spaceDimention[base.AxisY].Higher,
				Z: spaceDimention[base.AxisZ].Lower,
			},
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Higher,
				Y: spaceDimention[base.AxisY].Higher,
				Z: spaceDimention[base.AxisZ].Higher,
			},
		)
		moveDirection = datatypes.SideRight
	case base.AxisZ:
		points = field.GetMonomersWithin(
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Lower,
				Y: spaceDimention[base.AxisY].Lower,
				Z: spaceDimention[base.AxisZ].Higher,
			},
			base.Vector3DF{
				X: spaceDimention[base.AxisX].Higher,
				Y: spaceDimention[base.AxisY].Higher,
				Z: spaceDimention[base.AxisZ].Higher,
			},
		)
		moveDirection = datatypes.SideDown
	default:
		return nil, 0.0
	}

	return moveSurfaceAlong(&points, moveDirection, axisAlong)
}

func moveSurfaceAlong(points *[]*datatypes.Monomer, moveDirection datatypes.Side, axisAlong base.Axis) ([]*datatypes.Monomer, float64) {
	var border *float64
	resultingPoints := make([]*datatypes.Monomer, 0)
	for len(*points) > 0 {
		for i := len(*points) - 1; i >= 0; i-- {
			point := (*points)[i]
			if point.IsNotTypeOf(datatypes.MonomerTypeUndefined) {
				resultingPoints = append(resultingPoints, point)
				*points = append((*points)[:i], (*points)[i+1:]...)
				if border == nil {
					border = new(float64)
					switch axisAlong {
					case base.AxisX:
						*border = point.Coords().X
					case base.AxisY:
						*border = point.Coords().Y
					default:
						*border = point.Coords().Z
					}
				}
			} else {
				nextPoint, err := point.GetSibling(moveDirection)
				if err == nil {
					(*points)[i] = nextPoint
				} else {
					*points = append((*points)[:i], (*points)[i+1:]...)
				}
			}
		}
	}
	return resultingPoints, *border
}

func (analyzer *CyclesAnalyzer) analyzeSurfaceIntersections() {
	analyzer.file.WriteString("LOAD BEARING BONDS CROSSING X,Y,Z CUTTING PLANES:\n")
	for _, axis := range analyzer.axises {
		startingMonomers, start, end := analyzer.getStartingPointsForCycles(axis)
		moveDirection := getMoveDirection(axis)
		if moveDirection == datatypes.SideUndefined {
			analyzer.printer.PrintflnError("could not define the move direction of the axis %s", axis.ToString())
			continue
		}
		surfaceIntersections := make(map[float64]int)
		step := getDimentions(axis, startingMonomers, moveDirection)
		start += step / 2
		for coordOnAxis := start; coordOnAxis < end; coordOnAxis += step {
			// Calculate the intersections
			surfaceIntersections[coordOnAxis] = getIntersectionsCount(axis, coordOnAxis, &startingMonomers, moveDirection)
			// Move points forwards
			for i := len(startingMonomers) - 1; i >= 0; i-- {
				nextMonomer, err := startingMonomers[i].GetSibling(moveDirection)
				if err != nil {
					startingMonomers = append(startingMonomers[:i], startingMonomers[i+1:]...)
					continue
				}
				startingMonomers[i] = nextMonomer
			}
		}

		// Print the results
		keys := make([]float64, len(surfaceIntersections))
		i := 0
		for key := range surfaceIntersections {
			keys[i] = key
			i++
		}
		slices.Sort(keys)
		analyzer.file.WriteString(fmt.Sprintf("For %s:\n", axis.ToString()))
		for i, key := range keys {
			analyzer.file.WriteString(fmt.Sprintf("\t%d: %d\n", i+1, surfaceIntersections[key]))
		}
	}
}

func getMoveDirection(axis base.Axis) datatypes.Side {
	switch axis {
	case base.AxisX:
		return datatypes.SideForward
	case base.AxisY:
		return datatypes.SideLeft
	case base.AxisZ:
		return datatypes.SideUp
	default:
		return datatypes.SideUndefined
	}
}

func getDimentions(axisAlong base.Axis, startingPoints []*datatypes.Monomer, moveDirection datatypes.Side) (step float64) {
	for _, startingPoint := range startingPoints {
		nextPoint, err := startingPoint.GetSibling(moveDirection)
		if err != nil {
			continue
		}
		switch axisAlong {
		case base.AxisX:
			step = nextPoint.Coords().X - startingPoint.Coords().X
		case base.AxisY:
			step = nextPoint.Coords().Y - startingPoint.Coords().Y
		case base.AxisZ:
			step = nextPoint.Coords().Z - startingPoint.Coords().Z
		}
		break
	}
	return
}

func getIntersectionsCount(axis base.Axis, coordOnAxis float64, prevPoints *[]*datatypes.Monomer, moveDirection datatypes.Side) int {
	intersectionsCount := 0
	for i := len(*prevPoints) - 1; i >= 0; i-- {
		prevPoint := (*prevPoints)[i]
		nextPoint, err := (*prevPoints)[i].GetSibling(moveDirection)
		if err != nil { // If it's the deadend, remove the point due to being reluctant
			*prevPoints = append((*prevPoints)[:i], (*prevPoints)[i+1:]...)
			continue
		}
		switch axis {
		case base.AxisX:
			if prevPoint.Coords().X < coordOnAxis && coordOnAxis < nextPoint.Coords().X {
				intersectionsCount++
			}
		case base.AxisY:
			if prevPoint.Coords().Y < coordOnAxis && coordOnAxis < nextPoint.Coords().Y {
				intersectionsCount++
			}
		case base.AxisZ:
			if prevPoint.Coords().Z < coordOnAxis && coordOnAxis < nextPoint.Coords().Z {
				intersectionsCount++
			}
		}
	}
	return intersectionsCount
}
