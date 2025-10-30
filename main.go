package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"polymers/build_globula"
	interp "polymers/command_interpreter"
	"polymers/datatypes"
	"polymers/global_data"
	"polymers/output_format"
	"polymers/savers"
	"polymers/views"
	"strconv"
	"time"
)

var inputData build_globula.CalcAlgInputData
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
	fileNumber := 1
	output_format.PrintInfo("Welcome to the Polymer Builder 2.0. Please, type the space dimention: ")
	var spaceDimention global_data.SpaceDimention
	spaceDimention.X = 10
	spaceDimention.Y = 10
	spaceDimention.Z = 100
	//fmt.Scanln(&spaceDimention)
	output_format.PrintfInfo("The space dimention set by user is %d\n", spaceDimention)

	output_format.PrintlnInfo("Configuring the global data")
	global_data.ConfigureGlobalData(spaceDimention)
	output_format.PrintlnInfo("Configuring the global data finished")
	inputData.PolymersCount = 5
	inputData.AcceptThreshold = 0.1
	inputData.MaxMonomersCount = 2000
	inputData.SphereRadius = 10
	//inputData.SphereRadius = 100
	/*
		inputData.PolymersCount = 5
		inputData.AcceptThreshold = 0.1
		inputData.MaxMonomersCount = 40
		inputData.SphereRadius = 20
	*/
	output_format.PrintlnInfo("The preparations are done! Now you may set up the input data and run the algorithm.")
	cmdReader := bufio.NewReader(os.Stdin)
	commands := make([]string, 0)
	scriptPtr := flag.String("script", "", "define a script file")
	flag.Parse()
	if scriptPtr != nil {
		commands = getCommandsFromScript(*scriptPtr)
	}
	isWorking := true
	for isWorking {
		fmt.Print("> ")
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
			output_format.PrintlnError(data.(string))
		case interp.COMMAND_HELP:
			PrintHelp()
		case interp.COMMAND_SET_POLYMERS_COUNT:
			inputData.PolymersCount = data.(int)
		case interp.COMMAND_SET_ACCEPT_THRESHOLD:
			inputData.AcceptThreshold = data.(float64)
		case interp.COMMAND_SET_MAX_MONOMERS_COUNT:
			inputData.MaxMonomersCount = data.(int)
		case interp.COMMAND_SET_SPHERE_RADIUS:
			rad := data.(int)
			if rad > int(spaceDimention.X) || rad > int(spaceDimention.Y) || rad > int(spaceDimention.Z) {
				output_format.PrintflnError("Sphere radius must be less or equal space dimention %d", spaceDimention)
			} else {
				inputData.SphereRadius = rad
			}
		case interp.COMMAND_SHOW_PARAMETERS:
			PrintParams()
		case interp.COMMAND_BUILD:
			m := data.(map[string]interface{})
			buildGlobula(m["alg"].(build_globula.AlgType), m["params"].([]string))
		case interp.COMMAND_SHOW_GLOBULAS_LIST:
			for _, globula := range globulas {
				fmt.Println(globula.Name())
			}
		case interp.COMMAND_SHOW_GLOBULA:
			globulaName := data.(string)
			var globula *views.GlobulaView = getGlobulaByName(globulaName)
			if globula != nil {
				PrintGlobulaInfo(globula)
			} else {
				output_format.PrintlnError("There is no globula called \"" + globulaName + "\"")
			}
		case interp.COMMAND_SAVE_GLOBULA:
			globulaName := data.(string)
			var globula *views.GlobulaView = getGlobulaByName(globulaName)
			if globula == nil {
				output_format.PrintlnError("There is no globula called \"" + globulaName + "\"")
				break
			}
			content, _ := savers.SaveToLammps(globula)
			f, err := os.Create(globulaName + strconv.Itoa(fileNumber) + ".data")
			if err != nil {
				output_format.PrintlnError(err.Error())
				return
			}
			defer f.Close()
			f.Write([]byte(content))
			/*f, err = os.Create(globulaName + strconv.Itoa(fileNumber) + ".json")
			if err != nil {
				output_format.PrintlnError(err.Error())
				break
			}
			if err := json.NewEncoder(f).Encode(globula); err != nil {
				output_format.PrintlnError(err.Error())
			}*/
			fileNumber++
		case interp.COMMAND_HIGHLIGHT_CLUSTERS_ALL:
			globulaName := data.(string)
			var originalGlobula *views.GlobulaView = getGlobulaByName(globulaName)
			globula := originalGlobula.DeepCopy(originalGlobula.Name() + "_all_clusters")
			globula.Name()
			globulas = append(globulas, globula)

			fmt.Println("Start highlighting clusters")
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
				output_format.PrintlnError("There is no globula named " + globulaName)
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
			globula := originGlobula.DeepCopy(originGlobula.Name() + "_patterned")
			fileName := data["fileName"]
			file, err := os.Open(fileName)
			if err != nil {
				output_format.PrintlnError(err.Error())
				break
			}
			defer func() {
				if closeErr := file.Close(); closeErr != nil {
					output_format.PrintlnError(closeErr.Error())
				}
			}()
			scanner := bufio.NewScanner(file)
			var pattern string
			if scanner.Scan() {
				pattern = scanner.Text()
			} else {
				output_format.PrintlnError("The given file does not contain a pattern")
				break
			}

			var anyNotExist bool = false
			for _, letter := range pattern {
				l := string(letter)
				if globula.GetMonomerTypeByLiteral(l) == datatypes.MONOMER_TYPE_UNDEFINED {
					output_format.PrintlnError(l + " does not have its decryption")
					anyNotExist = true
				}
			}
			if anyNotExist {
				output_format.PrintlnError("Please, define the missing decryptions to continue")
				break
			}

			views.ForEachPolymer(globula, func(pv *views.PolymerView) {
				currentLetterNumber := 0
				views.ForEachMonomer(pv, func(m *datatypes.Monomer) {
					m.MonomerType = globula.GetMonomerTypeByLiteral(string(pattern[currentLetterNumber]))
					currentLetterNumber = (currentLetterNumber + 1) % len(pattern)
				})
			})

			globulas = append(globulas, globula)

		case interp.COMMAND_EXIT:
			isWorking = false

		case interp.COMMAND_SCRIPT:
			filename, err := data.(string)
			if !err {
				output_format.PrintlnError("Usage: script <filename>")
			}
			commands = getCommandsFromScript(filename)

		default:
			output_format.PrintlnError("'" + line[:len(line)-1] + "' is not supported")
		}
	}
}

