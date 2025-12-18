package atomistic

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/output_format"
	"polymers/views"
	"strconv"
	"strings"
	"unicode"
)

var printer output_format.IPrint

func MakeAtomistic(globula *views.GlobulaView) {
	// Retrieve the polymer from globula
	polymer := getPolymer(globula)
	// Resize it to an apropriate size
	resizePolymerByScale(polymer, 2.9)
	// Place molecules into their places
	placeMolecules(polymer)
	// Make connections between molecules
	connectMonomers(polymer)
	// Fill ends of the polymer
	fillEnds(polymer)
	// Save it into the file
	savePolymer(polymer)
}

func getPolymer(globula *views.GlobulaView) *_Polymer {
	pattern := NewPolymer()
	views.ForEachPolymer_If(globula, func(pv *views.PolymerView) bool {
		views.ForEachMonomer(pv, func(m *datatypes.Monomer) bool {
			coords := m.Coords()
			pattern.Monomers = append(pattern.Monomers, &_Monomer{
				Atoms: []_Atom{
					{
						Label:  globula.GetLiteral(m.MonomerType),
						Coords: base.Vector3D_To_Vector3DF(&coords),
					},
				},
			})
			return true
		})
		return false
	})
	return pattern
}

func resizePolymerByScale(pattern *_Polymer, scale float64) {
	pivot := &pattern.Monomers[0].Atoms[0].Coords
	for i := 1; i < len(pattern.Monomers); i++ {
		firstAtom := &pattern.Monomers[i].Atoms[0]
		direction := base.SubtractVecF(&firstAtom.Coords, pivot)
		direction.MultiplyByConstantF(scale)
		firstAtom.Coords = *base.AddVecF(pivot, direction)
	}
}

func getMoleculeFromFile(filePath string) *_Monomer {
	file, err := os.Open(filePath)
	printer = output_format.GetPrint()
	if err != nil {
		printer.PrintlnError(err.Error())
		return nil
	}
	defer file.Close()

	// Skip all the lines before molecule info
	scanner := bufio.NewScanner(file)
	for scanner.Scan() && scanner.Text() != "@<TRIPOS>MOLECULE" {
	}
	if !scanner.Scan() {
		printer.PrintlnError("Wrong structure")
		return nil
	}

	molecule := &_Monomer{}
	// 1. Fill Molecule info
	bondsCount, err := fillMoleculeData(molecule, scanner)
	if err != nil {
		printer.PrintlnError(err.Error())
		return nil
	}

	// Skip until ATOMS section
	for scanner.Scan() && scanner.Text() != "@<TRIPOS>ATOM" {
	}
	// 2. Fill Atoms info
	if err := fillAtomsInfo(molecule, scanner); err != nil {
		printer.PrintlnError(err.Error())
		return nil
	}

	// Skip until BONDS section
	for scanner.Scan() && scanner.Text() != "@<TRIPOS>BOND" {
	}
	// 3. Fill Bonds info
	if err := fillBondsInfo(molecule, scanner, bondsCount); err != nil {
		printer.PrintlnError(err.Error())
		return nil
	}
	return molecule
}

func fillMoleculeData(molecule *_Monomer, scanner *bufio.Scanner) (int, error) {
	// Save the molecule name
	molecule.Name = scanner.Text()

	scanner.Scan()
	fields := strings.Fields(scanner.Text())
	atomsCount, err := strconv.Atoi(fields[0])
	if err != nil {
		printer.PrintlnError(err.Error())
		return 0, nil
	}

	bondsCount, err := strconv.Atoi(fields[1])
	if err != nil {
		printer.PrintlnError(err.Error())
		return 0, nil
	}

	molecule.Atoms = make([]_Atom, atomsCount)
	molecule.Bonds = make(map[_AtomNumber]map[_AtomNumber]_BondValence)
	return int(bondsCount), nil
}

func fillAtomsInfo(molecule *_Monomer, scanner *bufio.Scanner) error {
	// Assume we're at the first line of ATOMS section
	for i := 0; i < len(molecule.Atoms) && scanner.Scan(); i++ {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 9 {
			return errors.New("wrong ATOMS section line format")
		}
		atom := &molecule.Atoms[i]
		atom.Number, _ = strconv.Atoi(fields[0])
		atom.Label = getLabel(fields)
		atom.Coords = getCoords(fields)
	}
	return nil
}

