package command_interpreter

import (
	"errors"
	"polymers/build_globula"
	"strconv"
	"unicode"
)

const (
	COMMAND_HELP_STR              = "help"
	COMMAND_BUILD_STR             = "build"
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
	COMMAND_FILE_STR              = "file"
	COMMAND_DFS_STR               = "dfs"
	COMMAND_ATOMISTIC_STR         = "atomistic"
	COMMAND_LOADER_STR            = "load"
)

type Command = int

const (
	COMMAND_UNDEFINED = -1
	COMMAND_HELP      = iota
	COMMAND_BUILD
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
	COMMAND_ATOMISTIC
	COMMAND_LOADER
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
	var token string
	var char rune = rune(getCurrChar())
	for char == ' ' {
		moveForward()
		if finished() {
			return "", errors.New("no command was found")
		}
		char = rune(getCurrChar())
	}
	for unicode.IsLetter(char) || unicode.IsNumber(char) ||
		char == '_' ||
		char == '.' ||
		char == '%' ||
		char == '(' || char == ')' ||
		char == '*' ||
		char == ',' ||
		char == '-' {
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

	if f, ok := map[string]func() (Command, interface{}){
		COMMAND_HELP_STR:              func() (Command, interface{}) { return COMMAND_HELP, nil },
		COMMAND_PATTERN_STR:           pattern,
		COMMAND_BUILD_STR:             build,
		COMMAND_SHOW_STR:              show,
		COMMAND_SAVE_STR:              save,
		COMMAND_CLUSTERS_STR:          clusters,
		COMMAND_AGE_STR:               age,
		COMMAND_HIGHLIGHT_BORDERS_STR: borders,
		COMMAND_RESET_STR:             resetGlobula,
		COMMAND_EXIT_STR:              func() (Command, interface{}) { return COMMAND_EXIT, nil },
		COMMAND_WATERIZE_STR:          waterize,
		COMMAND_TRUNK_STR:             trunk,
		COMMAND_SCRIPT_STR:            script,
		COMMAND_COMMON_STATS_STR:      commonStats,
		COMMAND_ATOMISTIC_STR:         atomistic,
		COMMAND_LOADER_STR:            load,
	}[token]; ok {
		return f()
	} else {
		return COMMAND_UNDEFINED, "Undefined command: " + token
	}
}

func pattern() (Command, interface{}) {
	m := make(map[string]string)

	t, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, "Wrong source type"
	}
	switch t {
	case COMMAND_FILE_STR:
		fileName, err := getParameterAsString()
		if err != nil {
			return COMMAND_UNDEFINED, "Wrong usage"
		}
		m["fileName"] = fileName
	case COMMAND_PATTERN_STR:
		patt, err := getParameterAsString()
		if err != nil {
			return COMMAND_UNDEFINED, "Wrong usage"
		}
		if patt, err = getPattern(patt); err == nil {
			m["pattern"] = patt
		} else {
			return COMMAND_UNDEFINED, err.Error()
		}
	default:
		return COMMAND_UNDEFINED, "Wrong source type"
	}

	outputName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}

	m["outputName"] = outputName

	return COMMAND_PATTERN, m
}

func build() (Command, interface{}) {
	objective, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}

	predefinedParams := make([]string, 0)
	for getCurrChar() != '"' {
		p, err := getNextToken()
		if err != nil {
			continue
		}
		predefinedParams = append(predefinedParams, p)
	}

	m := make(map[string]interface{})
	m["params"] = predefinedParams

	name, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, "Specify the name of a particle"
	}
	m["name"] = name

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
	return COMMAND_SHOW_GLOBULA, nil
}

func save() (Command, interface{}) {
	getNextToken() // skip all the empty spaces until " or the end
	if finished() {
		return COMMAND_UNDEFINED, "Wrong usage"
	}

	if string(getCurrChar()) != "\"" {
		return COMMAND_UNDEFINED, "Wrong usage"
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
		}
	}
	return COMMAND_UNDEFINED, "Globula name must be wrapped with \"\""
}

