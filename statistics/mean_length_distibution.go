package statistics

import (
	"polymers/views"
)

func CalculateMeanLengthOfChains(globula *views.GlobulaView) float64 {
	atomsCount := globula.GetAtomsCount()
	return float64(atomsCount) / float64(globula.Len())
}
