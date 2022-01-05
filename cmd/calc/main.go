package main

import (
	"encoding/xml"
	"fmt"
	"log"
)

func main() {
	var v map[string]interface{} = map[string]interface{}{"a": 104, "b": "hello"}
	b, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
