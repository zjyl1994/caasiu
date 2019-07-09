package caasiu

import (
	"encoding/json"
	"fmt"
	"regexp"
)

const (
	regexInteger = "^\\d+$"
)

var builtinRules map[string]func(string, string, interface{}) (bool, string)

func init() {
	builtinRules = map[string]func(string, string, interface{}) (bool, string){
		"string":  ruleString,
		"integer": ruleInteger,
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
