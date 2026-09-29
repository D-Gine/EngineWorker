package main

import (
	"fmt"

	ece "github.com/D-Gine/ECE"
)

func main() {
	reg := ece.NewRegistry()
	reg.RegisterComponents[int]()
	err := reg.AddComponent(0, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	int, err := reg.GetComponents[int]().Get(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("entity0 int = '%d'", int)
}
