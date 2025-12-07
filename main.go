package main

import (
	"bufio"
	"flag"
	"math/rand"
	"os"
	"polymers/build_globula"
	interp "polymers/command_interpreter"
	"polymers/datatypes"
	"polymers/dfs"
	"polymers/global_data"
	"polymers/output_format"
	"polymers/savers"
	"polymers/views"
	"slices"
	"strconv"
	"time"
)

var globulas []*views.GlobulaView

func getGlobulaByName(name string) *views.GlobulaView {
	for _, glob := range globulas {
		if glob.Name() == name {
			return glob
		}
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	output_format.SetPrint(&output_format.ColoredConsolePrint{})
	fileNumber := 1
	output_format.GetPrint().PrintInfo("Welcome to the Polymer Builder 2.0. Please, type the space dimention: ")
	var spaceDimention global_data.SpaceDimention
	spaceDimention.X = 3
	spaceDimention.Y = 3
	spaceDimention.Z = 1
	//output_format.GetPrint().Readln(&spaceDimention.X, &spaceDimention.Y, &spaceDimention.Z)
	output_format.GetPrint().PrintfInfo("The space dimention set by user is %d\n", spaceDimention)

	output_format.GetPrint().PrintlnInfo("Configuring the global data")
	global_data.ConfigureGlobalData(spaceDimention)
	output_format.GetPrint().PrintlnInfo("Configuring the global data finished")
	output_format.GetPrint().PrintlnInfo("The preparations are done! Now you may set up the input data and run the algorithm.")
	cmdReader := bufio.NewReader(os.Stdin)
	commands := make([]string, 0)
	scriptPtr := flag.String("script", "", "define a script file")
	flag.Parse()
	if scriptPtr != nil {
		commands = getCommandsFromScript(*scriptPtr)
	}
	isWorking := true
	for isWorking {
		output_format.GetPrint().Print("> ")
		var line string
		if len(commands) == 0 {
			line, _ = cmdReader.ReadString('\n')
		} else {
			line = commands[0]
			commands = commands[1:]
		}
		command, data := interp.Interpret(line)
		switch command {
		case interp.COMMAND_UNDEFINED:
			output_format.GetPrint().PrintlnError(data.(string))
		case interp.COMMAND_HELP:
			interp.PrintHelp()
		case interp.COMMAND_BUILD:
			m := data.(map[string]interface{})
			buildGlobula(m["alg"].(build_globula.AlgType), m["params"].([]string), m["name"].(string))
		case interp.COMMAND_SHOW_GLOBULAS_LIST:
			for _, globula := range globulas {
				output_format.GetPrint().PrintlnInfo(globula.Name())
			}
		case interp.COMMAND_SHOW_GLOBULA:
			globulaName := data.(string)
			var globula *views.GlobulaView = getGlobulaByName(globulaName)
			if globula != nil {
				PrintGlobulaInfo(globula)
			} else {
				output_format.GetPrint().PrintlnError("There is no globula called \"" + globulaName + "\"")
			}
		case interp.COMMAND_SAVE_GLOBULA:
			globulaName := data.(string)
			var globula *views.GlobulaView = getGlobulaByName(globulaName)
			if globula == nil {
				output_format.GetPrint().PrintlnError("There is no globula called \"" + globulaName + "\"")
				break
			}
			content, _ := savers.SaveToLammps(globula)
			f, err := os.Create(globulaName + strconv.Itoa(fileNumber) + ".data")
			if err != nil {
				output_format.GetPrint().PrintlnError(err.Error())
				return
			}
			defer f.Close()
			f.Write([]byte(content))
			/*f, err = os.Create(globulaName + strconv.Itoa(fileNumber) + ".json")
			if err != nil {
				output_format.GetPrint().PrintlnError(err.Error())
				break
			}
			if err := json.NewEncoder(f).Encode(globula); err != nil {
				output_format.GetPrint().PrintlnError(err.Error())
			}*/
			fileNumber++
		case interp.COMMAND_HIGHLIGHT_CLUSTERS_ALL:
			globulaName := data.(string)
			var originalGlobula *views.GlobulaView = getGlobulaByName(globulaName)
			globula := originalGlobula.DeepCopy(originalGlobula.Name() + "_all_clusters")
			globula.Name()
			globulas = append(globulas, globula)

			output_format.GetPrint().Println("Start highlighting clusters")
			xClusters, yClusters, zClusters := globula.CommonClusters()
			if xClusters != nil {
				xClusters.Colorize(false)
			}
			if yClusters != nil {
				yClusters.Colorize(false)
			}
			if zClusters != nil {
				zClusters.Colorize(false)
			}

		case interp.COMMAND_AGE:
			data := data.(map[string]interface{})
			groupsCountStr := data["count"].(string)
			globulaName := data["globula"].(string)
			doCrosslinks := data["make_crosslinks"].(bool)
			originalGlobula := getGlobulaByName(globulaName)
			if originalGlobula == nil {
				output_format.GetPrint().PrintlnError("There is no globula named " + globulaName)
				break
			}
			globula := originalGlobula.DeepCopy(originalGlobula.Name() + "_aged_" + strconv.Itoa(len(globulas)))
			globulas = append(globulas, globula)
			groupsCount := 0
			if groupsCountStr[len(groupsCountStr)-1] == '%' {
				percent, _ := strconv.ParseFloat(groupsCountStr[:len(groupsCountStr)-1], 64)
				groupsCount = int(float64(globula.GetAtomsCount()) * float64(percent) / 100)
			} else {
				groupsCount, _ = strconv.Atoi(groupsCountStr)
			}
			globula.DoAging2(groupsCount, doCrosslinks)

		case interp.COMMAND_RESET:
			data := data.(map[string]interface{})
			globulaName := data["globula"].(string)
			globula := getGlobulaByName(globulaName)
			globula.Reset()

		case interp.COMMAND_RESET_FULL:
			data := data.(map[string]interface{})
			globulaName := data["globula"].(string)
			globula := getGlobulaByName(globulaName)
			globula.FullReset()

		case interp.COMMAND_HIGHLIGHT_BORDERS:
			data := data.(map[string]interface{})
			globulaName := data["globula"].(string)
			originGlobula := getGlobulaByName(globulaName)
			globula := originGlobula.DeepCopy(globulaName + "highlighted_borders")
			globula.HighlightBorders()

		case interp.COMMAND_WATERIZE:
			data := data.(map[string]interface{})
			globulaName := data["globula"].(string)
			originGlobula := getGlobulaByName(globulaName)
			globula := originGlobula.DeepCopy(globulaName + "_waterized")
			globula.Waterize()
			globulas = append(globulas, globula)

		case interp.COMMAND_TRUNK:
			data := data.(map[string]interface{})
			globulaName := data["globula"].(string)
			originGlobula := getGlobulaByName(globulaName)
			newSize, ok := data["new_size"]
			var globula *views.GlobulaView
			if ok {
				globula = originGlobula.DeepCopy(globulaName + "_trunk_custom")
				globula.MakeHomogenousAsCustom(newSize.(int))
			} else {
				globula = originGlobula.DeepCopy(globulaName + "_trunk_shortest")
				globula.MakeHomogenousAsShortest()
			}
			globulas = append(globulas, globula)

		case interp.COMMAND_PATTERN:
			data := data.(map[string]string)
			globulaName := data["globulaName"]
			originGlobula := getGlobulaByName(globulaName)
			outputName := data["outputName"]
			globula := originGlobula.DeepCopy(outputName)
			var pattern string
			fileName, ok := data["fileName"]
			if ok {
				file, err := os.Open(fileName)
				if err != nil {
					output_format.GetPrint().PrintlnError(err.Error())
					break
				}
				defer func() {
					if closeErr := file.Close(); closeErr != nil {
						output_format.GetPrint().PrintlnError(closeErr.Error())
					}
				}()
				scanner := bufio.NewScanner(file)
				if scanner.Scan() {
					pattern = scanner.Text()
				} else {
					output_format.GetPrint().PrintlnError("The given file does not contain a pattern")
					break
				}
			} else if p, ok := data["pattern"]; ok {
				pattern = p
			} else {
				output_format.GetPrint().PrintflnError("See usage of this command")
				break
			}

			var anyNotExist bool = false
			for _, letter := range pattern {
				l := string(letter)
				if globula.GetMonomerTypeByLiteral(l) == datatypes.MONOMER_TYPE_UNDEFINED {
					output_format.GetPrint().PrintlnError(l + " does not have its decryption")
					anyNotExist = true
				}
			}
			if anyNotExist {
				output_format.GetPrint().PrintlnError("Please, define the missing decryptions to continue")
				break
			}

			if globula.Is(views.GLOBULA_GLOBULA_TYPE) {
				views.ForEachPolymer(globula, func(pv *views.PolymerView) {
					currentLetterNumber := 0
					views.ForEachMonomer(pv, func(m *datatypes.Monomer) bool {
						m.MonomerType = globula.GetMonomerTypeByLiteral(string(pattern[currentLetterNumber]))
						currentLetterNumber = (currentLetterNumber + 1) % len(pattern)
						return true
					})
				})
			} else {
				literalsTable := globula.GetLiterals()
				literals := make([]datatypes.MonomerType, 0)
				for m := range *literalsTable {
					literals = append(literals, m)
				}
				slices.SortFunc(literals, func(a, b datatypes.MonomerType) int {
					if int(a) < int(b) {
						return -1
					} else if int(a) == int(b) {
						return 0
					} else {
						return 1
					}
				})
				threshold := len(*literalsTable) / 2
				latestType := threshold + 1
				views.ForEachPolymer(globula, func(pv *views.PolymerView) {
					currentLetterNumber := 0
					views.ForEachMonomer(pv, func(m *datatypes.Monomer) bool {
						for _, monType := range literals {
							if string(pattern[currentLetterNumber]) == (*literalsTable)[monType] {
								m.MonomerType = monType
								if latestType < threshold {
									m.MonomerType = m.MonomerType + datatypes.MonomerType(threshold)
								}
								latestType = int(m.MonomerType)
								currentLetterNumber = currentLetterNumber + 1
								return true
							}
						}
						return false
					})
				})
			}

			globulas = append(globulas, globula)

		case interp.COMMAND_EXIT:
			isWorking = false

		case interp.COMMAND_SCRIPT:
			filename, err := data.(string)
			if !err {
				output_format.GetPrint().PrintlnError("Usage: script <filename>")
			}
			commands = getCommandsFromScript(filename)

		case interp.COMMAND_COMMON_STATS:
			data := data.(map[string]string)
			globulaName := data["globula"]
			globula := getGlobulaByName(globulaName)
			text := globula.GetStatistics()
			f, err := os.Create(data["filename"])
			if err != nil {
				output_format.GetPrint().PrintlnError(err.Error())
				return
			}
			defer f.Close()
			f.Write([]byte(text))

		case interp.COMMAND_DFS:
			globula := dfs.DoDFS()
			globulas = append(globulas, globula)

		default:
			output_format.GetPrint().PrintlnError("'" + line[:len(line)-1] + "' is not supported")
		}
	}
}

func buildGlobula(algType build_globula.AlgType, predefinedParams []string, particleName string) {
	inputDataBuilder := build_globula.CreateInputDataBuilder(algType)
	inputData_, err := inputDataBuilder.CreateInputData(algType, predefinedParams, particleName)
	if err != nil {
		output_format.GetPrint().PrintlnError(err.Error())
		return
	}
	calcAlg := build_globula.CreateCalcAlg(inputData_, algType)
	finishedPolymers := calcAlg.Calc()
	if finishedPolymers == nil {
		output_format.GetPrint().PrintlnError("The result of building is nil")
	} else {
		globula := views.NewGlobulaView(inputData_.GetName(), finishedPolymers, inputData_.GetGlobulaType())
		globula.SetLiterals(inputData_.GetLiterals())
		globulas = append(globulas, globula)
	}
}

func PrintGlobulaInfo(globula *views.GlobulaView) {
	output_format.GetPrint().Println("\n\tGlobula Name: " + globula.Name())
	output_format.GetPrint().Println("\tPolymers Count: " + strconv.Itoa(globula.Len()))
	output_format.GetPrint().Println("\tPolymers:")
	var monomersCount int
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		monomersCount += pol.Len()
		output_format.GetPrint().Println("\t\tPolymer Name: " + pol.Name())
		output_format.GetPrint().Println("\t\tMonomers Count: " + strconv.Itoa(pol.Len()))
		output_format.GetPrint().PrintEmptyLine()
	})
	output_format.GetPrint().Println("\t\tMonomers in total: " + strconv.Itoa(monomersCount))
}

func getCommandsFromScript(filename string) []string {
	commands := make([]string, 0)

	file, err := os.Open(filename)
	if err != nil {
		output_format.GetPrint().PrintlnError(err.Error())
		return commands
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		commands = append(commands, line)
	}

	return commands
}
