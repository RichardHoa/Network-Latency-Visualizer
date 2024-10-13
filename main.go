package main

import (
	"fmt"
	"github.com/RichardHoa/Network-Latency-Visualizer/cronjob"
	"github.com/RichardHoa/Network-Latency-Visualizer/ping"
	"github.com/joho/godotenv"
	"github.com/nexidian/gocliselect"
	"log"
	"os"
	"os/exec"
	"strings"
	// "time"
)

func main() {
	loadENV := godotenv.Load()
	if loadENV != nil {
		log.Fatal("Error loading .env file")
	}

	WORKING_DIR := os.Getenv("WORKING_DIR")
	fmt.Printf("Working dir is: %s\n", WORKING_DIR)

	fmt.Println("Welcome to the Network-Latency-Visualizer!")
	fmt.Println("--------------------------------------------")

	if len(os.Args) > 1 && os.Args[1] == "-edit" {
		fmt.Println("Editing mode")

		menu := gocliselect.NewMenu("Choose your action")

		menu.AddItem("Remove cronjob", "remove cronjob")
		menu.AddItem("Help me", "help me")

		choice := menu.Display()

		switch choice {
		case "remove cronjob":
			fmt.Println("yeah")
		case "help me":
			fmt.Println("I will help you")

		}

	} else {

		output, listCronjobErr := exec.Command("crontab", "-l").Output()
		if listCronjobErr != nil {
			log.Fatal(listCronjobErr)
		}

		if strings.Contains(string(output), "scanning") {
			fmt.Println("We already set up cronjob")
			fmt.Println("Perform automatic scanning")

			err := ping.PingScanning(WORKING_DIR)
			if err != nil {
				log.Fatal(err)
			}

		} else {
			fmt.Println("We haven't set up cronjob")
			err := cronjob.SetUpCronJob(WORKING_DIR)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}


