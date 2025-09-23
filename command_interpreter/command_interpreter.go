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
	Command_pols_count_str        = "pols_count"
	Command_threshold_str         = "threshold"
	Command_max_mon_count_str     = "max_mon_count"
	Command_sphere_rad_str        = "sphere_rad"
	Command_build_str             = "build"
	COMMAND_EXIT_STR              = "exit"
	COMMAND_SHOW_STR              = "show"
	COMMAND_PARAMETERS_STR        = "params"
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
)

type Command = int

const (
	COMMAND_UNDEFINED = -1
	COMMAND_HELP      = iota
	COMMAND_SET_POLYMERS_COUNT
	COMMAND_SET_ACCEPT_THRESHOLD
	COMMAND_SET_MAX_MONOMERS_COUNT
	COMMAND_SET_SPHERE_RADIUS
	COMMAND_BUILD_GLOBULA
	COMMAND_BUILD_THREAD_GLOBULA
	COMMAND_SHOW_PARAMETERS
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
	for unicode.IsLetter(char) || unicode.IsNumber(char) || char == '_' {
		token += string(getCurrChar())
		moveForward()
		if finished() {
			break
		}
		char = rune(getCurrChar())
	}
	return token, nil
}

func getGlobulaName() (string, error) {
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

	if token == COMMAND_EXIT_STR {
		return COMMAND_EXIT, nil
	}

	if token == COMMAND_WATERIZE_STR {
		return waterize()
	}

	if token == COMMAND_TRUNK_STR {
		return trunk()
	}

	return COMMAND_UNDEFINED, "Undefined command: " + token
}

func set() (Command, interface{}) {
	token, error := getNextToken()
	if error != nil {
		return COMMAND_UNDEFINED, error.Error()
	}

	if token == Command_pols_count_str {
		return interpret_with_num(COMMAND_SET_POLYMERS_COUNT, token)
	}

	if token == Command_max_mon_count_str {
		return interpret_with_num(COMMAND_SET_MAX_MONOMERS_COUNT, token)
	}

	if token == Command_sphere_rad_str {
		return interpret_with_num(COMMAND_SET_SPHERE_RADIUS, token)
	}

	if token == Command_threshold_str {
		token, error = getNextToken()
		if error != nil {
			return COMMAND_UNDEFINED, error.Error()
		}
		if _, error := strconv.Atoi(token); error != nil {
			return COMMAND_UNDEFINED, error.Error()
		}
		dot := getCurrChar()
		if dot != '.' {
			return COMMAND_UNDEFINED, errors.New("Threshold must be float")
		}
		moveForward()
		d_num_str := token + string(dot)
		token, error = getNextToken()
		if error != nil {
			return COMMAND_UNDEFINED, error.Error()
		}
		if num, error := strconv.Atoi(token); error != nil {
			return COMMAND_UNDEFINED, error.Error()
		} else if num < 0 {
			return COMMAND_UNDEFINED, Command_threshold_str + " expects a non-negative number"
		}
		d_num_str += token
		if d_num, error := strconv.ParseFloat(d_num_str, 64); error != nil {
			return COMMAND_UNDEFINED, error.Error()
		} else {
			return COMMAND_SET_ACCEPT_THRESHOLD, d_num
		}
	}

	return COMMAND_UNDEFINED, "Undefined parameter: " + token
}

func interpret_with_num(comm Command, usage string) (Command, interface{}) {
	token, error := getNextToken()
	if error != nil {
		return COMMAND_UNDEFINED, error.Error()
	}
	if !finished() {
		return COMMAND_UNDEFINED, "Usage: set " + usage + " <integer>"
	}
	if num, error := strconv.Atoi(token); error == nil {
		if num < 0 {
			return COMMAND_UNDEFINED, usage + " expects a non-negative number"
		} else {
			return comm, num
		}
	} else {
		return COMMAND_UNDEFINED, "Couldn't parse as int: " + token
	}
}

func build() (Command, interface{}) {
	token, error := getNextToken()
	if error != nil {
		return COMMAND_UNDEFINED, error.Error()
	}

	if token == COMMAND_GLOBULA_STR {
		return COMMAND_BUILD_GLOBULA, build_globula.GlobulaBuildAlg
	}

	if token == COMMAND_THREAD_STR {
		return COMMAND_BUILD_THREAD_GLOBULA, build_globula.ThreadBuildAlg
	}

	return getUndefinedCommand(token)
}

func show() (Command, interface{}) {
	token, error := getNextToken()
	if error != nil {
		return COMMAND_UNDEFINED, error.Error()
	}

	if token == COMMAND_PARAMETERS_STR {
		return COMMAND_SHOW_PARAMETERS, nil
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
	globulaName, err := getGlobulaName()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	if finished() {
		return COMMAND_UNDEFINED, string("Usage: age <globula_name> <groups_count>")
	}
	getNextToken() // skip empty spaces
	if finished() {
		return COMMAND_UNDEFINED, string("Usage: age <globula_name> <groups_count>")
	}
	moveForward()
	if finished() {
		return COMMAND_UNDEFINED, string("Usage: age <globula_name> <groups_count>")
	}
	token, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err
	}

	groupCount, err := strconv.Atoi(token)
	if err != nil {
		return COMMAND_UNDEFINED, err
	}
	m := make(map[string]interface{})
	m["globula"] = globulaName
	m["count"] = groupCount
	return COMMAND_AGE, m
}

func borders() (Command, interface{}) {
	globulaName, err := getGlobulaName()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	m := make(map[string]interface{})
	m["globula"] = globulaName
	return COMMAND_HIGHLIGHT_BORDERS, m
}

func resetGlobula() (Command, interface{}) {
	globulaName, err := getGlobulaName()
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
	globulaName, err := getGlobulaName()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	m := make(map[string]interface{})
	m["globula"] = globulaName
	return COMMAND_WATERIZE, m
}

func trunk() (Command, interface{}) {
	globulaName, err := getGlobulaName()
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
