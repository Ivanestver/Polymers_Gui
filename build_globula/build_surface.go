package build_globula

import (
	"math"
	"polymers/datatypes"
	"polymers/global_data"
	"polymers/views"
)

type SurfaceAlgInputData struct {
	Xlength, Ylength, Zlength int
}

func (inputData *SurfaceAlgInputData) GetGlobulaType() views.GlobulaProperty {
	return views.GLOBULA_SURFACE_TYPE
}

func (inputData *SurfaceAlgInputData) GetLiterals() map[datatypes.MonomerType]string {
	return BuildThreadAlgInputData{}.GetLiterals()
}

type SurfaceInputDataBuilder struct {
}

func (builder *SurfaceInputDataBuilder) CreateInputData(algType AlgType, defaultParams []string, particleName string) (ICalcAlgInputData, error) {
	params := append([]string{"0", "0", "0"}, defaultParams...)
	inputData, err := BuildThreadAlgInputDataBuilder{}.CreateInputData(algType, params, particleName)
	threadInputData := inputData.(BuildThreadAlgInputData)
	return &SurfaceAlgInputData{
		Xlength: int(threadInputData.ThreadRadius) * 2,
		Ylength: int(threadInputData.ThreadLength),
		Zlength: threadInputData.MaxPolymersCount,
	}, err
}

type SurfaceCalcAlg struct {
	inputData *SurfaceAlgInputData
}

func (alg *SurfaceCalcAlg) Calc() []*datatypes.Polymer {
	xDiv2 := alg.inputData.Xlength / 2
	yDiv2 := alg.inputData.Ylength / 2
	radius := math.Sqrt(
		float64(xDiv2)*float64(xDiv2) +
			float64(yDiv2)*float64(yDiv2))
	threadAlg := BuildThreadAlg{inputData: BuildThreadAlgInputData{
		Cell: struct {
			Lx int
			Ly int
			Lz int
		}{
			Lx: 0,
			Ly: 0,
			Lz: 0,
		},
		ThreadRadius:     radius,
		ThreadLength:     float64(alg.inputData.Zlength),
		MaxPolymersCount: alg.inputData.Xlength * alg.inputData.Ylength,
	}}
	spaceDimention := global_data.GetGlobalData().SpaceDimention
	global_data.ConfigureGlobalData(global_data.SpaceDimention{
		{Lower: 0.0, Higher: float64(alg.inputData.Xlength)},
		{Lower: 0.0, Higher: float64(alg.inputData.Ylength)},
		{Lower: 0.0, Higher: float64(alg.inputData.Zlength)},
	})
	polymers := threadAlg.Calc()
	global_data.ConfigureGlobalData(spaceDimention)
	return polymers
}
