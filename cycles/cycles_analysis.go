package cycles

import (
	"polymers/base"
	"polymers/datatypes"
	"polymers/global_data"
	"polymers/output_format"
	"polymers/views"
	"slices"
)

type CyclesAnalyzer struct {
	globula *views.GlobulaView
	axises  []base.Axis
	printer output_format.IPrint
}

func Analyze(globula *views.GlobulaView, axises []base.Axis) {
	analyzer := CyclesAnalyzer{
		globula: globula,
		axises:  axises,
		printer: output_format.GetPrint(),
	}
	analyzer.analyzeCounts()
	analyzer.analyzeSurfaceIntersections()
}

func (analyzer *CyclesAnalyzer) analyzeCounts() {
	cyclesMap := make(map[base.Axis][][]*datatypes.Monomer)
	for _, axis := range analyzer.axises {
		dfs := makeDFSAlg(axis)
		startingPoints := analyzer.getStartingPointsForCycles(axis)
		if startingPoints == nil && len(startingPoints) == 0 {
			analyzer.printer.PrintlnError("no starting points have been defined")
			continue
		}
		cycles := make([][]*datatypes.Monomer, 0)
		for _, startingPoint := range startingPoints {
			cycles = append(cycles, dfs.FindCycles(startingPoint)...)
		}
		cyclesMap[axis] = cycles
	}

	// Print axises
	analyzer.printer.PrintlnInfo("The number of cycles for each side:")
	for _, axis := range []base.Axis{base.AxisX, base.AxisY, base.AxisZ} {
		analyzer.printer.PrintflnInfo("\t%s: %d", axis.ToString(), len(cyclesMap[axis]))
	}
}

func (analyzer *CyclesAnalyzer) getStartingPointsForCycles(axisAlong base.Axis) []*datatypes.Monomer {
	var field datatypes.IField
	views.ForEachPolymerIf(analyzer.globula, func(pv *views.PolymerView) bool {
		field = pv.GetUnderlinedField()
		return false
	})
	spaceDimention := global_data.GetGlobalData().SpaceDimention
	startingPoints := startingPointsDefiner(axisAlong, spaceDimention, field)
	finishingPoints := finishingPointsDefiner(axisAlong, spaceDimention, field)

	return slices.DeleteFunc(startingPoints, func(p *datatypes.Monomer) bool {
		return p.IsTypeOf(datatypes.MonomerTypeUndefined) || slices.ContainsFunc(finishingPoints, func(m *datatypes.Monomer) bool {
			return datatypes.MonomersAreEqual(p, m)
		})
	})
}

func startingPointsDefiner(axisAlong base.Axis, spaceDimention global_data.SpaceDimention, field datatypes.IField) []*datatypes.Monomer {
	switch axisAlong {
	case base.AxisX:
		return field.GetMonomersWithin(
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
	case base.AxisY:
		return field.GetMonomersWithin(
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
	case base.AxisZ:
		return field.GetMonomersWithin(
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
	default:
		return nil
	}
}

func finishingPointsDefiner(axisAlong base.Axis, spaceDimention global_data.SpaceDimention, field datatypes.IField) []*datatypes.Monomer {
	switch axisAlong {
	case base.AxisX:
		return field.GetMonomersWithin(
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
	case base.AxisY:
		return field.GetMonomersWithin(
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
	case base.AxisZ:
		return field.GetMonomersWithin(
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
	default:
		return nil
	}
}

func (analyzer *CyclesAnalyzer) analyzeSurfaceIntersections() {
	for _, axis := range analyzer.axises {
		startingMonomers := analyzer.getStartingPointsForCycles(axis)
		if len(startingMonomers) == 0 {
			continue
		}
		moveDirection := getMoveDirection(axis)
		if moveDirection == datatypes.SideUndefined {
			analyzer.printer.PrintflnError("could not define the move direction of the axis %s", axis.ToString())
			continue
		}
		surfaceIntersections := make(map[float64]int)
		start, end, step := getDimentions(axis, startingMonomers, moveDirection)
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
		analyzer.printer.Printfln("For %s", axis.ToString())
		for _, key := range keys {
			analyzer.printer.Printfln("%f.2: %d", key, surfaceIntersections[key])
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

func getDimentions(axisAlong base.Axis, startingPoints []*datatypes.Monomer, moveDirection datatypes.Side) (start, finiish, step float64) {
	spaceDimention := global_data.GetGlobalData().SpaceDimention
	start = spaceDimention[axisAlong].Lower
	finiish = spaceDimention[axisAlong].Higher
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
	}
	start += step / 2
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
