package loaders

import (
	"errors"
	"polymers/datatypes"
	"polymers/views"
)

type ILoader interface {
	Load(filename string, field datatypes.FieldType) (*views.GlobulaView, error)
}

func NewLoader(loaderType string, field datatypes.FieldType) (ILoader, error) {
	if loader, ok := map[string]ILoader{
		"lammps": &_LammpsLoader{},
	}[loaderType]; ok {
		return loader, nil
	} else {
		return nil, errors.New("not a registered type")
	}
}
