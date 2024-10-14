package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/RichardHoa/Network-Latency-Visualizer/cronjob"
	"github.com/RichardHoa/Network-Latency-Visualizer/chart"
	"github.com/RichardHoa/Network-Latency-Visualizer/ping"
	// "github.com/joho/godotenv"
	"github.com/nexidian/gocliselect"
	// "time"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] != "-a" {
		fmt.Println("Invalid command line")
		os.Exit(1)
	}

	// WORKING_DIR here, important input for the whole program
	WORKING_DIR := "/Users/hoathaidang/Documents/bootdev/go-networking"
	// Testing print
	fmt.Printf("Working dir is: %s\n", WORKING_DIR)
	// Welcome message
	fmt.Println("Welcome to the Network-Latency-Visualizer!")
	fmt.Println("--------------------------------------------")

	// Get into advanced mode by using -a
	if len(os.Args) > 1 && os.Args[1] == "-a" {
		fmt.Println("advanced mode")

		// Create a terminal menu for the user
		menu := gocliselect.NewMenu("Choose your action")

		// Create option for the user
		menu.AddItem("Remove cronjob", "remove cronjob")
		menu.AddItem("Chart", "chart")
		menu.AddItem("Help me", "help me")

		// Get the choice from the user
		choice := menu.Display()

		switch choice {
		// Remove existing cronjob
		case "remove cronjob":
			cronjob.SaveCronJob("", WORKING_DIR, "remove")
		case "chart":
			chart.LineExamples{}.Examples()

		// Display help message
		case "help me":
			fmt.Println("To edit a cronjob, remove it first then do ./scanning")
		}

		os.Exit(0)

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
		fmt.Println("Perform automatic scanning")

		err := ping.PingScanning(WORKING_DIR)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		fmt.Println("Look like we haven't set up the cronjob, let's do it now!")
		err := cronjob.SetUpCronJob(WORKING_DIR)
		if err != nil {
			log.Fatal(err)
		}
	}

}
