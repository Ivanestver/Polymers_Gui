package cycles

import (
	"math"
	"polymers/base"
	"polymers/datatypes"
	"polymers/global_data"
	"polymers/output_format"
)

type _DFSAlg struct {
	sidesToMove []datatypes.Side
	usedPoints  map[int64]struct{}
	printer     output_format.IPrint
	axisAlong   base.Axis
}

func makeDFSAlg(axisAlong base.Axis) _DFSAlg {
	return _DFSAlg{
		sidesToMove: datatypes.GetNormalByAxis(axisAlong),
		usedPoints:  make(map[int64]struct{}),
		printer:     output_format.GetPrint(),
		axisAlong:   axisAlong,
	}
}

func (dfs *_DFSAlg) FindCycles(startPoint *datatypes.Monomer) [][]*datatypes.Monomer {
	// Let's start with the X axis for test and debug purposes
	printer := output_format.GetPrint()
	stack := base.Stack{}
	stack.Push(startPoint)
	paths := make([][]*datatypes.Monomer, 0)
	currPath := make([]*datatypes.Monomer, 0)
	for !stack.IsEmpty() {
		obj, ok := stack.Peek()
		if !ok {
			printer.PrintlnError("something wrong with the dpt")
		}
		currMonomer := obj.(*datatypes.Monomer)
		if dfs.isNewPoint(currMonomer) {
			// for len(currPath) > 0 && !areSiblisgs(currMonomer, currPath[len(currPath)-1]) {
			// 	last := currPath[len(currPath)-1]
			// 	delete(dfs.usedPoints, last.Number)
			// 	currPath = currPath[:len(currPath)-1]
			// }
			dfs.usedPoints[currMonomer.Number] = struct{}{}
			currPath = append(currPath, currMonomer)
			if dfs.isEdge(currMonomer) {
				addNewPath(&paths, currPath)
			} else {
				dfs.updateStack(currMonomer, &stack)
			}
		} else {
			if _, ok := stack.Pop(); !ok {
				printer.PrintlnError("when push popping from the stack")
			} else {
				last := currPath[len(currPath)-1]
				delete(dfs.usedPoints, last.Number)
				currPath = currPath[:len(currPath)-1]
			}
		}
	}
	return paths
}

func (dfs *_DFSAlg) isNewPoint(currMonomer *datatypes.Monomer) bool {
	if _, ok := dfs.usedPoints[currMonomer.Number]; ok || currMonomer.IsTypeOf(datatypes.MonomerTypeUndefined) {
		return false
	} else {
		return true
	}
}

func (dfs *_DFSAlg) updateStack(currMonomer *datatypes.Monomer, stack *base.Stack) {
	for _, sideToMove := range dfs.sidesToMove {
		sibling, err := currMonomer.GetSibling(sideToMove)
		if err != nil {
			dfs.printer.PrintlnError(err.Error())
			continue
		}
		if dfs.isNewPoint(sibling) {
			stack.Push(sibling)
		} else {
			continue
		}
	}
}

func (dfs *_DFSAlg) isEdge(currMonomer *datatypes.Monomer) bool {
	globalData := global_data.GetGlobalData()
	edgeValue := globalData.SpaceDimention[dfs.axisAlong].Higher
	return base.CompareFloat(edgeValue, currMonomer.Coords().X) ||
		base.CompareFloat(edgeValue, currMonomer.Coords().Y) ||
		base.CompareFloat(edgeValue, currMonomer.Coords().Z)
}

func addNewPath(paths *[][]*datatypes.Monomer, candidate []*datatypes.Monomer) {
	*paths = append(*paths, candidate)
}

func areSiblisgs(mon1, mon2 *datatypes.Monomer) bool {
	c1 := mon1.Coords()
	c2 := mon2.Coords()
	return math.Abs(c1.X-c2.X) == 1 && c1.Y == c2.Y && c1.Z == c2.Z ||
		c1.X == c2.X && math.Abs(c1.Y-c2.Y) == 1 && c1.Z == c2.Z ||
		c1.X == c2.X && c1.Y == c2.Y && math.Abs(c1.Z-c2.Z) == 1
}
