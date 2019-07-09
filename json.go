package caasiu

import (
	"fmt"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

type JSON struct {
	RawJson       []byte                                                      // JSON 原文
	RegisterRules map[string]func(string, string, interface{}) (bool, string) // 已注册的规则
}

func NewJSON(json []byte, rules map[string]func(string, string, interface{}) (bool, string)) *JSON {
	var result JSON
	result.RawJson = json
	result.RegisterRules = builtinRules
	for ruleName, ruleFunc := range rules {
		result.RegisterRules[ruleName] = ruleFunc
	}
	return &result
}
func (j *JSON) Valid(rule map[string][]string) (bool, []string) {
	sjson, err := simplejson.NewJson(j.RawJson)
	if err != nil {
		return false, []string{err.Error()}
	}
	var errMsg []string
	for fieldName, rulesOnField := range rule {
		fieldPaths := strings.Split(fieldName, ".")
		currentJsonLevel := sjson.GetPath(fieldPaths...)
		if currentJsonLevel.Interface() == nil {
			if stringInArray("required", rulesOnField) {
				errMsg = append(errMsg, fmt.Sprintf(`field "%s" is required`, fieldName))
			}
			continue
		}
		for _, oneRule := range rulesOnField {
			ruleCommand := strings.Split(oneRule, ":")
			if ruleFunc, ok := j.RegisterRules[ruleCommand[0]]; ok {
				valid, errMessage := ruleFunc(oneRule, fieldName, currentJsonLevel.Interface())
				if !valid {
					errMsg = append(errMsg, errMessage)
				}
			}
		}
	}
	if len(errMsg) > 0 {
		return false, errMsg
	} else {
		return true, nil
	}
}
