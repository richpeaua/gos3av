package scan

import (
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
var mockScanResult ScanResult

func TestScanResults(t *testing.T) {
	testScenarios := []struct {
		name, input string
		expectedOut string
	}{
		{
			name: "ExpectedOutClean",
			input: mockScannerOutputClean,
			expectedOut: mockScanResult{
				Status: "CLEAN",
				Signature: "",
				Duration: 18.16,
			},
		},
		{
			name: "ExpectedOutInfected",
			input: mockScannerOutputInfected,
			expectedOut: mockScanResult{
				Status: "INFECTED",
				Signature: "Win.Test.EICAR_HDB-1",
				Duration: 16.34,
			},
		},

	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {

			// Given a parsed scan output
			pr := parseScanOutput(scenario.input)

			// 1. ExpectedOut
			if want, got := scenario.expectedOut, pr; got != want {
				t.Errorf("wanted `%s`, but got `%s`", want, got)
			}
		})
	}
}
