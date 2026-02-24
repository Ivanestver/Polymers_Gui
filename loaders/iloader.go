package loaders

import (
	"errors"
	"polymers/views"
)

type FieldType = int

const (
	FIELD_TYPE_REAL FieldType = iota
	FIELD_TYPE_LATTICE
)

type ILoader interface {
	Load(filename string, field FieldType) (*views.GlobulaView, error)
}

func NewLoader(loaderType string, field FieldType) (ILoader, error) {
	if loader, ok := map[string]ILoader{
		"lammps": &_LammpsLoader{},
	}[loaderType]; ok {
		return loader, nil
	} else {
		return nil, errors.New("not a registered type")
	}
}
