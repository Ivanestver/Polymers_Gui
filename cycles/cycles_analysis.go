package cycles

import (
	"fmt"
	"os"
	"polymers/base"
	"polymers/datatypes"
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
	defer analyzer.file.Close()
	analyzer.analyzeSurfaceIntersections()
	analyzer.analyzePaths()
}

func (analyzer *CyclesAnalyzer) analyzePaths() {
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
		fmt.Fprintf(analyzer.file, "\t%s: %d\n", axis.ToString(), len(cyclesMap[axis]))
	}
}

func (analyzer *CyclesAnalyzer) getStartingPointsForCycles(axisAlong base.Axis) ([]*datatypes.Monomer, float64, float64) {
	var field datatypes.IField
	views.ForEachPolymerIf(analyzer.globula, func(pv *views.PolymerView) bool {
		field = pv.GetUnderlinedField()
		return false
	})
	startingPoints, leftBorder := defineStartingPoints(axisAlong, field)
	finishingPoints, rightBorder := defineFinishingPoints(axisAlong, field)

	return slices.DeleteFunc(startingPoints, func(p *datatypes.Monomer) bool {
		return p.IsTypeOf(datatypes.MonomerTypeUndefined) || slices.ContainsFunc(finishingPoints, func(m *datatypes.Monomer) bool {
			return datatypes.MonomersAreEqual(p, m)
		})
	}), leftBorder, rightBorder
}

func defineStartingPoints(axisAlong base.Axis, field datatypes.IField) ([]*datatypes.Monomer, float64) {
	points := field.GetMinMonomersByAxis(axisAlong)
	leftBorder := 0.0
	if len(points) > 0 {
		leftBorder = points[0].Coords()[axisAlong]
	}

	return points, leftBorder
}

func defineFinishingPoints(axisAlong base.Axis, field datatypes.IField) ([]*datatypes.Monomer, float64) {
	points := field.GetMaxMonomersByAxis(axisAlong)
	leftBorder := 0.0
	if len(points) > 0 {
		leftBorder = points[0].Coords()[axisAlong]
	}

	return points, leftBorder
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
					*border = point.Coords()[axisAlong]
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
	if border == nil {
		return nil, 0.0
	} else {
		return resultingPoints, *border
	}
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
		step := (end - start) * 0.05
		start += step / 2
		for coordOnAxis := start; coordOnAxis < end; coordOnAxis += step {
			// Calculate the intersections
			surfaceIntersections[coordOnAxis] = analyzer.getIntersectionsCount(axis, coordOnAxis, &startingMonomers, moveDirection)
		}

		// Print the results
		keys := make([]float64, len(surfaceIntersections))
		i := 0
		for key := range surfaceIntersections {
			keys[i] = key
			i++
		}
		slices.Sort(keys)
		fmt.Fprintf(analyzer.file, "For %s:\n", axis.ToString())
		for i, key := range keys {
			fmt.Fprintf(analyzer.file, "\t%d: %d\n", i+1, surfaceIntersections[key])
		}
	}
}

func getMoveDirection(axis base.Axis) datatypes.Side {
	return [base.AxisCount]datatypes.Side{
		datatypes.SideForward,
		datatypes.SideLeft,
		datatypes.SideUp,
	}[axis]
}

func getDimentions(axisAlong base.Axis, startingPoints []*datatypes.Monomer, moveDirection datatypes.Side) (step float64) {
	for _, startingPoint := range startingPoints {
		nextPoint, err := startingPoint.GetSibling(moveDirection)
		if err != nil {
			continue
		}
		step = nextPoint.Coords()[axisAlong] - startingPoint.Coords()[axisAlong]
		break
	}
	return
}

func (analyzer *CyclesAnalyzer) getIntersectionsCount(axis base.Axis, coordOnAxis float64, prevPoints *[]*datatypes.Monomer, moveDirection datatypes.Side) int {
	intersectionsCount := 0
	views.ForEachPolymer(analyzer.globula, func(p *views.PolymerView) {
		views.ForEachMonomer(p, func(m *datatypes.Monomer) bool {
			siblings := m.GetSiblingsAlongAxis(axis)
			mCoord := m.Coords()[axis]
			if mCoord > coordOnAxis || base.CompareFloat(mCoord, coordOnAxis) {
				return true
			}
			for _, sibling := range siblings {
				siblingCoord := sibling.Coords()[axis]
				if coordOnAxis < siblingCoord {
					intersectionsCount++
				}
			}
			return true
		})
	})
	return intersectionsCount
}
