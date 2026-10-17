package statistics

import (
	"fmt"
	"polymers/views"
	"slices"
	"strings"
)

func ShowLengthDistribution(globula *views.GlobulaView) string {
	builder := strings.Builder{}
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
	lengths := make([]int, len(distribution))
	i := 0
	for length := range distribution {
		lengths[i] = length
		i++
	}
	slices.Sort(lengths)
	builder.WriteString("9. Распределение по длинам цепей\n")
	builder.WriteString("Длина: Количество\n")
	for _, length := range lengths {
		count := distribution[length]
		fmt.Fprintf(&builder, "%d: %d\n", length, count)
	}
	builder.WriteString("\n")
	return builder.String()
}
