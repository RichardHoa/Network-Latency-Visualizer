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
)

func main() {

	// Get working directory by env
	WORKING_DIR := os.Getenv("WORKING_DIR")

	// If env is not available, set the working directory to current working dir
	if WORKING_DIR == "" {
		pwdCommand := exec.Command("pwd")

		pwdByte, err := pwdCommand.Output()
		if err != nil {
			log.Fatal(err)
		}

		WORKING_DIR = strings.TrimSpace(string(pwdByte))
	}

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

	// If user do not input custom -a flag, then switch to collecting data mode

	// List existing cronjobs
	cronjobList, listCronjobErr := exec.Command("crontab", "-l").Output()
	if listCronjobErr != nil {
		log.Fatal(listCronjobErr)
	}

	// If we already set up cronjob, then we will collect data 
	if strings.Contains(string(cronjobList), "scanning") {
		fmt.Println("We already set up cronjob")
		fmt.Println("Perform automatic scanning, record network data and record Upload Speed and Download Speed")

		err := ping.RecordPingData(WORKING_DIR)
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
