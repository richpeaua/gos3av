package main

import (
	"fmt"
)

func helloName(n string) (string, error) {
	greeting := fmt.Sprintf("Hello, %s!", n)
	return greeting, nil
}

func main() {
}