package build_globula

import "polymers/datatypes"

type BuildThreadAlgInputData struct {
}

func (alg BuildThreadAlgInputData) GetName() string {
	return "Thread"
}

type BuildThreadAlg struct {
	inputData BuildThreadAlgInputData
}

func (alg *BuildThreadAlg) Calc() []*datatypes.Polymer {
	return nil
}