func getLabel(fields []string) string {
	var builder strings.Builder
	for _, c := range fields[1] {
		if !unicode.IsLetter(c) {
			break
		}
		builder.WriteRune(c)
	}
	return builder.String()
}

func getCoords(fields []string) base.Vector3DF {
	x, _ := strconv.ParseFloat(fields[2], 64)
	y, _ := strconv.ParseFloat(fields[3], 64)
	z, _ := strconv.ParseFloat(fields[4], 64)
	return base.Vector3DF{
		X: x,
		Y: y,
		Z: z,
	}
}

func fillBondsInfo(molecule *_Monomer, scanner *bufio.Scanner, bondsCount int) error {
	// Assume we're at the first line of BONDS section
	for i := 0; i < bondsCount && scanner.Scan(); i++ {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 4 {
			return errors.New("wrong BONDS section line format")
		}
		startMonomerNumber, err := strconv.Atoi(fields[1])
		if err != nil {
			return err
		}
		endMonomerNumber, err := strconv.Atoi(fields[2])
		if err != nil {
			return err
		}
		bondType, err := strconv.Atoi(fields[3])
		if err != nil {
			return err
		}

		part, ok := molecule.Bonds[startMonomerNumber]
		if !ok {
			molecule.Bonds[startMonomerNumber] = make(map[_AtomNumber]_BondValence)
			part = molecule.Bonds[startMonomerNumber]
		}
		part[endMonomerNumber] = bondType
	}

	return nil
}

func placeMolecules(polymer *_Polymer) {
	labelToPrototypeMap := make(map[string]*_Monomer)
	labelToPrototypeMap["O"] = getMoleculeFromFile("C2ch.mol2")
	labelToPrototypeMap["N"] = getMoleculeFromFile("C2Och.mol2")
	labelToPrototypeMap["C"] = getMoleculeFromFile("C2Ach.mol2")
	for i := 0; i < len(polymer.Monomers); i++ {
		prototype := labelToPrototypeMap[polymer.Monomers[i].Atoms[0].Label]
		molecule := prototype.Copy()
		molecule.MoveTo(&polymer.Monomers[i].Atoms[0].Coords)
		polymer.Monomers[i] = molecule
	}

	reNumberAtoms(polymer)
}

func reNumberAtoms(polymer *_Polymer) {
	atomNumber := 2
	for _, monomer := range polymer.Monomers {
		replacementMap := make(map[_AtomNumber]_AtomNumber)
		for i := 0; i < len(monomer.Atoms); i++ {
			atom := &monomer.Atoms[i]
			replacementMap[atom.Number] = atomNumber
			atom.Number = atomNumber
			atomNumber++
		}

		newBonds := make(map[_AtomNumber]map[_AtomNumber]_BondValence)
		for startAtomNumber := range monomer.Bonds {
			newBonds[replacementMap[startAtomNumber]] = make(map[_AtomNumber]_BondValence)
			m := newBonds[replacementMap[startAtomNumber]]
			for endAtomNumber, valence := range monomer.Bonds[startAtomNumber] {
				m[replacementMap[endAtomNumber]] = valence
			}
		}
		monomer.Bonds = newBonds
	}
}

func savePolymer(polymer *_Polymer) {
	bytesData := []byte(makeFileContent(polymer))
	if err := os.WriteFile("output_data.mol2", bytesData, 0644); err != nil {
		printer.PrintlnError(err.Error())
	}
}

func makeFileContent(polymer *_Polymer) string {
	var builder strings.Builder
	builder.WriteString("@<TRIPOS>MOLECULE\n")
	builder.WriteString("Test\n")
	builder.WriteString(fmt.Sprintf("%d %d 0 0 0\n", getAtomsCount(polymer), getBondsCount(polymer)))
	builder.WriteString("SMALL\n")
	builder.WriteString("USER_CHARGES\n")
	builder.WriteString("\n")
	builder.WriteString("\n")

	var builderBonds strings.Builder
	builder.WriteString("@<TRIPOS>ATOM\n")
	bondNumber := 1
	for _, monomer := range polymer.Monomers {
		for _, atom := range monomer.Atoms {
			builder.WriteString(turnAtomToString(&atom))
		}
		saveBonds(&monomer.Bonds, &builderBonds, &bondNumber)
	}
	saveBonds(&polymer.Bonds, &builderBonds, &bondNumber)

	builder.WriteString("@<TRIPOS>BOND\n")
	builder.WriteString(builderBonds.String())
	return builder.String()
}

