package main

import (
	"bufio"
	"flag"
	"math/rand"
	"os"
	"polymers/atomistic"
	"polymers/base"
	"polymers/build_globula"
	interp "polymers/command_interpreter"
	"polymers/dfs"
	"polymers/global_data"
	"polymers/loaders"
	"polymers/output_format"
	pattern_lib "polymers/pattern"
	"polymers/savers"
	"polymers/views"
	"strconv"
	"strings"
	"time"
)

var globulas []*views.GlobulaView
var printer output_format.IPrint

func getGlobulaByName(name string) *views.GlobulaView {
	for _, glob := range globulas {
		if glob.Name() == name {
			return glob
		}
	}
	return nil
}

func setUpSpaceDimention(commands *[]string) global_data.SpaceDimention {
	var spaceDimention global_data.SpaceDimention
	if len(*commands) == 0 {
		printer.PrintInfo("Please, type the space dimention: ")
		printer.Readln(&spaceDimention.X, &spaceDimention.Y, &spaceDimention.Z)
		return spaceDimention
	}
	line := (*commands)[0]
	parts := strings.Split(line, " ")
	if len(parts) != 4 || parts[0] != "space" {
		printer.PrintInfo("No space dimention definition found. Please, type the space dimention: ")
		printer.Readln(&spaceDimention.X, &spaceDimention.Y, &spaceDimention.Z)
	} else {
		*commands = (*commands)[1:]
		if x, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
			spaceDimention.X = x
		} else {
			printer.PrintlnError(err.Error())
		}
		if y, err := strconv.ParseInt(parts[2], 10, 64); err == nil {
			spaceDimention.Y = y
		} else {
			printer.PrintlnError(err.Error())
		}
		if z, err := strconv.ParseInt(parts[3], 10, 64); err == nil {
			spaceDimention.Z = z
		} else {
			printer.PrintlnError(err.Error())
		}
	}
	printer.PrintfInfo("The space dimention set by user is %d\n", spaceDimention)
	return spaceDimention
}

func main() {
	rand.Seed(time.Now().UnixNano())
	output_format.SetPrint(&output_format.ColoredConsolePrint{})
	printer = output_format.GetPrint()
	printer.PrintlnInfo("Welcome to the Polymer Builder 2.0")
	commands := make([]string, 0)
	scriptPtr := flag.String("script", "", "define a script file")
	modePtr := flag.String("mode", "", "define mode to manage execution")
	flag.Parse()
	if len(*scriptPtr) > 0 {
		commands = getCommandsFromScript(*scriptPtr, *modePtr)
	}
	spaceDimention := setUpSpaceDimention(&commands)
	printer.PrintlnInfo("Configuring the global data")
	global_data.ConfigureGlobalData(spaceDimention)
	printer.PrintlnInfo("Configuring the global data finished")
	printer.PrintlnInfo("The preparations are done! Now you may set up the input data and run the algorithm.")
	cmdReader := bufio.NewReader(os.Stdin)
	isWorking := true
	for isWorking {
		printer.Print("> ")
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
			printer.PrintlnError(data.(string))
		case interp.COMMAND_HELP:
			interp.PrintHelp()
		case interp.COMMAND_BUILD:
			m := data.(map[string]interface{})
			buildGlobula(m["alg"].(build_globula.AlgType), m["params"].([]string), m["name"].(string))
		case interp.COMMAND_SHOW_GLOBULAS_LIST:
			for _, globula := range globulas {
				printer.PrintlnInfo(globula.Name())
			}
		case interp.COMMAND_SHOW_GLOBULA:
			globulaName := data.(string)
			var globula *views.GlobulaView = getGlobulaByName(globulaName)
			if globula != nil {
				PrintGlobulaInfo(globula)
			} else {
				printer.PrintlnError("There is no globula called \"" + globulaName + "\"")
			}
		case interp.COMMAND_SAVE_GLOBULA:
			globulaName := data.(string)
			var globula *views.GlobulaView = getGlobulaByName(globulaName)
			if globula == nil {
				printer.PrintlnError("There is no globula called \"" + globulaName + "\"")
				break
			}
			content, _ := savers.SaveToLammps(globula)
			f, err := os.Create(globulaName + ".data")
			if err != nil {
				printer.PrintlnError(err.Error())
				return
			}
			defer f.Close()
			f.Write([]byte(content))
			/*f, err = os.Create(globulaName + strconv.Itoa(fileNumber) + ".json")
			if err != nil {
				printer.PrintlnError(err.Error())
				break
			}
			if err := json.NewEncoder(f).Encode(globula); err != nil {
				printer.PrintlnError(err.Error())
			}*/
		case interp.COMMAND_HIGHLIGHT_CLUSTERS_ALL:
			globulaName := data.(string)
			var originalGlobula *views.GlobulaView = getGlobulaByName(globulaName)
			globula := originalGlobula.DeepCopy(originalGlobula.Name() + "_all_clusters")
			globula.Name()
			globulas = append(globulas, globula)

			printer.Println("Start highlighting clusters")
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
			algType := data["alg_type"].(int)
			newGlobulaName := data["new_globula_name"].(string)
			originalGlobula := getGlobulaByName(globulaName)
			if originalGlobula == nil {
				printer.PrintlnError("There is no globula named " + globulaName)
				break
			}
			globula := originalGlobula.DeepCopy(newGlobulaName)
			globulas = append(globulas, globula)
			groupsCount := 0
			if groupsCountStr[len(groupsCountStr)-1] == '%' {
				percent, _ := strconv.ParseFloat(groupsCountStr[:len(groupsCountStr)-1], 64)
				groupsCount = int(float64(globula.GetAtomsCount()) * float64(percent) / 100)
			} else {
				groupsCount, _ = strconv.Atoi(groupsCountStr)
			}
			if algType == 1 {
				globula.DoAging1(groupsCount)
			} else if algType == 2 {
				globula.DoAging2(groupsCount, doCrosslinks)
			} else if algType == 3 {
				ncut, ok := data["ncut"].(int)
				if ok {
					continue
				}
				nOContaining, ok := data["nOContaining"].(int)
				if ok {
					continue
				}
				ncross, ok := data["ncross"].(int)
				if ok {
					continue
				}
				globula.DoAging3(ncut, nOContaining, ncross)
			}

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
			if globula := processPattern(data.(map[string]string)); globula != nil {
				globulas = append(globulas, globula)
			}

		case interp.COMMAND_EXIT:
			isWorking = false

		case interp.COMMAND_SCRIPT:
			filename, err := data.(string)
			if !err {
				printer.PrintlnError("Usage: script <filename>")
			}
			commands = getCommandsFromScript(filename, "")

		case interp.COMMAND_COMMON_STATS:
			data := data.(map[string]string)
			globulaName := data["globula"]
			globula := getGlobulaByName(globulaName)
			text := globula.GetStatistics()
			f, err := os.Create(data["filename"])
			if err != nil {
				printer.PrintlnError(err.Error())
				return
			}
			defer f.Close()
			f.Write([]byte(text))

		case interp.COMMAND_DFS:
			globula := dfs.DoDFS()
			globulas = append(globulas, globula)

		case interp.COMMAND_ATOMISTIC:
			data := data.(map[string]string)
			globulaName := data["globula"]
			globula := getGlobulaByName(globulaName)
			if config, ok := data["config"]; ok && globula != nil {
				atomistic.MakeAtomistic(globula, config)
			} else {
				printer.PrintflnError("Could not find a globula of the name %s", globulaName)
			}

		case interp.COMMAND_LOADER:
			data := data.(map[string]string)
			filetype := data["filetype"]
			filename := data["filename"]
			globulaName := data["globulaName"]
			if loader, err := loaders.NewLoader(filetype); err == nil {
				if newGlobula, err := loader.Load(filename, globulaName); err == nil {
					globulas = append(globulas, newGlobula)
				} else {
					printer.PrintflnError("When loading: %s", err.Error())
				}
			} else {
				printer.PrintflnError("When loading: %s", err.Error())
			}

		default:
			printer.PrintlnError("'" + line[:len(line)-1] + "' is not supported")
		}
	}
}

