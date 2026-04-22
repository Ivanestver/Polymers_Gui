package cycles

import (
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
				analyzer.nodes[int(mon.Number-1)] = mon.Coords()
			}
		}
	}
	return analyzer
}

func Analyze(globula *views.GlobulaView, axises []base.Axis) {
	file, err := os.Create("cycles.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	oldPrinter := outputformat.GetPrint()
	outputformat.SetPrint(outputformat.NewFilePrint(file))
	analyzer := NewCyclesAnalyzer(globula, axises)
	analyzer.analyzeOrigin()
	analyzer.analyzePreprocessed()
	//analyzer.analyzeSurfaceIntersections()
	analyzer.analyzePaths()
	outputformat.SetPrint(oldPrinter)
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
	for _, c := range clusters {
		if len(c) == 1 {
			analyzer.graph.Disable(c[0])
		} else if !analyzer.isPassingCluster(c) {
			for _, i := range c {
				analyzer.graph.Disable(i)
			}
		}
	}
}

func (analyzer *CyclesAnalyzer) isPassingCluster(c cluster) bool {
	for _, axis := range []base.Axis{base.AxisX, base.AxisY, base.AxisZ} {
		minClusters, maxClusters := analyzer.getMinMaxOfCluster(c, axis)
		if len(minClusters) > 0 && len(maxClusters) > 0 {
			return true
		}
	}
	return false
}

func (analyzer *CyclesAnalyzer) getMinMaxOfCluster(c cluster, axis base.Axis) (minClusters, maxClusters []int) {
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	eX := (spaceDimention[axis].Higher - spaceDimention[axis].Lower) * 0.02
	minClusters = make([]int, 0)
	maxClusters = make([]int, 0)
	for _, i := range c {
		coords := analyzer.nodes[i]
		if spaceDimention[axis].Lower > coords[axis] ||
			base.CompareFloatWithE(coords[axis], spaceDimention[axis].Lower, eX) {
			minClusters = append(minClusters, i)
		}

		if coords[axis] > spaceDimention[axis].Higher ||
			base.CompareFloatWithE(coords[axis], spaceDimention[axis].Higher, eX) {
			maxClusters = append(maxClusters, i)
		}
	}
	return
}

func (analyzer *CyclesAnalyzer) findClusters() []cluster {
	clusters := make([]cluster, 0)
	nodes := analyzer.graph.GetAvailableNodes()
	visitedPoints := make(map[int]struct{})
	for _, node := range nodes {
		if _, ok := visitedPoints[node]; ok {
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
	analyzer.analyzeDensity()
	analyzer.analyzeSubchains()
}

func (analyzer *CyclesAnalyzer) analyzeDensity() {
	spaceDim := globaldata.GetGlobalData().SpaceDimention
	N := analyzer.graph.GetAvailableNodesCount()
	V := spaceDim.GetV()
	analyzer.printer.Printfln("Плотность сетки: %f", float64(N)/V)
}

type subchain []int

func (analyzer *CyclesAnalyzer) analyzeSubchains() {
	// The strategy is to find a subchain and replace it
	availableNodes := analyzer.graph.GetAvailableNodes()
	subchains := make([]subchain, 0)
	// First, find the subchains
	for _, node := range availableNodes {
		connectedNodes := analyzer.graph.GetConnectedOf(node)
		if len(connectedNodes) < 3 {
			continue
		}
		calculatedSubchains := analyzer.findSubchainsFrom(node)
		for _, s := range calculatedSubchains {
			if len(s) < 2 || slices.ContainsFunc(subchains, func(sub subchain) bool {
				if len(sub) != len(s) {
					return false
				}
				for _, item := range sub {
					if !slices.Contains(s, item) {
						return false
					}
				}
				return true
			}) {
				continue
			}
			subchains = append(subchains, s)
		}
	}

	analyzer.printer.Printfln("Число субцепей: %d", len(subchains))

	distribution := make(map[int]int)
	for _, subc := range subchains {
		distribution[len(subc)] += 1
	}
	analyzer.printer.Println("Распределение субцепей по длинам субцепи (длина-количество):")
	lengths := make([]int, len(distribution))
	for l := range distribution {
		lengths = append(lengths, l)
	}
	minLength := slices.Min(lengths)
	maxLength := slices.Max(lengths)
	for l := minLength; l <= maxLength; l++ {
		analyzer.printer.Printfln("\t%d: %d", l, distribution[l])
	}

	// Second, reduce them
	for _, s := range subchains {
		if len(s) < 2 {
			continue
		}
		for i := 0; i < len(s)-1; i++ {
			analyzer.graph.BreakConnection(s[i], s[i+1])
			if 0 < i && i < len(s)-1 {
				analyzer.graph.Disable(s[i])
			}
		}
		analyzer.graph.MakeConnection(s[0], s[len(s)-1])
	}
}

func (analyzer *CyclesAnalyzer) findSubchainsFrom(i int) []subchain {
	stack := base.Stack{}
	visited := make(map[int]struct{})
	visited[i] = struct{}{}
	subchains := make([]subchain, 0)
	s := subchain{i}
	connectedNodes := analyzer.graph.GetConnectedOf(i)
	for _, j := range connectedNodes {
		stack.Push(j)
	}
	for !stack.IsEmpty() {
		item, ok := stack.Pop()
		if !ok {
			continue
		}
		node := item.(int)
		if len(s) > 0 && s[len(s)-1] == node {
			s = s[:len(s)-1]
		}
		if _, ok := visited[node]; ok {
			continue
		}
		stack.Push(node)
		connectedNodes = analyzer.graph.GetConnectedOf(node)
		s = append(s, node)
		visited[node] = struct{}{}
		if len(connectedNodes) == 2 {
			for _, n := range connectedNodes {
				stack.Push(n)
			}
		} else {
			temp := make([]int, len(s))
			copy(temp, s)
			subchains = append(subchains, temp)
		}
	}
	return subchains
}

func (analyzer *CyclesAnalyzer) analyzePaths() {
	cyclesMap := make(map[base.Axis][][]int)
	for _, axis := range analyzer.axises {
		// startingPoints, leftBorder, rightBorder := analyzer.getStartingPointsForCycles(axis)
		// if startingPoints == nil && len(startingPoints) == 0 {
		// 	analyzer.printer.PrintlnError("no starting points have been defined")
		// 	continue
		// }
		// cycles := make([][]*datatypes.Monomer, 0)
		// dfs := makeDFSAlg(axis, leftBorder, rightBorder)
		// for _, startingPoint := range startingPoints {
		// 	cycles = append(cycles, dfs.FindCycles(startingPoint)...)
		// }
		// cyclesMap[axis] = cycles
		clusters := analyzer.findClusters()
		paths := make([][]int, 0)
		currPath := make([]int, 0)
		for _, c := range clusters {
			minClusters1, maxClusters1 := analyzer.getMinMaxOfCluster(c, axis)
			minClusters := slices.DeleteFunc(minClusters1, func(i int) bool {
				return slices.Contains(maxClusters1, i)
			})
			maxClusters := slices.DeleteFunc(maxClusters1, func(i int) bool {
				return slices.Contains(minClusters1, i)
			})
			for _, startNode := range minClusters {
				visited := make(map[int]struct{})
				stack := base.Stack{startNode}
				for !stack.IsEmpty() {
					item, _ := stack.Peek()
					node := item.(int)
					if _, ok := visited[node]; !ok {
						visited[node] = struct{}{}
						currPath = append(currPath, node)
						if slices.Contains(maxClusters, node) {
							p := make([]int, len(currPath))
							copy(p, currPath)
							paths = append(paths, p)
						} else {
							connectedNodes := analyzer.graph.GetConnectedOf(node)
							for _, connectedNode := range connectedNodes {
								if _, ok := visited[connectedNode]; !ok {
									stack.Push(connectedNode)
								}
							}
						}
					} else {
						stack.Pop()
						if node == currPath[len(currPath)-1] {
							last := currPath[len(currPath)-1]
							delete(visited, last)
							currPath = currPath[:len(currPath)-1]
						}
					}
				}
			}
		}
		cyclesMap[axis] = paths
	}

	// Print axises
	analyzer.printer.Println("Across cutting planes X,Y,Z :")
	for _, axis := range []base.Axis{base.AxisX, base.AxisY, base.AxisZ} {
		analyzer.printer.Printfln("\t%s: %d", axis.ToString(), len(cyclesMap[axis]))
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
	analyzer.printer.Printfln("LOAD BEARING BONDS CROSSING X,Y,Z CUTTING PLANES:")
	stepsCount := 20.0
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	for _, axis := range analyzer.axises {
		start := spaceDimention[axis].Lower
		end := spaceDimention[axis].Higher
		surfaceIntersections := make(map[float64]int)
		step := (end - start) / stepsCount
		start += step / 2
		for coordOnAxis := start; coordOnAxis < end; coordOnAxis += step {
			// Calculate the intersections
			surfaceIntersections[coordOnAxis] = analyzer.getIntersectionsCount(axis, coordOnAxis)
		}

		// Print the results
		keys := make([]float64, len(surfaceIntersections))
		i := 0
		for key := range surfaceIntersections {
			keys[i] = key
			i++
		}
		slices.Sort(keys)
		analyzer.printer.Printfln("For %s:", axis.ToString())
		for i, key := range keys {
			analyzer.printer.Printfln("\t%d: %d", i+1, surfaceIntersections[key])
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

func (analyzer *CyclesAnalyzer) getIntersectionsCount(axis base.Axis, coordOnAxis float64) int {
	intersectionsCount := 0
	nodesCount := analyzer.graph.GetNodesCount()
	for i := 0; i < nodesCount-1; i++ {
		if !analyzer.graph.IsAvailable(i) {
			continue
		}
		for j := i + 1; j < nodesCount; j++ {
			if analyzer.graph.IsAvailable(j) && analyzer.graph.AreConnected(i, j) {
				xi := analyzer.nodes[i][axis]
				xj := analyzer.nodes[j][axis]
				isIn := func(left, right float64) bool {
					return (base.CompareFloat(left, coordOnAxis) || left < coordOnAxis) &&
						(base.CompareFloat(coordOnAxis, right) || coordOnAxis < right)
				}
				if isIn(xi, xj) || isIn(xj, xi) {
					intersectionsCount++
				}
			}
		}
	}
	return intersectionsCount
}
