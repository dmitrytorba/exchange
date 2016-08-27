package main

import (
	"fmt"
)

func main() {
	e, err := createExchange()
	if err != nil {
		panic(err)
	}

	fmt.Println(e)
}
