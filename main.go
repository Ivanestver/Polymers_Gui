package main

import (
	"bufio"
	dt "data_visualizer/datatypes"
	"data_visualizer/visualizers"
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

func readAtomsSection(scanner *bufio.Scanner) ([]dt.Atom, error) {
	atoms := make([]dt.Atom, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			return atoms, nil
		}

		parts := strings.Split(line, " ")
		if len(parts) != ATOMS_INFO_COUNT {
			return []dt.Atom{}, errors.New("(In Atoms): '" + line + "' was not apropriately been built")
		}

		number, err := strconv.Atoi(parts[ATOM_NUMBER_POSITION])
		if err != nil {
			return []dt.Atom{}, errors.New("(In Atoms): " + err.Error())
		}
		atomType, err := strconv.Atoi(parts[ATOM_TYPE_POSITION])
		if err != nil {
			return []dt.Atom{}, errors.New("(In Atoms): " + err.Error())
		}
		pointX, err := strconv.Atoi(parts[ATOM_POINT_X])
		if err != nil {
			return []dt.Atom{}, errors.New("(In Atoms): " + err.Error())
		}
		pointY, err := strconv.Atoi(parts[ATOM_POINT_Y])
		if err != nil {
			return []dt.Atom{}, errors.New("(In Atoms): " + err.Error())
		}
		pointZ, err := strconv.Atoi(parts[ATOM_POINT_Z])
		if err != nil {
			return []dt.Atom{}, errors.New("(In Atoms): " + err.Error())
		}

		atoms = append(atoms, dt.Atom{
			M_Number:   number - 1,
			M_AtomType: atomType,
			M_Point: dt.Point{
				M_X: pointX,
				M_Y: pointY,
				M_Z: pointZ,
			},
		})
	}

	return []dt.Atom{}, errors.New("(In Atoms): Unexpected end of file")
}

func readBondsSection(scanner *bufio.Scanner) (dt.Bonds, error) {
	bonds := make(dt.Bonds)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			return bonds, nil
		}
		parts := strings.Split(line, " ")
		if len(parts) != BOND_INFO_COUNT {
			return dt.Bonds{}, errors.New("(In Bonds): '" + line + "' was not apropriately been built")
		}

		sourceAtomNumber, err := strconv.Atoi(parts[BOND_SOURCE_ATOM])
		if err != nil {
			return dt.Bonds{}, errors.New("(In Bonds): " + err.Error())
		}
		sourceAtomNumber--

		destinationAtomNumber, err := strconv.Atoi(parts[BOND_DESTINATION_ATOM])
		if err != nil {
			return dt.Bonds{}, errors.New("(In Bonds): " + err.Error())
		}
		destinationAtomNumber--

		bonds[sourceAtomNumber] = append(bonds[sourceAtomNumber], destinationAtomNumber)
		bonds[destinationAtomNumber] = append(bonds[destinationAtomNumber], sourceAtomNumber)
	}

	return bonds, nil
}

func extractInfoFromFile(f *os.File) ([]dt.Atom, dt.Bonds, error) {
	scanner := bufio.NewScanner(f)
	// Read until the Atoms section
	readUntilSection(scanner, "Atoms")

	atoms, err := readAtomsSection(scanner)
	if err != nil {
		return atoms, dt.Bonds{}, err
	}

	readUntilSection(scanner, "Bonds")

	bonds, err := readBondsSection(scanner)
	return atoms, bonds, err
}

func main() {
	// Get the file name specified
	inputFileName, err := getInputFileName()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Open the specified file
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Parse the file and extract the atoms and bonds
	atoms, bonds, err := extractInfoFromFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Data has successfully been loaded.")
	for {
		fmt.Println("Available visualizers:")
		visualizersDescriptions := visualizers.GetAllVisualizersDescription()
		for _, desc := range visualizersDescriptions {
			fmt.Println("\t" + desc)
		}
		fmt.Println("Please, specify the visualizer:")
		choice := 0
		_, err := fmt.Scanf("%i", &choice)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		visualizer := visualizers.NewVisualizer(choice)
		visualizer.Visualize(atoms, bonds)
	}
}
