package buildglobula

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/globaldata"
	"polymers/outputformat"
	"polymers/views"
	"strconv"
	"unicode/utf8"
)

type PatternInputDataBuilder struct {
}

func (builder PatternInputDataBuilder) CreateInputData(defaultParams []string, particleName string) (ICalcAlgInputData, error) {
	if len(particleName) == 0 {
		return nil, errors.New("не задан файл с паттерном")
	}
	printStatistics := true
	polymersCount := 1
	if len(defaultParams) > 0 {
		printStatistics = defaultParams[0] == "true"
	}
	if len(defaultParams) > 1 {
		if c, err := strconv.Atoi(defaultParams[1]); err == nil {
			polymersCount = c
		} else {
			return nil, err
		}
	}

	return PatternInputData{
		patternFileName: particleName,
		printStatistics: printStatistics,
		polymersCount:   polymersCount,
	}, nil
}

type PatternInputData struct {
	patternFileName string
	printStatistics bool
	polymersCount   int
}

func (inputData PatternInputData) GetGlobulaType() views.GlobulaProperty {
	return views.GlobulaPatternType
}

type PatternCalcAlg AbstractAlg[PatternInputData]

/*
Y
^
|
|*-* *-* *-*
|| | | | | |
|* * * * * * *
|| | | | | | |
|* *-* *-* *-*
 ----------------------------------> X
*/

func (alg *PatternCalcAlg) Calc() []*datatypes.Polymer {
	file, err := os.Open(alg.inputData.patternFileName)
	if err != nil {
		outputformat.GetPrint().PrintlnError(err.Error())
		return nil
	}
	spaceDimention := globaldata.GetGlobalData().SpaceDimention
	field := datatypes.NewField(uint64(
		max(
			spaceDimention[base.AxisX].Higher,
			spaceDimention[base.AxisX].Lower,
			spaceDimention[base.AxisY].Higher,
			spaceDimention[base.AxisY].Lower,
			spaceDimention[base.AxisZ].Higher,
			spaceDimention[base.AxisZ].Lower)))
	polymers := make([]*datatypes.Polymer, alg.inputData.polymersCount)
	for i := range polymers {
		polymers[i] = datatypes.NewPolymer(field, int64(i))
	}
	direction := base.AxisYVec
	currPoint := base.AxisYVecReversed
	scanner := bufio.NewScanner(file)
	metElements := make(map[base.MendeleevTableElement]float64)
	allElementsCount := 0.0
	for scanner.Scan() && scanner.Err() == nil {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		for _, beadType := range line {
			allElementsCount += 1.0
			currPoint = base.AddVecF(currPoint, direction)
			mon := field.GetMonomerByCoords(currPoint)
			mon.MonomerType = base.MendeleevTableElement(string(beadType))
			metElements[mon.MonomerType] += 1
			polymers[0].AddMonomer(mon)
		}
		currPoint = base.AddVecF(base.AddVecF(currPoint, direction), base.AxisXVec)
		direction.MultiplyByConstantF(-1.0)
	}
	if alg.inputData.printStatistics {
		beautifulPrint(polymers, metElements)
	}
	return polymers
}

func beautifulPrint(polymers []*datatypes.Polymer, metElements map[base.MendeleevTableElement]float64) {
	allElementsCount := 0.0
	for _, p := range polymers {
		allElementsCount += float64(p.Len())
	}
	fmt.Println("Распределение по встреченным буквам:")
	fmt.Println("|-----------------------------------------|")
	fmt.Println("|Элемент|Количество|Процентное соотношение|")
	fmt.Println("|-------|----------|----------------------|")
	for element, count := range metElements {
		fmt.Printf("|%s|%s|%s|\n",
			fillToLength(string(element), utf8.RuneCountInString("Элемент")),
			fillToLength(strconv.Itoa(int(count)), utf8.RuneCountInString("Количество")),
			fillToLength(strconv.FormatFloat(count/allElementsCount*100.0, 'f', 2, 64)+"%", utf8.RuneCountInString("Процентное соотношение")),
		)
		fmt.Println("|-------|----------|----------------------|")
	}
	fmt.Printf("|Общее количество мономеров|%s|\n",
		fillToLength(strconv.Itoa(int(allElementsCount)), utf8.RuneCountInString("--------------")))
	fmt.Println("|-----------------------------------------|")
}

func fillToLength(s string, upToLength int) string {
	if len(s) >= upToLength {
		return s
	}
	newS := s
	for i := len(newS); i < upToLength; i++ {
		newS += " "
	}
	return newS
}
