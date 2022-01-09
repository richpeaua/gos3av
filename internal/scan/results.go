package scan

import (
	"math"
	"strconv"
	"strings"
)

type ScanResult struct {
	Status, Signature string
	Duration float64
}

func parseScanOutput(file string, so scanOutput) ScanResult {
	resultItems := make(map[string]string)

	outputLines := strings.Split(string(so), "\n")

	for _, line := range outputLines {
		if strings.Contains(line, ":") {
			lineItems := strings.Split(line, ":")
			k, v := lineItems[0], lineItems[1]
			resultItems[k] = strings.TrimSpace(v)
		}
	}

	stat := "INFECTED"
	sig  := ""
	dur, _ := strconv.ParseFloat(strings.Split(resultItems["Time"], " ")[0], 32)

	if resultItems["Infected files"] == "0" {
		stat = "CLEAN"
	}

	if stat == "INFECTED" {
		sig = strings.Replace(resultItems[file], "FOUND", "", 1)
	}

	return ScanResult{
		Status: stat,
		Signature: strings.TrimSpace(sig),
		Duration: math.Floor(dur * 100)/100,
	}
} 