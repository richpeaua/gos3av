package main

import (
	"errors"
	"fmt"
)

func helloName(n string) (string, error) {
	greeting := fmt.Sprintf("Hello, %s!", n)

	switch n {
	case "":
		return greeting, errors.New("`name` cannot be blank")
	}

	return greeting, nil
}

func main() {
}