package main

import (
	"fmt"
	"log"
	"os/exec"
	"github.com/joho/godotenv"
	"strings"
	"github.com/RichardHoa/Network-Latency-Visualizer/cronjob"
)

func main() {
	WORKING_DIR := "/Users/hoathaidang/Documents/bootdev/go-networking"

	fmt.Printf("Working dir is: %s\n", WORKING_DIR)

	setupCronJobErr := godotenv.Load()

	if setupCronJobErr != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Welcome to the Network-Latency-Visualizer!")
	fmt.Println("--------------------------------------------")


	output, listCronjobErr := exec.Command("crontab", "-l").Output()
	if listCronjobErr != nil {
		log.Fatal(listCronjobErr)
	}

	fmt.Printf("Your current cronjobs are: %s\n", string(output))



	if strings.Contains(string(output), "networking") {
		fmt.Println("We already set up cronjob")


	} else {
		fmt.Println("We didn't set up cronjob")
		err := cronjob.SetUpCronJob(WORKING_DIR)
		if err != nil {
			log.Fatal(err)
		}
	}

}


