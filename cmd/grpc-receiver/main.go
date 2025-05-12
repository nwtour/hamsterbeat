package main

import (
	"fmt"
	"hamsterbeat/internal/hamsterbeat"
)

func main() {
	err := hamsterbeat.Server()
	if err != nil {
		fmt.Printf("Start failed %s", err)
	}
}
