package statistics

import (
	"polymers/views"
)

func CalculateMeanChainLength(globula *views.GlobulaView) float64 {
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
	})() / float64(len(lengths))
}
