package main

import "fmt"
import "github.com/zjyl1994/caasiu"

func main() {
	fmt.Println("Caasiu example")
	jsonStr := `
	{
		"string":"sdfsdfdsdsfsd",
		"data1":{
			"data2":{
				"data3":45.7
			},
			"data4":67,
			"data5":"7430893"
		}
	}
	`
	cs := caasiu.NewJSON([]byte(jsonStr), nil)
	valid, errMsg := cs.Valid(map[string][]string{
		"string":            []string{"string"},
		"data1.data2.data3": []string{"integer"},
		"data1.data4":       []string{"integer"},
		"data1.data5":       []string{"string", "integer"},
		"data6.data7":       []string{"required"},
	})
	fmt.Println("VALID", valid)
	fmt.Println("ERROR", errMsg)
}
