package loaders

import (
	"errors"
	"polymers/views"
)

type ILoader interface {
	Load(filename string) (*views.GlobulaView, error)
}

func NewLoader(loaderType string) (ILoader, error) {
	return nil, errors.New("not a registered type")
}
