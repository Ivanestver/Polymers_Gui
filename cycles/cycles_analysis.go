package cycles

import (
	"math"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/globaldata"
	"polymers/outputformat"
	"polymers/views"
	"slices"
)

type CyclesAnalyzer struct {
	graph      *Graph
	globula    *views.GlobulaView
	axises     []base.Axis
	printer    outputformat.IPrint
	nodes      map[int]base.Vector3DF
	stepsCount int
}

func NewCyclesAnalyzer(globula *views.GlobulaView, axises []base.Axis, stepsCount int) *CyclesAnalyzer {
	analyzer := &CyclesAnalyzer{}
	analyzer.graph = NewGraph(globula)
	analyzer.globula = globula
	analyzer.axises = axises
	analyzer.printer = outputformat.GetPrint()
	analyzer.nodes = make(map[int]base.Vector3DF)
	analyzer.stepsCount = stepsCount
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

func Analyze(globula *views.GlobulaView, axises []base.Axis, stepsCount int) {
	file, err := os.Create("cycles.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	oldPrinter := outputformat.GetPrint()
	outputformat.SetPrint(outputformat.NewFilePrint(file))
	analyzer := NewCyclesAnalyzer(globula, axises, stepsCount)
	analyzer.analyzeOrigin()
	analyzer.analyzePreprocessed()
	analyzer.analyzeSurfaceIntersections()
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
	clusters := analyzer.findClusters()
	analyzer.printer.Println("Across cutting planes X,Y,Z :")
	for _, axis := range analyzer.axises {
		paths := make([][]int, 0)
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
				queue := base.Queue{startNode}
				parents := make(map[int]*base.Set)
				lengths := make(map[int]int)
				for _, node := range analyzer.graph.GetAvailableNodes() {
					lengths[node] = math.MaxInt
				}
				lengths[startNode] = 0
				visited[startNode] = struct{}{}
				for !queue.IsEmpty() {
					item, _ := queue.Pop()
					node := item.(int)
					if slices.Contains(maxClusters, node) {
						continue
					}
					connectedNodes := analyzer.graph.GetConnectedOf(node)
					for _, connectedNode := range connectedNodes {
						newLength := lengths[node] + 1
						if _, ok := visited[connectedNode]; !ok {
							queue.Push(connectedNode)
						}
						if newLength < lengths[connectedNode] {
							lengths[connectedNode] = newLength
							s := &base.Set{}
							s.Insert(node)
							parents[connectedNode] = s
						} else if newLength == lengths[connectedNode] {
							s := parents[connectedNode]
							s.Insert(node)
						}
					}
					visited[node] = struct{}{}
				}

				for _, node := range maxClusters {
					if _, ok := parents[node]; !ok {
						continue
					}
					currNode := node
					path := []int{}
					stack := base.Stack{currNode}
					for !stack.IsEmpty() {
						i := stack.PeekNotSafe().(int)
						if len(path) > 0 && i == path[len(path)-1] {
							stack.Pop()
							path = path[:len(path)-1]
							continue
						}
						path = append(path, i)
						if i == startNode {
							temp := make([]int, len(path))
							copy(temp, path)
							paths = append(paths, temp)
							slices.Reverse(temp)
							if !slices.ContainsFunc(paths, func(p []int) bool {
								return slices.Equal(p, temp)
							}) {
								paths = append(paths, temp)
							}
						} else {
							parentsOfI := parents[i]
							for p := range *parentsOfI {
								stack.Push(p)
							}
						}
					}
				}
			}
		}
		analyzer.printer.Printfln("\t%s: %d", axis.ToString(), len(paths))
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
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	for _, axis := range analyzer.axises {
		start := spaceDimention[axis].Lower
		end := spaceDimention[axis].Higher
		surfaceIntersections := make(map[float64]int)
		step := (end - start) / float64(analyzer.stepsCount)
		start += step / 2
		currMonomers, _ := analyzer.getMinMaxOfCluster(analyzer.graph.GetAvailableNodes(), axis)
		for coordOnAxis := start; coordOnAxis < end; coordOnAxis += step {
			// Calculate the intersections
			surfaceIntersections[coordOnAxis] = analyzer.getIntersectionsCount(axis, coordOnAxis, &currMonomers)
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
			analyzer.printer.Printfln("\t%d. %f: %d", i+1, key, surfaceIntersections[key])
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

func (analyzer *CyclesAnalyzer) getIntersectionsCount(axis base.Axis, coordOnAxis float64, currNodes *[]int) int {
	intersectionsCount := 0
	for {
		nodesToRemove := make([]int, 0)
		for _, currNode := range *currNodes {
			connectedNodes := analyzer.graph.GetConnectedOf(currNode)
			currCoords := analyzer.nodes[currNode]
			if currCoords[axis] > coordOnAxis {
				continue
			}
			outOfBounds := 0
			for _, connectedNode := range connectedNodes {
				connectedCoords := analyzer.nodes[connectedNode]
				if connectedCoords[axis]-currCoords[axis] < 0.0 {
					continue
				}
				if connectedCoords[axis] < coordOnAxis {
					outOfBounds++
					continue
				}
				intersectionsCount++
			}
			if outOfBounds == len(connectedNodes) {
				nodesToRemove = append(nodesToRemove, currNode)
			}
		}
		if len(nodesToRemove) == 0 {
			break
		}
		nodesToAdd := make([]int, 0)
		*currNodes = slices.DeleteFunc(*currNodes, func(currNode int) bool {
			if slices.Contains(nodesToRemove, currNode) {
				connectedNodes := analyzer.graph.GetConnectedOf(currNode)
				nodeToRemoveCoords := analyzer.nodes[currNode]
				connectedNodes = slices.DeleteFunc(connectedNodes, func(connectedNode int) bool {
					connectedNodeCoords := analyzer.nodes[connectedNode]
					return connectedNodeCoords[axis]-nodeToRemoveCoords[axis] < 0.0
				})
				nodesToAdd = append(nodesToAdd, connectedNodes...)
				return true
			} else {
				return false
			}
		})
		*currNodes = append(*currNodes, nodesToAdd...)
	}
	return intersectionsCount
}
