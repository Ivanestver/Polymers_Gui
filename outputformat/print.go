package outputformat

var iprint IPrint

func SetPrint(p IPrint) {
	iprint = p
}

func GetPrint() IPrint {
	if iprint == nil {
		panic("iprint is nil")
	}
	return iprint
}
