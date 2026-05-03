package cristallinity

import (
	"errors"
	"math"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/outputformat"
	"polymers/savers"
	"polymers/views"
)

type _CristallizedStick []*datatypes.Monomer

func (stick *_CristallizedStick) GetDirection() base.Vector3DF {
	if len(*stick) < 3 {
		return base.InvalidVectorF()
	}
	return base.SubtractVecF((*stick)[2].Coords(), (*stick)[0].Coords())
}

func (stick *_CristallizedStick) GetCenterOfMasses() base.Vector3DF {
	ret := base.IndentityVectorF()
	for _, m := range *stick {
		v := m.Coords()
		ret.AddF(&v)
	}
	ret.MultiplyByConstantF(1.0 / float64(len(*stick)))
	return ret
}

func (stick *_CristallizedStick) GetLengthTo(other *_CristallizedStick) float64 {
	v1 := stick.GetCenterOfMasses()
	direction1 := stick.GetDirection()
	v2 := other.GetCenterOfMasses()
	vVector := base.SubtractVecF(v1, v2)
	vecMultiplication := base.VectorProduct(vVector, direction1)
	return vecMultiplication.Len() / direction1.Len()
}

type _CristallizedDomain []_CristallizedStick

func areCodirectional(v1, v2 base.Vector3DF) bool {
	cosOfVectors := base.GetCos(v1, v2)
	return base.CompareFloatWithE(1.0, cosOfVectors, 0.1)
}

func areColinear(v1, v2 base.Vector3DF) bool {
	cosOfVectors := math.Abs(base.GetCos(v1, v2))
	return base.CompareFloatWithE(1.0, cosOfVectors, 0.1)
}

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
	analyzer.carbonSkeleton = append(analyzer.carbonSkeleton, CH2Carbon)
	// 2. Go back and forth to recover the carbon skeleton
	backMonomer, forthMonomer := analyzer.getDirectingMonomers(CH2Carbon)
	if backMonomer == nil && forthMonomer == nil {
		panic("It's impossible if -CH2- exists but its sibings don't")
	}
	if backMonomer != nil {
		// 2.1. Go back
		analyzer.carbonSkeleton = append([]*datatypes.Monomer{backMonomer}, analyzer.carbonSkeleton...)
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
	}
	if forthMonomer != nil {
		// 2.2. Go forth
		analyzer.carbonSkeleton = append(analyzer.carbonSkeleton, forthMonomer)
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

func (analyzer *_CristallinityAnalyzer) findCristallizedParts() []_CristallizedStick {
	const offset = 2
	const invalidMonomerNumber = -1
	initDirection := base.InvalidVectorF()
	sticks := make([]_CristallizedStick, 0)
	startMonomerNumber := invalidMonomerNumber
	endMonomerNumber := invalidMonomerNumber
	for i := 0; i < len(analyzer.carbonSkeleton)-offset; i++ {
		prev := analyzer.carbonSkeleton[i]
		curr := analyzer.carbonSkeleton[i+offset]
		directionVector := base.SubtractVecF(prev.Coords(), curr.Coords())
		if startMonomerNumber == invalidMonomerNumber {
			startMonomerNumber = i
			endMonomerNumber = i + offset
			initDirection = directionVector
		} else {
			if areCodirectional(initDirection, directionVector) {
				endMonomerNumber = i + offset
			} else {
				lenOfStick := endMonomerNumber - startMonomerNumber + 1
				if (lenOfStick-2)/2 < 2 {
					continue
				}
				currStick := make(_CristallizedStick, lenOfStick)
				for j := startMonomerNumber; j <= endMonomerNumber; j++ {
					currStick[j-startMonomerNumber] = analyzer.carbonSkeleton[j]
				}
				sticks = append(sticks, currStick)
				startMonomerNumber = invalidMonomerNumber
			}
		}
	}
	return sticks
}

func (analyzer *_CristallinityAnalyzer) debugSticks(sticks []_CristallizedStick, monomerType datatypes.MonomerType) {
	for _, stick := range sticks {
		for _, monInStick := range stick {
			monInStick.MonomerType = monomerType
		}
	}
	if s, err := savers.SaveToLammps(analyzer.globula); err == nil {
		file, err := os.Create("cristall.data")
		if err == nil {
			defer file.Close()
			file.WriteString(s)
		} else {
			panic(err.Error())
		}
	}
}

func (analyzer *_CristallinityAnalyzer) joinSticksToDomains(sticks []_CristallizedStick) []_CristallizedDomain {
	sets := make([]base.UnorderedSet[int], len(sticks))
	for i := range sets {
		sets[i] = base.UnorderedSet[int]{}
		sets[i].Insert(i)
	}
	for {
		sets1 := make([]base.UnorderedSet[int], 0)
		used := base.UnorderedSet[int]{}
		for i := 0; i < len(sets); i++ {
			if used.Contains(i) {
				continue
			}
			for j := 0; j < len(sets); j++ {
				if used.Contains(j) || i == j {
					continue
				}
				sticks1 := sets[i]
				sticks2 := sets[j]
				if domainsAreClose(sticks1, sticks2, sticks) {
					used.Insert(i)
					used.Insert(j)
					set1 := base.UnorderedSet[int]{}
					for s := range sticks1 {
						set1.Insert(s)
					}
					for s := range sticks2 {
						set1.Insert(s)
					}
					sets1 = append(sets1, set1)
				}
			}
		}
		if len(sets) == len(sets1) || len(sets1) == 0 {
			break
		}
		sets = sets1
	}
	domains := make([]_CristallizedDomain, len(sets))
	for i, s := range sets {
		domains[i] = _CristallizedDomain{}
		for value := range s {
			domains[i] = append(domains[i], sticks[value])
		}
	}
	return domains
}

func domainsAreClose(sticks1, sticks2 base.UnorderedSet[int], sticks []_CristallizedStick) bool {
	for s1 := range sticks1 {
		for s2 := range sticks2 {
			if areClose(sticks[s1], sticks[s2]) {
				return true
			}
		}
	}
	return false
}

func areClose(stick1, stick2 _CristallizedStick) bool {
	return areColinear(
		stick1.GetDirection(),
		stick2.GetDirection(),
	) && stick1.GetLengthTo(&stick2) < 1.44
}

func Analyze(globula *views.GlobulaView) {
	printer := outputformat.GetPrint()
	analyzer, err := makeCristallinityAnalyzer(globula)
	if err != nil {
		printer.PrintflnError("%v", err)
	}
	sticks := analyzer.findCristallizedParts()
	analyzer.joinSticksToDomains(sticks)
}
