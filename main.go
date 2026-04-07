/*
Package main is the starting point
*/
package main

import (
	"bufio"
	"errors"
	"flag"
	"os"
	"polymers/atomistic"
	"polymers/base"
	"polymers/buildglobula"
	interp "polymers/commandinterpreter"
	"polymers/cycles"
	"polymers/datatypes"
	"polymers/globaldata"
	"polymers/loaders"
	"polymers/outputformat"
	"polymers/patternlib"
	"polymers/savers"
	"polymers/views"
	"strconv"
	"strings"
)

var globula *views.GlobulaView
var printer outputformat.IPrint

func setUpSpaceDimention(commands *[]string) globaldata.SpaceDimention {
	var spaceDimention globaldata.SpaceDimention
	spaceDim := struct{ X, Y, Z, Xl, Yl, Zl float64 }{}
	spaceDim.Xl = 0.0
	spaceDim.Xl = 0.0
	spaceDim.Xl = 0.0
	if len(*commands) == 0 {
		printer.PrintInfo("Please, type the space dimention: ")
		printer.Readln(&spaceDim.X, &spaceDim.Y, &spaceDim.Z)
		spaceDimention[base.AxisX].Lower = 0.0
		spaceDimention[base.AxisY].Lower = 0.0
		spaceDimention[base.AxisZ].Lower = 0.0
		spaceDimention[base.AxisX].Higher = spaceDim.X
		spaceDimention[base.AxisY].Higher = spaceDim.Y
		spaceDimention[base.AxisZ].Higher = spaceDim.Z
		return spaceDimention
	}
	line := (*commands)[0]
	parts := strings.Split(line, " ")
	if (len(parts) != 4 && len(parts) != 7) || parts[0] != "space" {
		printer.PrintInfo("No space dimention definition found. Please, type the space dimention: ")
		printer.Readln(&spaceDim.X, &spaceDim.Y, &spaceDim.Z)
	} else {
		*commands = (*commands)[1:]
		if x, err := strconv.ParseFloat(parts[1], 64); err == nil {
			spaceDim.X = x
		} else {
			printer.PrintlnError(err.Error())
		}
		if y, err := strconv.ParseFloat(parts[2], 64); err == nil {
			spaceDim.Y = y
		} else {
			printer.PrintlnError(err.Error())
		}
		if z, err := strconv.ParseFloat(parts[3], 64); err == nil {
			spaceDim.Z = z
		} else {
			printer.PrintlnError(err.Error())
		}
		if len(parts) == 7 {
			if x, err := strconv.ParseFloat(parts[4], 64); err == nil {
				spaceDim.Xl = x
			} else {
				printer.PrintlnError(err.Error())
			}
			if y, err := strconv.ParseFloat(parts[5], 64); err == nil {
				spaceDim.Yl = y
			} else {
				printer.PrintlnError(err.Error())
			}
			if z, err := strconv.ParseFloat(parts[6], 64); err == nil {
				spaceDim.Zl = z
			} else {
				printer.PrintlnError(err.Error())
			}
		}
	}
	spaceDimention[base.AxisX].Lower = spaceDim.Xl
	spaceDimention[base.AxisY].Lower = spaceDim.Yl
	spaceDimention[base.AxisZ].Lower = spaceDim.Zl
	spaceDimention[base.AxisX].Higher = spaceDim.X
	spaceDimention[base.AxisY].Higher = spaceDim.Y
	spaceDimention[base.AxisZ].Higher = spaceDim.Z
	printer.PrintfInfo("The space dimention set by user is %d\n", spaceDimention)
	return spaceDimention
}

