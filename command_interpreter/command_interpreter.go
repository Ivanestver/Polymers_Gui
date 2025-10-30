package command_interpreter

import (
	"errors"
	"polymers/build_globula"
	"strconv"
	"unicode"
)

const (
	Command_help_str              = "help"
	Command_set_str               = "set"
	Command_build_str             = "build"
	COMMAND_EXIT_STR              = "exit"
	COMMAND_SHOW_STR              = "show"
	COMMAND_GLOBULA_STR           = "globula"
	COMMAND_SAVE_STR              = "save"
	COMMAND_CLUSTERS_STR          = "clusters"
	COMMAND_CLUSTERS_ALL_STR      = "all"
	COMMAND_AGE_STR               = "age"
	COMMAND_RESET_STR             = "reset"
	COMMAND_FULL_STR              = "full"
	COMMAND_HIGHLIGHT_BORDERS_STR = "highlight_borders"
	COMMAND_THREAD_STR            = "thread"
	COMMAND_WATERIZE_STR          = "waterize"
	COMMAND_TRUNK_STR             = "trunk"
	COMMAND_PATTERN_STR           = "pattern"
	COMMAND_SCRIPT_STR            = "script"
	COMMAND_COMMON_STATS_STR      = "common_stats"
)

type Command = int

const (
	COMMAND_UNDEFINED = -1
	COMMAND_HELP      = iota
	COMMAND_BUILD
	COMMAND_SHOW_GLOBULAS_LIST
	COMMAND_SHOW_GLOBULA
	COMMAND_SAVE_GLOBULA
	COMMAND_HIGHLIGHT_CLUSTERS_ALL
	COMMAND_AGE
	COMMAND_RESET
	COMMAND_RESET_FULL
	COMMAND_EXIT
	COMMAND_READ_DATA
	COMMAND_HIGHLIGHT_BORDERS
	COMMAND_WATERIZE
	COMMAND_TRUNK
	COMMAND_PATTERN
	COMMAND_SCRIPT
	COMMAND_COMMON_STATS
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
		return "", errors.New("Uncompleted command")
	}
	var token string
	var char rune = rune(getCurrChar())
	for char == ' ' {
		moveForward()
		if finished() {
			return "", errors.New("No command was found")
		}
		char = rune(getCurrChar())
	}
	for unicode.IsLetter(char) || unicode.IsNumber(char) || char == '_' || char == '.' || char == '%' {
		token += string(getCurrChar())
		moveForward()
		if finished() {
			break
		}
		char = rune(getCurrChar())
	}
	return token, nil
}

func getParameterAsString() (string, error) {
	getNextToken() // skip all the empty spaces until " or the end
	if finished() {
		return "", nil
	}

	if string(getCurrChar()) != "\"" {
		return "", errors.New("Usage: show globula \"<globula name>\"")
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
	return "", errors.New("Globula name must be wrapped with \"\"")
}

func getUndefinedCommand(token string) (Command, string) {
	return COMMAND_UNDEFINED, "Undefined parameter '" + token + "'"
}

func Interpret(program string) (Command, interface{}) {
	reset()
	currProgram = program
	return s()
}

func s() (Command, interface{}) {
	token, error := getNextToken()
	if error != nil {
		return COMMAND_UNDEFINED, error.Error()
	}

	if token == Command_help_str {
		return COMMAND_HELP, nil
	}

	if token == Command_set_str {
		return set()
	}

	if token == Command_build_str {
		return build()
	}

	if token == COMMAND_SHOW_STR {
		return show()
	}

	if token == COMMAND_SAVE_STR {
		return save()
	}

	if token == COMMAND_CLUSTERS_STR {
		return clusters()
	}

	if token == COMMAND_AGE_STR {
		return age()
	}

	if token == COMMAND_HIGHLIGHT_BORDERS_STR {
		return borders()
	}

	if token == COMMAND_RESET_STR {
		return resetGlobula()
	}

	if token == COMMAND_EXIT_STR {
		return COMMAND_EXIT, nil
	}

	if token == COMMAND_WATERIZE_STR {
		return waterize()
	}

	if token == COMMAND_TRUNK_STR {
		return trunk()
	}

	if token == COMMAND_SCRIPT_STR {
		return script()
	}

	if token == COMMAND_COMMON_STATS_STR {
		return commonStats()
	}

	return COMMAND_UNDEFINED, "Undefined command: " + token
}

func set() (Command, interface{}) {
	token, error := getNextToken()
	if error != nil {
		return COMMAND_UNDEFINED, error.Error()
	}

	if token == COMMAND_PATTERN_STR {
		return pattern()
	}

	return COMMAND_UNDEFINED, "Undefined parameter: " + token
}

func pattern() (Command, interface{}) {
	globulaName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, "Usage: set pattern <globula_name> <file_name>"
	}

	fileName, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, "Usage: set pattern <globula_name> <file_name>"
	}

	m := make(map[string]string)
	m["globulaName"] = globulaName
	m["fileName"] = fileName
	return COMMAND_PATTERN, m
}

func build() (Command, interface{}) {
	objective, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}

	predefinedParams := make([]string, 0)
	for !finished() {
		p, err := getNextToken()
		if err != nil {
			continue
		}
		predefinedParams = append(predefinedParams, p)
	}

	m := make(map[string]interface{})
	m["params"] = predefinedParams

	if objective == COMMAND_GLOBULA_STR {
		m["alg"] = build_globula.GlobulaBuildAlg
		return COMMAND_BUILD, m
	}

	if objective == COMMAND_THREAD_STR {
		m["alg"] = build_globula.ThreadBuildAlg
		return COMMAND_BUILD, m
	}

	return getUndefinedCommand(objective)
}