func getAtomsCount(polymer *_Polymer) int {
	atomsCount := 0
	for _, molecule := range polymer.Monomers {
		atomsCount += len(molecule.Atoms)
	}
	return atomsCount
}

func getBondsCount(polymer *_Polymer) int {
	bondsCount := 0
	for _, monomer := range polymer.Monomers {
		bondsCount += monomer.GetBondsCount()
	}
	return bondsCount + len(polymer.Bonds)
}

func turnAtomToString(atom *_Atom) string {
	return fmt.Sprintf("%d %s %f %f %f %s 0 ***** 0\n", atom.Number, atom.Label, atom.Coords.X, atom.Coords.Y, atom.Coords.Z, atom.Label)
}

func saveBonds(bonds *map[_AtomNumber]map[_AtomNumber]_BondValence, builderBonds *strings.Builder, bondNumber *int) {
	for startMonomerNumber, m := range *bonds {
		for endMonomerNumber, valence := range m {
			builderBonds.WriteString(
				fmt.Sprintf(
					"%d %d %d %d\n",
					*bondNumber,
					startMonomerNumber,
					endMonomerNumber,
					valence))
			(*bondNumber)++
		}
	}
}

func connectMonomers(polymer *_Polymer) {
	if len(polymer.Monomers) < 2 {
		return
	}

	for i := 0; i < len(polymer.Monomers)-1; i++ {
		left := polymer.Monomers[i]
		right := polymer.Monomers[i+1]

		_, nextLeft := left.GetConnectionAtoms()
		prevRight, _ := right.GetConnectionAtoms()
		if _, ok := polymer.Bonds[nextLeft.Number]; !ok {
			polymer.Bonds[nextLeft.Number] = make(map[_AtomNumber]_BondValence)
		}
		polymer.Bonds[nextLeft.Number][prevRight.Number] = 1
	}
}

func fillEnds(polymer *_Polymer) {
	fillStartingMonomer(polymer)
	fillTerminatingMonomer(polymer)
}

func fillStartingMonomer(polymer *_Polymer) {
	polymer.Monomers = append([]*_Monomer{
		{
			Name: "Starting",
			Atoms: []_Atom{
				{
					Number:  1,
					Label:   "H",
					Mass:    1,
					Charge:  0,
					Valence: 1,
					Coords:  getStartingMonomerCoords(polymer),
				},
			},
		},
	}, polymer.Monomers...)

	if prev, _ := polymer.Monomers[1].GetConnectionAtoms(); prev != nil {
		polymer.Bonds[1] = make(map[_AtomNumber]_BondValence)
		polymer.Bonds[1][prev.Number] = 1
	}
}

func fillTerminatingMonomer(polymer *_Polymer) {
	polymer.Monomers = append(polymer.Monomers,
		&_Monomer{
			Name: "Terminating",
			Atoms: []_Atom{
				{
					Number:  getAtomsCount(polymer) + 1,
					Label:   "H",
					Mass:    1,
					Charge:  0,
					Valence: 1,
					Coords:  getTerminatingMonomerCoords(polymer),
				},
			},
		},
	)

	lastMonomerNumber := len(polymer.Monomers) - 1
	if _, next := polymer.Monomers[lastMonomerNumber-1].GetConnectionAtoms(); next != nil {
		polymer.Bonds[next.Number] = make(map[_AtomNumber]_BondValence)
		polymer.Bonds[next.Number][polymer.Monomers[lastMonomerNumber].Atoms[0].Number] = 1
	}
}

func getStartingMonomerCoords(polymer *_Polymer) base.Vector3DF {
	firstMonomerMassCenter := polymer.Monomers[0].GetMassCenter()
	secondMonomerMassCenter := polymer.Monomers[1].GetMassCenter()

	direction := base.SubtractVecF(&firstMonomerMassCenter, &secondMonomerMassCenter)

	firstMonomerMassCenter.AddF(direction)
	return firstMonomerMassCenter
}

func getTerminatingMonomerCoords(polymer *_Polymer) base.Vector3DF {
	monomersCount := len(polymer.Monomers)
	prelastMonomerMassCenter := polymer.Monomers[monomersCount-2].GetMassCenter()
	lastMonomerMassCenter := polymer.Monomers[monomersCount-1].GetMassCenter()

	direction := base.SubtractVecF(&lastMonomerMassCenter, &prelastMonomerMassCenter)

	lastMonomerMassCenter.AddF(direction)
	return lastMonomerMassCenter
}
