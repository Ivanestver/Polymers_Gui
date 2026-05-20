package cristallinity

import (
	"errors"
	"fmt"
	"math"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/globaldata"
	"polymers/outputformat"
	"polymers/savers"
	"polymers/views"
	"slices"

	"gonum.org/v1/gonum/mat"
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
	return base.SubtractVecF((*stick)[len(*stick)-1].Coords(), (*stick)[0].Coords())
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

func (stick *_CristallizedStick) GetVectorSkeleton() []base.Vector3DF {
	vectors := make([]base.Vector3DF, 0)
	for i := 1; i < len(*stick); i++ {
		vectors = append(vectors,
			base.SubtractVecF(
				(*stick)[i].Coords(),
				(*stick)[i-1].Coords(),
			),
		)
	}
	return vectors
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

type _CristallinityCarbonSkeletonBuilder struct {
	globula *views.GlobulaView
	level   ScaleLevel
	carbons map[int64]*datatypes.Monomer
}

func makeCarbonSkeletonBuilder(globula *views.GlobulaView, level ScaleLevel) *_CristallinityCarbonSkeletonBuilder {
	return &_CristallinityCarbonSkeletonBuilder{
		globula: globula,
		level:   level,
		carbons: map[int64]*datatypes.Monomer{},
	}
}

func (builder *_CristallinityCarbonSkeletonBuilder) Build() ([]_CarbonSkeleton, error) {
	if builder.level == Molecular {
		carbonSkeleton := make([]_CarbonSkeleton, builder.globula.Len())
		for polNumber := 0; polNumber < builder.globula.Len(); polNumber++ {
			polymer := builder.globula.GetPolymerByIdx(polNumber)
			carbonSkeleton[polNumber] = make(_CarbonSkeleton, polymer.Len())
			for monNumber := 0; monNumber < polymer.Len(); monNumber++ {
				carbonSkeleton[polNumber][monNumber] = polymer.GetMonomerByIdx(monNumber)
			}
		}
		return carbonSkeleton, nil
	}
	// 1. Find a CH2 to realize where a carbon skeleton is
	builder.defineAllCH2Carbons()
	if len(builder.carbons) == 0 {
		return nil, errors.New("отсутствуют мономеры -CH2-")
	}
	carbonSkeletons := make([]_CarbonSkeleton, 0)
	keys := make([]int64, 0)
	for key := range builder.carbons {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	for len(keys) > 0 {
		// 2. Go back and forth to recover the carbon skeleton
		key := keys[0]
		CH2Carbon := builder.carbons[key]
		delete(builder.carbons, key)
		keys = keys[1:]
		carbonSkeleton := make(_CarbonSkeleton, 0)
		carbonSkeleton = append(carbonSkeleton, CH2Carbon)
		backMonomer, forthMonomer := builder.getDirectingMonomers(CH2Carbon)
		if backMonomer == nil && forthMonomer == nil {
			panic("It's impossible if -CH2- exists but its sibings don't")
		}
		if backMonomer != nil {
			// 2.1. Go back
			carbonSkeleton = append([]*datatypes.Monomer{backMonomer}, carbonSkeleton...)
			if err := builder.defineSkeleton(&carbonSkeleton,
				func(s *_CarbonSkeleton) *datatypes.Monomer {
					return (*s)[0]
				},
				func(s *_CarbonSkeleton) *datatypes.Monomer {
					return (*s)[1]
				},
				func(s *_CarbonSkeleton, newMon *datatypes.Monomer) {
					*s = append([]*datatypes.Monomer{newMon}, (*s)...)
					delete(builder.carbons, newMon.Number)
					if idx := slices.Index(keys, newMon.Number); idx != -1 {
						keys = append(keys[:idx], keys[idx+1:]...)
					}
				}); err != nil {
				return nil, err
			}
		}
		if forthMonomer != nil {
			// 2.2. Go forth
			carbonSkeleton = append(carbonSkeleton, forthMonomer)
			if err := builder.defineSkeleton(&carbonSkeleton, func(s *_CarbonSkeleton) *datatypes.Monomer {
				return (*s)[len(*s)-1]
			},
				func(s *_CarbonSkeleton) *datatypes.Monomer {
					return (*s)[len(*s)-2]
				},
				func(s *_CarbonSkeleton, newMon *datatypes.Monomer) {
					*s = append(*s, newMon)
					delete(builder.carbons, newMon.Number)
					if idx := slices.Index(keys, newMon.Number); idx != -1 {
						keys = append(keys[:idx], keys[idx+1:]...)
					}
				}); err != nil {
				return nil, err
			}
		}
		carbonSkeletons = append(carbonSkeletons, carbonSkeleton)
		carbonSkeleton = make(_CarbonSkeleton, 0)
	}
	return carbonSkeletons, nil
}

func (builder *_CristallinityCarbonSkeletonBuilder) defineAllCH2Carbons() {
	for polymernNumber := 0; polymernNumber < builder.globula.Len(); polymernNumber++ {
		polymer := builder.globula.GetPolymerByIdx(polymernNumber)
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
			if HCount >= 2 {
				builder.carbons[monomer.Number] = monomer
			}
		}
	}
}

func (builder *_CristallinityCarbonSkeletonBuilder) defineSkeleton(carbonSkeleton *_CarbonSkeleton, fCurrMon, fPrevMon func(s *_CarbonSkeleton) *datatypes.Monomer, fAdd func(s *_CarbonSkeleton, newMon *datatypes.Monomer)) error {
	// Assume startPoint and directingMonomer are already in the skeleton
	// Where to move
	canMove := true
	for canMove {
		currMonomer := fCurrMon(carbonSkeleton)
		if currMonomer == nil {
			return errors.New("атом в углеродном скелете не может быть пустым местом")
		}
		canMove = false
		for _, sibling := range currMonomer.GetSiblings() {
			if sibling == nil || sibling.IsNotTypeOf(base.Carbon) {
				continue
			}
			if sibling == fPrevMon(carbonSkeleton) {
				continue
			}
			fAdd(carbonSkeleton, sibling)
			canMove = true
			break
		}
	}
	return nil
}

func (builder *_CristallinityCarbonSkeletonBuilder) getDirectingMonomers(startMonomer *datatypes.Monomer) (backMonomer, forthMonomer *datatypes.Monomer) {
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

type _CristallinityAnalyzer struct {
	globula        *views.GlobulaView
	carbonSkeleton []_CarbonSkeleton
	offset         int
	outputFile     *os.File
	level          ScaleLevel
	baseStick      []*datatypes.Monomer
}

func makeCristallinityAnalyzer(globula *views.GlobulaView, offset int, outputFilename string, level ScaleLevel, baseElement base.MendeleevTableElement) (*_CristallinityAnalyzer, error) {
	analyzer := &_CristallinityAnalyzer{
		globula:        globula,
		carbonSkeleton: make([]_CarbonSkeleton, 0),
		offset:         offset,
		level:          level,
		baseStick:      nil,
	}
	builder := makeCarbonSkeletonBuilder(globula, level)
	if skeletons, err := builder.Build(); err != nil {
		return nil, err
	} else {
		analyzer.carbonSkeleton = skeletons
	}
	if file, err := os.Create(outputFilename); err == nil {
		analyzer.outputFile = file
	} else {
		return nil, err
	}
	analyzer.baseStick = defineBaseStick(globula, baseElement)
	return analyzer, nil
}

func defineBaseStick(globula *views.GlobulaView, baseElement base.MendeleevTableElement) []*datatypes.Monomer {
	if baseElement == base.MendeleevTableElementUndefined {
		return nil
	}

	baseStick := make([]*datatypes.Monomer, 0)
	for polNumber := 0; polNumber < globula.Len(); polNumber++ {
		polymer := globula.GetPolymerByIdx(polNumber)
		if polymer.GetMonomerByIdx(0).MonomerType == baseElement {
			for monNumber := 0; monNumber < polymer.Len(); monNumber++ {
				baseStick = append(baseStick, polymer.GetMonomerByIdx(monNumber))
			}
		}
	}
	return baseStick
}

// func (analyzer *_CristallinityAnalyzer) defineSkeleton(carbonSkeleton *_CarbonSkeleton, fCurrMon, fPrevMon func(s *_CarbonSkeleton) *datatypes.Monomer, fAdd func(s *_CarbonSkeleton, newMon *datatypes.Monomer)) error {
// 	// Assume startPoint and directingMonomer are already in the skeleton
// 	// Where to move
// 	canMove := true
// 	for canMove {
// 		currMonomer := fCurrMon(carbonSkeleton)
// 		if currMonomer == nil {
// 			return errors.New("атом в углеродном скелете не может быть пустым местом")
// 		}
// 		canMove = false
// 		for _, sibling := range currMonomer.GetSiblings() {
// 			if sibling == nil || sibling.IsNotTypeOf(base.Carbon) {
// 				continue
// 			}
// 			if sibling == fPrevMon(carbonSkeleton) {
// 				continue
// 			}
// 			fAdd(carbonSkeleton, sibling)
// 			canMove = true
// 			break
// 		}
// 	}
// 	return nil
// }

func (analyzer *_CristallinityAnalyzer) findCristallizedParts() []_CristallizedStick {
	sticks := make([]_CristallizedStick, 0)
	offset := analyzer.offset
	for _, carbonSkeleton := range analyzer.carbonSkeleton {
		const invalidMonomerNumber = -1
		initDirection := base.InvalidVectorF()
		startMonomerNumber := invalidMonomerNumber
		endMonomerNumber := invalidMonomerNumber
		for i := 0; i < len(carbonSkeleton)-offset; i += offset {
			prev := carbonSkeleton[i]
			curr := carbonSkeleton[i+offset]
			directionVector := base.SubtractVecF(curr.Coords(), prev.Coords())
			if startMonomerNumber == invalidMonomerNumber {
				startMonomerNumber = i
				endMonomerNumber = i + offset
				initDirection = directionVector
			} else {
				if areCodirectional(initDirection, directionVector) {
					endMonomerNumber = i + offset
					initDirection = directionVector
					if i < len(carbonSkeleton)-2-offset {
						continue
					}
				}
				lenOfStick := endMonomerNumber - startMonomerNumber + 1
				if lenOfStick/2 < 2 {
					startMonomerNumber = invalidMonomerNumber
					continue
				}
				currStick := make(_CristallizedStick, lenOfStick)
				for j := startMonomerNumber; j <= endMonomerNumber; j++ {
					currStick[j-startMonomerNumber] = carbonSkeleton[j]
				}
				sticks = append(sticks, currStick)
				startMonomerNumber = i
				endMonomerNumber = i + offset
				initDirection = directionVector
			}
		}
	}
	return sticks
}

func (analyzer *_CristallinityAnalyzer) debugSticks(sticks []_CristallizedStick, monomerType base.MendeleevTableElement, filename string) {
	for _, stick := range sticks {
		for _, monInStick := range stick {
			monInStick.MonomerType = monomerType
		}
	}
	if len(filename) == 0 {
		return
	}
	if s, err := savers.SaveToLammps(analyzer.globula); err == nil {
		file, err := os.Create(filename)
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
	fmt.Fprintln(analyzer.outputFile, "Orientation")
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
	S = (3.0*S - 1.0) / 2.0
	fmt.Fprintf(analyzer.outputFile, "S = %f", S)
}

func (analyzer *_CristallinityAnalyzer) getDirector(sticks []_CristallizedStick) base.Vector3DF {
	if analyzer.baseStick == nil {
	director := base.IdentityVectorF()
	for _, stick := range sticks {
		stickDirection := stick.GetDirection()
		director.AddF(&stickDirection)
	}
	return director
	} else {
		return base.SubtractVecF(
			analyzer.baseStick[len(analyzer.baseStick)-1].Coords(),
			analyzer.baseStick[0].Coords())
	}
}

func (analyzer *_CristallinityAnalyzer) calculateS() error {
	vectors := make([]base.Vector3DF, 0)
	offset := analyzer.offset
	spacedim := globaldata.GetGlobalData().SpaceDimention
	field := datatypes.NewRealField([3][2]float64{
		{spacedim[base.AxisX].Lower, spacedim[base.AxisX].Higher},
		{spacedim[base.AxisY].Lower, spacedim[base.AxisY].Higher},
		{spacedim[base.AxisZ].Lower, spacedim[base.AxisZ].Higher},
	})
	sticks := make([]datatypes.IPolymer, len(analyzer.carbonSkeleton))
	for j, carbonSkeleton := range analyzer.carbonSkeleton {
		pol := datatypes.NewRealPolymer(field, int64(j))
		pol.AddMonomer(field.GetMonomerByCoords(carbonSkeleton[0].Coords()))
		for i := 0; i < len(carbonSkeleton)-offset; i += offset {
			prev := carbonSkeleton[i]
			curr := carbonSkeleton[i+offset]
			pol.AddMonomer(field.GetMonomerByCoords(curr.Coords()))
			directionVector := base.SubtractVecF(curr.Coords(), prev.Coords())
			vectors = append(vectors, directionVector)
		}
		sticks[j] = pol
	}
	newGlobula := views.NewGlobulaView(sticks, views.GlobulaGlobulaType)
	if s, err := savers.SaveToLammps(newGlobula); err == nil {
		file, _ := os.Create("cristall_vectors.data")
		defer file.Close()
		file.WriteString(s)
	} else {
		fmt.Printf("%v", err)
	}
	meanCosTheta := 0.0
	for i := 0; i < len(vectors)-1; i++ {
		for j := i + 1; j < len(vectors); j++ {
			cosTheta := base.GetCos(vectors[i], vectors[j])
			meanCosTheta += cosTheta * cosTheta
		}
	}
	n := float64(len(vectors))
	meanCosTheta = meanCosTheta / (n * (n - 1) / 2.0)
	fmt.Fprintf(analyzer.outputFile, "meanCosTheta = %f\n", meanCosTheta)
	meanCosTheta = (3*meanCosTheta - 1) / 2
	fmt.Fprintf(analyzer.outputFile, "S = %f\n", meanCosTheta)
	return nil
}

func (analyzer *_CristallinityAnalyzer) getVectors() []base.Vector3DF {
	vectors := make([]base.Vector3DF, 0)
	offset := analyzer.offset
	for _, carbonSkeleton := range analyzer.carbonSkeleton {
		for i := 0; i < len(carbonSkeleton)-offset; i += offset {
			prev := carbonSkeleton[i]
			curr := carbonSkeleton[i+offset]
			directionVector := base.SubtractVecF(curr.Coords(), prev.Coords())
			vectors = append(vectors, directionVector.Normalized())
		}
	}
	return vectors
}

func (analyzer *_CristallinityAnalyzer) analyzeOrientationViaTensor(vectors []base.Vector3DF) error {
	n := float64(len(vectors))
	outerProduct := mat.NewSymDense(int(base.AxisCount), nil)
	for _, v := range vectors {
		ux := v[base.AxisX]
		uy := v[base.AxisY]
		uz := v[base.AxisZ]
		data := outerProduct.RawSymmetric().Data
		data[0] += ux * ux
		data[1] += ux * uy
		data[2] += ux * uz
		data[3] += uy * ux
		data[4] += uy * uy
		data[5] += uy * uz
		data[6] += uz * ux
		data[7] += uz * uy
		data[8] += uz * uz
	}

	outerProduct.ScaleSym(1.0/n, outerProduct)
	for i := range int(base.AxisCount) {
		val := outerProduct.At(i, i)
		outerProduct.SetSym(i, i, val-1.0/3)
	}

	var eig mat.EigenSym
	if ok := eig.Factorize(outerProduct, true); !ok {
		return errors.New("не удалось разложить матрицу")
	}

	values := eig.Values(nil)

	maxIdx := 0
	for i := 1; i < len(values); i++ {
		if values[i] > values[maxIdx] {
			maxIdx = i
			break
		}
	}
	S := 1.5 * values[maxIdx]
	var eVecs mat.Dense
	eig.VectorsTo(&eVecs)
	director := base.Vector3DF{
		eVecs.At(0, maxIdx),
		eVecs.At(1, maxIdx),
		eVecs.At(2, maxIdx),
	}

	fmt.Fprintln(analyzer.outputFile, "Orientation")
	fmt.Fprintf(analyzer.outputFile, "S = %.4f\n", S)
	fmt.Fprintf(analyzer.outputFile, "director = %v", director)
	return nil
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

func (analyzer *_CristallinityAnalyzer) toVectorsFlat(sticks []_CristallizedStick) []base.Vector3DF {
	vectors := make([]base.Vector3DF, 0)
	for _, stick := range sticks {
		for i := analyzer.offset; i < len(stick); i += analyzer.offset {
			v := base.SubtractVecF(
				stick[i].Coords(),
				stick[i-analyzer.offset].Coords(),
			)
			vectors = append(vectors, v)
		}
	}
	return vectors
}

func Analyze(globula *views.GlobulaView, offset int, outputFilename string, level ScaleLevel, baseElem base.MendeleevTableElement, topPercent float64) {
	printer := outputformat.GetPrint()
	if err := validateInputParams(offset, outputFilename, level); err != nil {
		printer.PrintlnError(err.Error())
		return
	}
	analyzer, err := makeCristallinityAnalyzer(globula, offset, outputFilename, level, baseElem)
	if err != nil {
		printer.PrintflnError("%v", err)
		return
	}
	defer analyzer.outputFile.Close()
	// if err := analyzer.calculateS(); err != nil {
	// 	fmt.Printf("%v\n", err)
	// }
	// sticks := make([]_CristallizedStick, len(analyzer.carbonSkeleton))
	// for i, skeleton := range analyzer.carbonSkeleton {
	// 	sticks[i] = make(_CristallizedStick, len(skeleton))
	// 	copy(sticks[i], skeleton)
	// }
	// analyzer.debugSticks(sticks, base.Oxygen, "cristall.data")
	// analyzer.debugSticks(sticks, base.Carbon, "")
	sticks1 := analyzer.findCristallizedParts()
	if len(sticks1) == 0 {
		printer.PrintflnWarning("Отсутствуют кристаллические домены")
		return
	}
	slices.SortFunc(sticks1, func(cs1, cs2 _CristallizedStick) int {
		if len(cs1) < len(cs2) {
			return 1
		} else if len(cs1) == len(cs2) {
			return 0
		} else {
			return -1
		}
	})
	partOf := int(float64(len(sticks1)) * topPercent)
	sticks1 = sticks1[:partOf]
	analyzer.analyzeOrientations(sticks1)
	analyzer.debugSticks(sticks1, base.Fluorine, "cristall_orientation.data")
	// domains := analyzer.joinSticksToDomains(sticks)
	// analyzer.analyzeDomains(domains)
	// vectors := analyzer.getVectors()
	// vectors := analyzer.toVectorsFlat(sticks1)
	// analyzer.analyzeOrientationViaTensor(vectors)
}
