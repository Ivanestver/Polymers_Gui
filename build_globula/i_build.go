package build_globula

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
	GetLiterals() map[datatypes.MonomerType]string
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

func GetLiteralsTable() map[datatypes.MonomerType]string {
	m := make(map[datatypes.MonomerType]string)
	m[datatypes.MONOMER_TYPE_UNDEFINED] = ""
	m[datatypes.MONOMER_TYPE_USUAL] = "O"
	m[datatypes.MONOMER_TYPE_O_CONTAINING] = "N"
	m[datatypes.MONOMER_TYPE_VYNIL] = "C"
	m[datatypes.MONOMER_TYPE_FWISE] = "F"
	m[datatypes.MONOMER_TYPE_CLWISE] = "Cl"
	m[datatypes.MONOMER_TYPE_CROSSLINKED] = "H"
	m[datatypes.MONOMER_TYPE_WATER] = "I"
	return m
}
