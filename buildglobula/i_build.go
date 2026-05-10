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
)

type ICalcAlg interface {
	Calc() []*datatypes.Polymer
}

type ICalcAlgInputData interface {
	GetGlobulaType() views.GlobulaProperty
}

type IInputDataBuilder interface {
	CreateInputData(algType AlgType, defaultParams []string, particleName string) (ICalcAlgInputData, error)
}

func CreateInputDataBuilder(algType AlgType) IInputDataBuilder {
	switch algType {
	case GlobulaBuildAlg:
		return &CalcAlgInputDataBuilder{}
	case ThreadBuildAlg:
		return &BuildThreadAlgInputDataBuilder{}
	case SurfaceBuildAlg:
		return &SurfaceInputDataBuilder{}
	default:
		return nil
	}
}

func CreateCalcAlg(inputData ICalcAlgInputData, algType AlgType) ICalcAlg {
	switch algType {
	case GlobulaBuildAlg:
		inp, ok := inputData.(CalcAlgInputData)
		if ok {
			return &CalcAlg{
				inputData: inp,
			}
		}
	case ThreadBuildAlg:
		inp, ok := inputData.(BuildThreadAlgInputData)
		if ok {
			return &BuildThreadAlg{
				inputData: inp,
			}
		}
	case SurfaceBuildAlg:
		inp, ok := inputData.(*SurfaceAlgInputData)
		if ok {
			return &SurfaceCalcAlg{
				inputData: inp,
			}
		}
	default:
		return nil
	}
	return nil
}
