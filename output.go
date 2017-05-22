package main

import (
	"fmt"
)

func Println(message interface{}, verbose bool) {
	if verbose {
		fmt.Println(message)
	}
}
