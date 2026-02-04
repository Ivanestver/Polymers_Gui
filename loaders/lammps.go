package loaders

import (
	"os"
	"polymers/base"
	"polymers/build_globula"
	"polymers/datatypes"
	"polymers/global_data"
	"polymers/views"
	"strings"

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
	jsonStruct, err := lammps_parser.Parse(string(content), globulaname)
	if err != nil {
		return nil, err
	}
	return parseFromJson(&jsonStruct, strings.Split(filename, ".")[0])
}

func parseFromJson(jsonStruct *lammps_structs.LammpsStruct, globulaName string) (*views.GlobulaView, error) {
	field := makeField()
	polymers := makePolymers(jsonStruct, field)
	globula := views.NewGlobulaView(globulaName, polymers, views.GLOBULA_THREAD_TYPE, build_globula.BuildThreadAlgInputData{}.GetLiterals())
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
		moleculeID := atom.MoleculeID
		polymer, ok := polymersMap[moleculeID]
		if !ok {
			polymersMap[moleculeID] = datatypes.NewPolymer(field, int64(moleculeID))
			polymer = polymersMap[moleculeID]
		}
		polymer.AddMonomer(datatypes.NewMonomer(base.Vector3DF{
			X: atom.X,
			Y: atom.Y,
			Z: atom.Z,
		}, datatypes.MonomerType(atom.AtomType)))
	}
	polymers := make([]*datatypes.Polymer, len(polymersMap))
	for i, polymer := range polymersMap {
		polymers[i-1] = polymer
	}
	return polymers
}
