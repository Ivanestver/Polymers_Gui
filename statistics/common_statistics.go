package statistics

import (
	"math"
	"polymers/outputformat"
	"polymers/views"
	"strconv"
	"strings"
)

func GetStatistics(globula *views.GlobulaView) string {
	builder := strings.Builder{}
	builder.WriteString(showNumberOfParticles(globula))
	builder.WriteString(showNumberOfChains(globula))
	builder.WriteString(showTheoreticalAgeStatistics(globula))
	builder.WriteString(ShowActualAgeStatistics(globula))
	builder.WriteString(ShowMeanLengthOfChains(globula))
	builder.WriteString(ShowLengthDistribution(globula))
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
