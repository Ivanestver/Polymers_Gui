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
	Sc                             MendeleevTableElement = "Sc"
	Ti                             MendeleevTableElement = "Ti"
	V                              MendeleevTableElement = "V"
	Cr                             MendeleevTableElement = "Cr"
	Mn                             MendeleevTableElement = "Mn"
	Fe                             MendeleevTableElement = "Fe"
	Co                             MendeleevTableElement = "Co"
	Ni                             MendeleevTableElement = "Ni"
	Cu                             MendeleevTableElement = "Cu"
	Zn                             MendeleevTableElement = "Zn"
	Ga                             MendeleevTableElement = "Ga"
	Ge                             MendeleevTableElement = "Ge"
	As                             MendeleevTableElement = "As"
	Se                             MendeleevTableElement = "Se"
	Br                             MendeleevTableElement = "Br"
	Kr                             MendeleevTableElement = "Kr"
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
	table[14] = P
	table[15] = S
	table[16] = Cl
	table[17] = Ar
	table[18] = K
	table[19] = Ca
	table[20] = Sc
	table[21] = Ti
	table[22] = V
	table[23] = Cr
	table[24] = Mn
	table[25] = Fe
	table[26] = Co
	table[27] = Ni
	table[28] = Cu
	table[29] = Zn
	table[30] = Ga
	table[31] = Ge
	table[32] = As
	table[33] = Se
	table[34] = Br
	table[35] = Kr
	return table
}

func GetMendeleevTableReversed() map[MendeleevTableElement]int {
	m := make(map[MendeleevTableElement]int)
	table := GetMendeleevTable()
	for i, e := range table {
		m[e] = i
	}
	return m
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
