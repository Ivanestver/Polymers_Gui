package base

import "strings"

type MendeleevTableElement string

const (
	MendeleevTableElementUndefined MendeleevTableElement = ""
	Hidrogen                       MendeleevTableElement = "H"
	Carbon                         MendeleevTableElement = "C"
	Oxygen                         MendeleevTableElement = "O"
	Nitrogen                       MendeleevTableElement = "N"
)

func GetMendeleevTable() [118]MendeleevTableElement {
	table := [118]MendeleevTableElement{}
	table[0] = Hidrogen
	table[5] = Carbon
	table[6] = Nitrogen
	table[7] = Oxygen
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
		if len(s) == len(elem) && MendeleevTableElement(s) == elem {
			return elem
		}
	}
	return MendeleevTableElementUndefined
}
