package dfs

import (
	"polymers/base"
	dt "polymers/datatypes"
	"polymers/global_data"
	"polymers/output_format"
	"polymers/views"
	"slices"
)

func DoDFS() *views.GlobulaView {
	field := dt.NewField(0)
	polymers := make([]*dt.Polymer, 0)
	for i := 0; i < 0; i++ {
		polymers[i] = dt.NewPolymer(field, int64(i))
	}

	globalData := global_data.GetGlobalData()

	printInf := output_format.GetPrint()

	printInf.Printfln("Enter the outcoming axis: %d - X, %d - Y, %d - Z", dt.X_AXIS, dt.Y_AXIS, dt.Z_AXIS)
	//var outcomingAxis dt.Axis = dt.X_AXIS
	//printInf.Readln(&outcomingAxis)

	printInf.Printfln("Enter the direction: %d - X, %d - Y, %d - Z", dt.X_AXIS, dt.Y_AXIS, dt.Z_AXIS)
	var directionAxis dt.Axis = dt.Y_AXIS
	//printInf.Readln(&directionAxis)

	direction := dt.DIRECTION_FORWARD
	outcomingSide := dt.GetSide(directionAxis, direction)
	polymerNumber := int64(1)
	for x := int64(0); x < globalData.SpaceDimention.X; x++ {
		for z := int64(0); z < globalData.SpaceDimention.Z; z++ {
			if cycles := makeForwardingCycles(base.Vector3D{X: x, Y: 0, Z: z}, outcomingSide, field, &polymerNumber); cycles != nil {
				minLengthPolymer := slices.MinFunc(*cycles, func(c1, c2 *dt.Polymer) int {
					if c1.Len() < c2.Len() {
						return -1
					} else if c1.Len() == c2.Len() {
						return 0
					} else {
						return 1
					}
				})

				polymers = append(polymers, slices.DeleteFunc(*cycles, func(c *dt.Polymer) bool { return c.Len() > minLengthPolymer.Len() })...)
			}
		}
	}

	globula := views.NewGlobulaView("DFS", polymers, views.GLOBULA_GLOBULA_TYPE)
	m := make(map[dt.MonomerType]string)
	m[dt.MONOMER_TYPE_USUAL] = "O"
	globula.SetLiterals(m)
	return globula
}

func makeForwardingCycles(point base.Vector3D, outcomingSide dt.Side, field *dt.Field, polymerNumber *int64) *[]*dt.Polymer {
	startMonomer := field.GetMonomerByCoords(point)
	if startMonomer == nil {
		return nil
	}
	movingSides := dt.GetMovementSides()
	reversedSide := dt.GetReversedSide(outcomingSide)
	movingSides = slices.DeleteFunc(movingSides, func(side dt.Side) bool { return side == reversedSide })
	//movingSides = []dt.Side{dt.SIDE_Forward, dt.SIDE_Left, dt.SIDE_Right}
	var stack base.Stack
	stack.Push(startMonomer)
	visitedMons := make(map[base.Vector3D]bool)
	currPath := make([]base.Vector3D, 0)
	(*polymerNumber)++
	pathsCoords := make([][]base.Vector3D, 0)
	for !stack.IsEmpty() {
		inf, ok := stack.Peek()
		if !ok {
			return nil
		}
		currMon := inf.(*dt.Monomer)
		if visited, ok := visitedMons[currMon.Coords()]; !ok || !visited {
			visitedMons[currMon.Coords()] = true
		} else {
			if visited {
				currCoords := currMon.Coords()
				if base.VectorsAreEqual(&currPath[len(currPath)-1], &currCoords) {
					currPath = slices.Delete(currPath, len(currPath)-1, len(currPath))
					visitedMons[currMon.Coords()] = false
				}
				stack.Pop()
				continue
			}
		}
		currPath = append(currPath, currMon.Coords())
		if isEdge(currMon, outcomingSide) {
			newPath := slices.Clone(currPath)
			pathsCoords = append(pathsCoords, newPath)
			continue
		}
		for _, movingSide := range movingSides {
			if nextMon, err := currMon.GetSibling(movingSide); err == nil {
				stack.Push(nextMon)
			}
		}
	}
	polymers := new([]*dt.Polymer)
	for _, path := range pathsCoords {
		newPolymer := dt.NewPolymer(field, *polymerNumber)
		for _, point := range path {
			newPolymer.AddMonomer(field.GetMonomerByCoords(point))
		}
		(*polymerNumber)++
		*polymers = append(*polymers, newPolymer)
	}
	return polymers
}

func isEdge(mon *dt.Monomer, outComingSide dt.Side) bool {
	spaceDimention := &global_data.GetGlobalData().SpaceDimention
	edgeX := spaceDimention.X - 1
	edgeY := spaceDimention.Y - 1
	edgeZ := spaceDimention.Z - 1
	point := mon.Coords()
	return (outComingSide == dt.SIDE_Forward && point.X == edgeX) ||
		(outComingSide == dt.SIDE_Backward && point.X == 0) ||
		(outComingSide == dt.SIDE_Right && point.Y == 0) ||
		(outComingSide == dt.SIDE_Left && point.Y == edgeY) ||
		(outComingSide == dt.SIDE_Up && point.Z == edgeZ) ||
		(outComingSide == dt.SIDE_Down && point.Z == 0)
}
