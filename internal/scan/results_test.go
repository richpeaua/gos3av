package scan

import (
	// "fmt"
	"testing"
)

// Stdout of scan of clean file by TestClamscanProcess
const mockScannerOutputClean = `
testfiles/clean.txt: OK

----------- SCAN SUMMARY -----------
Known viruses: 8586950
Engine version: 0.104.1
Scanned directories: 0
Scanned files: 1
Infected files: 0
Data scanned: 0.00 MB
Data read: 0.00 MB (ratio 0.00:1)
Time: 18.165 sec (0 m 18 s)
Start Date: 2021:12:31 00:07:59
End Date:   2021:12:31 00:08:17
`
// Stdout of scan of infected file by TestClamscanProcess
const mockScannerOutputInfected = `
testfiles/infected.txt: Win.Test.EICAR_HDB-1 FOUND

----------- SCAN SUMMARY -----------
Known viruses: 8586950
Engine version: 0.104.1
Scanned directories: 0
Scanned files: 1
Infected files: 1
Data scanned: 0.00 MB
Data read: 0.00 MB (ratio 0.00:1)
Time: 16.342 sec (0 m 16 s)
Start Date: 2021:12:31 00:09:55
End Date:   2021:12:31 00:10:12	
`

func TestScanResults(t *testing.T) {
	testScenarios := []struct {
		name, file  string
		input       scanOutput
		status      string
		signature   string
	}{
		{
			name: "ExpectedOutClean",
			input: mockScannerOutputClean,
			file: "testfiles/clean.txt",
			status: "CLEAN",
		},
		{
			name: "ExpectedOutInfected",
			input: mockScannerOutputInfected,
			file: "testfiles/infected.txt",
			status: "INFECTED",
			signature: "Win.Test.EICAR_HDB-1",
		},

	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {

			// Given a parsed scan output
			pr := parseScanOutput(scenario.file, scenario.input)

			// 1. ExpectedValue: Clean file
			if want, got := scenario.status, pr.Status; got != want {
				t.Errorf("wanted `%v`, but got `%v`", want, got)
			}

			// 2. ExpectedValue: Infected file
			if want, got := scenario.status, pr.Status; got != want {
				t.Errorf("wanted `%v`, but got `%v`", want, got)
			}

			// 3. ExpectedValue: Infected file virus signature
			if  want, got := scenario.signature, pr.Signature; got != want {
				t.Errorf("wanted signature `%v`, but got `%v`", want, got)
			}
		})
	}
}
