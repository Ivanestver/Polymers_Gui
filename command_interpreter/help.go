package command_interpreter

import "polymers/output_format"

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
		output_format.GetPrint().PrintEmptyLine()
	}
}

func printHelp() {
	output_format.GetPrint().Printfln("%s - print this help", COMMAND_HELP_STR)
}

func printPattern() {
	output_format.GetPrint().Printfln("%s %s <pattern_source> %s - set pattern to globula, where", COMMAND_PATTERN_STR, globulaName, outputName)
	output_format.GetPrint().Printfln("\t%s - the name of an existing globula the pattern is applied to", globulaName)
	output_format.GetPrint().Printfln("\t<pattern_source> - the source of a pattern. The values are")
	output_format.GetPrint().Printfln("\t\t%s %s - from a file specified by %s", COMMAND_FILE_STR, filename, filename)
	output_format.GetPrint().Printfln("\t\t%s %s - directly specify the pattern", COMMAND_PATTERN_STR, "<pattern>")
	output_format.GetPrint().Printfln("\t%s - the the name of a resulting globula", outputName)
}

func printBuild() {
	paramsList := "<params_list>"
	output_format.GetPrint().Printfln("%s %s %s %s - build a globula of a type specified with a given name, where", COMMAND_BUILD_STR, typeStr, paramsList, globulaName)
	output_format.GetPrint().Printfln("\t%s - the type of a globula:", typeStr)
	output_format.GetPrint().Printfln("\t\t%s - a spherical globula", COMMAND_GLOBULA_STR)
	output_format.GetPrint().Printfln("\t\t%s - a thread", COMMAND_THREAD_STR)
	output_format.GetPrint().Printfln("\t%s - the list of parameters specific for a given type:", paramsList)
	output_format.GetPrint().Printfln("\t\t%s - globula count, polymers count in each globula, accept threshold, maximum monomers count, sphere radius", COMMAND_GLOBULA_STR)
	output_format.GetPrint().Printfln("\t\t%s - cell size along X, cell size along Y, cell size along Z, diameter, length, max polymers count", COMMAND_THREAD_STR)
	output_format.GetPrint().Printfln("\t%s - the name of a built globula", globulaName)
}

func printShow() {
	output_format.GetPrint().Printfln("%s <what_to_show> - show information specified by <what_to_show>", COMMAND_SHOW_STR)
	output_format.GetPrint().Println("\t'' (nothing) - show all globulas built")
	output_format.GetPrint().Printfln("\tglobula \"%s\" - show globula parameters", globulaName)
}

func printSave() {
	output_format.GetPrint().Printfln("%s %s - save the globula", COMMAND_SAVE_STR, globulaName)
}

func printClusters() {
	params := "<params>"
	output_format.GetPrint().Printfln("%s %s %s - highlight clusters in a globula", COMMAND_CLUSTERS_STR, globulaName, params)
	output_format.GetPrint().Printfln("\t%s - a globula where to highlight", globulaName)
	output_format.GetPrint().Printfln("\t%s - what clusters to highlight:", params)
	output_format.GetPrint().Printfln("\t\t%s - all clusters", COMMAND_CLUSTERS_ALL_STR)
}

func printAge() {
	groupsCount := "<groupsCount>"
	makeCrosslinks := "<make_crosslinks>"
	output_format.GetPrint().Printfln("%s %s %s %s - age a given globula. The number of aged groups is specified. It's possible to manage whether to make crosslinks", COMMAND_AGE_STR, globulaName, groupsCount, makeCrosslinks)
	output_format.GetPrint().Printfln("\t%s - specify the age ratio - either the number of aged groups or the percent of them", groupsCount)
	output_format.GetPrint().Printfln("\t\ttrue - to make crosslinks")
	output_format.GetPrint().Printfln("\t\tfalse - not to make crosslinks")
}

func printBorders() {
	output_format.GetPrint().Printfln("%s %s - highlight border. NOTE: the feature is still in development", COMMAND_HIGHLIGHT_BORDERS_STR, globulaName)
}

func printResetGlobula() {
	full := "<" + COMMAND_FULL_STR + ">"
	output_format.GetPrint().Printfln("%s %s %s - reset globula to its initial state", COMMAND_RESET_STR, globulaName, full)
	output_format.GetPrint().Printfln("\t%s - specify whether to make full reset:", full)
	output_format.GetPrint().Printfln("\t\t'' (nothing) - make shallow reset")
	output_format.GetPrint().Printfln("\t\t%s - make full reset", COMMAND_FULL_STR)
}

func printWaterize() {
	output_format.GetPrint().Printfln("%s %s - drawn the globula into water", COMMAND_WATERIZE_STR, globulaName)
}

func printTrunk() {
	newSize := "<new_size>"
	output_format.GetPrint().Printfln("%s %s %s - trunkate all the polymers in globula to a given size", COMMAND_TRUNK_STR, globulaName, newSize)
	output_format.GetPrint().Printfln("\t%s - specify the new size as integer. Leave empty if to trunkate to the shortest polymer's size", newSize)
}

func printScript() {
	output_format.GetPrint().Printfln("%s %s - use script specified by %s", COMMAND_SCRIPT_STR, filename, filename)
}

func printCommonStats() {
	output_format.GetPrint().Printfln("%s %s %s - calculate the common statistics of a globula and save into a given file. The following information is showed", COMMAND_COMMON_STATS_STR, globulaName, outputName)
	output_format.GetPrint().Println("\t1. The number of monomers")
	output_format.GetPrint().Println("\t2. The initial number of chains")
	output_format.GetPrint().Println("\t3. The expected age ratio")
	output_format.GetPrint().Println("\t4. The expected breaks count")
	output_format.GetPrint().Println("\t   The expected crosslinks count")
	output_format.GetPrint().Println("\t5. The expected age groups distribution")
	output_format.GetPrint().Println("\t	C")
	output_format.GetPrint().Println("\t	N")
	output_format.GetPrint().Println("\t	H")
	output_format.GetPrint().Println("\t6. The actual breaks count")
	output_format.GetPrint().Println("\t   The actual crosslinks count")
	output_format.GetPrint().Println("\t7. The actual age groups distribution")
	output_format.GetPrint().Println("\t	C")
	output_format.GetPrint().Println("\t	N")
	output_format.GetPrint().Println("\t	H")
	output_format.GetPrint().Println("\t   The actual age ratio")
	output_format.GetPrint().Println("\t8. The mean chain length")
}
