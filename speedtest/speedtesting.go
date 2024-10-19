package speedtest

import (
	"fmt"
	"github.com/showwin/speedtest-go/speedtest"
	"os"
	"strings"
	"time"
)

func RecordSpeedTestData(WORKING_DIR string) error {
	workingDirReport := WORKING_DIR + "/speedtest/speedtest.txt"
	var speedtestClient = speedtest.New()

	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	var DLSpeed speedtest.ByteRate
	var ULSpeed speedtest.ByteRate

	for _, s := range targets {
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()
		DLSpeed = s.DLSpeed
		ULSpeed = s.ULSpeed

		s.Context.Reset()
	}

	DLSpeed = speedtest.ByteRate(DLSpeed.Mbps())
	ULSpeed = speedtest.ByteRate(ULSpeed.Mbps())

	// Open the file for appending
	file, openFileErr := os.OpenFile(workingDirReport, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openFileErr != nil {
		return openFileErr
	}
	defer file.Close()

	// Create custom text for the report
	resultString := fmt.Sprintf("%.2f MB/s | %.2f MB/s", DLSpeed, ULSpeed)
	resultString += " | " + time.Now().Format("2006-01-02 15:04:05") + "\n"

	_, writeFileErr := file.WriteString(resultString)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil

}

func ReadSpeedTestReport(reportPath string) (DLSpeed []string, ULSpeed []string, timeString []string, err error) {

	report, openFileErr := os.ReadFile(reportPath)
	if openFileErr != nil {
		return nil, nil, nil, openFileErr
	}
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
