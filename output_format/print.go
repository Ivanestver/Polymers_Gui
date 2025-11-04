package output_format

var iprint IPrint

func SetPrint(p IPrint) {
	iprint = p
}

func GetPrint() IPrint {
	return iprint
}
