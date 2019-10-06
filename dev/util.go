package dev

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(vs ...interface{}) {
	for _, v := range vs {
		bytes, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Println(string(bytes))
		}
	}
}

func PrintJSONIndent(vs ...interface{}) {
	for _, v := range vs {
		bytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			fmt.Println(string(bytes))
		}
	}
}
