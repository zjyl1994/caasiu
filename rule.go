package caasiu

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	regexInteger = "^\\d+$"
)

var builtinRules map[string]func(string, string, interface{}) (bool, string)

func init() {
	builtinRules = map[string]func(string, string, interface{}) (bool, string){
		"string":  ruleString,
		"integer": ruleInteger,
		"regex":   ruleRegexp,
		"in":      ruleIn,
	}
}

func ruleString(ruleName string, fieldName string, value interface{}) (bool, string) {
	_, ok := value.(string)
	if ok {
		return true, ""
	} else {
		return false, fmt.Sprintf(`field "%s" must string`, fieldName)
	}
}

func ruleInteger(ruleName string, fieldName string, value interface{}) (bool, string) {
	errMsg := fmt.Sprintf(`field "%s" must integer`, fieldName)
	var valueString string
	switch t := value.(type) {
	case string:
		valueString = t
	case json.Number:
		valueString = t.String()
	default:
		return false, errMsg
	}
	if regexp.MustCompile(regexInteger).Match([]byte(valueString)) {
		return true, ""
	} else {
		return false, errMsg
	}
}

func ruleRegexp(ruleName string, fieldName string, value interface{}) (bool, string) {
	ruleParamStartAt := strings.Index(ruleName, ":")
	if ruleParamStartAt == -1 {
		return false, "no param available"
	}
	param := ruleName[ruleParamStartAt+1:]
	var valueString string
	switch t := value.(type) {
	case string:
		valueString = t
	case json.Number:
		valueString = t.String()
	default:
		return false, fmt.Sprintf(`field "%s" can't cast to string`, fieldName)
	}
	regex, err := regexp.Compile(param)
	if err != nil {
		return false, "illegal regexp"
	}
	if regex.Match([]byte(valueString)) {
		return true, ""
	} else {
		return false, fmt.Sprintf(`field "%s" not match "%s"`, fieldName, param)
	}
}

func ruleIn(ruleName string, fieldName string, value interface{}) (bool, string) {
	ruleParamStartAt := strings.Index(ruleName, ":")
	if ruleParamStartAt == -1 {
		return false, "no param available"
	}
	param := ruleName[ruleParamStartAt+1:]
	params := strings.Split(param, ",")
	var valueString string
	switch t := value.(type) {
	case string:
		valueString = t
	case json.Number:
		valueString = t.String()
	default:
		return false, fmt.Sprintf(`field "%s" can't cast to string`, fieldName)
	}
	if stringInArray(valueString, params) {
		return true, ""
	} else {
		return false, fmt.Sprintf(`field "%s" not in [%s]`, fieldName, param)
	}
}
