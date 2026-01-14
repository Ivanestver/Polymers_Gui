package loaders

import (
	"bufio"
	"os"
	"polymers/views"
)

type _LammpsMetadata struct {
	atomsCount     int
	atomTypesCount int
	bondsCount     int
	bondTypeCount  int
}

type _LammpsLoader struct {
	_LammpsMetadata
	builtGlobula *views.GlobulaView
	scanner      *bufio.Scanner
}

func Load(filename string) (*views.GlobulaView, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	loader := _LammpsLoader{
		scanner: bufio.NewScanner(file),
	}
	loader.load()
	return loader.builtGlobula, nil
}

func (loader *_LammpsLoader) load() error {
}

func (loader *_LammpsLoader) loadMetadata() error {

}
