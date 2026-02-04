package loaders

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"polymers/base"
	"polymers/build_globula"
	"polymers/datatypes"
	"polymers/global_data"
	"polymers/views"
	"strconv"
	"strings"
)

type _LammpsMetadata struct {
	atomsCount     int
	atomTypesCount int
	bondsCount     int
	bondTypesCount int
	atomTypes      map[string]struct {
		Mass  int
		Label string
	}
	bondTypes map[string]struct {
		sth1 int
		sth2 float64
	}
	polymers map[int]*datatypes.Polymer
	atoms    map[int64]struct {
		Monomer       *datatypes.Monomer
		PolymerNumber int
	}
}

type _LammpsLoader struct {
	_LammpsMetadata
	builtGlobula *views.GlobulaView
	scanner      *bufio.Scanner
	filename     string
}

func (loader *_LammpsLoader) Load(filename, globulaname string) (*views.GlobulaView, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	loader.filename = filename

	loader.scanner = bufio.NewScanner(file)
	if err := loader.load(globulaname); err != nil {
		return nil, err
	}
	loader.scanner = nil
	return loader.builtGlobula, nil
}

func (loader *_LammpsLoader) load(globulaName string) error {
	if err := loader.loadMetadata(); err != nil {
		return err
	}
	if err := loader.loadMasses(); err != nil {
		return err
	}
	if err := loader.loadBondTypes(); err != nil {
		return err
	}
	if err := loader.loadAtoms(); err != nil {
		return err
	}
	if err := loader.loadBonds(); err != nil {
		return err
	}
	polymers := make([]*datatypes.Polymer, len(loader.polymers))
	for i := 0; i < len(polymers); i++ {
		polymers[i] = loader.polymers[i]
	}
	literalsTable := build_globula.BuildThreadAlgInputData{}.GetLiterals()
	loader.builtGlobula = views.NewGlobulaView(globulaName, polymers, views.GLOBULA_THREAD_TYPE, literalsTable)
	return nil
}

func (loader *_LammpsLoader) loadMetadata() error {
	for loader.scanner.Scan() {
		line := loader.scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] < '0' || line[0] > '9' {
			continue
		}
		// read atoms count section
		if count, err := getNumber(line); err == nil {
			loader.atomsCount = count
			break
		}
	}
	loader.atoms = make(map[int64]struct {
		Monomer       *datatypes.Monomer
		PolymerNumber int
	})

	// read atoms types section
	if err := writeMetadata(loader, &loader.atomTypesCount, "atom types"); err != nil {
		return err
	}
	loader.atomTypes = make(map[string]struct {
		Mass  int
		Label string
	})

	// read bonds section
	if err := writeMetadata(loader, &loader.bondsCount, "bonds counts"); err != nil {
		return err
	}
	loader.bondTypes = make(map[string]struct {
		sth1 int
		sth2 float64
	})

	// read bonds types section
	if err := writeMetadata(loader, &loader.bondTypesCount, "bonds types"); err != nil {
		return err
	}

	return nil
}

func getNumber(s string) (int, error) {
	if len(s) == 0 {
		return 0, errors.New("string is empty")
	}

	isNumber := func(b byte) bool { return '0' <= b && b <= '9' }

	var builder strings.Builder
	currentSymbol := 0
	for isNumber(s[currentSymbol]) && currentSymbol < len(s) {
		builder.WriteByte(s[currentSymbol])
		currentSymbol++
	}
	return strconv.Atoi(builder.String())
}

func writeMetadata(loader *_LammpsLoader, metadata *int, sectionName string) error {
	if loader.scanner.Scan() {
		line := loader.scanner.Text()
		if len(line) == 0 {
			return errors.New("Could not find " + sectionName + " section")
		}
		if line[0] < '0' || line[0] > '9' {
			return errors.New("Could not find " + sectionName + " section")
		}
		if count, err := getNumber(line); err == nil {
			*metadata = count
		}
	}
	return nil
}

func (loader *_LammpsLoader) loadMasses() error {
	for loader.scanner.Scan() {
		if loader.scanner.Text() == "Masses" {
			break
		}
	}
	loader.scanner.Scan()

	for atomTypeLineNumber := 0; atomTypeLineNumber < loader.atomTypesCount && loader.scanner.Scan(); atomTypeLineNumber++ {
		line := loader.scanner.Text()
		if len(line) == 0 {
			break
		}
		parts := strings.Split(line, " ")
		if len(parts) != 4 {
			return fmt.Errorf("wrong line in the Masses section (line number in there: %d)", atomTypeLineNumber+1)
		}
		number := parts[0]
		mass, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		label := parts[3]
		loader.atomTypes[number] = struct {
			Mass  int
			Label string
		}{
			Mass:  mass,
			Label: label,
		}
	}

	return nil
}

