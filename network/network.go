package network

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type NetworkData struct {
	ProcessName string
	ReceivedMB  []string
	SentMB      []string
	Time        []string
}

// Function to record Network Data to file
func RecordNetworkData(WORKING_DIR string) error {

	workingDirReport := WORKING_DIR + "/network/network.txt"

	// Open the file for appending
	file, openFileErr := os.OpenFile(workingDirReport, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openFileErr != nil {
		return openFileErr
	}
	defer file.Close()

	// Scan using nettop
	networkcmd, err := exec.Command("nettop", "-l", "1", "-P", "-x").Output()
	networkcmd = networkcmd[:len(networkcmd)-1]
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(networkcmd), "\n")
	lines = lines[1:]

	// Using regex to extract data
	re := regexp.MustCompile(`\d+:\d+:\d+\.\d+\s+([\w\s\(\)\.]+)\.(\d+)\s+(\d+)\s+(\d+)`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)

		// Get the network consumption in byte
		receivedByte, err := strconv.ParseFloat(matches[3], 64)
		if err != nil {
			fmt.Println(err)
		}

		// Get the sent network in byte
		sentByte, err := strconv.ParseFloat(matches[4], 64)
		if err != nil {
			fmt.Println(err)
		}

		// Convert from byte to MB
		receivedMB := receivedByte / float64(1000000)
		sentMB := sentByte / float64(1000000)

		// Convert to string
		receivedMBString := fmt.Sprintf("%.5f", receivedMB)
		sentMBString := fmt.Sprintf("%.5f", sentMB)

		finalResult := matches[1] + " | " + receivedMBString + " | " + sentMBString + " | " + time.Now().Format("2006-01-02 15:04:05")

		// Only record process that has received and sent network data
		if receivedByte > 0 && sentByte > 0 {
			_, writeFileErr := file.WriteString(finalResult + "\n")
			if writeFileErr != nil {
				return writeFileErr
			}
		}
	}

	return nil
}

// Function to read the network report and return the stats, ready for chart building
func ReadNetworkData(WORKING_DIR string) (networkMap map[string]NetworkData, err error) {
	filePath := WORKING_DIR + "/network/network.txt"

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Get all the lines in slice format
	lines := strings.Split(string(file), "\n")
	lines = lines[:len(lines)-1]

	var networkDataMap = make(map[string]NetworkData)

	for _, line := range lines {
		slice := strings.Split(line, " | ")
		processName := slice[0]

		// If the process is already in the map, update its network data
		if networkData, ok := networkDataMap[processName]; ok {
			// Update existing slice
			receivedMB := append(networkData.ReceivedMB, slice[1])
			sentMB := append(networkData.SentMB, slice[2])
			time := append(networkData.Time, slice[3])

			// Replace the old network data by a new network data
			networkDataMap[processName] = NetworkData{
				ProcessName: processName,
				ReceivedMB:  receivedMB,
				SentMB:      sentMB,
				Time:        time,
			}
		// If the process is not in the map, add it
		} else {
			networkDataMap[processName] = NetworkData{
				ProcessName: processName,
				ReceivedMB:  []string{slice[1]},
				SentMB:      []string{slice[2]},
				Time:        []string{slice[3]},
			}
		}
	}
	return networkDataMap, nil

}

// Sort the map in descending order
func SortNetworkDataMap(networkDataMap map[string]NetworkData, sortedByReceivedData bool) (keysSortedInDesc []string) {

	// Initialize a slice of string containing all the keys sorted in descending order
	keysDesc := make([]string, 0, len(networkDataMap))

	// Add all the key of the map to the slice
	for key := range networkDataMap {
		keysDesc = append(keysDesc, key)
	}

	// Sorted the key based on the requirement
	sort.SliceStable(keysDesc, func(i, j int) bool {
		var (
			totalMBI float64
			totalMBJ float64
			MBSliceI    []string
			MBSliceJ    []string
		)

		// if sorted by incoming network is true, sort by incoming network
		if sortedByReceivedData {
			MBSliceI = networkDataMap[keysDesc[i]].ReceivedMB
			MBSliceJ = networkDataMap[keysDesc[j]].ReceivedMB
		// if sorted by incoming network is false, sort by outgoing network
		} else {
			MBSliceI = networkDataMap[keysDesc[i]].SentMB
			MBSliceJ = networkDataMap[keysDesc[j]].SentMB
		}

		for _, value := range MBSliceI {
			valueFloat, _ := strconv.ParseFloat(value, 64)
			totalMBI += valueFloat
		}

		for _, value := range MBSliceJ {
			valueFloat, _ := strconv.ParseFloat(value, 64)
			totalMBJ += valueFloat
		}

		avgMBI := totalMBI / float64(len(MBSliceI))
		avgMBJ := totalMBJ / float64(len(MBSliceJ))

		return avgMBI > avgMBJ
	})

	return keysDesc

}

// Function to get the top N keys in descending order
func GetTopDesc(keysSorted []string, topNumber int) (topKeysInDesc []string) {
	topKeys := make([]string, 0, 3)
	for i := 0; i < topNumber; i++ {
		topKeys = append(topKeys, keysSorted[i])
	}

	return topKeys

}
