package caasiu

import (
	"fmt"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

type JSON struct {
	RawJson       []byte                                              // JSON 原文
	RegisterRules map[string]func(string, interface{}) (bool, string) // 已注册的规则
}

func NewJSON(json []byte, rules map[string]func(string, interface{}) (bool, string)) *JSON {
	var result JSON
	result.RawJson = json
	result.RegisterRules = rules
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
		var currentJsonLevel *simplejson.Json
		for _, oneField := range fieldPaths {
			currentJsonLevel = sjson.Get(oneField)
			if currentJsonLevel == nil {
				break
			}
		}
		if currentJsonLevel == nil {
			if stringInArray("required", rulesOnField) {
				errMsg = append(errMsg, fmt.Sprintf(`field "%s" is required`, fieldName))
			}
			continue
		}
		for _, oneRule := range rulesOnField {
			ruleCommand := strings.Split(oneRule, ":")
			if ruleFunc, ok := j.RegisterRules[ruleCommand[0]]; ok {
				valid, errMessage := ruleFunc(oneRule, currentJsonLevel.Interface{})
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