func (loader *_LammpsLoader) loadBondTypes() error {
	for loader.scanner.Scan() && loader.scanner.Text() != "Bond Coeffs # harmonic" {
	}
	loader.scanner.Scan()

	for bondTypeLineNumber := 0; bondTypeLineNumber < loader.bondTypesCount && loader.scanner.Scan(); bondTypeLineNumber++ {
		parts := strings.Split(loader.scanner.Text(), " ")
		if len(parts) != 3 {
			return fmt.Errorf("wrong line in the Bond Coeffs section (line number in there: %d)", bondTypeLineNumber+1)
		}
		number := parts[0]
		sth1, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		sth2, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return err
		}
		loader.bondTypes[number] = struct {
			sth1 int
			sth2 float64
		}{
			sth1: sth1,
			sth2: sth2,
		}
	}

	return nil
}

func (loader *_LammpsLoader) loadAtoms() error {
	for loader.scanner.Scan() && loader.scanner.Text() != "Atoms # full" {
	}
	loader.scanner.Scan()

	literals := build_globula.BuildThreadAlgInputData{}.GetLiterals()
	loader.polymers = make(map[int]*datatypes.Polymer)
	field := datatypes.NewField(uint64(global_data.GetGlobalData().SpaceDimention.X))
	for atomLineNumber := 0; atomLineNumber < loader.atomsCount && loader.scanner.Scan(); atomLineNumber++ {
		parts := strings.Split(loader.scanner.Text(), " ")
		if len(parts) != 10 {
			return fmt.Errorf("wrong line in the Atoms section (line number in there: %d)", atomLineNumber+1)
		}

		atomID, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return err
		}

		polymerID, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		polymerID--

		atomTypeNumber := parts[2]

		// Miss the charge because of needing it

		x, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			return err
		}

		y, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return err
		}

		z, err := strconv.ParseFloat(parts[6], 64)
		if err != nil {
			return err
		}

		monomer :=
			field.GetMonomerByCoords(base.Vector3DF{
				X: x,
				Y: y,
				Z: z,
			})
		monomer.Number = atomID
		monomer.MonomerType = getMonomerTypeByLiteral(literals, loader.atomTypes[atomTypeNumber].Label)

		loader.atoms[atomID] = struct {
			Monomer       *datatypes.Monomer
			PolymerNumber int
		}{Monomer: monomer, PolymerNumber: polymerID}
	}

	for i := 1; i <= len(loader.atoms); i++ {
		atom, ok := loader.atoms[int64(i)]
		if !ok {
			return fmt.Errorf("missing line in the Atoms section: %d", i)
		}
		polymer, ok := loader.polymers[atom.PolymerNumber]
		if !ok {
			polymer = datatypes.NewPolymer(field, int64(atom.PolymerNumber))
			loader.polymers[atom.PolymerNumber] = polymer
		}
		polymer.AddMonomer(atom.Monomer)
	}
	return nil
}

func getMonomerTypeByLiteral(literals map[datatypes.MonomerType]string, literal string) datatypes.MonomerType {
	for key, value := range literals {
		if value == literal {
			return key
		}
	}
	return datatypes.MONOMER_TYPE_UNDEFINED
}

func (loader *_LammpsLoader) loadBonds() error {
	for loader.scanner.Scan() && loader.scanner.Text() != "Bonds" {
	}
	loader.scanner.Scan()

	for atomLineNumber := 0; atomLineNumber < loader.atomsCount && loader.scanner.Scan(); atomLineNumber++ {
		parts := strings.Split(loader.scanner.Text(), " ")
		if len(parts) != 4 {
			return fmt.Errorf("wrong line in the Bonds section (line number in there: %d)", atomLineNumber+1)
		}

		connectionType, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		firstAtomNumber, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return err
		}

		secondAtomNumber, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return err
		}

		firstAtom := loader.atoms[firstAtomNumber]
		secondAtom := loader.atoms[secondAtomNumber]

		if err := datatypes.MakeConnection(firstAtom.Monomer, secondAtom.Monomer, datatypes.ConnectionType(connectionType)); err != nil {
			return err
		}
	}
	return nil
}