func buildGlobula(algType build_globula.AlgType, predefinedParams []string, particleName string) {
	inputDataBuilder := build_globula.CreateInputDataBuilder(algType)
	inputData_, err := inputDataBuilder.CreateInputData(algType, predefinedParams, particleName)
	if err != nil {
		printer.PrintlnError(err.Error())
		return
	}
	calcAlg := build_globula.CreateCalcAlg(inputData_, algType)
	finishedPolymers := calcAlg.Calc()
	if finishedPolymers == nil {
		printer.PrintlnError("The result of building is nil")
	} else {
		globula := views.NewGlobulaView(inputData_.GetName(), finishedPolymers, inputData_.GetGlobulaType(), inputData_.GetLiterals())
		globulas = append(globulas, globula)
	}
}

func PrintGlobulaInfo(globula *views.GlobulaView) {
	printer.Println("\n\tGlobula Name: " + globula.Name())
	printer.Println("\tPolymers Count: " + strconv.Itoa(globula.Len()))
	printer.Println("\tPolymers:")
	var monomersCount int
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		monomersCount += pol.Len()
		printer.Println("\t\tPolymer Name: " + pol.Name())
		printer.Println("\t\tMonomers Count: " + strconv.Itoa(pol.Len()))
		printer.PrintEmptyLine()
	})
	printer.Println("\t\tMonomers in total: " + strconv.Itoa(monomersCount))
}

func getCommandsFromScript(filename string, mode string) []string {
	commands := make([]string, 0)

	file, err := os.Open(filename)
	if err != nil {
		printer.PrintlnError(err.Error())
		return commands
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	modeSectionStart := base.Stack{}
	isTakeOnlyMode := (len(mode) > 0)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if isTakeOnlyMode { // if --mode has been defined
			if line[len(line)-1] == '{' { // if it's a label
				modeSectionStart.Push(line) // mark the following commands are within this mode
			} else { // it's not a label
				if !modeSectionStart.IsEmpty() { // the mode is activated
					if line == "}" { // the section is over
						modeSectionStart.Pop()
					} else if modeSectionStart.PeekNotSafe() == mode+"{" {
						commands = append(commands, line) // this is valid command
					} else {
						continue
					}
				} else {
					commands = append(commands, line)
				}
			}
		} else {
			commands = append(commands, line)
		}
	}

	return commands
}

func processPattern(data map[string]string) *views.GlobulaView {
	globulaName := data["globulaName"]
	originGlobula := getGlobulaByName(globulaName)
	outputName := data["outputName"]
	globula := originGlobula.DeepCopy(outputName)
	pattern, ok := pattern_lib.GetPattern(data)
	if !ok {
		return nil
	}

	if pattern_lib.AnyLetterIsUndefined(pattern, globula) {
		printer.PrintlnError("Please, define the missing decryptions to continue")
		return nil
	}

	if globula.Is(views.GLOBULA_GLOBULA_TYPE) {
		pattern_lib.ApplyAsGlobula(globula, pattern)
	} else {
		pattern_lib.ApplyAsThread(globula, pattern)
	}

	return globula
}
