package visualizers

import (
	dt "data_visualizer/datatypes"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type sStack struct {
	s []int
}

func makeStack() sStack {
	return sStack{
		s: make([]int, 0),
	}
}

func (s *sStack) push(value int) {
	s.s = append(s.s, value)
}

func (s *sStack) top() int {
	return s.s[len(s.s)-1]
}

func (s *sStack) pop() int {
	value := s.top()
	s.s = s.s[:len(s.s)-1]
	return value
}

func (s *sStack) isEmpty() bool {
	return len(s.s) == 0
}

func (s *sStack) isNotEmpty() bool {
	return len(s.s) != 0
}

const STEPS_COUNT = 1000

type sAgeMassCountVisualizer struct {
	atoms *[]dt.Atom
	bonds *dt.Bonds
}

func (visualizer *sAgeMassCountVisualizer) Visualize(atoms []dt.Atom, bonds dt.Bonds, outputName string) {
	// cache the values to use them across the methods
	visualizer.atoms = &atoms
	visualizer.bonds = &bonds
	partsMasses := visualizer.dst()
	visualizer.showMasses(partsMasses, outputName)
}

func allIsTrue(usedAtoms *map[int]bool) bool {
	for _, b := range *usedAtoms {
		if !b {
			return false
		}
	}
	return true
}

func (visualizer *sAgeMassCountVisualizer) dst() []int {
	// Create the table of used atoms
	usedAtoms := make(map[int]bool)
	for i := 0; i < len(*visualizer.atoms); i++ {
		usedAtoms[i] = false
	}

	masses := make([]int, 0)
	// Do until all atoms have been visited
	for !allIsTrue(&usedAtoms) {
		// Firstly, we need to find a point to start from
		startPoint := 0
		for ; startPoint < len(*visualizer.atoms); startPoint++ {
			if !usedAtoms[startPoint] {
				break
			}
		}
		visitedPointsCount := 0
		stack := makeStack()
		stack.push(startPoint)
		for stack.isNotEmpty() {
			currentPoint := stack.pop()
			visitedPointsCount++
			usedAtoms[currentPoint] = true
			for _, nextPoint := range (*visualizer.bonds)[currentPoint] {
				if !usedAtoms[nextPoint] {
					stack.push(nextPoint)
				}
			}
		}
		masses = append(masses, visitedPointsCount)
	}
	return masses
}

func (visualizer *sAgeMassCountVisualizer) makeDistributions(masses []int) map[int]int {
	minMass := math.MaxInt
	maxMass := 0
	for _, mass := range masses {
		if mass < minMass {
			minMass = mass
		}
		if mass > maxMass {
			maxMass = mass
		}
	}
	step := int(math.Ceil(float64(len(*visualizer.atoms)) / STEPS_COUNT))
	distribution := make(map[int]int)
	if step == 0 {
		return distribution
	}
	for i := 0; i <= STEPS_COUNT; i++ {
		distribution[minMass+step*i] = 0
	}

	for _, mass := range masses {
		i := int(math.Ceil(float64(mass) / float64(step)))
		distribution[minMass+step*i] += 1
	}
	keys := make([]int, 0)
	for k := range distribution {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	finalDistribution := make(map[int]int)
	for _, key := range keys {
		if distribution[key] != 0 {
			finalDistribution[key] = distribution[key]
		}
	}

	return finalDistribution
}

func getXYValues(distribution *map[int]int) ([]string, []opts.BarData) {
	keys := make([]int, 0)
	for k := range *distribution {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	xValues := make([]string, 0)
	yValues := make([]opts.BarData, 0)
	for _, key := range keys {
		xValues = append(xValues, strconv.Itoa(key))
		yValues = append(yValues, opts.BarData{Value: (*distribution)[key]})
	}
	return xValues, yValues
}

func (visualizer *sAgeMassCountVisualizer) showMasses(masses []int, outputName string) {
	distribution := visualizer.makeDistributions(masses)
	bar := charts.NewBar()
	xValues, yValues := getXYValues(&distribution)
	bar.SetXAxis(xValues).AddSeries("Placeholder", yValues)
	f, _ := os.Create(outputName + ".html")
	bar.Render(f)
}
