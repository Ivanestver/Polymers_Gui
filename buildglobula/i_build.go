package buildglobula

import (
	"polymers/datatypes"
	"polymers/views"
)

type AlgType int

const (
	GlobulaBuildAlg AlgType = iota
	ThreadBuildAlg
	SurfaceBuildAlg
	PatternBuildAlg
	CrystallBuildAlg
)

type ICalcAlg interface {
	Calc() []*datatypes.Polymer
}

type ICalcAlgInputData interface {
	GetGlobulaType() views.GlobulaProperty
}

type ICalcAlgInputDataBuilder interface {
	CreateInputData(algType AlgType, defaultParams []string, particleName string) (ICalcAlgInputData, error)
}

type AbstractAlg[T ICalcAlgInputData] struct {
	inputData T
}

func CreateInputDataBuilder(algType AlgType) ICalcAlgInputDataBuilder {
	switch algType {
	case GlobulaBuildAlg:
		return GlobulaInputDataBuilder{}
	case ThreadBuildAlg:
		return ThreadInputDataBuilder{}
	case SurfaceBuildAlg:
		return SurfaceInputDataBuilder{}
	case PatternBuildAlg:
		return PatternInputDataBuilder{}
	case CrystallBuildAlg:
		return CrystallInputDataBuilder{}
	default:
		return nil
	}
}

func CreateCalcAlg(inputData ICalcAlgInputData, algType AlgType) ICalcAlg {
	switch algType {
	case GlobulaBuildAlg:
		inp, ok := inputData.(GlobulaInputData)
		if ok {
			return &GlobulaCalcAlg{
				inputData: inp,
			}
		}
	case ThreadBuildAlg:
		inp, ok := inputData.(ThreadInputData)
		if ok {
			return &ThreadCalcAlg{
				inputData: inp,
			}
		}
	case SurfaceBuildAlg:
		inp, ok := inputData.(SurfaceInputData)
		if ok {
			return &SurfaceCalcAlg{
				inputData: inp,
			}
		}
	case PatternBuildAlg:
		inp, ok := inputData.(PatternInputData)
		if ok {
			return &PatternCalcAlg{
				inputData: inp,
			}
		}
	case CrystallBuildAlg:
		inp, ok := inputData.(CrystallBuildInputData)
		if ok {
			return &CrystallCalcAlg{
				inputData: inp,
			}
		}
	default:
		return nil
	}
	return nil
}
