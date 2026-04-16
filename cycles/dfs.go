package cycles

import (
	"math"
	"polymers/base"
	"polymers/datatypes"
	"polymers/outputformat"
)

type _DFSAlg struct {
	sidesToMove []datatypes.Side
	usedPoints  map[int64]struct{}
	printer     outputformat.IPrint
	axisAlong   base.Axis
	leftBorder  float64
	rightBorder float64
}

func makeDFSAlg(axisAlong base.Axis, leftBorder, rightBorder float64) _DFSAlg {
	return _DFSAlg{
		sidesToMove: datatypes.GetNormalByAxis(axisAlong),
		usedPoints:  make(map[int64]struct{}),
		printer:     outputformat.GetPrint(),
		axisAlong:   axisAlong,
		leftBorder:  leftBorder,
		rightBorder: rightBorder,
	}
}

func (dfs *_DFSAlg) FindCycles(startPoint *datatypes.Monomer) [][]*datatypes.Monomer {
	// Let's start with the X axis for test and debug purposes
	printer := outputformat.GetPrint()
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
	siblings := currMonomer.GetSiblingsAlongAxis(dfs.axisAlong)
	for _, sibling := range siblings {
		if dfs.isNewPoint(sibling) {
			stack.Push(sibling)
		} else {
			continue
		}
	}
}

func (dfs *_DFSAlg) isEdge(currMonomer *datatypes.Monomer) bool {
	return len(currMonomer.GetSiblingsAlongAxisStrict(dfs.axisAlong)) == 0
}

func addNewPath(paths *[][]*datatypes.Monomer, candidate []*datatypes.Monomer) {
	*paths = append(*paths, candidate)
}

func areSiblisgs(mon1, mon2 *datatypes.Monomer) bool {
	c1 := mon1.Coords()
	c2 := mon2.Coords()
	return math.Abs(c1[base.AxisX]-c2[base.AxisX]) == 1 && c1[base.AxisY] == c2[base.AxisY] && c1[base.AxisZ] == c2[base.AxisZ] ||
		c1[base.AxisX] == c2[base.AxisX] && math.Abs(c1[base.AxisY]-c2[base.AxisY]) == 1 && c1[base.AxisZ] == c2[base.AxisZ] ||
		c1[base.AxisX] == c2[base.AxisX] && c1[base.AxisY] == c2[base.AxisY] && math.Abs(c1[base.AxisZ]-c2[base.AxisZ]) == 1
}
