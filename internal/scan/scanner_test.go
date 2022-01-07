package scan

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/rs/zerolog"
)

// Stdout of scan of clean file by TestClamscanProcess
const testScannerOutputClean = `
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
const testScannerOutputInfected = `
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

var testScenario string

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	os.Exit(m.Run())
}

// fakeExecCommandClamav is a function that initialises a new exec.Cmd, one which will
// simply call a mock process, TestClamavProc, rather than the command it is provided. It will
// also pass through the command and its arguments as an argument to the helperProcess
func fakeExecCommandClamav(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestClamscanProc", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	ts := "TEST_SCENARIO=" + testScenario
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", ts}
	return cmd
}

// fakeExecCommandFreshclam is a function that initialises a new exec.Cmd, one which will
// simply call a mock process, TestFreschlamProc, rather than the command it is provided. It will
// also pass through the command and its arguments as an argument to the helperProcess
func fakeExecCommandFreshclam(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestFreshclamProc", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	ts := "TEST_SCENARIO=" + testScenario
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", ts}
	return cmd
}

// TestClamscanProc is a method that is called as a substitute for the clamscan shell command,
// the GO_TEST_PROCESS flag ensures that if it is called as part of the test suite, it is
// skipped.
func TestClamscanProc(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	
	// ensure command is passed
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	// behaviour per test scenario
	switch os.Getenv("TEST_SCENARIO") {
	case "ExpectedOutClean":
		fmt.Fprint(os.Stdout, testScannerOutputClean)
	case "ExpectedOutInfected":
		fmt.Fprint(os.Stdout, testScannerOutputInfected)
		os.Exit(1)
	default:
		fmt.Fprint(os.Stderr, "error")
		os.Exit(2)
	}
}

// TestFreshclaProc is a method that is called as a substitute for the freshclam shell command,
// the GO_TEST_PROCESS flag ensures that if it is called as part of the test suite, it is
// skipped.
func TestFreshclamProc(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	
	// ensure command is passed
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	// behaviour per test scenario
	switch os.Getenv("TEST_SCENARIO") {
	case "ExpectedPass":
		fmt.Fprint(os.Stdout, "updated")
	default:
		fmt.Fprint(os.Stderr, "error")
		os.Exit(2)
	}

}

func TestUpdateDB(t *testing.T) {
	testScenarios := []struct {
		name     string
		dbpath   string
		confpath string
	}{
		{
			name: "ExpectedPass",
			dbpath: "/db/path",
			confpath: "/conf/path",
		},
		{
			name: "ExpectedErrNoDbPath",
			dbpath: "",
			confpath: "/conf/path",
		},
		{
			name: "ExpectedErrNoConfPath",
			dbpath: "/db/path",
			confpath: "",
		},

	}

	tvs := NewVirusScanner()

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
		
			// 1. Fake cmd setup	
			testScenario = scenario.name
			execCommand = fakeExecCommandFreshclam
			defer func(){ execCommand = exec.Command }()

			// 2. DB Update
			tvs.DBPath = scenario.dbpath
			tvs.FreshclamConfPath = scenario.confpath

			err := tvs.UpdateDB()

			// 3. Update without no errors
			if scenario.name == "ExpectedPass" && err != nil {
				t.Errorf("Expected `nil` error, got %v", err)
			}

			// 4.1 ExpectedErrors: no db path
			if tvs.DBPath == "" && err == nil {
				t.Errorf("ExpectedErr missing db path, got %v", err)
			}

			// 4.2 ExpectedErrors: no freshclam conf path
			if tvs.FreshclamConfPath == "" && err == nil {
				t.Errorf("ExpectedErr missing freshclam conf path, got %v", err)
			}
		})
	}


}

func TestFileScan(t *testing.T) {
	testScenarios := []struct {
		name  string
		fpath string
		fcont []byte
	}{
		{
			name: "ExpectedOutClean",
			fpath: "./clean_file.txt",
		},
		{
			name: "ExpectedOutInfected",
			fpath: "./infected_file.txt",
		},
		{
			name: "ExpectedErr",
		},
	}
	
	tvs := NewVirusScanner()

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {

			// 1. Fake cmd setup
			testScenario = scenario.name
			execCommand = fakeExecCommandClamav
			defer func(){ execCommand = exec.Command }()

			// 2. Scan file
			output, scanerr := tvs.ScanFile(scenario.fpath)
			

			// 3. ExpectedOutClean
			if output == "" && scanerr == nil {
				t.Errorf("expected output, but got `%v`", output)
			}

			// 4. ExpectedErr
			if scenario.name == "ExpectedErr" && scanerr == nil {
				t.Errorf("expected an error, but got `%v`", scanerr)
			}	

		})
	}
}