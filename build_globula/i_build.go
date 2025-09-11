package build_globula

import "polymers/datatypes"

type AlgType int

const (
	GlobulaBuildAlg AlgType = iota
	ThreadBuildAlg
)

type ICalcAlg interface {
	Calc() []*datatypes.Polymer
	GetLiteralsTable() map[datatypes.MonomerType]string
}

type ICalcAlgInputData interface {
	GetName() string
}

type IInputDataBuilder interface {
	CreateInputData(algType AlgType) (ICalcAlgInputData, error)
}

func CreateInputDataBuilder(algType AlgType) IInputDataBuilder {
	if algType == GlobulaBuildAlg {
		return &CalcAlgInputDataBuilder{}
	} else {
		return &BuildThreadAlgInputDataBuilder{}
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
	default:
		return nil
	}
	return nil
}
