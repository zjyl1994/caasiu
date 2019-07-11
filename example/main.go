package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zjyl1994/caasiu"
)

/* test request information
URL http://127.0.0.1:8080/example?param1=jhsjdfh&param2=573&param3=623dsa%2064
Body
{
	"string":"sdfSD56fd+sdsfsd",
	"data1":{
		"data2":{
			"data3":45.7
		},
		"data4":67,
		"data5":"7430893"
	},
	"data6":"3",
	"data7":5,
	"zero":0,
	"test": [
		{
			"a": "c"
		},
		{
			"b": 4
		}
	],
	"bool":true
}
*/

func myHandler(w http.ResponseWriter, r *http.Request) {
	c, err := caasiu.New(r)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	valid, errMsg := c.Valid(caasiu.Rules{
		QueryString: caasiu.Rule{
			"param1": []string{"required", "string", "alpha"},
			"param2": []string{"required", "string", "integer"},
			"param3": []string{"required", "string", "alphanum"},
			"param4": []string{"required", "string", "in:true,false"},
			"param5": []string{"string"},
		},
		Body: caasiu.Rule{
			"string":            []string{"string", "regex:^[0-9a-zA-Z_]+$"},
			"data1.data2.data3": []string{"integer"},
			"data1.data4":       []string{"integer"},
			"data1.data5":       []string{"string", "integer"},
			"data6.data7":       []string{"required"},
			"data6":             []string{"in:3,5,7,9"},
			"data7":             []string{"in:1,2,4,6"},
			"zero":              []string{"required"},
			"test":              []string{"array"},
			"bool":              []string{"bool"},
		},
	})
	var result map[string]interface{}
	qs := c.QueryString().Data()
	if c.IsJSONBody() {
		js := c.Body().Data()
		result = map[string]interface{}{
			"valid":  valid,
			"errMsg": errMsg,
			"data": map[string]interface{}{
				"param3": qs["param3"],
				"data7":  js.GetPath("data7").MustInt(),
			},
		}
	} else {
		result = map[string]interface{}{
			"valid":  valid,
			"errMsg": errMsg,
			"data": map[string]interface{}{
				"param3": qs["param3"],
				"data7":  nil,
			},
		}
	}
	b, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", b)
}

func main() {
	fmt.Println("caasiu example")
	http.HandleFunc("/example", myHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
