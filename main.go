package main

import (
	"errors"
	"fmt"
)

func helloName(n string) (string, error) {
	greeting := fmt.Sprintf("Hello, %s!", n)

	switch {
	case n == "":
		return greeting, errors.New("name cannot be blank")
	case len(n) > 64:
		return greeting, errors.New("name cannot exceed 64 chars")
	}
	return greeting, nil
}

func main() {
}