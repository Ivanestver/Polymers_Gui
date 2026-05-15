package cristallinity

import (
	"errors"
	"fmt"
	"math"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/outputformat"
	"polymers/savers"
	"polymers/views"
	"slices"
)

type ScaleLevel = string // Атомистический, молекулярный и т.д.

const (
	Atomistic ScaleLevel = "atomistic"
	Molecular ScaleLevel = "molecular"
)

type _CristallizedStick []*datatypes.Monomer

func (stick *_CristallizedStick) GetDirection() base.Vector3DF {
	if len(*stick) < 3 {
		return base.InvalidVectorF()
	}
	return base.SubtractVecF((*stick)[2].Coords(), (*stick)[0].Coords())
}

func (stick *_CristallizedStick) GetCenterOfMasses() base.Vector3DF {
	ret := base.IdentityVectorF()
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

type _CarbonSkeleton []*datatypes.Monomer

type _CristallinityAnalyzer struct {
	globula        *views.GlobulaView
	carbonSkeleton []_CarbonSkeleton
	offset         int
	outputFile     *os.File
	level          ScaleLevel
}

func makeCristallinityAnalyzer(globula *views.GlobulaView, offset int, outputFilename string, level ScaleLevel) (*_CristallinityAnalyzer, error) {
	analyzer := &_CristallinityAnalyzer{
		globula:        globula,
		carbonSkeleton: make([]_CarbonSkeleton, 0),
		offset:         offset,
		level:          level,
	}
	if err := analyzer.defineCarbonSkeleton(); err != nil {
		return nil, err
	}
	if file, err := os.Create(outputFilename); err == nil {
		analyzer.outputFile = file
	} else {
		return nil, err
	}
	return analyzer, nil
}

func (analyzer *_CristallinityAnalyzer) defineCarbonSkeleton() error {
	// 0. Check whether it's an atomistic representation or a molecular one
	if analyzer.level == Molecular {
		analyzer.carbonSkeleton = make([]_CarbonSkeleton, analyzer.globula.Len())
		for polNumber := 0; polNumber < analyzer.globula.Len(); polNumber++ {
			polymer := analyzer.globula.GetPolymerByIdx(polNumber)
			analyzer.carbonSkeleton[polNumber] = make(_CarbonSkeleton, polymer.Len())
			for monNumber := 0; monNumber < polymer.Len(); monNumber++ {
				analyzer.carbonSkeleton[polNumber][monNumber] = polymer.GetMonomerByIdx(monNumber)
			}
		}
		return nil
	}
	// 1. Find a CH2 to realize where a carbon skeleton is
	CH2Carbon := analyzer.getCH2Cabron()
	if CH2Carbon == nil {
		return errors.New("отсутствуют мономеры -CH2-")
	}
	carbonSkeleton := make(_CarbonSkeleton, 0)
	carbonSkeleton = append(carbonSkeleton, CH2Carbon)
	// 2. Go back and forth to recover the carbon skeleton
	backMonomer, forthMonomer := analyzer.getDirectingMonomers(CH2Carbon)
	if backMonomer == nil && forthMonomer == nil {
		panic("It's impossible if -CH2- exists but its sibings don't")
	}
	if backMonomer != nil {
		// 2.1. Go back
		carbonSkeleton = append([]*datatypes.Monomer{backMonomer}, carbonSkeleton...)
		if err := analyzer.defineSkeleton(func(s *_CarbonSkeleton) *datatypes.Monomer {
			return (*s)[0]
		},
			func(s *_CarbonSkeleton) *datatypes.Monomer {
				return (*s)[1]
			},
			func(s *_CarbonSkeleton, newMon *datatypes.Monomer) {
				*s = append([]*datatypes.Monomer{newMon}, (*s)...)
			}); err != nil {
			return err
		}
	}
	if forthMonomer != nil {
		// 2.2. Go forth
		carbonSkeleton = append(carbonSkeleton, forthMonomer)
		if err := analyzer.defineSkeleton(func(s *_CarbonSkeleton) *datatypes.Monomer {
			return (*s)[len(*s)-1]
		},
			func(s *_CarbonSkeleton) *datatypes.Monomer {
				return (*s)[len(*s)-2]
			},
			func(s *_CarbonSkeleton, newMon *datatypes.Monomer) {
				*s = append(*s, newMon)
			}); err != nil {
			return err
		}
	}
	analyzer.carbonSkeleton = append(analyzer.carbonSkeleton, carbonSkeleton)
	return nil
}

func (analyzer *_CristallinityAnalyzer) getCH2Cabron() *datatypes.Monomer {
	for polymernNumber := 0; polymernNumber < analyzer.globula.Len(); polymernNumber++ {
		polymer := analyzer.globula.GetPolymerByIdx(polymernNumber)
		for monomerIdx := 0; monomerIdx < polymer.Len(); monomerIdx++ {
			monomer := polymer.GetMonomerByIdx(monomerIdx)
			if monomer.IsNotTypeOf(base.Carbon) {
				continue
			}
			siblings := monomer.GetSiblings()
			HCount := 0
			for _, sibling := range siblings {
				if sibling == nil {
					continue
				}
				if sibling.IsTypeOf(base.Hydrogen) {
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
		if sibling.IsTypeOf(base.Carbon) {
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

func (analyzer *_CristallinityAnalyzer) defineSkeleton(fCurrMon, fPrevMon func(s *_CarbonSkeleton) *datatypes.Monomer, fAdd func(s *_CarbonSkeleton, newMon *datatypes.Monomer)) error {
	// Assume startPoint and directingMonomer are already in the skeleton
	// Where to move
	for skeletonNumber := 0; skeletonNumber < len(analyzer.carbonSkeleton); skeletonNumber++ {
		carbonSkeleton := analyzer.carbonSkeleton[skeletonNumber]
		canMove := true
		for canMove {
			currMonomer := fCurrMon(&carbonSkeleton)
			if currMonomer == nil {
				return errors.New("атом в углеродном скелете не может быть пустым местом")
			}
			canMove = false
			for _, sibling := range currMonomer.GetSiblings() {
				if sibling == nil || sibling.IsNotTypeOf(base.Carbon) {
					continue
				}
				if sibling == fPrevMon(&carbonSkeleton) {
					continue
				}
				fAdd(&carbonSkeleton, sibling)
				canMove = true
				break
			}
		}
	}
	return nil
}

func (analyzer *_CristallinityAnalyzer) findCristallizedParts() []_CristallizedStick {
	sticks := make([]_CristallizedStick, 0)
	offset := analyzer.offset
	for _, carbonSkeleton := range analyzer.carbonSkeleton {
		const invalidMonomerNumber = -1
		initDirection := base.InvalidVectorF()
		startMonomerNumber := invalidMonomerNumber
		endMonomerNumber := invalidMonomerNumber
		for i := 0; i < len(carbonSkeleton)-offset; i++ {
			prev := carbonSkeleton[i]
			curr := carbonSkeleton[i+offset]
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
						currStick[j-startMonomerNumber] = carbonSkeleton[j]
					}
					sticks = append(sticks, currStick)
					startMonomerNumber = invalidMonomerNumber
				}
			}
		}
	}
	return sticks
}

func (analyzer *_CristallinityAnalyzer) debugSticks(sticks []_CristallizedStick, monomerType base.MendeleevTableElement) {
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

func (analyzer *_CristallinityAnalyzer) analyzeDomains(domains []_CristallizedDomain) {
	analyzer.analyzeCristallinity(domains, analyzer.outputFile)
}

func (analyzer *_CristallinityAnalyzer) analyzeCristallinity(domains []_CristallizedDomain, file *os.File) {
	domainMonomersCount := 0
	for _, domain := range domains {
		for _, stick := range domain {
			domainMonomersCount += len(stick)
		}
	}

	fmt.Fprintf(file, "Степень кристалличности: %f\n", float64(domainMonomersCount)/float64(analyzer.globula.GetAtomsCount()))
}

func (analyzer *_CristallinityAnalyzer) analyzeOrientations(sticks []_CristallizedStick) {
	director := analyzer.getDirector(sticks)
	S := 0.0
	for _, stick := range sticks {
		stickDirection := stick.GetDirection()
		cosTheta := base.GetCos(stickDirection, director)
		fmt.Println(cosTheta)
		S += cosTheta * cosTheta
	}
	S = S / float64(len(sticks))
	fmt.Println(S)
	fmt.Println(math.Acos(math.Sqrt(S)))
	S = (3.0*(S/float64(len(sticks))) - 1.0) / 2.0
	fmt.Fprintf(analyzer.outputFile, "S = %f", S)
}

func (analyzer *_CristallinityAnalyzer) getDirector(sticks []_CristallizedStick) base.Vector3DF {
	// director := base.IdentityVectorF()
	// for _, stick := range sticks {
	// 	stickDirection := stick.GetDirection()
	// 	director.AddF(&stickDirection)
	// }
	// return director
	director := slices.MaxFunc(sticks, func(s1, s2 _CristallizedStick) int {
		if len(s1) < len(s2) {
			return -1
		} else if len(s1) == len(s2) {
			return 0
		} else {
			return 1
		}
	})
	return director.GetDirection()
}

func validateInputParams(offset int, outputFilename string, level ScaleLevel) error {
	if offset < 1 {
		return errors.New("offset должен быть больше 0")
	}
	if len(outputFilename) == 0 {
		return errors.New("название файла не должно быть пустым")
	}
	if level != Atomistic && level != Molecular {
		return fmt.Errorf("неверный уровень: %s", string(level))
	}
	return nil
}

func Analyze(globula *views.GlobulaView, offset int, outputFilename string, level ScaleLevel) {
	printer := outputformat.GetPrint()
	if err := validateInputParams(offset, outputFilename, level); err != nil {
		printer.PrintlnError(err.Error())
		return
	}
	analyzer, err := makeCristallinityAnalyzer(globula, offset, outputFilename, level)
	if err != nil {
		printer.PrintflnError("%v", err)
		return
	}
	defer analyzer.outputFile.Close()
	sticks := analyzer.findCristallizedParts()
	if len(sticks) == 0 {
		printer.PrintflnWarning("Отсутствуют кристаллические домены")
		return
	}
	analyzer.analyzeOrientations(sticks)
	// analyzer.debugSticks(sticks, base.Hydrogen)
	// domains := analyzer.joinSticksToDomains(sticks)
	// analyzer.analyzeDomains(domains)
}
