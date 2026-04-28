package cristallinity

import (
	"errors"
	"math"
	"polymers/base"
	"polymers/datatypes"
	"polymers/outputformat"
	"polymers/views"
)

type _CristallinityAnalyzer struct {
	globula        *views.GlobulaView
	carbonSkeleton []*datatypes.Monomer
}

func makeCristallinityAnalyzer(globula *views.GlobulaView) (*_CristallinityAnalyzer, error) {
	analyzer := &_CristallinityAnalyzer{
		globula:        globula,
		carbonSkeleton: make([]*datatypes.Monomer, 0),
	}
	if err := analyzer.defineCarbonSkeleton(); err != nil {
		return nil, err
	}
	return analyzer, nil
}

func (analyzer *_CristallinityAnalyzer) defineCarbonSkeleton() error {
	// 1. Find a CH2 to realize where a carbon skeleton is
	CH2Carbon := analyzer.getCH2Cabron()
	if CH2Carbon == nil {
		return errors.New("отсутствуют мономеры -CH2-")
	}
	// 2. Go back and forth to recover the carbon skeleton
	backMonomer, forthMonomer := analyzer.getDirectingMonomers(CH2Carbon)
	if backMonomer == nil || forthMonomer == nil {
		panic("It's impossible if -CH2- exists but its sibings don't")
	}
	// 2.1. Go back
	if err := analyzer.defineSkeleton(func(s *[]*datatypes.Monomer) *datatypes.Monomer {
		return (*s)[0]
	},
		func(s *[]*datatypes.Monomer) *datatypes.Monomer {
			return (*s)[1]
		},
		func(s *[]*datatypes.Monomer, newMon *datatypes.Monomer) {
			*s = append([]*datatypes.Monomer{newMon}, (*s)...)
		}); err != nil {
		return err
	}
	// 2.2. Go forth
	if err := analyzer.defineSkeleton(func(s *[]*datatypes.Monomer) *datatypes.Monomer {
		return (*s)[len(*s)-1]
	},
		func(s *[]*datatypes.Monomer) *datatypes.Monomer {
			return (*s)[len(*s)-2]
		},
		func(s *[]*datatypes.Monomer, newMon *datatypes.Monomer) {
			*s = append(*s, newMon)
		}); err != nil {
		return err
	}
	return nil
}

func (analyzer *_CristallinityAnalyzer) getCH2Cabron() *datatypes.Monomer {
	for polymernNumber := 0; polymernNumber < analyzer.globula.Len(); polymernNumber++ {
		polymer := analyzer.globula.GetPolymerByIdx(polymernNumber)
		for monomerIdx := 0; monomerIdx < polymer.Len(); monomerIdx++ {
			monomer := polymer.GetMonomerByIdx(monomerIdx)
			if monomer.MonomerType != datatypes.MonomerTypeUsual {
				continue
			}
			siblings := monomer.GetSiblings()
			HCount := 0
			for _, sibling := range siblings {
				if sibling == nil {
					continue
				}
				if sibling.MonomerType == datatypes.MonomerTypeCrosslinked {
					HCount++
				}
			}
			if HCount == 2 {
				return monomer
			}
		}
	}
	return nil
}

func (analyzer *_CristallinityAnalyzer) getDirectingMonomers(startMonomer *datatypes.Monomer) (backMonomer, forthMonomer *datatypes.Monomer) {
	siblings := startMonomer.GetSiblings()
	for _, sibling := range siblings {
		if sibling.MonomerType == datatypes.MonomerTypeUsual {
			if backMonomer == nil {
				backMonomer = sibling
			} else if forthMonomer == nil {
				forthMonomer = sibling
			} else {
				break
			}
		}
	}
	return
}

func (analyzer *_CristallinityAnalyzer) defineSkeleton(fCurrMon, fPrevMon func(s *[]*datatypes.Monomer) *datatypes.Monomer, fAdd func(s *[]*datatypes.Monomer, newMon *datatypes.Monomer)) error {
	// Assume startPoint and directingMonomer are already in the skeleton
	// Where to move
	canMove := true
	for canMove {
		currMonomer := fCurrMon(&analyzer.carbonSkeleton)
		if currMonomer == nil {
			return errors.New("атом в углеродном скелете не может быть пустым местом")
		}
		canMove = false
		for _, sibling := range currMonomer.GetSiblings() {
			if sibling == nil || sibling.MonomerType != datatypes.MonomerTypeUsual {
				continue
			}
			if sibling == fPrevMon(&analyzer.carbonSkeleton) {
				continue
			}
				fAdd(&analyzer.carbonSkeleton, sibling)
				canMove = true
				break
		}
	}
	return nil
}

func Analyze(globula *views.GlobulaView) {
	printer := outputformat.GetPrint()
	_, err := makeCristallinityAnalyzer(globula)
	if err != nil {
		printer.PrintflnError("%v", err)
	}
}
