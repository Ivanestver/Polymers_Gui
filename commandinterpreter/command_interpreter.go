package commandinterpreter

import (
	"errors"
	"fmt"
	"polymers/base"
	"polymers/buildglobula"
	"strconv"
	"strings"
	"unicode"
)

const (
	CommandHelpSTR             = "help"
	CommandBuildSTR            = "build"
	CommandExitSTR             = "exit"
	CommandShowSTR             = "show"
	CommandGlobulaSTR          = "globula"
	CommandSaveSTR             = "save"
	CommandClustersSTR         = "clusters"
	CommandClustersAllATR      = "all"
	CommandAgeSTR              = "age"
	CommandResetSTR            = "reset"
	CommandFullSTR             = "full"
	CommandHighlightBordersSTR = "highlight_borders"
	CommandThreadSTR           = "thread"
	CommandSurfaceSTR          = "surface"
	CommandWaterizeSTR         = "waterize"
	CommandTrunkSTR            = "trunk"
	CommandPatternSTR          = "pattern"
	CommandScriptSTR           = "script"
	CommandCommonStatsSTR      = "common_stats"
	CommandFileSTR             = "file"
	CommandCyclesSTR           = "cycles"
	CommandAtomisticSTR        = "atomistic"
	CommandLoaderSTR           = "load"
	CommandLatticeSTR          = "lattice"
	CommandRealSTR             = "real"
	CommandSpaceSTR            = "space"
	CommandCristallinitySTR    = "cristallinity"
	CommandCommentSTR          = '#'
	CommandTrajectoriesSTR     = "traj"
	CommandColorizeSTR         = "colorize"
)

type Command = int

const (
	CommandUndefined = -1
	CommandHelp      = iota
	CommandBuild
	CommandShowGlobula
	CommandSaveGlobula
	CommandHighlightClustersAll
	CommandAge
	CommandReset
	CommandResetFull
	CommandExit
	CommandReadData
	CommandHighlightBorders
	CommandWaterize
	CommandTrunk
	CommandPattern
	CommandScript
	CommandCommonStats
	CommandAtomistic
	CommandLoader
	CommandCycles
	CommandSpace
	CommandCristallinity
	CommandTrajectories
	CommandColorize
)

var currProgram string
var currChar int = 0

func moveForward() {
	currChar++
}

func finished() bool {
	return currChar == len(currProgram)
}

func getCurrChar() byte {
	return currProgram[currChar]
}

func reset() {
	currChar = 0
}

func getNextToken() (string, error) {
	if finished() {
		return "", errors.New("incompleted command")
	}
	var token strings.Builder
	char := rune(getCurrChar())
	for char == ' ' {
		moveForward()
		if finished() {
			return "", errors.New("no command was found")
		}
		char = rune(getCurrChar())
	}
	for {
		if unicode.IsLetter(char) || unicode.IsNumber(char) ||
			char == '_' ||
			char == '.' ||
			char == '%' ||
			char == '(' || char == ')' ||
			char == '*' ||
			char == ',' ||
			char == '-' ||
			char == '/' {
			token.WriteString(string(getCurrChar()))
			moveForward()
			if finished() {
				break
			}
			char = rune(getCurrChar())
		} else if char == CommandCommentSTR {
			for !finished() {
				moveForward()
			}
			break
		} else {
			break
		}
	}
	return token.String(), nil
}

func getParameterAsString() (string, error) {
	getNextToken() // skip all the empty spaces until " or the end
	if finished() {
		return "", nil
	}

	if string(getCurrChar()) != "\"" {
		return "", errors.New("wrong parameter")
	}
	moveForward()
	var globulaName string
	for !finished() {
		token, err := getNextToken()
		if err != nil {
			return "", err
		}
		if len(globulaName) == 0 {
			globulaName += token
		} else {
			globulaName += " " + token
		}
		if string(getCurrChar()) == "\"" {
			moveForward()
			return globulaName, nil
		}
	}
	return "", errors.New("globula name must be wrapped with \"\"")
}

func getUndefinedCommand(token string) (Command, string) {
	return CommandUndefined, "Undefined parameter '" + token + "'"
}

func Interpret(program string) (Command, any) {
	reset()
	currProgram = program
	return s()
}

