package main

import (
	"fmt"
	"testing"
	"strings"
	"strconv"
)

func TestHelloName(t *testing.T) {
	testScenarios := []struct {
		name 		 string
		expectedErrorMsg string	
	}{
		{
			name: "rich",
		},
		{
			name: "",
			expectedErrorMsg: "name must not be blank",
		},
		{
			name: strings.Repeat("a", 64),
			expectedErrorMsg: "name must be under 64 chars",
		},
	}

	for i, scenario := range testScenarios {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			greeting, errmsg := helloName(scenario.name)
			
			// 1.1 Greeting: expected name
			if want, got := fmt.Sprint("Hello, %s!", scenario.name), greeting; want != got {
				t.Errorf("expected %s, but got %s", want, got)
				return
			}

			// 1.2 Greeting: expected errors
			if want, got := scenario.expectedErrorMsg, errmsg; want != got {
				t.Errorf("expected error %s, but got %s", want, got)
				return
			}
		})
	}
}