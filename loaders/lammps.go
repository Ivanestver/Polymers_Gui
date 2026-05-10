package loaders

import (
	"fmt"
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/views"

	lammps_parser "github.com/Ivanestver/lammps-file-parser/deserialize"
	lammps_structs "github.com/Ivanestver/lammps-file-parser/structs"
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
}

func (loader *_LammpsLoader) Load(filename string, fieldType datatypes.FieldType) (*views.GlobulaView, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	jsonStruct, err := lammps_parser.Deserialize(string(content), filename)
	if err != nil {
		return nil, err
	}
	return parseFromJSON(jsonStruct, fieldType)
}

func parseFromJSON(jsonStruct *lammps_structs.LammpsStruct, fieldType datatypes.FieldType) (*views.GlobulaView, error) {
	field := makeField(fieldType, jsonStruct.SpaceDimention)
	polymers, err := makePolymers(jsonStruct, field, fieldType)
	if err != nil {
		return nil, err
	}
	globula := views.NewGlobulaView(polymers, views.GlobulaThreadType)
	return globula, nil
}

func makeField(fieldType datatypes.FieldType, spaceDimention [3][2]float64) datatypes.IField {
	//return datatypes.NewField(getMaxDimention())
	return datatypes.CreateField(fieldType, spaceDimention)
}

func makePolymers(lammpsStruct *lammps_structs.LammpsStruct, field datatypes.IField, fieldType datatypes.FieldType) ([]datatypes.IPolymer, error) {
	polymersMap := make(map[int]datatypes.IPolymer)
	for _, atom := range lammpsStruct.Atoms {
		moleculeID := atom.MoleculeID - 1
		polymer, ok := polymersMap[moleculeID]
		var a *datatypes.Monomer
		if !ok {
			polymersMap[moleculeID] = datatypes.NewIPolymer(fieldType, field, int64(moleculeID))
			polymer = polymersMap[moleculeID]
		}
		a = field.GetMonomerByCoords(base.Vector3DF{
			atom.X,
			atom.Y,
			atom.Z,
		})
		if a == nil {
			panic("Monomer is nil")
		}
		atomTypeStruct := lammpsStruct.AtomTypes[atom.AtomType]
		mendeleevTableElement := base.RecognizeElement(atomTypeStruct.AtomLabel)
		if mendeleevTableElement == base.MendeleevTableElementUndefined {
			panic(fmt.Sprintf("Неизвестный элемент: %s", atomTypeStruct.AtomLabel))
		}
		a.MonomerType = mendeleevTableElement
		polymer.AddMonomer(a)
	}

	for _, bond := range lammpsStruct.Bonds {
		atom1 := lammpsStruct.Atoms[bond.Ends[0]-1]
		atom2 := lammpsStruct.Atoms[bond.Ends[1]-1]
		mon1 := field.GetMonomerByCoords(base.Vector3DF{atom1.X, atom1.Y, atom1.Z})
		mon2 := field.GetMonomerByCoords(base.Vector3DF{atom2.X, atom2.Y, atom2.Z})
		if err := datatypes.MakeConnection(mon1, mon2, datatypes.ConnectionTypeOne); err != nil {
			return nil, err
		}
	}

	polymers := make([]datatypes.IPolymer, len(polymersMap))
	for i, polymer := range polymersMap {
		polymers[i] = polymer
	}
	return polymers, nil
}
