package scan

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
)

// VirusScanner object encapsulates simple clamav virus scanner functionalities such as
// scanning a target file or updating the virus database
type VirusScanner struct {

	// DBPath is the path of the Clamav virus databases used to perform the scan
	DBPath 	          string

	// FreshclamConfPath is the path of the config file for the freshclam command, which
	// provides the database mirror information as well as database file compression options
	FreshclamConfPath string
}

// NewVirusScanner returns a VirusScanner struct
func NewVirusScanner() VirusScanner {
	return VirusScanner{}
}

// UpdateDB pulls down the latest clamav virus databases
func (vs VirusScanner) UpdateDB() error {
	logFuncName := "VirusScanner.UpdateDB"

	var err error

	switch {
	case vs.DBPath == "":
		err = errors.New("missing DB path")
		log.Error().Str("func", logFuncName).AnErr("DbPathMissing", err)
		return err
	case vs.FreshclamConfPath == "":
		err = errors.New("missing freshclam conf path")
		log.Error().Str("func", logFuncName).AnErr("FreshclamConfPathMissing", err)
		return err
	}

	args := fmt.Sprintf("--config-file=%s '--datadir=%s",
				vs.FreshclamConfPath, vs.DBPath)
	cmd := execCommand("freshclam", args)

	log.Info().Str("func", logFuncName).Msg("Downloading clamav virus database files to: " + vs.DBPath)

	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Error().Str("func", logFuncName).AnErr("DbUpdateErr", err)
		return err
	}

	log.Info().Str("func", logFuncName).Msg("Database updated")

	return nil
}

// execCommand var which holds a mockable exec.Command function for running the Clamscan program
var execCommand = exec.Command

// ScanFile scans a file at a given filepath and returns the output
func (vs VirusScanner) ScanFile(fpath string) (string, error) {
	logFuncName := "VirusScanner.ScanFile"

	cmd := execCommand("clamscan", fpath)
	
	log.Info().Str("func", logFuncName).Msgf("Scanning file: %s ...", fpath)

	start := time.Now()
	out, err := cmd.CombinedOutput()
	duration := time.Since(start)

	exitCode := cmd.ProcessState.ExitCode()
	// clamscan exits with exitcode 1 if infected file is found so exitcode 1 should not be considered an error
	if err != nil && exitCode > 1 {
		log.Error().Str("func", logFuncName).AnErr("ScanFileErr", err).Msg("")
		return "", err
	}

	log.Info().Str("func", logFuncName).Msgf("Scan completed: %.2f secs", duration.Seconds())

	return string(out), nil
}