func show() (Command, interface{}) {
	token, error := getNextToken()
	if error != nil {
		return COMMAND_UNDEFINED, error.Error()
	}

	if token == COMMAND_GLOBULA_STR {
		return showGlobula()
	}

	return COMMAND_UNDEFINED, "Undefined parameter '" + token + "'"
}

func showGlobula() (Command, interface{}) {
	getNextToken() // skip all the empty spaces until " or the end
	if finished() {
		return COMMAND_SHOW_GLOBULAS_LIST, nil
	}

	if string(getCurrChar()) != "\"" {
		return COMMAND_UNDEFINED, "Usage: show globula \"<globula name>\""
	}
	moveForward()
	var globulaName string
	for !finished() {
		token, error := getNextToken()
		if error != nil {
			return COMMAND_UNDEFINED, error.Error()
		}
		if len(globulaName) == 0 {
			globulaName += token
		} else {
			globulaName += " " + token
		}
		if string(getCurrChar()) == "\"" {
			return COMMAND_SHOW_GLOBULA, globulaName
		}
	}
	return COMMAND_UNDEFINED, "Globula name must be wrapped with \"\""
}

func save() (Command, interface{}) {
	getNextToken() // skip all the empty spaces until " or the end
	if finished() {
		return COMMAND_UNDEFINED, "Usage: show globula \"<globula name>\""
	}

	if string(getCurrChar()) != "\"" {
		return COMMAND_UNDEFINED, "Usage: show globula \"<globula name>\""
	}
	moveForward()
	var globulaName string
	for !finished() {
		token, error := getNextToken()
		if error != nil {
			return COMMAND_UNDEFINED, error.Error()
		}
		if len(globulaName) == 0 {
			globulaName += token
		} else {
			globulaName += " " + token
		}
		if string(getCurrChar()) == "\"" {
			return COMMAND_SAVE_GLOBULA, globulaName
		}
	}
	return COMMAND_UNDEFINED, "Globula name must be wrapped with \"\""
}

func clusters() (Command, interface{}) {
	getNextToken() // skip all the empty spaces until " or the end
	if finished() {
		return COMMAND_UNDEFINED, "Usage: clusters \"<globula name>\" <param>"
	}

	if string(getCurrChar()) != "\"" {
		return COMMAND_UNDEFINED, "Usage: clusters \"<globula name>\" <param>"
	}
	moveForward()
	var globulaName string
	for !finished() {
		token, error := getNextToken()
		if error != nil {
			return COMMAND_UNDEFINED, error.Error()
		}
		if len(globulaName) == 0 {
			globulaName += token
		} else {
			globulaName += " " + token
		}
		if string(getCurrChar()) == "\"" {
			moveForward()
			break
		}
	}

	tokenAll, err := getNextToken()
	print(tokenAll)
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}

	if tokenAll == COMMAND_CLUSTERS_ALL_STR {
		return COMMAND_HIGHLIGHT_CLUSTERS_ALL, globulaName
	}

	return COMMAND_UNDEFINED, "Usage: clusters \"<globula name>\" <param>"
}

func age() (Command, interface{}) {
	globulaName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	if finished() {
		return COMMAND_UNDEFINED, string("Usage: age <globula_name> <groups_count> <make_crosslinks>")
	}
	groupCount, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err
	}

	moveForward()
	if finished() {
		return COMMAND_UNDEFINED, string("Usage: age <globula_name> <groups_count> <make_crosslinks>")
	}
	token, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err
	}

	if token != "true" && token != "false" {
		return COMMAND_UNDEFINED, string("Error: make_crosslinks parameter must be either \"true\" or \"false\"")
	}

	m := make(map[string]interface{})
	m["globula"] = globulaName
	m["count"] = groupCount
	m["make_crosslinks"] = (token == "true")
	return COMMAND_AGE, m
}

func borders() (Command, interface{}) {
	globulaName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	m := make(map[string]interface{})
	m["globula"] = globulaName
	return COMMAND_HIGHLIGHT_BORDERS, m
}

func resetGlobula() (Command, interface{}) {
	globulaName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	getNextToken() // skip empty spaces
	if finished() {
		m := make(map[string]interface{})
		m["globula"] = globulaName
		m["full"] = false
		return COMMAND_RESET, m
	}
	moveForward()
	token, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}

	if token == COMMAND_FULL_STR {
		m := make(map[string]interface{})
		m["globula"] = globulaName
		m["full"] = false
		return COMMAND_RESET_FULL, m
	}

	return COMMAND_UNDEFINED, string("Usage: reset <globula_name> full")
}

func waterize() (Command, interface{}) {
	globulaName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	m := make(map[string]interface{})
	m["globula"] = globulaName
	return COMMAND_WATERIZE, m
}

func trunk() (Command, interface{}) {
	globulaName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	getNextToken()
	m := make(map[string]interface{})
	m["globula"] = globulaName
	if finished() { // No size has been provided, therefore, use the shortest
		return COMMAND_TRUNK, m
	}

	token, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}

	newSize, err := strconv.Atoi(token)
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}

	m["new_size"] = newSize
	return COMMAND_TRUNK, m
}

func script() (Command, interface{}) {
	filename, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	return COMMAND_SCRIPT, filename
}

func commonStats() (Command, interface{}) {
	globulaName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	getNextToken()
	m := make(map[string]string)
	m["globula"] = globulaName

	filename, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, string("Usage: command_stats <Globula Name> <Output File Name>")
	}

	m["filename"] = filename
	return COMMAND_COMMON_STATS, m
}
