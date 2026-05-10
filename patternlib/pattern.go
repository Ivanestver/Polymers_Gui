package patternlib

import (
	"bufio"
	"os"
	"polymers/base"
	"polymers/outputformat"
	"polymers/views"
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
	hasMonomerWithElement := func(l string) bool {
		polymersCount := globula.Len()
		for i := 0; i < polymersCount; i++ {
			polymer := globula.GetPolymerByIdx(i)
			for j := 0; j < polymer.Len(); j++ {
				monomer := polymer.GetMonomerByIdx(j)
				if string(l) == string(monomer.MonomerType) {
					return true
				}
			}
		}
		return false
	}
	for _, l := range pattern {
		if !hasMonomerWithElement(string(l)) {
			printer.PrintflnError("%s does not have its decryption", string(l))
			return true
		}
	}
	return false
}

func ApplyAsGlobula(globula *views.GlobulaView, pattern string) {
	for polymerNumber := 0; polymerNumber < globula.Len(); polymerNumber++ {
		currentLetterNumber := 0
		polymer := globula.GetPolymerByIdx(polymerNumber)
		for monomerNumber := 0; monomerNumber < polymer.Len(); monomerNumber++ {
			m := polymer.GetMonomerByIdx(monomerNumber)
			m.MonomerType = base.RecognizeElement(string(pattern[currentLetterNumber]))
			currentLetterNumber = (currentLetterNumber + 1) % len(pattern)
		}
	}
}
