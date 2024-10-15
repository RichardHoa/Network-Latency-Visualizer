package network

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type NetworkData struct {
	processName string
	MBIn        []string
	MBOut       []string
	time        []string
}

func RecordNetworkData(WORKING_DIR string) error {
	fmt.Println("Recording network data...")

	workingDirReport := WORKING_DIR + "/network/network.txt"

	// Open the file for appending
	file, openFileErr := os.OpenFile(workingDirReport, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openFileErr != nil {
		return openFileErr
	}
	defer file.Close()

	networkcmd, err := exec.Command("nettop", "-l", "1", "-P", "-x").Output()
	networkcmd = networkcmd[:len(networkcmd)-1]
	if err != nil {
		log.Fatal(err)
	}
	stringSlice := strings.Split(string(networkcmd), "\n")
	stringSlice = stringSlice[1:]

	re := regexp.MustCompile(`\d+:\d+:\d+\.\d+\s+([\w\s\(\)\.]+)\.(\d+)\s+(\d+)\s+(\d+)`)

	for _, s := range stringSlice {
		matches := re.FindStringSubmatch(s)
		byteIn, err := strconv.ParseFloat(matches[3], 64)
		if err != nil {
			fmt.Println(err)
		}
		MBIn := byteIn / float64(1000000)
		MBInString := fmt.Sprintf("%.5f", MBIn)

		byteOut, err := strconv.ParseFloat(matches[4], 64)
		if err != nil {
			fmt.Println(err)
		}
		MBOut := byteOut / float64(1000000)
		MBOutString := fmt.Sprintf("%.5f", MBOut)

		formattedString := matches[1] + " | " + MBInString + " | " + MBOutString + " | " + time.Now().Format("2006-01-02 15:04:05")

		if byteIn > 0 && byteOut > 0 {

			_, writeFileErr := file.WriteString(formattedString + "\n")
			if writeFileErr != nil {
				return writeFileErr
			}
		}

	}

	return nil
}

func ReadNetworkData(WORKING_DIR string) error {
	filePath := WORKING_DIR + "/network/network.txt"

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(file), "\n")
	lines = lines[:len(lines)-1]
	var networkDataMap = make(map[string]NetworkData)

	for _, line := range lines {
		slice := strings.Split(line, " | ")

		name := slice[0]

		if networkData, exists := networkDataMap[name]; exists {

			MBIn := append(networkData.MBIn, slice[1])
			MBOut := append(networkData.MBOut, slice[2])
			time := append(networkData.time, slice[3])

			networkDataMap[name] = NetworkData{
				processName: name,
				MBIn:        MBIn,
				MBOut:       MBOut,
				time:        time,
			}

		} else {
			networkDataMap[name] = NetworkData{
				processName: name,
				MBIn:        []string{slice[1]},
				MBOut:       []string{slice[2]},
				time:        []string{slice[3]},
			}

		}

	}

	for name, networkData := range networkDataMap {
		fmt.Println("--------------------------")
		fmt.Printf("Process name: %s\n", name)
		fmt.Printf("Length MBIN: %d\n", len(networkData.MBIn))
		fmt.Printf("Length MBOUT: %d\n", len(networkData.MBOut))
		fmt.Printf("Length Time: %d\n", len(networkData.time))
		// fmt.Printf("MB In: %s\n", strings.Join(networkData.MBIn, ", "))
		// fmt.Printf("MB Out: %s\n", strings.Join(networkData.MBOut, ", "))
		// fmt.Printf("Time: %s\n", strings.Join(networkData.time, ", "))
		fmt.Println()
	}

	return nil

}