func s() (Command, any) {
	token, error := getNextToken()
	if error != nil {
		return CommandUndefined, error.Error()
	}

	if f, ok := map[string]func() (Command, any){
		CommandHelpSTR:             func() (Command, any) { return CommandHelp, nil },
		CommandPatternSTR:          pattern,
		CommandBuildSTR:            build,
		CommandShowSTR:             show,
		CommandSaveSTR:             save,
		CommandClustersSTR:         clusters,
		CommandAgeSTR:              age,
		CommandHighlightBordersSTR: borders,
		CommandResetSTR:            resetGlobula,
		CommandExitSTR:             func() (Command, any) { return CommandExit, nil },
		CommandWaterizeSTR:         waterize,
		CommandTrunkSTR:            trunk,
		CommandScriptSTR:           script,
		CommandCommonStatsSTR:      commonStats,
		CommandAtomisticSTR:        atomistic,
		CommandLoaderSTR:           load,
		CommandCyclesSTR:           cycles,
		CommandSpaceSTR:            space,
		CommandCristallinitySTR:    cristallinity,
		CommandTrajectoriesSTR:     trajectories,
		CommandColorizeSTR:         colorize,
	}[token]; ok {
		return f()
	} else {
		return CommandUndefined, "Undefined command: " + token
	}
}

func pattern() (Command, any) {
	m := make(map[string]string)

	t, err := getNextToken()
	if err != nil {
		return CommandUndefined, "Wrong source type"
	}
	switch t {
	case CommandFileSTR:
		fileName, err := getParameterAsString()
		if err != nil {
			return CommandUndefined, "Wrong usage"
		}
		m["fileName"] = fileName
	case CommandPatternSTR:
		patt, err := getParameterAsString()
		if err != nil {
			return CommandUndefined, "Wrong usage"
		}
		if patt, err = getPattern(patt); err == nil {
			m["pattern"] = patt
		} else {
			return CommandUndefined, err.Error()
		}
	default:
		return CommandUndefined, "Wrong source type"
	}

	outputName, err := getParameterAsString()
	if err != nil {
		return CommandUndefined, err.Error()
	}

	m["outputName"] = outputName

	return CommandPattern, m
}

func build() (Command, any) {
	objective, err := getNextToken()
	if err != nil {
		return CommandUndefined, err.Error()
	}

	predefinedParams := make([]string, 0)
	for !finished() {
		p, err := getNextToken()
		if err != nil {
			continue
		}
		predefinedParams = append(predefinedParams, p)
	}

	m := make(map[string]any)
	m["params"] = predefinedParams

	name, err := getParameterAsString()
	if err != nil {
		return CommandUndefined, "Specify the name of a particle"
	}
	m["name"] = name

	if objective == CommandGlobulaSTR {
		m["alg"] = buildglobula.GlobulaBuildAlg
		return CommandBuild, m
	}

	if objective == CommandThreadSTR {
		m["alg"] = buildglobula.ThreadBuildAlg
		return CommandBuild, m
	}

	if objective == CommandSurfaceSTR {
		m["alg"] = buildglobula.SurfaceBuildAlg
		return CommandBuild, m
	}

	return getUndefinedCommand(objective)
}

func show() (Command, any) {
	token, error := getNextToken()
	if error != nil {
		return CommandUndefined, error.Error()
	}

	if token == CommandGlobulaSTR {
		return showGlobula()
	}

	return CommandUndefined, "Undefined parameter '" + token + "'"
}

func showGlobula() (Command, any) {
	return CommandShowGlobula, nil
}

func save() (Command, any) {
	filename, err := getParameterAsString()
	if err != nil {
		return CommandUndefined, err.Error()
	}
	return CommandSaveGlobula, filename
}

func clusters() (Command, any) {
	getNextToken() // skip all the empty spaces until " or the end
	if finished() {
		return CommandUndefined, "Wrong usage"
	}

	tokenAll, err := getNextToken()
	if err != nil {
		return CommandUndefined, err.Error()
	}

	if tokenAll == CommandClustersAllATR {
		return CommandHighlightClustersAll, globulaName
	}

	return CommandUndefined, "Wrong usage"
}

