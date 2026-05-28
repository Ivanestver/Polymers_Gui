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
	d1 := stick.GetDirection()
	d2 := other.GetDirection()
	vectorMul := base.VectorProduct(d1, d2)
	vectorMulLength := vectorMul.Len()
	massCenter1 := stick.GetCenterOfMasses()
	massCenter2 := other.GetCenterOfMasses()
	vecConn := base.SubtractVecF(massCenter1, massCenter2)
	nominator := base.DotProduct(vecConn, vectorMul)
	return nominator / vectorMulLength
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
	printer        outputformat.IPrint
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
		analyzer.printer = outputformat.NewFilePrint(file)
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

func (analyzer *_CristallinityAnalyzer) findCristallizedSticks() []_CristallizedStick {
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
				if int(math.Ceil(float64(lenOfStick)/2)) < 2 {
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
	sets := make(map[int]base.UnorderedSet[int])
	for i := range sticks {
		s := base.UnorderedSet[int]{}
		s.Insert(i)
		sets[i] = s
	}
	for {
		noMoreJoins := true
		joinMap := make(map[int]base.UnorderedSet[int])
		keys := make([]int, 0)
		for i := range sets {
			keys = append(keys, i)
		}
		slices.SortFunc(keys, func(a1, a2 int) int {
			if a1 < a2 {
				return -1
			} else if a1 == a2 {
				return 0
			} else {
				return 1
			}
		})
		for i := 0; i < len(keys); i++ {
			key := keys[i]
			if _, ok := joinMap[key]; ok {
				continue
			} else {
				joinMap[key] = sets[key]
			}
			for j := i + 1; j < len(keys); j++ {
				keyNext := keys[j]
				s := joinMap[key]
				if s.Contains(keyNext) {
					continue
				}
				sticks1 := sets[key]
				sticks2 := sets[keyNext]
				if domainsAreClose(sticks1, sticks2, sticks) {
					s.Insert(keyNext)
					noMoreJoins = false
				}
			}
		}
		if noMoreJoins {
			break
		}
		used := base.UnorderedSet[int]{}
		sets = make(map[int]base.UnorderedSet[int])
		for k := 0; k < len(keys); k++ {
			key := keys[k]
			if used.Contains(key) {
				continue
			}
			stack := base.Stack{key}
			usedInCluster := base.UnorderedSet[int]{}
			newCluster := base.UnorderedSet[int]{}
			for !stack.IsEmpty() {
				in, ok := stack.Pop()
				if !ok {
					continue
				}
				curr := in.(int)
				newCluster.Insert(curr)
				if usedInCluster.Contains(curr) {
					continue
				}
				if s, ok := joinMap[curr]; ok {
					used.Insert(curr)
					usedInCluster.Insert(curr)
					for next := range s {
						stack.Push(next)
					}
				}
			}
			sets[key] = newCluster
		}
	}
	domains := make([]_CristallizedDomain, len(sets))
	i := 0
	for _, s := range sets {
		domains[i] = _CristallizedDomain{}
		for value := range s {
			domains[i] = append(domains[i], sticks[value])
		}
		i++
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
	d1 := stick1.GetDirection()
	d2 := stick2.GetDirection()
	cosTheta := math.Abs(base.GetCos(d1, d2))
	return cosTheta > 0.99 && stick1.GetLengthTo(&stick2) < 2
}

func (analyzer *_CristallinityAnalyzer) analyzeDomains(domains []_CristallizedDomain) {
	analyzer.analyzeCristallinity(domains)
}

func (analyzer *_CristallinityAnalyzer) analyzeCristallinity(domains []_CristallizedDomain) {
	domainMonomersCount := 0
	for _, domain := range domains {
		for _, stick := range domain {
			domainMonomersCount += len(stick)
		}
	}

	analyzer.printer.Printfln("Степень кристалличности: %f\n", float64(domainMonomersCount)/float64(analyzer.globula.GetAtomsCount()))
}

func (analyzer *_CristallinityAnalyzer) analyzeOrientations(sticks []_CristallizedStick) {
	analyzer.printer.Println("Ориентация по формуле <3*cosTheta-1>/2")
	director := analyzer.getDirector(sticks)
	director = director.Normalized()
	S := 0.0
	for _, stick := range sticks {
		stickDirection := stick.GetDirection()
		stickDirection = stickDirection.Normalized()
		cosTheta := base.GetCos(stickDirection, director)
		S += cosTheta * cosTheta
	}
	S = S / float64(len(sticks))
	S = (3.0*S - 1.0) / 2.0
	analyzer.printer.Printfln("S = %f", S)
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
			vectors = append(vectors, directionVector.Normalized())
		}
		sticks[j] = pol
	}
	newGlobula := views.NewGlobulaView(sticks, views.GlobulaGlobulaType)
	if s, err := savers.SaveToLammps(newGlobula); err == nil {
		file, _ := os.Create("cristall_vectors.data")
		defer file.Close()
		file.WriteString(s)
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
	analyzer.printer.Printfln("meanCosTheta = %f\n", meanCosTheta)
	meanCosTheta = (3*meanCosTheta - 1) / 2
	analyzer.printer.Printfln("S = %f\n", meanCosTheta)
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

func (analyzer *_CristallinityAnalyzer) analyzeOrientationViaTensor(sticks []_CristallizedStick) error {
	vectors := func() []base.Vector3DF {
		vectors := make([]base.Vector3DF, len(sticks))
		for i, stick := range sticks {
			vectors[i] = stick.GetDirection()
			vectors[i] = vectors[i].Normalized()
		}
		return vectors
	}()
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

	analyzer.printer.Println("Ориентация с помощью формулы с тензором")
	analyzer.printer.Printfln("S = %.4f\n", S)
	analyzer.printer.Printfln("director = %v", director)
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

func (analyzer *_CristallinityAnalyzer) debugSticksStartEnd(sticks []_CristallizedStick) {
	spacedim := globaldata.GetGlobalData().SpaceDimention
	field := datatypes.NewRealField([3][2]float64{
		{spacedim[base.AxisX].Lower, spacedim[base.AxisX].Higher},
		{spacedim[base.AxisY].Lower, spacedim[base.AxisY].Higher},
		{spacedim[base.AxisZ].Lower, spacedim[base.AxisZ].Higher},
	})
	polymers := make([]datatypes.IPolymer, len(sticks))
	for i, stick := range sticks {
		pol := datatypes.NewRealPolymer(field, int64(i))
		pol.AddMonomer(field.GetMonomerByCoords(stick[0].Coords()))
		pol.AddMonomer(field.GetMonomerByCoords(stick[len(stick)-1].Coords()))
		polymers[i] = pol
	}
	newGlobula := views.NewGlobulaView(polymers, views.GlobulaGlobulaType)
	if s, err := savers.SaveToLammps(newGlobula); err == nil {
		file, _ := os.Create("cristall_vectors_molecular.data")
		defer file.Close()
		file.WriteString(s)
	}
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
	analyzeJoinDomains(analyzer, topPercent)
}

func analyzeWithPercent(analyzer *_CristallinityAnalyzer, topPercent float64) {
	sticks := analyzer.findCristallizedSticks()
	if len(sticks) == 0 {
		analyzer.printer.PrintflnWarning("Отсутствуют кристаллические домены")
		return
	}
	slices.SortFunc(sticks, func(cs1, cs2 _CristallizedStick) int {
		if len(cs1) < len(cs2) {
			return 1
		} else if len(cs1) == len(cs2) {
			return 0
		} else {
			return -1
		}
	})
	partOf := int(float64(len(sticks)) * topPercent)
	sticks = sticks[:partOf]
	analyzer.analyzeOrientations(sticks)
	analyzer.analyzeOrientationViaTensor(sticks)
}

func analyzeJoinDomains(analyzer *_CristallinityAnalyzer, topPercent float64) {
	sticks := analyzer.findCristallizedSticks()
	domains := analyzer.joinSticksToDomains(sticks)
	slices.SortFunc(domains, func(d1, d2 _CristallizedDomain) int {
		lenD1 := len(d1)
		lenD2 := len(d2)
		if lenD1 < lenD2 {
			return 1
		} else if lenD1 == lenD2 {
			return 0
		} else {
			return -1
		}
	})
	firstTop := int(math.Ceil(float64(len(domains)) * topPercent))
	domains = domains[:firstTop]
	sticks = func() []_CristallizedStick {
		sticks := make([]_CristallizedStick, 0)
		for _, domain := range domains {
			sticks = append(sticks, domain[0])
		}
		return sticks
	}()
	analyzer.debugSticks(sticks, base.Oxygen, "domains.dataj")
	analyzer.analyzeOrientations(sticks)
	analyzer.analyzeOrientationViaTensor(sticks)
}
