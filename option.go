package caasiu

type Option struct {
	ExtraRule map[string]func(string, string, interface{}) (bool, string)
}

type Rules struct {
	QueryString Rule
	Body        Rule
}

type Rule = map[string][]string
