package base

import (
	"strings"
)

type MendeleevTableElement string

const (
	MendeleevTableElementUndefined MendeleevTableElement = ""
	Hydrogen                       MendeleevTableElement = "H"
	Carbon                         MendeleevTableElement = "C"
	Oxygen                         MendeleevTableElement = "O"
	Nitrogen                       MendeleevTableElement = "N"
	Fluorine                       MendeleevTableElement = "F"
)

func GetMendeleevTable() [118]MendeleevTableElement {
	table := [118]MendeleevTableElement{}
	table[0] = Hydrogen
	table[5] = Carbon
	table[6] = Nitrogen
	table[7] = Oxygen
	table[8] = Fluorine
	return table
}

func RecognizeElement(str string) MendeleevTableElement {
	s := ""
	if len(str) == 0 {
		return MendeleevTableElementUndefined
	} else if len(str) == 1 {
		s = str
	} else {
		s = str[:2]
	}
	s = strings.ToUpper(s)
	for _, elem := range GetMendeleevTable() {
		if elem != MendeleevTableElementUndefined && strings.HasPrefix(s, string(elem)) {
			return elem
		}
	}
	return MendeleevTableElementUndefined
}
