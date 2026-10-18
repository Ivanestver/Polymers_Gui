package statistics

import (
	"errors"
	"polymers/views"
)

func CalculateMeanLengthOfChains(globula *views.GlobulaView) float64 {
	sum := 0.0
	for polNumber := range globula.Len() {
		sum += float64(globula.GetPolymerByIdx(polNumber).Len())
	}
	return sum / float64(globula.Len())
}

func CalculateMeanLengthOfCrystall(globula *views.GlobulaView) (float64, error) {
	if !globula.Is(views.GlobulaCrystalType) {
		return 0.0, errors.New("глобула не является кристаллической")
	}
	// Now we assume that globula consists of only one polymer
	lengths := make([]float64, 0)
	for polNumber := range globula.Len() {
		pol := globula.GetPolymerByIdx(polNumber)
		length := 0.0
		for monNumber := range pol.Len() {
			mon := pol.GetMonomerByIdx(monNumber)
			if mon.NextMonomer != nil {
				length += 1.0
			} else {
				lengths = append(lengths, length)
				length = 0.0
			}
		}
	}
	return (func() float64 {
		sum := lengths[0]
		for i := 1; i < len(lengths); i++ {
			sum += lengths[i]
		}
		return sum
	})() / float64(len(lengths)), nil
}
