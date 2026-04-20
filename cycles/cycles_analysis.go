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
	graph   *Graph
	globula *views.GlobulaView
	axises  []base.Axis
	printer outputformat.IPrint
	file    *os.File
	nodes   map[int]base.Vector3DF
}

func NewCyclesAnalyzer(globula *views.GlobulaView, axises []base.Axis) *CyclesAnalyzer {
	analyzer := &CyclesAnalyzer{}
	analyzer.graph = NewGraph(globula)
	analyzer.globula = globula
	analyzer.axises = axises
	analyzer.printer = outputformat.GetPrint()
	analyzer.nodes = make(map[int]base.Vector3DF)
	for polNum := 0; polNum < globula.Len(); polNum++ {
		pol := globula.GetPolymerByID(polNum)
		if pol == nil {
			continue
		}
		for monNumber := 0; monNumber < pol.Len(); monNumber++ {
			if mon := pol.GetMonomerByID(monNumber); mon != nil {
				analyzer.nodes[int(mon.Number)] = mon.Coords()
			}
		}
	}
	return analyzer
}

func Analyze(globula *views.GlobulaView, axises []base.Axis) {
	analyzer := NewCyclesAnalyzer(globula, axises)
	analyzer.file, _ = os.Create("cycles.log")
	defer analyzer.file.Close()
	analyzer.analyzeOrigin()
	analyzer.analyzePreprocessed()
	// analyzer.analyzeSurfaceIntersections()
	// analyzer.analyzePaths()
}

func (analyzer *CyclesAnalyzer) analyzeOrigin() {
	analyzer.analyzeParticlesCount()
	analyzer.analyzeNumberOfVertexes()
	analyzer.analyzeTrees()
	analyzer.analyzeClusters()
}

func (analyzer *CyclesAnalyzer) analyzeParticlesCount() {
	analyzer.printer.Printfln("Общее число вершин N: %d", analyzer.graph.GetAvailableNodesCount())
}

func (analyzer *CyclesAnalyzer) analyzeNumberOfVertexes() {
	numberOfVertexes := make(map[int]int)
	nodesCount := analyzer.graph.GetNodesCount()
	for i := 0; i < nodesCount; i++ {
		if !analyzer.graph.IsAvailable(i) {
			continue
		}
		connectedCount := 0
		for j := 0; j < nodesCount; j++ {
			if analyzer.graph.AreConnected(i, j) {
				connectedCount++
			}
		}
		numberOfVertexes[connectedCount] += 1
	}

	analyzer.printer.Println("Распределение по количеству связей")
	for i := 0; i < 6; i++ {
		if number, ok := numberOfVertexes[i]; ok {
			analyzer.printer.Printfln("\t%d %d", i, number)
		} else {
			analyzer.printer.Printfln("\t%d 0", i)
		}
	}
}

func (analyzer *CyclesAnalyzer) analyzeTrees() {
	nodesCount := analyzer.graph.GetNodesCount()
	nodesInTreesCount := 0
	for {
		nodesToDisable := make([]int, 0)
		for i := 0; i < nodesCount; i++ {
			connectedNodes := analyzer.graph.GetConnectedOf(i)
			if len(connectedNodes) == 1 {
				nodesToDisable = append(nodesToDisable, i)
			}
		}
		if len(nodesToDisable) == 0 {
			break
		}
		nodesInTreesCount += len(nodesToDisable)
		for _, nodeToDisable := range nodesToDisable {
			analyzer.graph.Disable(nodeToDisable)
		}
	}

	analyzer.printer.Printfln("Количество узлов в деревьях: %d", nodesInTreesCount)
}

type cluster []int

func (analyzer *CyclesAnalyzer) analyzeClusters() {
	// Find clusters
	clusters := analyzer.findClusters()
	analyzer.printer.Println("Кластерный анализ:")
	analyzer.printer.Printfln("\tКоличество кластеров: %d", len(clusters))
	nodesInClustersCount := 0
	for _, c := range clusters {
		nodesInClustersCount += len(c)
	}
	analyzer.printer.Printfln("\tКоличество узлов в кластерах: %d", nodesInClustersCount)

	// Remove non-periodic
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	eX := (spaceDimention[base.AxisX].Higher - spaceDimention[base.AxisX].Lower) * 0.02
	eY := (spaceDimention[base.AxisY].Higher - spaceDimention[base.AxisY].Lower) * 0.02
	eZ := (spaceDimention[base.AxisZ].Higher - spaceDimention[base.AxisZ].Lower) * 0.02
	for _, c := range clusters {
		hasMin := false
		hasMax := false
		for _, i := range c {
			coords := analyzer.nodes[i]
			hasMin = coords[base.AxisX] < spaceDimention[base.AxisX].Lower || base.CompareFloatWithE(coords[base.AxisX], spaceDimention[base.AxisX].Lower, eX) ||
				coords[base.AxisY] < spaceDimention[base.AxisY].Lower || base.CompareFloatWithE(coords[base.AxisY], spaceDimention[base.AxisY].Lower, eY) ||
				coords[base.AxisZ] < spaceDimention[base.AxisZ].Lower || base.CompareFloatWithE(coords[base.AxisZ], spaceDimention[base.AxisZ].Lower, eZ)

			hasMax = coords[base.AxisX] > spaceDimention[base.AxisX].Higher || base.CompareFloatWithE(coords[base.AxisX], spaceDimention[base.AxisX].Higher, eX) ||
				coords[base.AxisY] > spaceDimention[base.AxisY].Higher || base.CompareFloatWithE(coords[base.AxisY], spaceDimention[base.AxisY].Higher, eY) ||
				coords[base.AxisZ] > spaceDimention[base.AxisZ].Higher || base.CompareFloatWithE(coords[base.AxisZ], spaceDimention[base.AxisZ].Higher, eZ)
			if hasMin && hasMax {
				break
			}
		}
		if !hasMin || !hasMax {
			for _, i := range c {
				analyzer.graph.Disable(i)
			}
		} else if len(c) == 1 {
			analyzer.graph.Disable(c[0])
		}
	}
}

func (analyzer *CyclesAnalyzer) findClusters() []cluster {
	clusters := make([]cluster, 0)
	nodesCount := analyzer.graph.GetNodesCount()
	visitedPoints := make(map[int]struct{})
	for node := 0; node < nodesCount; node++ {
		if _, ok := visitedPoints[node]; ok {
			continue
		}
		if !analyzer.graph.IsAvailable(node) {
			continue
		}
		c := analyzer.findCluster(node, visitedPoints)
		if len(c) > 0 {
			clusters = append(clusters, c)
		}
	}
	return clusters
}

func (analyzer *CyclesAnalyzer) findCluster(startPoint int, visitedPoints map[int]struct{}) cluster {
	if !analyzer.graph.IsAvailable(startPoint) {
		return cluster{}
	}
	stack := base.Stack{}
	c := cluster{}
	stack.Push(startPoint)
	for !stack.IsEmpty() {
		item, _ := stack.Pop()
		point := item.(int)
		if _, ok := visitedPoints[point]; ok {
			continue
		}
		c = append(c, point)
		visitedPoints[point] = struct{}{}
		for _, connected := range analyzer.graph.GetConnectedOf(point) {
			if _, ok := visitedPoints[connected]; !ok {
				stack.Push(connected)
			}
		}
	}
	return c
}

func (analyzer *CyclesAnalyzer) analyzePreprocessed() {
	analyzer.analyzeParticlesCount()
	analyzer.analyzeNumberOfVertexes()
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
