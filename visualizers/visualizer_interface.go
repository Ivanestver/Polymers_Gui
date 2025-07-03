package visualizers

import (
	dt "data_visualizer/datatypes"
	"strconv"
)

const (
	VISUALIZER_TYPE_UNKNOWN    = -1
	VISUALIZER_TYPE_MASS_COUNT = iota
	VISUALIER_TYPE_COUNT
)

type VisualizerType = int

var visualizerTypesStrings map[VisualizerType]string = map[VisualizerType]string{
	VISUALIZER_TYPE_MASS_COUNT: "Parts Masses-Count Distribution after Aging",
}

func GetAllVisualizersDescription() []string {
	description := make([]string, VISUALIER_TYPE_COUNT-VISUALIZER_TYPE_MASS_COUNT)
	for t, s := range visualizerTypesStrings {
		description = append(description, strconv.Itoa(t)+" - "+s)
	}
	return description
}

type IVisualizer interface {
	Visualize(atoms []dt.Atom, bonds dt.Bonds, outputName string)
}

func NewVisualizer(visualizerType VisualizerType) IVisualizer {
	switch visualizerType {
	case VISUALIZER_TYPE_MASS_COUNT:
		return new(sAgeMassCountVisualizer)
	}
	return nil
}