func age() (Command, any) {
	groupCount, err := getNextToken()
	if err != nil {
		return CommandUndefined, err
	}

	moveForward()
	if finished() {
		return CommandUndefined, string("Wrong usage")
	}
	doCrosslinks, err := getNextToken()
	if err != nil {
		return CommandUndefined, err
	}

	if doCrosslinks != "true" && doCrosslinks != "false" {
		return CommandUndefined, string("Error: make_crosslinks parameter must be either \"true\" or \"false\"")
	}

	token, err := getNextToken()
	if err != nil {
		return CommandUndefined, err
	}
	ageAlgType, err := strconv.Atoi(token)
	if err != nil {
		return CommandUndefined, err
	}

	m := make(map[string]any)
	// For Age Algorithm Type 3 and 4
	if !finished() && (ageAlgType == 3 || ageAlgType == 4) {
		ncut, err := getNextToken()
		if err != nil {
			return CommandUndefined, err
		}
		m["ncut"] = ncut
		if finished() {
			return CommandUndefined, errors.New("wrong usage of the age command")
		}
		nOContaining, err := getNextToken()
		if err != nil {
			return CommandUndefined, err
		}
		m["OContaining"] = nOContaining
		if finished() {
			return CommandUndefined, errors.New("wrong usage of the age command")
		}
		ncross, err := getNextToken()
		if err != nil {
			return CommandUndefined, err
		}
		m["ncross"] = ncross
	}

	m["count"] = groupCount
	m["make_crosslinks"] = (doCrosslinks == "true")
	m["alg_type"] = ageAlgType
	return CommandAge, m
}

func borders() (Command, any) {
	return CommandHighlightBorders, nil
}

func resetGlobula() (Command, any) {
	getNextToken() // skip empty spaces
	if finished() {
		m := make(map[string]any)
		m["full"] = false
		return CommandReset, m
	}
	moveForward()
	token, err := getNextToken()
	if err != nil {
		return CommandUndefined, err.Error()
	}

	if token == CommandFullSTR {
		m := make(map[string]any)
		m["full"] = true
		return CommandResetFull, m
	}

	return CommandUndefined, string("Wrong usage")
}

func waterize() (Command, any) {
	return CommandWaterize, nil
}

func trunk() (Command, any) {
	m := make(map[string]any)
	if finished() { // No size has been provided, therefore, use the shortest
		return CommandTrunk, m
	}

	token, err := getNextToken()
	if err != nil {
		return CommandUndefined, err.Error()
	}

	newSize, err := strconv.Atoi(token)
	if err != nil {
		return CommandUndefined, err.Error()
	}

	m["new_size"] = newSize
	return CommandTrunk, m
}

func script() (Command, any) {
	filename, err := getNextToken()
	if err != nil {
		return CommandUndefined, err.Error()
	}
	return CommandScript, filename
}

func commonStats() (Command, any) {
	m := make(map[string]string)

	filename, err := getParameterAsString()
	if err != nil {
		return CommandUndefined, string("Wrong usage")
	}

	m["filename"] = filename
	return CommandCommonStats, m
}

func getPattern(raw string) (string, error) {
	dst := ""
	curr := 0
	for curr < len(raw) {
		c := raw[curr]
		if 'A' <= c && c <= 'Z' {
			dst += string(c)
			curr++
		} else if c == '(' {
			curr++
			if err := getPatternS2(raw, &curr, &dst); err != nil {
				return "", nil
			}
		} else {
			return "", errors.New("The pattern doesn't follow the rules: " + raw)
		}
	}
	return dst, nil
}

func getPatternS2(raw string, curr *int, dst *string) error {
	if *curr >= len(raw) {
		return errors.New("The pattern doesn't follow the rules: " + raw)
	}
	start := *curr
	for {
		c := raw[*curr]
		if 'A' <= c && c <= 'Z' {
			*curr++
		} else if c == ')' && start < *curr {
			*curr++
			return getPatternS3(raw, curr, start, dst)
		} else {
			return errors.New("The pattern doesn't follow the rules: " + raw)
		}
	}
}

func getPatternS3(raw string, curr *int, start int, dst *string) error {
	if *curr >= len(raw) {
		return errors.New("The pattern doesn't follow the rules: " + raw)
	}
	c := raw[*curr]
	if 'A' <= c && c <= 'Z' {
		*dst += string(c)
		*curr++
		return nil
	} else if c == '*' {
		*curr++
		return getPatternS4(raw, curr, start, dst)
	} else {
		return errors.New("The pattern doesn't follow the rules: " + raw)
	}
}

func getPatternS4(raw string, curr *int, start int, dst *string) error {
	startNum := *curr
	for *curr < len(raw) {
		c := raw[*curr]
		if '0' <= c && c <= '9' {
			*curr++
		} else {
			break
		}
	}
	if startNum == *curr {
		return errors.New("The pattern doesn't follow the rules: " + raw)
	} else {
		if itersCount, err := strconv.Atoi(raw[startNum:*curr]); err == nil {
			for i := 0; i < itersCount; i++ {
				*dst += raw[start : startNum-2]
			}
			return nil
		} else {
			return err
		}
	}
}