func clusters() (Command, interface{}) {
	getNextToken() // skip all the empty spaces until " or the end
	if finished() {
		return COMMAND_UNDEFINED, "Wrong usage"
	}

	tokenAll, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}

	if tokenAll == COMMAND_CLUSTERS_ALL_STR {
		return COMMAND_HIGHLIGHT_CLUSTERS_ALL, globulaName
	}

	return COMMAND_UNDEFINED, "Wrong usage"
}

func age() (Command, interface{}) {
	globulaName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	if finished() {
		return COMMAND_UNDEFINED, string("Wrong usage")
	}
	groupCount, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err
	}

	moveForward()
	if finished() {
		return COMMAND_UNDEFINED, string("Wrong usage")
	}
	token, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err
	}

	if token != "true" && token != "false" {
		return COMMAND_UNDEFINED, string("Error: make_crosslinks parameter must be either \"true\" or \"false\"")
	}

	token, err = getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, err
	}
	ageAlgType, err := strconv.Atoi(token)
	if err != nil {
		return COMMAND_UNDEFINED, err
	}

	m := make(map[string]interface{})
	// For Age Algorithm Type 3
	if !finished() && ageAlgType == 3 {
		ncut, err := getNextToken()
		if err != nil {
			return COMMAND_UNDEFINED, err
		}
		m["ncut"] = ncut
		if finished() {
			return COMMAND_UNDEFINED, errors.New("Wrong usage of the age command")
		}
		nOContaining, err := getNextToken()
		if err != nil {
			return COMMAND_UNDEFINED, err
		}
		m["OContaining"] = nOContaining
		if finished() {
			return COMMAND_UNDEFINED, errors.New("Wrong usage of the age command")
		}
		ncross, err := getNextToken()
		if err != nil {
			return COMMAND_UNDEFINED, err
		}
		m["ncross"] = ncross
	}

	m["globula"] = globulaName
	m["count"] = groupCount
	m["make_crosslinks"] = (token == "true")
	m["alg_type"] = ageAlgType
	return COMMAND_AGE, m
}

func borders() (Command, interface{}) {
	return COMMAND_HIGHLIGHT_BORDERS, nil
}

func resetGlobula() (Command, interface{}) {
	getNextToken() // skip empty spaces
	if finished() {
		m := make(map[string]interface{})
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
		m["full"] = true
		return COMMAND_RESET_FULL, m
	}

	return COMMAND_UNDEFINED, string("Wrong usage")
}

func waterize() (Command, interface{}) {
	return COMMAND_WATERIZE, nil
}

func trunk() (Command, interface{}) {
	m := make(map[string]interface{})
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
	m := make(map[string]string)

	filename, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, string("Wrong usage")
	}

	m["filename"] = filename
	return COMMAND_COMMON_STATS, m
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
			if err := getPattern_S2(raw, &curr, &dst); err != nil {
				return "", nil
			}
		} else {
			return "", errors.New("The pattern doesn't follow the rules: " + raw)
		}
	}
	return dst, nil
}

func getPattern_S2(raw string, curr *int, dst *string) error {
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
			return getPattern_S3(raw, curr, start, dst)
		} else {
			return errors.New("The pattern doesn't follow the rules: " + raw)
		}
	}
}

func getPattern_S3(raw string, curr *int, start int, dst *string) error {
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
		return getPattern_S4(raw, curr, start, dst)
	} else {
		return errors.New("The pattern doesn't follow the rules: " + raw)
	}
}

func getPattern_S4(raw string, curr *int, start int, dst *string) error {
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

func atomistic() (Command, interface{}) {
	m := make(map[string]string)

	fileName, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, err.Error()
	}
	m["config"] = fileName

	return COMMAND_ATOMISTIC, m
}

func load() (Command, interface{}) {
	filetype, err := getNextToken()
	if err != nil {
		return COMMAND_UNDEFINED, "no filetype specified"
	}

	m := make(map[string]string)
	m["filetype"] = filetype

	filename, err := getParameterAsString()
	if err != nil {
		return COMMAND_UNDEFINED, "no filename specified"
	}
	m["filename"] = filename

	return COMMAND_LOADER, m
}
