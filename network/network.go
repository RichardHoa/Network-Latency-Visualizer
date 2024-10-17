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
	MBIn        []string
	MBOut       []string
	Time        []string
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

func ReadNetworkData(WORKING_DIR string) (networkMap map[string]NetworkData, err error) {
	filePath := WORKING_DIR + "/network/network.txt"

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
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
			time := append(networkData.Time, slice[3])

			networkDataMap[name] = NetworkData{
				ProcessName: name,
				MBIn:        MBIn,
				MBOut:       MBOut,
				Time:        time,
			}

		} else {
			networkDataMap[name] = NetworkData{
				ProcessName: name,
				MBIn:        []string{slice[1]},
				MBOut:       []string{slice[2]},
				Time:        []string{slice[3]},
			}

		}

	}

	return networkDataMap, nil

}

func SortNetworkDataMap(networkDataMap map[string]NetworkData, isMBIn bool) (keysSortedInDesc []string) {

	keysDesc := make([]string, 0, len(networkDataMap))

	for key := range networkDataMap {
		keysDesc = append(keysDesc, key)
	}

	sort.SliceStable(keysDesc, func(i, j int) bool {
		var (
			totalMBI float64
			totalMBJ float64
			MBInI    []string
			MBInJ    []string
		)

		if isMBIn {
			MBInI = networkDataMap[keysDesc[i]].MBIn
			MBInJ = networkDataMap[keysDesc[j]].MBIn
		} else {
			MBInI = networkDataMap[keysDesc[i]].MBOut
			MBInJ = networkDataMap[keysDesc[j]].MBOut
		}

		for _, v := range MBInI {
			vFloat64, _ := strconv.ParseFloat(v, 64)
			totalMBI += vFloat64
		}

		for _, v := range MBInJ {
			vFloat64, _ := strconv.ParseFloat(v, 64)
			totalMBJ += vFloat64
		}

		averageMbInI := totalMBI / float64(len(MBInI))
		averageMBInJ := totalMBJ / float64(len(MBInJ))

		return averageMbInI > averageMBInJ
	})

	return keysDesc

}

func GetTopDesc(keysMB []string, topNumber int) (topKeysInDesc []string) {
	topKeys := make([]string, 0, 3)
	for i := 0; i < topNumber; i++ {
		topKeys = append(topKeys, keysMB[i])
	}

	return topKeys

}
