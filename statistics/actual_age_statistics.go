package statistics

import (
	"polymers/base"
	dt "polymers/datatypes"
	"polymers/views"
)

func CalculateActualAgeStatistics(globula *views.GlobulaView) (NCount, CCount, cutsCount int, crossCount map[int64]int) {
	NCount = 0
	CCount = 0
	cutsCount = 0
	crossCount = make(map[int64]int)
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

	return
}
