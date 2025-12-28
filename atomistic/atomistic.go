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

type _Config struct {
	Scale         float64
	Substitutions map[string]*_Subtitution
	SaveFile      string
}

func _NewConfig() *_Config {
	return &_Config{
		Scale:         0.0,
		Substitutions: make(map[string]*_Subtitution),
	}
}

func MakeAtomistic(globula *views.GlobulaView, configFile string) {
	printer = output_format.GetPrint()
	config, err := createConfig(configFile)
	if err != nil {
		printer.PrintflnError("When atomistic: %s", err.Error())
		return
	}
	// Retrieve the polymer from globula
	polymer := getPolymer(globula)
	// Resize it to an apropriate size
	resizePolymerByScale(polymer, config.Scale)
	// Place molecules into their places
	placeMolecules(polymer, config)
	// Fill ends of the polymer
	//fillEnds(polymer)
	// Make connections between molecules
	connectMonomers(polymer)
	// Save it into the file
	savePolymer(polymer, config)
}

func createConfig(configFileName string) (*_Config, error) {
	file, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	config := _NewConfig()
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			printer.PrintlnWarning("Line is empty")
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			printer.PrintlnWarning("The line has wrong format: " + line)
			continue
		}

		switch parts[0] {
		case "scale":
			writeScale(config, parts[1:], line)
		case "monomer":
			writeMonomer(config, parts[1:], line)
		case "save":
			writeSave(config, parts[1:], line)
		}
	}
	return config, nil
}

func writeScale(config *_Config, parts []string, line string) {
	if scale, err := strconv.ParseFloat(parts[0], 64); err == nil {
		config.Scale = scale
	} else {
		printer.PrintflnWarning("The following error occured: %s (in line: %s)", err.Error(), line)
	}
}

func writeMonomer(config *_Config, parts []string, line string) {
	label := parts[0]
	variants, err := strconv.Atoi(parts[1])
	if err != nil {
		printer.PrintflnWarning("The error occured: %s (in line '%s')", err.Error(), line)
		return
	}
	if variants < 1 {
		printer.PrintflnWarning("The number of variants must be not less than 1 (in line %s)", line)
		return
	}
	substitution := make(_Subtitution, variants)
	for i := 0; i < variants; i++ {
		substitution[i] = getMoleculeFromFile(parts[2+i])
	}
	config.Substitutions[label] = &substitution
}

func writeSave(config *_Config, parts []string, line string) {
	if len(parts) == 0 {
		printer.PrintflnError("The file name must be specified (in line: %s)", line)
	}
	config.SaveFile = parts[0]
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
		if isHead(fields) {
			molecule.Head = atom
		} else if isTail(fields) {
			molecule.Tail = atom
		} else {
			continue
		}
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

func isHead(fields []string) bool {
	return isEnd(fields, 1)
}

func isTail(fields []string) bool {
	return isEnd(fields, 2)
}

func isEnd(fields []string, markExpected int) bool {
	markReal := getEndMark(fields)
	return markExpected == markReal
}

func getEndMark(fields []string) int {
	mark := string(fields[len(fields)-1][0])
	if number, err := strconv.ParseInt(mark, 10, 32); err == nil {
		return int(number)
	} else {
		printer.PrintflnError("Could not parse the end mark as int (%s) due to the error: %s", mark, err.Error())
		return -1
	}
}

func fillBondsInfo(monomer *_Monomer, scanner *bufio.Scanner, bondsCount int) error {
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

		part, ok := monomer.Bonds[startMonomerNumber]
		if !ok {
			monomer.Bonds[startMonomerNumber] = make(map[_AtomNumber]_BondValence)
			part = monomer.Bonds[startMonomerNumber]
		}
		part[endMonomerNumber] = bondType
	}

	return nil
}

func placeMolecules(polymer *_Polymer, config *_Config) {
	const (
		isFirst = iota
		isSecond
		isThird
		COUNT
	)
	prevMonomerState := isThird
	for i := 0; i < len(polymer.Monomers); i++ {
		label := polymer.Monomers[i].Atoms[0].Label
		prototype := config.Substitutions[label]
		var molecule *_Monomer
		molecule = prototype.GetNext(prevMonomerState).Copy()
		prevMonomerState = (prevMonomerState + 1) % COUNT
		molecule.MoveTo(&polymer.Monomers[i].Atoms[0].Coords)
		polymer.Monomers[i] = molecule
	}

	reNumberAtoms(polymer)
}

func reNumberAtoms(polymer *_Polymer) {
	atomNumber := 1
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

func savePolymer(polymer *_Polymer, config *_Config) {
	bytesData := []byte(makeFileContent(polymer))
	var filename string
	if len(config.SaveFile) > 0 {
		filename = config.SaveFile
	} else {
		filename = "output_data.mol2"
	}
	if err := os.WriteFile(filename, bytesData, 0644); err != nil {
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

		if _, ok := polymer.Bonds[left.Head.Number]; !ok {
			polymer.Bonds[left.Head.Number] = make(map[_AtomNumber]_BondValence)
		}
		polymer.Bonds[left.Head.Number][right.Tail.Number] = 1
	}
}

func fillEnds(polymer *_Polymer) {
	fillStartingMonomer(polymer)
	fillTerminatingMonomer(polymer)
}

func fillStartingMonomer(polymer *_Polymer) {
	startingMonomer := &_Monomer{
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
	}
	startingMonomer.Head = &startingMonomer.Atoms[0]
	startingMonomer.Tail = &startingMonomer.Atoms[0]
	polymer.Monomers = append([]*_Monomer{startingMonomer}, polymer.Monomers...)
}

func fillTerminatingMonomer(polymer *_Polymer) {
	terminatingMonomer := &_Monomer{
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
	}
	terminatingMonomer.Head = &terminatingMonomer.Atoms[0]
	terminatingMonomer.Tail = &terminatingMonomer.Atoms[0]
	polymer.Monomers = append(polymer.Monomers, terminatingMonomer)
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
