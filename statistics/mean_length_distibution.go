package statistics

import (
	"polymers/views"
	"strconv"
	"strings"
)

func ShowMeanLengthOfChains(globula *views.GlobulaView) string {
	builder := strings.Builder{}
	builder.WriteString("8. Средняя длина цепи: ")
	atomsCount := globula.GetAtomsCount()
	builder.WriteString(strconv.FormatFloat(float64(atomsCount)/float64(globula.Len()), 'f', 2, 64))
	builder.WriteString("\n")
	return builder.String()
}
