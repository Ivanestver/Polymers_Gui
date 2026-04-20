package cycles

import "polymers/views"

type Graph struct {
	graph   [][]int
	enabled []bool
}

func NewGraph(globula *views.GlobulaView) *Graph {
	nodesCount := globula.GetAtomsCount()
	graph := &Graph{
		graph:   make([][]int, nodesCount),
		enabled: make([]bool, nodesCount),
	}
	for i := 0; i < nodesCount; i++ {
		graph.graph[i] = make([]int, nodesCount)
		graph.enabled[i] = true
	}

	for polymerNumber := 0; polymerNumber < globula.Len(); polymerNumber++ {
		polymer := globula.GetPolymerByID(polymerNumber)
		if polymer == nil {
			continue
		}
		for monomerNumber := 0; monomerNumber < polymer.Len(); monomerNumber++ {
			monomer := polymer.GetMonomerByID(monomerNumber)
			if monomer != nil {
				siblings := monomer.GetSiblings()
				for _, sibling := range siblings {
					graph.graph[monomer.Number-1][sibling.Number-1] = 1
				}
			}
		}
	}

	return graph
}

func (graph *Graph) GetNodesCount() int {
	return len(graph.graph)
}

func (graph *Graph) GetAvailableNodesCount() int {
	enabledCount := 0
	for _, enabled := range graph.enabled {
		if enabled {
			enabledCount++
		}
	}
	return enabledCount
}

func (graph *Graph) AreConnected(i, j int) bool {
	return graph.IsAvailable(i) &&
		graph.IsAvailable(j) &&
		graph.graph[i][j] == 1
}

func (graph *Graph) GetConnectedOf(i int) []int {
	connected := make([]int, 0)
	if graph.IsAvailable(i) {
		for j := 0; j < len(graph.graph[i]); j++ {
			if graph.enabled[j] && graph.graph[i][j] == 1 {
				connected = append(connected, j)
			}
		}
	}
	return connected
}

func (graph *Graph) Disable(i int) {
	if graph.inBounds(i) {
		graph.enabled[i] = false
	}
}

func (graph *Graph) outOfBounds(i int) bool {
	return i < 0 || i >= graph.GetNodesCount()
}

func (graph *Graph) inBounds(i int) bool {
	return !graph.outOfBounds(i)
}

func (graph *Graph) IsAvailable(i int) bool {
	return graph.inBounds(i) && graph.enabled[i]
}
