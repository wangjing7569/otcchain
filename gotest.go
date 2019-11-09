package main

import (
	"fmt"
	"encoding/json"
)

type Foo struct {
	Name string
	Age  int
}

func main() {
	data := []byte(`{"Name": "John Doe", "Age": 25}`)
	fmt.Println(len(data))
	var f Foo

	json.Unmarshal(data, &f)
	fmt.Printf("%s is %d years old.\n", f.Name, f.Age)

	output, _ := json.Marshal(&f)
	fmt.Println(string(output))
}