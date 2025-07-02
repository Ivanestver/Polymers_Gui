package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getInputFileName() (string, error) {
	if len(os.Args) < 3 {
		return "", errors.New("Not enough parameters. Minimum 3 required")
	}
	inputFileNamePosition := -1
	for i := 0; i < len(os.Args); i++ {
		if strings.EqualFold(os.Args[i], "--ifile") {
			inputFileNamePosition = i + 1
		}
	}
	if inputFileNamePosition >= len(os.Args) {
		return "", errors.New("The --ifile requires a file .data specified")
	}
	return os.Args[inputFileNamePosition], nil
}

type Point struct {
	m_X int
	m_Y int
	m_Z int
}

type Atom struct {
	m_Number   int
	m_AtomType int
	m_Point    Point
}

type Bonds map[int][]int

const (
	ATOMS_INFO_COUNT     = 10 // see the Atoms section structure in a .data file
	ATOM_NUMBER_POSITION = 0
	ATOM_TYPE_POSITION   = 2
	ATOM_POINT_X         = 4
	ATOM_POINT_Y         = 5
	ATOM_POINT_Z         = 6

	BOND_INFO_COUNT       = 4
	BOND_SOURCE_ATOM      = 2
	BOND_DESTINATION_ATOM = 3
)

func readUntilSection(scanner *bufio.Scanner, sectionName string) {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, sectionName) {
			scanner.Scan()
			return
		}
	}
}

func readAtomsSection(scanner *bufio.Scanner) ([]Atom, error) {
	atoms := make([]Atom, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			return atoms, nil
		}

		parts := strings.Split(line, " ")
		if len(parts) != ATOMS_INFO_COUNT {
			return []Atom{}, errors.New("(In Atoms): '" + line + "' was not apropriately been built")
		}

		number, err := strconv.Atoi(parts[ATOM_NUMBER_POSITION])
		if err != nil {
			return []Atom{}, errors.New("(In Atoms): " + err.Error())
		}
		atomType, err := strconv.Atoi(parts[ATOM_TYPE_POSITION])
		if err != nil {
			return []Atom{}, errors.New("(In Atoms): " + err.Error())
		}
		pointX, err := strconv.Atoi(parts[ATOM_POINT_X])
		if err != nil {
			return []Atom{}, errors.New("(In Atoms): " + err.Error())
		}
		pointY, err := strconv.Atoi(parts[ATOM_POINT_Y])
		if err != nil {
			return []Atom{}, errors.New("(In Atoms): " + err.Error())
		}
		pointZ, err := strconv.Atoi(parts[ATOM_POINT_Z])
		if err != nil {
			return []Atom{}, errors.New("(In Atoms): " + err.Error())
		}

		atoms = append(atoms, Atom{
			m_Number:   number,
			m_AtomType: atomType,
			m_Point: Point{
				m_X: pointX,
				m_Y: pointY,
				m_Z: pointZ,
			},
		})
	}

	return []Atom{}, errors.New("(In Atoms): Unexpected end of file")
}

func readBondsSection(scanner *bufio.Scanner) (Bonds, error) {
	bonds := make(Bonds)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			return bonds, nil
		}
		parts := strings.Split(line, " ")
		if len(parts) != BOND_INFO_COUNT {
			return Bonds{}, errors.New("(In Bonds): '" + line + "' was not apropriately been built")
		}

		sourceAtomNumber, err := strconv.Atoi(parts[BOND_SOURCE_ATOM])
		if err != nil {
			return Bonds{}, errors.New("(In Bonds): " + err.Error())
		}
		destinationAtomNumber, err := strconv.Atoi(parts[BOND_SOURCE_ATOM])
		if err != nil {
			return Bonds{}, errors.New("(In Bonds): " + err.Error())
		}
		bonds[sourceAtomNumber] = append(bonds[sourceAtomNumber], destinationAtomNumber)
	}

	return bonds, nil
}

func extractInfoFromFile(f *os.File) ([]Atom, Bonds, error) {
	scanner := bufio.NewScanner(f)
	// Read until the Atoms section
	readUntilSection(scanner, "Atoms")

	atoms, err := readAtomsSection(scanner)
	if err != nil {
		return atoms, Bonds{}, err
	}

	readUntilSection(scanner, "Bonds")

	bonds, err := readBondsSection(scanner)
	return atoms, bonds, err
}

func main() {
	inputFileName, err := getInputFileName()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	atoms, bonds, err := extractInfoFromFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(len(atoms))
	fmt.Println(len(bonds))
}
