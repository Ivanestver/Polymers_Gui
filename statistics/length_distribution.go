package statistics

import (
	"polymers/views"
)

func CalculateLengthDistribution(globula *views.GlobulaView) map[int]int {
	distribution := make(map[int]int) // [length]count_of_such_chains
	for polNumber := range globula.Len() {
		polymer := globula.GetPolymerByIdx(polNumber)
		currChainLength := 0
		for monNumber := range polymer.Len() {
			currMon := polymer.GetMonomerByIdx(monNumber)
			currChainLength++
			if currMon.NextMonomer == nil {
				distribution[currChainLength]++
				currChainLength = 0
			}
		}
	}
	return distribution
}
