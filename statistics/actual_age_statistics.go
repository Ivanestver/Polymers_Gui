package statistics

import (
	"polymers/base"
	dt "polymers/datatypes"
	"polymers/views"
	"strconv"
	"strings"
)

func ShowActualAgeStatistics(globula *views.GlobulaView) string {
	NCount := 0
	CCount := 0
	cutsCount := 0
	crossCount := make(map[int64]int)
	for polNumber := range globula.Len() {
		pol := globula.GetPolymerByIdx(polNumber)
		for monNumber := range pol.Len() {
			mon := pol.GetMonomerByIdx(monNumber)
			if mon.NextMonomer != nil &&
				((mon.MonomerType == base.N || mon.MonomerType == base.O) &&
					(mon.NextMonomer.MonomerType == base.N ||
						mon.NextMonomer.MonomerType == base.O)) {
				cutsCount++
			}
			if mon.MonomerType == base.N {
				CCount++
			} else if mon.MonomerType == base.O {
				NCount++
			} else if _, ok := crossCount[mon.Number]; !ok && mon.MonomerType == base.H {
				crossCount[mon.Number] = 1
				for _, side := range dt.GetMovementSides() {
					connectionType := mon.GetTypeOfConnectionWithSide(side)
					if connectionType == dt.ConnectionTypeCrosslinks {
						crossMon, _ := mon.GetSibling(side)
						crossCount[crossMon.Number] = 1
					}
				}
			}
		}
	}

	builder := strings.Builder{}
	ageGroupsCount := float64(NCount + CCount)
	builder.WriteString("6. Фактическое количество разрывов: ")
	builder.WriteString(strconv.Itoa(cutsCount))
	builder.WriteString("\n")
	builder.WriteString("   Фактическое количество сшивок: ")
	builder.WriteString(strconv.Itoa(len(crossCount) / 2))
	builder.WriteString("\n")
	builder.WriteString("7. Фактическое распределение по состаренным группам:\n")
	builder.WriteString("   C: ")
	builder.WriteString(strconv.Itoa(CCount))
	builder.WriteString(" (")
	if ageGroupsCount != 0 {
		builder.WriteString(strconv.FormatFloat(float64(CCount)/ageGroupsCount*100, 'f', 2, 64))
	} else {
		builder.WriteString("0")
	}
	builder.WriteString("%)")
	builder.WriteString("\n")
	builder.WriteString("   N: ")
	builder.WriteString(strconv.Itoa(NCount))
	builder.WriteString(" (")
	if ageGroupsCount != 0 {
		builder.WriteString(strconv.FormatFloat(float64(NCount)/ageGroupsCount*100, 'f', 2, 64))
	} else {
		builder.WriteString("0")
	}
	builder.WriteString("%)")
	builder.WriteString("\n")
	builder.WriteString("   H: ")
	builder.WriteString(strconv.Itoa(len(crossCount)))
	builder.WriteString("\n")
	builder.WriteString("   Фактическая степень старения: ")
	builder.WriteString(strconv.FormatFloat(ageGroupsCount/float64(globula.GetAtomsCount())*100.0, 'f', 2, 64))
	builder.WriteString("%\n")
	return builder.String()
}
