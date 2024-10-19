package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/RichardHoa/Network-Latency-Visualizer/cronjob"
	"github.com/RichardHoa/Network-Latency-Visualizer/network"
	"github.com/RichardHoa/Network-Latency-Visualizer/ping"
	"github.com/RichardHoa/Network-Latency-Visualizer/speedtest"

	// "github.com/RichardHoa/Network-Latency-Visualizer/table"
	"github.com/joho/godotenv"
	// "time"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	WORKING_DIR := os.Getenv("WORKING_DIR")

	if WORKING_DIR == "" {

		pwdCommand := exec.Command("pwd")

		pwdByte, err := pwdCommand.Output()
		if err != nil {
			log.Fatal(err)
		}

		WORKING_DIR = strings.TrimSpace(string(pwdByte))

	}

	// "/Users/hoathaidang/Documents/bootdev/go-networking"

	if len(os.Args) > 1 && os.Args[1] != "-a" {
		fmt.Println("Invalid command line")
		os.Exit(1)
	}

	// Welcome message
	fmt.Println("Welcome to the Network-Latency-Visualizer!")
	fmt.Println("--------------------------------------------")

	// Get into advanced mode by using -a
	if len(os.Args) > 1 && os.Args[1] == "-a" {
		RunTerminal(WORKING_DIR)
	}

	// If user do not input custom -a flag, then we will set up cronjob

	// List existing cronjobs
	cronjobList, listCronjobErr := exec.Command("crontab", "-l").Output()
	if listCronjobErr != nil {
		log.Fatal(listCronjobErr)
	}

	// If we already set up cronjob, then we will perform automatic scanning
	if strings.Contains(string(cronjobList), "scanning") {
		fmt.Println("We already set up cronjob")
		fmt.Println("Perform automatic scanning, record network data and record Upload Speed and Download Speed")

		err := ping.PingScanning(WORKING_DIR)
		if err != nil {
			log.Fatal(err)
		}

		recordErr := network.RecordNetworkData(WORKING_DIR)
		if recordErr != nil {
			log.Fatal(err)
		}

		speedtestErr := speedtest.RecordSpeedTestData(WORKING_DIR)
		if speedtestErr != nil {
			log.Fatal(err)
		}
		os.Exit(1)

	} else {
		fmt.Println("Look like we haven't set up the cronjob, let's do it now!")
		err := cronjob.SetUpCronJob(WORKING_DIR)
		if err != nil {
			log.Fatal(err)
		}
	}

}
