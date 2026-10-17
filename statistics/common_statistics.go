package statistics

import (
	"fmt"
	"math"
	"polymers/outputformat"
	"polymers/views"
	"slices"
	"strconv"
	"strings"
)

func GetCommonStatistics(globula *views.GlobulaView) string {
	builder := strings.Builder{}
	builder.WriteString(showNumberOfParticles(globula))
	builder.WriteString(showNumberOfChains(globula))
	builder.WriteString(showTheoreticalAgeStatistics(globula))
	builder.WriteString(showActualAgeStatistics(globula))
	builder.WriteString(showMeanLengthOfChains(globula))
	builder.WriteString(showLengthDistribution(globula))
	return builder.String()
}

func showNumberOfParticles(globula *views.GlobulaView) string {
	builder := strings.Builder{}
	builder.WriteString("1. Число частиц: ")
	builder.WriteString(strconv.Itoa(globula.GetAtomsCount()))
	builder.WriteString("\n")
	return builder.String()
}

func showNumberOfChains(globula *views.GlobulaView) string {
	builder := strings.Builder{}
	builder.WriteString("2. Исходное число цепей ")
	builder.WriteString(strconv.Itoa(globula.Len()))
	builder.WriteString("\n")
	return builder.String()
}

func showTheoreticalAgeStatistics(globula *views.GlobulaView) string {
	outputformat.GetPrint().Println("Please, specify the theoretical age ratio in percents")
	ageRatioPercent := 25
	ageRatio := float64(ageRatioPercent) * 0.01
	atomsCount := globula.GetAtomsCount()
	//outputformat.GetPrint().Readln(&ageRatioPercent)
	builder := strings.Builder{}

	builder.WriteString("3. Ожидаемая степень старения: ")
	builder.WriteString(strconv.Itoa(int(float64(atomsCount) * ageRatio)))
	builder.WriteString(" (")
	builder.WriteString(strconv.Itoa(ageRatioPercent))
	builder.WriteString("%)\n")

	agedParticlesCount := math.Ceil(float64(atomsCount) * ageRatio)
	cutsCount := int(math.Ceil(agedParticlesCount * 0.44))
	crossCount := int(agedParticlesCount * 0.25)

	builder.WriteString("4. Ожидаемое количество разрывов: ")
	builder.WriteString(strconv.Itoa(cutsCount))
	builder.WriteString("\n")

	builder.WriteString("   Ожидаемое количество сшивок: ")
	builder.WriteString(strconv.Itoa(crossCount))
	builder.WriteString("\n")

	CCount := int(agedParticlesCount * 0.59)
	NCount := int(agedParticlesCount) - CCount
	HCount := crossCount * 2

	builder.WriteString("5. Ожидаемое распределение по состаренным группам:\n")
	builder.WriteString("   C: ")
	builder.WriteString(strconv.Itoa(CCount))
	builder.WriteString("\n")
	builder.WriteString("   N: ")
	builder.WriteString(strconv.Itoa(NCount))
	builder.WriteString("\n")
	builder.WriteString("   H: ")
	builder.WriteString(strconv.Itoa(HCount))
	builder.WriteString("\n")
	return builder.String()
}

func showActualAgeStatistics(globula *views.GlobulaView) string {
	NCount, CCount, cutsCount, crossCount := CalculateActualAgeStatistics(globula)

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

func showMeanLengthOfChains(globula *views.GlobulaView) string {
	meanLengthOfChaings := CalculateMeanLengthOfChains(globula)
	builder := strings.Builder{}
	builder.WriteString("8. Средняя длина цепи: ")
	builder.WriteString(strconv.FormatFloat(meanLengthOfChaings, 'f', 2, 64))
	builder.WriteString("\n")
	return builder.String()
}

func showLengthDistribution(globula *views.GlobulaView) string {
	builder := strings.Builder{}
	distribution := CalculateLengthDistribution(globula) // [length]count_of_such_chains
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
