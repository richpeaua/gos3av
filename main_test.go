package main

import (
	"fmt"
	"errors"
	"testing"
	"strings"
)

func TestHelloName(t *testing.T) {
	testScenarios := []struct {
		name     	 string
		input 		 string
		expectedError    error
	}{
		{
			name: "expectedOut",
			input: "rich",
			expectedError: nil,
		},
		{
			name: "expectedErrBlank",
			input: "",
			expectedError: errors.New("name must not be blank"),
		},
		{
			name: "expectedErrNameLen",
			input: strings.Repeat("a", 65),
			expectedError: errors.New("name must be under 64 chars"),
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
			}

			// 1.3 Greeting: expected error - long name
			if len(scenario.input) > 64 && errmsg == nil {
				t.Errorf("expected name length error, but got %v", errmsg)
			}

		})
	}
}