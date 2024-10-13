package ping

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

func PingScanning(workingDir string) (scanningErr error) {

	result, pingScanningErr := exec.Command("ping", "google.com", "-c", "10").Output()
	if pingScanningErr != nil {
		return pingScanningErr
	}
	outputString := string(result)
	resultArray := strings.Split(outputString, "\n")

	workingDirReport := workingDir + "/report/report.txt"

	file, openFileErr := os.OpenFile(workingDirReport, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openFileErr != nil {
		return openFileErr
	}

	finalString := resultArray[len(resultArray)-2] + " | " + time.Now().Format("2006-01-02 15:04:05") + "\n"
	_, err := file.WriteString(finalString)
	if err != nil {
		return err
	}

	defer file.Close()

	return nil

}
