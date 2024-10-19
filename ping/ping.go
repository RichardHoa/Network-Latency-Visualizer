package ping

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type PingStats struct {
	Min        []string
	Avg        []string
	Max        []string
	Sttdev     []string
	TimeString []string
}

// Function to scan ping for network latency then save the stats to ping/ping.txt
func RecordPingData(workingDir string) (scanningErr error) {
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

// Function to read the report and return the stats, ready for chart building
func ReadPingReport(reportPath string) (pingStats PingStats, err error) {

	report, openFileErr := os.ReadFile(reportPath)
	if openFileErr != nil {
		return PingStats{}, openFileErr
	}

	// Extract data
	lines := strings.Split(string(report), "\n")
	lines = lines[:len(lines)-1]

	pingStats = PingStats{}

	for _, line := range lines {
		lineSlice := strings.Split(line, "|")
		time := lineSlice[1]
		data := lineSlice[0]
		stats := strings.Split(strings.Split(data, "=")[1], "/")

		pingStats.Min = append(pingStats.Min, stats[0])
		pingStats.Avg = append(pingStats.Avg, stats[1])
		pingStats.Max = append(pingStats.Max, stats[2])
		pingStats.Sttdev = append(pingStats.Sttdev, strings.Split(stats[3], " ")[0])
		pingStats.TimeString = append(pingStats.TimeString, time)
	}

	return pingStats, nil

}
