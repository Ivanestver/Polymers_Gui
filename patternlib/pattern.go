package patternlib

import (
	"bufio"
	"os"
	"polymers/datatypes"
	"polymers/outputformat"
	"polymers/views"
	"slices"
)

func GetPattern(data map[string]string) (string, bool) {
	var pattern string
	printer := outputformat.GetPrint()
	fileName, ok := data["fileName"]
	if ok {
		file, err := os.Open(fileName)
		if err != nil {
			printer.PrintlnError(err.Error())
			return "", false
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				printer.PrintlnError(closeErr.Error())
			}
		}()
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			pattern = scanner.Text()
		} else {
			printer.PrintlnError("The given file does not contain a pattern")
			return "", false
		}
	} else if p, ok := data["pattern"]; ok {
		pattern = p
	} else {
		printer.PrintflnError("See usage of this command")
		return "", false
	}
	return pattern, true
}

func AnyLetterIsUndefined(pattern string, globula *views.GlobulaView) bool {
	printer := outputformat.GetPrint()
	for _, letter := range pattern {
		l := string(letter)
		if globula.GetMonomerTypeByLiteral(l) == datatypes.MonomerTypeUndefined {
			printer.PrintlnError(l + " does not have its decryption")
			return true
		}
	}
	return false
}

func ApplyAsGlobula(globula *views.GlobulaView, pattern string) {
	views.ForEachPolymer(globula, func(pv *views.PolymerView) {
		currentLetterNumber := 0
		views.ForEachMonomer(pv, func(m *datatypes.Monomer) bool {
			m.MonomerType = globula.GetMonomerTypeByLiteral(string(pattern[currentLetterNumber]))
			currentLetterNumber = (currentLetterNumber + 1) % len(pattern)
			return true
		})
	})
}

func ApplyAsThread(globula *views.GlobulaView, pattern string) {
	literalsTable := globula.GetLiterals()
	literals := make([]datatypes.MonomerType, 0)
	for m := range *literalsTable {
		literals = append(literals, m)
	}
	slices.SortFunc(literals, func(a, b datatypes.MonomerType) int {
		if int(a) < int(b) {
			return -1
		} else if int(a) == int(b) {
			return 0
		} else {
			return 1
		}
	})
	views.ForEachPolymer(globula, func(pv *views.PolymerView) {
		currentLetterNumber := 0
		views.ForEachMonomer(pv, func(m *datatypes.Monomer) bool {
			for _, monType := range literals {
				if string(pattern[currentLetterNumber]) == (*literalsTable)[monType] {
					m.MonomerType = monType
					currentLetterNumber = (currentLetterNumber + 1) % len(pattern)
					return true
				}
			}
			return false
		})
	})

}