func PrintHelp() {
	output_format.PrintEmptyLine()
	fmt.Println("help - output_format.Print this article")

	output_format.PrintEmptyLine()
	fmt.Println("set <parameter> <args> - set specific args to a parameter, where <parameter>:")
	fmt.Printf("\t%s <integer> - set polymers count\n", interp.Command_pols_count_str)
	fmt.Printf("\t%s <float> - set threshold\n", interp.Command_threshold_str)
	fmt.Printf("\t%s <integer> - set max monomers count\n", interp.Command_max_mon_count_str)
	fmt.Printf("\t%s <integer> - set sphere radius\n", interp.Command_sphere_rad_str)
	fmt.Println("\twhere")
	fmt.Println("\t\t<integer> - any non-negative integer")
	fmt.Println("\t\t<float> - any non-negatve float-point number")

	output_format.PrintEmptyLine()
	fmt.Println("build - create globula")

	output_format.PrintEmptyLine()
	fmt.Println("show <options> - show different information, where <options>:")
	fmt.Printf("\t%s - show parameters of building\n", interp.COMMAND_PARAMETERS_STR)

	output_format.PrintEmptyLine()
}

func PrintParams() {
	fmt.Printf("\n\tSpace dimention: %d\n", global_data.GetGlobalData().SpaceDimention)
	fmt.Printf("\tPolymers count: %d\n", inputData.PolymersCount)
	fmt.Printf("\tThreshold: %f\n", inputData.AcceptThreshold)
	fmt.Printf("\tMaximum Monomers Count: %d\n", inputData.MaxMonomersCount)
	fmt.Printf("\tSphere Radius: %d\n\n", inputData.SphereRadius)
}

func buildGlobula(algType build_globula.AlgType, predefinedParams []string) {
	inputDataBuilder := build_globula.CreateInputDataBuilder(algType)
	inputData_, err := inputDataBuilder.CreateInputData(algType, predefinedParams)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	calcAlg := build_globula.CreateCalcAlg(inputData_, algType)
	finishedPolymers := calcAlg.Calc()
	if finishedPolymers == nil {
		fmt.Printf("The result of building is nil")
	} else {
		globula := views.NewGlobulaView(inputData_.GetName()+" "+strconv.Itoa(len(globulas)), finishedPolymers, inputData_.GetGlobulaType())
		globula.SetLiterals(build_globula.GetLiteralsTable())
		globulas = append(globulas, globula)
	}
}

func PrintGlobulaInfo(globula *views.GlobulaView) {
	fmt.Println("\n\tGlobula Name: " + globula.Name())
	fmt.Println("\tPolymers Count: " + strconv.Itoa(globula.Len()))
	fmt.Println("\tPolymers:")
	var monomersCount int
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		monomersCount += pol.Len()
		fmt.Println("\t\tPolymer Name: " + pol.Name())
		fmt.Println("\t\tMonomers Count: " + strconv.Itoa(pol.Len()))
		output_format.PrintEmptyLine()
	})
	fmt.Println("\t\tMonomers in total: " + strconv.Itoa(monomersCount))
}

func getCommandsFromScript(filename string) []string {
	commands := make([]string, 0)

	file, err := os.Open(filename)
	if err != nil {
		output_format.PrintlnError(err.Error())
		return commands
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		commands = append(commands, line)
	}

	return commands
}
