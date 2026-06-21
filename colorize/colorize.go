package colorize

import (
	"polymers/base"
	"polymers/views"
)

func Colorize(globula *views.GlobulaView) error {
	table := base.GetMendeleevTable()
	for polNumber := range globula.Len() {
		polymer := globula.GetPolymerByIdx(polNumber)
		monomerType := table[polNumber]
		for monNumber := range globula.Len() {
			monomer := polymer.GetMonomerByIdx(monNumber)
			monomer.MonomerType = monomerType
		}
	}
	return nil
}
