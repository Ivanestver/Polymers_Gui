package main

import (
	"fmt"
	"math/rand"
	"os"
	"polymers/build_globula"
	interp "polymers/command_interpreter"
	"polymers/global_data"
	"polymers/savers"
	"polymers/views"
	"strconv"
	"time"
)

var inputData build_globula.CalcAlgInputData
var globulas []*views.GlobulaView

func main() {
	rand.Seed(time.Now().UnixNano())
	fileNumber := 1
	printInfo("Welcome to the Polymer Builder 2.0. Please, type the space dimention: ")
	var spaceDimention int64 = 10
	//fmt.Scanln(&spaceDimention)
	printfInfo("The space dimention set by user is %d\n", spaceDimention)

	printlnInfo("Configuring the global data")
	global_data.ConfigureGlobalData(spaceDimention)
	printlnInfo("Configuring the global data finished")
	inputData.PolymersCount = 5
	inputData.AcceptThreshold = 0.1
	inputData.MaxMonomersCount = 10000
	inputData.SphereRadius = 100
	/*
		inputData.PolymersCount = 3
		inputData.AcceptThreshold = 0.1
		inputData.MaxMonomersCount = 40
		inputData.SphereRadius = 10
	*/
	printlnInfo("The preparations are done! Now you may set up the input data and run the algorithm.")
	//cmdReader := bufio.NewReader(os.Stdin)
	isWorking := true
	commands := make([]string, 0)
	commands = append(commands, "build  ")
	commands = append(commands, "save \"Globula 0\"  ")
	commands = append(commands, "clusters \"Globula 0\" all  ")
	commands = append(commands, "save \"Globula 0\"  ")
	commands = append(commands, "exit  ")
	for isWorking {
		var line string = commands[0]
		commands = commands[1:]
		fmt.Print("> ")
		//line, _ = cmdReader.ReadString('\n')
		command, data := interp.Interpret(line[:len(line)-2])
		switch command {
		case interp.COMMAND_UNDEFINED:
			printlnError(data.(string))
			break
		case interp.COMMAND_HELP:
			printHelp()
			break
		case interp.COMMAND_SET_POLYMERS_COUNT:
			inputData.PolymersCount = data.(int)
			break
		case interp.COMMAND_SET_ACCEPT_THRESHOLD:
			inputData.AcceptThreshold = data.(float64)
			break
		case interp.COMMAND_SET_MAX_MONOMERS_COUNT:
			inputData.MaxMonomersCount = data.(int)
			break
		case interp.COMMAND_SET_SPHERE_RADIUS:
			rad := data.(int)
			if rad > int(spaceDimention) {
				printflnError("Sphere radius must be less or equal space dimention %d", spaceDimention)
			} else {
				inputData.SphereRadius = rad
			}
			break
		case interp.COMMAND_SHOW_PARAMETERS:
			printParams()
			break
		case interp.COMMAND_BUILD_GLOBULA:
			buildGlobula()
			break
		case interp.COMMAND_SHOW_GLOBULAS_LIST:
			for _, globula := range globulas {
				fmt.Println(globula.Name())
			}
			break
		case interp.COMMAND_SHOW_GLOBULA:
			globulaName := data.(string)
			var globula *views.GlobulaView
			for _, glob := range globulas {
				if glob.Name() == globulaName {
					globula = glob
					break
				}
			}
			if globula != nil {
				printGlobulaInfo(globula)
			} else {
				printlnError("There is no globula called \"" + globulaName + "\"")
			}
			break
		case interp.COMMAND_SAVE_GLOBULA:
			globulaName := data.(string)
			var globula *views.GlobulaView
			for _, glob := range globulas {
				if glob.Name() == globulaName {
					globula = glob
					break
				}
			}
			content, _ := savers.SaveToLammps(globula)
			f, err := os.Create(globulaName + strconv.Itoa(fileNumber) + ".data")
			fileNumber++
			if err != nil {
				printlnError(err.Error())
				break
			}
			f.Write([]byte(content))
			// if err := json.NewEncoder(os.Stdout).Encode(globula); err != nil {
			// 	printlnError(err.Error())
			// }
			break
		case interp.COMMAND_HIGHLIGHT_CLUSTERS_ALL:
			globulaName := data.(string)
			var globula *views.GlobulaView
			for _, glob := range globulas {
				if glob.Name() == globulaName {
					globula = glob
					break
				}
			}
			fmt.Println("Start highlighting clusters")
			xClusters, yClusters, zClusters := globula.CommonClusters()
			xClusters.Colorize(false)
			yClusters.Colorize(false)
			zClusters.Colorize(false)
		case interp.COMMAND_EXIT:
			isWorking = false
			break
		default:
			printlnError("'" + line[:len(line)-1] + "' is not supported")
		}
	}
}

func printHelp() {
	printEmptyLine()
	fmt.Println("help - print this article")

	printEmptyLine()
	fmt.Println("set <parameter> <args> - set specific args to a parameter, where <parameter>:")
	fmt.Printf("\t%s <integer> - set polymers count\n", interp.Command_pols_count_str)
	fmt.Printf("\t%s <float> - set threshold\n", interp.Command_threshold_str)
	fmt.Printf("\t%s <integer> - set max monomers count\n", interp.Command_max_mon_count_str)
	fmt.Printf("\t%s <integer> - set sphere radius\n", interp.Command_sphere_rad_str)
	fmt.Println("\twhere")
	fmt.Println("\t\t<integer> - any non-negative integer")
	fmt.Println("\t\t<float> - any non-negatve float-point number")

	printEmptyLine()
	fmt.Println("build - create globula")

	printEmptyLine()
	fmt.Println("show <options> - show different information, where <options>:")
	fmt.Printf("\t%s - show parameters of building\n", interp.COMMAND_SHOW_PARAMETERS_STR)

	printEmptyLine()
}

func printParams() {
	fmt.Printf("\n\tSpace dimention: %d\n", global_data.GetGlobalData().SpaceDimention)
	fmt.Printf("\tPolymers count: %d\n", inputData.PolymersCount)
	fmt.Printf("\tThreshold: %f\n", inputData.AcceptThreshold)
	fmt.Printf("\tMaximum Monomers Count: %d\n", inputData.MaxMonomersCount)
	fmt.Printf("\tSphere Radius: %d\n\n", inputData.SphereRadius)
}

func buildGlobula() {
	calcAlg := build_globula.CreateCalcAlg(inputData)
	finishedPolymers := calcAlg.Calc()
	globulas = append(globulas, views.NewGlobulaView("Globula "+strconv.Itoa(len(globulas)), finishedPolymers))
}

func printGlobulaInfo(globula *views.GlobulaView) {
	fmt.Println("\n\tGlobula Name: " + globula.Name())
	fmt.Println("\tPolymers Count: " + strconv.Itoa(globula.Len()))
	fmt.Println("\tPolymers:")
	var monomersCount int
	views.ForEachPolymer(globula, func(pol *views.PolymerView) {
		monomersCount += pol.Len()
		fmt.Println("\t\tPolymer Name: " + pol.Name())
		fmt.Println("\t\tMonomers Count: " + strconv.Itoa(pol.Len()))
		printEmptyLine()
	})
	fmt.Println("\t\tMonomers in total: " + strconv.Itoa(monomersCount))
}
