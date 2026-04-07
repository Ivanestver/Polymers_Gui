package commandinterpreter

import "polymers/outputformat"

const (
	globulaName = "<globula_name>"
	filename    = "<filename>"
	outputName  = "<output_name>"
	typeStr     = "<type>"
)

func PrintHelp() {
	printFuncs([]func(){
		printHelp,
		printPattern,
		printBuild,
		printShow,
		printSave,
		printClusters,
		printAge,
		printBorders,
		printResetGlobula,
		printWaterize,
		printTrunk,
		printScript,
		printCommonStats,
	})
}

func printFuncs(functions []func()) {
	for _, f := range functions {
		f()
		outputformat.GetPrint().PrintEmptyLine()
	}
}

func printHelp() {
	outputformat.GetPrint().Printfln("%s - print this help", CommandHelpSTR)
}

func printPattern() {
	outputformat.GetPrint().Printfln("%s %s <pattern_source> %s - set pattern to globula, where", CommandPatternSTR, globulaName, outputName)
	outputformat.GetPrint().Printfln("\t%s - the name of an existing globula the pattern is applied to", globulaName)
	outputformat.GetPrint().Printfln("\t<pattern_source> - the source of a pattern. The values are")
	outputformat.GetPrint().Printfln("\t\t%s %s - from a file specified by %s", CommandFileSTR, filename, filename)
	outputformat.GetPrint().Printfln("\t\t%s %s - directly specify the pattern", CommandPatternSTR, "<pattern>")
	outputformat.GetPrint().Printfln("\t%s - the the name of a resulting globula", outputName)
}

func printBuild() {
	paramsList := "<params_list>"
	outputformat.GetPrint().Printfln("%s %s %s %s - build a globula of a type specified with a given name, where", CommandBuildSTR, typeStr, paramsList, globulaName)
	outputformat.GetPrint().Printfln("\t%s - the type of a globula:", typeStr)
	outputformat.GetPrint().Printfln("\t\t%s - a spherical globula", CommandGlobulaSTR)
	outputformat.GetPrint().Printfln("\t\t%s - a thread", CommandThreadSTR)
	outputformat.GetPrint().Printfln("\t%s - the list of parameters specific for a given type:", paramsList)
	outputformat.GetPrint().Printfln("\t\t%s - globula count, polymers count in each globula, accept threshold, maximum monomers count, sphere radius", CommandGlobulaSTR)
	outputformat.GetPrint().Printfln("\t\t%s - cell size along X, cell size along Y, cell size along Z, diameter, length, max polymers count", CommandThreadSTR)
	outputformat.GetPrint().Printfln("\t%s - the name of a built globula", globulaName)
}

func printShow() {
	outputformat.GetPrint().Printfln("%s <what_to_show> - show information specified by <what_to_show>", CommandShowSTR)
	outputformat.GetPrint().Println("\t'' (nothing) - show all globulas built")
	outputformat.GetPrint().Printfln("\tglobula \"%s\" - show globula parameters", globulaName)
}

func printSave() {
	outputformat.GetPrint().Printfln("%s %s - save the globula", CommandSaveSTR, globulaName)
}

func printClusters() {
	params := "<params>"
	outputformat.GetPrint().Printfln("%s %s %s - highlight clusters in a globula", CommandClustersSTR, globulaName, params)
	outputformat.GetPrint().Printfln("\t%s - a globula where to highlight", globulaName)
	outputformat.GetPrint().Printfln("\t%s - what clusters to highlight:", params)
	outputformat.GetPrint().Printfln("\t\t%s - all clusters", CommandClustersAllATR)
}

func printAge() {
	groupsCount := "<groupsCount>"
	makeCrosslinks := "<make_crosslinks>"
	outputformat.GetPrint().Printfln("%s %s %s %s - age a given globula. The number of aged groups is specified. It's possible to manage whether to make crosslinks", CommandAgeSTR, globulaName, groupsCount, makeCrosslinks)
	outputformat.GetPrint().Printfln("\t%s - specify the age ratio - either the number of aged groups or the percent of them", groupsCount)
	outputformat.GetPrint().Printfln("\t\ttrue - to make crosslinks")
	outputformat.GetPrint().Printfln("\t\tfalse - not to make crosslinks")
}

func printBorders() {
	outputformat.GetPrint().Printfln("%s %s - highlight border. NOTE: the feature is still in development", CommandHighlightBordersSTR, globulaName)
}

func printResetGlobula() {
	full := "<" + CommandFullSTR + ">"
	outputformat.GetPrint().Printfln("%s %s %s - reset globula to its initial state", CommandResetSTR, globulaName, full)
	outputformat.GetPrint().Printfln("\t%s - specify whether to make full reset:", full)
	outputformat.GetPrint().Printfln("\t\t'' (nothing) - make shallow reset")
	outputformat.GetPrint().Printfln("\t\t%s - make full reset", CommandFullSTR)
}

func printWaterize() {
	outputformat.GetPrint().Printfln("%s %s - drawn the globula into water", CommandWaterizeSTR, globulaName)
}

func printTrunk() {
	newSize := "<new_size>"
	outputformat.GetPrint().Printfln("%s %s %s - trunkate all the polymers in globula to a given size", CommandTrunkSTR, globulaName, newSize)
	outputformat.GetPrint().Printfln("\t%s - specify the new size as integer. Leave empty if to trunkate to the shortest polymer's size", newSize)
}

func printScript() {
	outputformat.GetPrint().Printfln("%s %s - use script specified by %s", CommandScriptSTR, filename, filename)
}

func printCommonStats() {
	outputformat.GetPrint().Printfln("%s %s %s - calculate the common statistics of a globula and save into a given file. The following information is showed", CommandCommonStatsSTR, globulaName, outputName)
	outputformat.GetPrint().Println("\t1. The number of monomers")
	outputformat.GetPrint().Println("\t2. The initial number of chains")
	outputformat.GetPrint().Println("\t3. The expected age ratio")
	outputformat.GetPrint().Println("\t4. The expected breaks count")
	outputformat.GetPrint().Println("\t   The expected crosslinks count")
	outputformat.GetPrint().Println("\t5. The expected age groups distribution")
	outputformat.GetPrint().Println("\t	C")
	outputformat.GetPrint().Println("\t	N")
	outputformat.GetPrint().Println("\t	H")
	outputformat.GetPrint().Println("\t6. The actual breaks count")
	outputformat.GetPrint().Println("\t   The actual crosslinks count")
	outputformat.GetPrint().Println("\t7. The actual age groups distribution")
	outputformat.GetPrint().Println("\t	C")
	outputformat.GetPrint().Println("\t	N")
	outputformat.GetPrint().Println("\t	H")
	outputformat.GetPrint().Println("\t   The actual age ratio")
	outputformat.GetPrint().Println("\t8. The mean chain length")
}
