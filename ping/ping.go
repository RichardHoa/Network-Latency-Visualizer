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
	workingDirReport := workingDir + "/report/report.txt"
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
