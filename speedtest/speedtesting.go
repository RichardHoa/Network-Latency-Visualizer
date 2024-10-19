package speedtest

import (
	"fmt"
	"github.com/showwin/speedtest-go/speedtest"
	"os"
	"strings"
	"time"
)

// Record Speed Test Data
func RecordSpeedTestData(WORKING_DIR string) error {

	var (
		speedtestClient = speedtest.New()
		DLSpeed         speedtest.ByteRate
		ULSpeed         speedtest.ByteRate
	)
	// Get all the avail servers
	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})
	// Pick the first server from the list
	server := targets[0]
	// No ping test
	server.PingTest(nil)
	// Do download test and upload test
	server.DownloadTest()
	server.UploadTest()

	DLSpeed = server.DLSpeed
	ULSpeed = server.ULSpeed

	server.Context.Reset()

	// Convert from BYTE/s to MB/s
	DLSpeed = speedtest.ByteRate(DLSpeed.Mbps())
	ULSpeed = speedtest.ByteRate(ULSpeed.Mbps())

	workingDirReport := WORKING_DIR + "/speedtest/speedtest.txt"

	// Open the file for appending
	file, openFileErr := os.OpenFile(workingDirReport, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openFileErr != nil {
		return openFileErr
	}
	defer file.Close()

	// Create custom text for the report
	resultString := fmt.Sprintf("%.2f MB/s | %.2f MB/s", DLSpeed, ULSpeed)
	resultString += " | " + time.Now().Format("2006-01-02 15:04:05") + "\n"

	// Write the text to the file
	_, writeFileErr := file.WriteString(resultString)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil

}

// Function to read the speedtest report and return the stats, ready for chart building
func ReadSpeedTestReport(reportPath string) (DLSpeed []string, ULSpeed []string, timeString []string, err error) {
	// Read the file
	report, openFileErr := os.ReadFile(reportPath)
	if openFileErr != nil {
		return nil, nil, nil, openFileErr
	}

	// Extract data
	lines := strings.Split(string(report), "\n")
	lines = lines[:len(lines)-1]

	for _, line := range lines {
		sections := strings.Split(line, " | ")
		DL := strings.Split(sections[0], " ")[0]
		UL := strings.Split(sections[1], " ")[0]
		DLSpeed = append(DLSpeed, DL)
		ULSpeed = append(ULSpeed, UL)
		timeString = append(timeString, sections[2])
	}

	return DLSpeed, ULSpeed, timeString, nil

}
