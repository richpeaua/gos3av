package main

import (
	"fmt"
	"testing"
	"strings"
)

func TestHelloName(t *testing.T) {
	testScenarios := []struct {
		name     	 string
		input 		 string
	}{
		{
			name: "expectedOut",
			input: "rich",
		},
		{
			name: "expectedErrBlank",
			input: "",
		},
		{
			name: "expectedErrNameLen",
			input: strings.Repeat("a", 65),
		},
	}

	for _, scenario := range testScenarios {
			t.Run(scenario.name, func(t *testing.T) {

			greeting, errmsg := helloName(scenario.input)
			
			// 1.1 Greeting: expected name
			if want, got := fmt.Sprintf("Hello, %s!", scenario.input), greeting; want != got {
				t.Errorf("expected `%s`, but got `%s`", want, got)
				return
			}

			// 1.2 Greeting: expected error - blank
			if scenario.input == "" && errmsg == nil {
				t.Errorf("expected blank name error, but got %v", errmsg)
				return
			}

			// 1.3 Greeting: expected error - name too long
			if len(scenario.input) > 64 && errmsg == nil {
				t.Errorf("expected name length error, but got %v", errmsg)
				return
			}

		})
	}
}