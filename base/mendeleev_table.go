package base

import (
	"strings"
)

type MendeleevTableElement string

const (
	MendeleevTableElementUndefined MendeleevTableElement = ""
	H                              MendeleevTableElement = "H"
	He                             MendeleevTableElement = "He"
	Li                             MendeleevTableElement = "Li"
	Be                             MendeleevTableElement = "Be"
	B                              MendeleevTableElement = "B"
	C                              MendeleevTableElement = "C"
	N                              MendeleevTableElement = "N"
	O                              MendeleevTableElement = "O"
	F                              MendeleevTableElement = "F"
	Ne                             MendeleevTableElement = "Ne"
	Na                             MendeleevTableElement = "Na"
	Mg                             MendeleevTableElement = "Mg"
	Al                             MendeleevTableElement = "Al"
	Si                             MendeleevTableElement = "Si"
	P                              MendeleevTableElement = "P"
	S                              MendeleevTableElement = "S"
	Cl                             MendeleevTableElement = "Cl"
	Ar                             MendeleevTableElement = "Ar"
	K                              MendeleevTableElement = "K"
	Ca                             MendeleevTableElement = "Ca"
)

func GetMendeleevTable() [118]MendeleevTableElement {
	table := [118]MendeleevTableElement{}
	table[0] = H
	table[1] = He
	table[2] = Li
	table[3] = Be
	table[4] = B
	table[5] = C
	table[6] = N
	table[7] = O
	table[8] = F
	table[9] = Ne
	table[10] = Na
	table[11] = Mg
	table[12] = Al
	table[13] = Si
	table[14] = B
	table[15] = P
	table[16] = S
	table[17] = Cl
	table[18] = Ar
	table[19] = K
	table[20] = Ca
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
		elemUpper := elem
		if elemUpper != MendeleevTableElementUndefined && strings.HasPrefix(s, strings.ToUpper(string(elemUpper))) {
			return elem
		}
	}
	return MendeleevTableElementUndefined
}
