package ping

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Function to scan ping for network latency then save the stats to report/report.txt
func PingScanning(workingDir string) (scanningErr error) {
	// Initialize path of the report file
	workingDirReport := workingDir + "/ping/ping.txt"
	// Testing log
	fmt.Printf("Working dir ping: %s\n", workingDirReport)

	// Scan using ping
	scanResult, pingScanningErr := exec.Command("/sbin/ping", "google.com", "-c", "10").Output()
	if pingScanningErr != nil {
		return pingScanningErr
	}

	// Open the file for appending
	file, openFileErr := os.OpenFile(workingDirReport, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openFileErr != nil {
		return openFileErr
	}
	defer file.Close()

	// Create custom text for the report
	resultArray := strings.Split(string(scanResult), "\n")
	// Ex output: round-trip min/avg/max/stddev = 26.483/30.290/37.375/3.926 ms | 2024-10-13 20:35:10
	finalString := resultArray[len(resultArray)-2] + " | " + time.Now().Format("2006-01-02 15:04:05") + "\n"

	// Write the text to the file
	_, writeFileErr := file.WriteString(finalString)
	if writeFileErr != nil {
		return writeFileErr
	}

	// Return nil if no error
	return nil
}

func ReadPingReport(reportPath string) (min []string, avg []string, max []string, sttdev []string, timeString []string, err error) {

	report, openFileErr := os.ReadFile(reportPath)
	if openFileErr != nil {
		return nil, nil, nil, nil, nil, openFileErr
	}
	lines := strings.Split(string(report), "\n")
	lines = lines[:len(lines)-1]

	minSlice := make([]string, 0)
	avgSlice := make([]string, 0)
	maxSlice := make([]string, 0)
	stddevSlice := make([]string, 0)
	timeStringSlice := make([]string, 0)

	for _, line := range lines {
		lineSlice := strings.Split(line, "|")
		time := lineSlice[1]
		data := lineSlice[0]
		stats := strings.Split(strings.Split(data, "=")[1], "/")

		minSlice = append(minSlice, stats[0])
		avgSlice = append(avgSlice, stats[1])
		maxSlice = append(maxSlice, stats[2])
		stddevSlice = append(stddevSlice, strings.Split(stats[3], " ")[0])
		timeStringSlice = append(timeStringSlice, time)
	}

	return minSlice, avgSlice, maxSlice, stddevSlice, timeStringSlice, nil

}
