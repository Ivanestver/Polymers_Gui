package loaders

import (
	"os"
	"polymers/base"
	"polymers/datatypes"
	"polymers/global_data"
	"polymers/views"

	lammps_parser "github.com/Ivanestver/lammps-file-parser/parser"
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

func (loader *_LammpsLoader) Load(filename, globulaname string) (*views.GlobulaView, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	jsonStruct, err := lammps_parser.Parse(string(content), filename)
	if err != nil {
		return nil, err
	}
	return parseFromJson(&jsonStruct, globulaname)
}

func parseFromJson(jsonStruct *lammps_structs.LammpsStruct, globulaName string) (*views.GlobulaView, error) {
	field := makeField()
	polymers := makePolymers(jsonStruct, field)
	literals := make(map[datatypes.MonomerType]string)
	literals[0] = "O"
	literals[1] = "N"
	literals[2] = "C"
	literals[3] = "S"
	globula := views.NewGlobulaView(globulaName, polymers, views.GLOBULA_THREAD_TYPE, literals)
	return globula, nil
}

func makeField() *datatypes.Field {
	return datatypes.NewField(getMaxDimention())
}

func getMaxDimention() uint64 {
	globalData := global_data.GetGlobalData()
	return uint64(max(globalData.SpaceDimention.X, globalData.SpaceDimention.Y, globalData.SpaceDimention.Z))
}

func makePolymers(lammpsStruct *lammps_structs.LammpsStruct, field *datatypes.Field) []*datatypes.Polymer {
	polymersMap := make(map[int]*datatypes.Polymer)
	for _, atom := range lammpsStruct.Atoms {
		moleculeID := atom.MoleculeID - 1
		polymer, ok := polymersMap[moleculeID]
		if !ok {
			polymersMap[moleculeID] = datatypes.NewPolymer(field, int64(moleculeID))
			polymer = polymersMap[moleculeID]
		}
		polymer.AddMonomer(datatypes.NewMonomer(base.Vector3DF{
			X: atom.X,
			Y: atom.Y,
			Z: atom.Z,
		}, datatypes.MonomerType(atom.AtomType-1)))
	}
	polymers := make([]*datatypes.Polymer, len(polymersMap))
	for i, polymer := range polymersMap {
		polymers[i] = polymer
	}
	return polymers
}