func atomistic() (Command, any) {
	m := make(map[string]string)

	fileName, err := getParameterAsString()
	if err != nil {
		return CommandUndefined, err.Error()
	}
	m["config"] = fileName
	return CommandAtomistic, m
}

func load() (Command, any) {
	filetype, err := getNextToken()
	if err != nil {
		return CommandUndefined, "no filetype specified"
	}

	m := make(map[string]string)
	m["filetype"] = filetype

	filename, err := getParameterAsString()
	if err != nil {
		return CommandUndefined, "no filename specified"
	}
	m["filename"] = filename
	m["field"] = CommandRealSTR

	token, err := getNextToken()
	if err != nil {
		return CommandLoader, m
	}
	if token == CommandLatticeSTR {
		m["field"] = CommandLatticeSTR
	}

	return CommandLoader, m
}

func cycles() (Command, any) {
	m := make(map[string]any)
	token, err := getNextToken()
	if err != nil {
		return CommandCycles, m
	}
	steps, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return CommandUndefined, err.Error()
	}
	m["steps"] = int(steps)
	token, err = getNextToken()
	if err != nil {
		return CommandCycles, m
	}
	switch token {
	case "X":
		m["axis"] = base.AxisX
	case "Y":
		m["axis"] = base.AxisY
	case "Z":
		m["axis"] = base.AxisZ
	default:
		return CommandUndefined, fmt.Sprintf("axis can be only 'X', 'Y', 'Z', but given %s", token)
	}
	return CommandCycles, m
}

func space() (Command, any) {
	token, err := getNextToken()
	if err != nil {
		return CommandUndefined, err
	}
	xHigher, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return CommandUndefined, err
	}
	token, err = getNextToken()
	if err != nil {
		return CommandUndefined, err
	}
	yHigher, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return CommandUndefined, err
	}
	token, err = getNextToken()
	if err != nil {
		return CommandUndefined, err
	}
	zHigher, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return CommandUndefined, err
	}

	m := make(map[string]float64)
	m["x_higher"] = xHigher
	m["y_higher"] = yHigher
	m["z_higher"] = zHigher

	token, err = getNextToken()
	if err != nil {
		return CommandSpace, m
	}
	xLower, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return CommandSpace, m
	}
	token, err = getNextToken()
	if err != nil {
		return CommandSpace, m
	}
	yLower, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return CommandSpace, m
	}
	token, err = getNextToken()
	if err != nil {
		return CommandSpace, m
	}
	zLower, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return CommandSpace, m
	}
	m["x_lower"] = xLower
	m["y_lower"] = yLower
	m["z_lower"] = zLower
	return CommandSpace, m
}

func cristallinity() (Command, any) {
	token, err := getNextToken()
	if err != nil {
		return CommandCristallinity, struct{}{}
	}
	m := make(map[string]any)
	m["level"] = token

	token, err = getNextToken()
	if err != nil {
		return CommandCristallinity, m
	}
	offset, err := strconv.Atoi(token)
	if err != nil {
		return CommandCristallinity, m
	}
	m["offset"] = offset

	token, err = getParameterAsString()
	if err != nil {
		return CommandCristallinity, m
	}
	m["output_filename"] = token

	token, err = getNextToken()
	if err != nil {
		return CommandCristallinity, m
	}
	topPercent, err := strconv.Atoi(token)
	if err != nil {
		return CommandCristallinity, m
	}
	m["top_percent"] = topPercent

	token, err = getParameterAsString()
	if err != nil {
		return CommandCristallinity, m
	}
	m["base_elem"] = base.RecognizeElement(token)

	return CommandCristallinity, m
}

func trajectories() (Command, any) {
	files := make([]string, 0)
	token, err := getParameterAsString()
	for len(token) > 0 && err == nil {
		files = append(files, token)
		token, err = getParameterAsString()
	}
	if len(files) == 0 {
		return CommandUndefined, "необходимо указать хотя бы один файл траекторий"
	}
	m := make(map[string][]string)
	m["traj_filename"] = files
	return CommandTrajectories, m
}

func colorize() (Command, any) {
	return CommandColorize, nil
}