func main() {
	outputformat.SetPrint(&outputformat.ColoredConsolePrint{})
	printer = outputformat.GetPrint()
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
	globaldata.ConfigureGlobalData(spaceDimention)
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
		case interp.CommandUndefined:
			printer.PrintlnError(data.(string))
		case interp.CommandHelp:
			interp.PrintHelp()
		case interp.CommandBuild:
			m := data.(map[string]interface{})
			buildGlobula(m["alg"].(buildglobula.AlgType), m["params"].([]string), m["name"].(string))
		case interp.CommandShowGlobula:
			if globula != nil {
				PrintGlobulaInfo(globula)
			} else {
				printer.PrintlnError("There is no globula called")
			}
		case interp.CommandSaveGlobula:
			if globula == nil {
				printer.PrintlnError("There is no globula called")
				break
			}
			filename := data.(string)
			content, _ := savers.SaveToLammps(globula)
			f, err := os.Create(filename)
			if err != nil {
				printer.PrintlnError(err.Error())
				return
			}
			defer f.Close()
			f.Write([]byte(content))
		case interp.CommandHighlightClustersAll:
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

		case interp.CommandAge:
			data := data.(map[string]interface{})
			groupsCountStr := data["count"].(string)
			doCrosslinks := data["make_crosslinks"].(bool)
			algType := data["alg_type"].(int)
			if globula == nil {
				printer.PrintlnError("There is no globula")
				break
			}
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
			} else if algType == 3 || algType == 4 {
				ncut, err := strconv.Atoi(data["ncut"].(string))
				if err != nil {
					continue
				}
				nOContaining, err := strconv.Atoi(data["OContaining"].(string))
				if err != nil {
					continue
				}
				nOContaining = int(float64(ncut) * float64(nOContaining) / 100.0)
				ncross, err := strconv.Atoi(data["ncross"].(string))
				if err != nil {
					continue
				}
				switch algType {
				case 3:
					globula.DoAging3(ncut, nOContaining, ncross)
				case 4:
					globula.DoAgingSurface(ncut, nOContaining, ncross)
				}
			}

		case interp.CommandReset:
			globula.Reset()

		case interp.CommandResetFull:
			globula.FullReset()

		case interp.CommandHighlightBorders:
			globula.HighlightBorders()

		case interp.CommandWaterize:
			globula.Waterize()

		case interp.CommandTrunk:
			data := data.(map[string]interface{})
			newSize, ok := data["new_size"]
			if ok {
				globula.MakeHomogenousAsCustom(newSize.(int))
			} else {
				globula.MakeHomogenousAsShortest()
			}

		case interp.CommandPattern:
			if _, err := processPattern(data.(map[string]string)); err != nil {
				printer.PrintlnError(err.Error())
			}

		case interp.CommandExit:
			isWorking = false

		case interp.CommandScript:
			filename, err := data.(string)
			if !err {
				printer.PrintlnError("Usage: script <filename>")
			}
			commands = getCommandsFromScript(filename, "")

		case interp.CommandCommonStats:
			data := data.(map[string]string)
			text := globula.GetStatistics()
			f, err := os.Create(data["filename"])
			if err != nil {
				printer.PrintlnError(err.Error())
				return
			}
			defer f.Close()
			f.Write([]byte(text))

		case interp.CommandAtomistic:
			data := data.(map[string]string)
			if config, ok := data["config"]; ok && globula != nil {
				atomistic.MakeAtomistic(globula, config)
			} else {
				printer.PrintflnError("Build a globula first")
			}

		case interp.CommandLoader:
			data := data.(map[string]string)
			filetype := data["filetype"]
			filename := data["filename"]
			fieldTypeStr := data["field"]
			fieldType, ok := map[string]datatypes.FieldType{
				interp.CommandRealSTR:    datatypes.FieldTypeReal,
				interp.CommandLatticeSTR: datatypes.FieldTypeLattice,
			}[fieldTypeStr]
			if !ok {
				printer.PrintflnError("No such a field type: %s", fieldType)
			}
			if loader, err := loaders.NewLoader(filetype, fieldType); err == nil {
				if newGlobula, err := loader.Load(filename, fieldType); err == nil {
					globula = newGlobula
				} else {
					printer.PrintflnError("When loading: %s", err.Error())
				}
			} else {
				printer.PrintflnError("When loading: %s", err.Error())
			}
		case interp.CommandCycles:
			data := data.(map[string]base.Axis)
			var axises []base.Axis
			if axis, ok := data["axis"]; ok {
				axises = []base.Axis{axis}
			} else {
				axises = []base.Axis{base.AxisX, base.AxisY, base.AxisZ}
			}
			cycles.Analyze(globula, axises)

		default:
			printer.PrintlnError("'" + line[:len(line)-1] + "' is not supported")
		}
	}
}

func buildGlobula(algType buildglobula.AlgType, predefinedParams []string, particleName string) {
	inputDataBuilder := buildglobula.CreateInputDataBuilder(algType)
	inputData_, err := inputDataBuilder.CreateInputData(algType, predefinedParams, particleName)
	if err != nil {
		printer.PrintlnError(err.Error())
		return
	}
	calcAlg := buildglobula.CreateCalcAlg(inputData_, algType)
	finishedPolymers := calcAlg.Calc()
	if finishedPolymers == nil {
		printer.PrintlnError("The result of building is nil")
	} else {
		polymers := make([]datatypes.IPolymer, len(finishedPolymers))
		for i := 0; i < len(polymers); i++ {
			polymers[i] = finishedPolymers[i]
		}
		globula = views.NewGlobulaView(polymers, inputData_.GetGlobulaType(), inputData_.GetLiterals())
	}
}

func PrintGlobulaInfo(globula *views.GlobulaView) {
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

func processPattern(data map[string]string) (*views.GlobulaView, error) {
	pattern, ok := patternlib.GetPattern(data)
	if !ok {
		return nil, errors.New("couldn't retrieve pattern")
	}

	if patternlib.AnyLetterIsUndefined(pattern, globula) {
		printer.PrintlnError("Please, define the missing decryptions to continue")
		return nil, errors.New("please, define the missing decryptions to continue")
	}

	if globula.Is(views.GlobulaGlobulaType) {
		patternlib.ApplyAsGlobula(globula, pattern)
	} else {
		patternlib.ApplyAsThread(globula, pattern)
	}

	return globula, nil
}
