package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	// "time"

	"github.com/joho/godotenv"
	"strconv"
	"strings"
)

func main() {
	WORKING_DIR := "/Users/hoathaidang/Documents/bootdev/go-networking"
	fmt.Printf("Working dir is: %s\n", WORKING_DIR)

	setupCronJobErr := godotenv.Load()

	if setupCronJobErr != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Welcome to the Network-Latency-Visualizer!")

	fmt.Println("For the first step, please choose how often you want to check your network latency")
	fmt.Println("Minimum: 1 minutes, maximum: 1 day")

	timeInMinutes := askForTimeInput()

	fmt.Printf("You chose %d minutes\n", timeInMinutes)

	// output, setupCronJobErr := exec.Command("ping", "google.com", "-c", "10").Output()
	// if setupCronJobErr != nil {
	// 	log.Fatal(setupCronJobErr)
	// }

	// outputString := string(output)
	// outputArray := strings.Split(outputString, "\n")
	// fmt.Println(outputArray[len(outputArray)-2])

	// err := setUpCronJob("1 * * * *", WORKING_DIR)

	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func setUpCronJob(timeString string, WORKING_DIR string) error {

	cronjobWorkingDir := WORKING_DIR + "/cronjob/cron.txt"

	crontabJobs, setupCronJobErr := exec.Command("crontab", "-l").Output()
	if setupCronJobErr != nil {
		return setupCronJobErr
	}

	file, openFileErr := os.OpenFile(cronjobWorkingDir, os.O_RDWR|os.O_CREATE, 0777)
	if openFileErr != nil {
		return openFileErr
	}

	defer file.Close()

	cronjob := "* * * * * " + WORKING_DIR + "/networking" + "\n"

	_, writeExistingCronJobErr := file.Write(crontabJobs)
	if writeExistingCronJobErr != nil {
		return writeExistingCronJobErr
	}

	_, writeNewCronJobErr := file.WriteString(cronjob)
	if writeNewCronJobErr != nil {
		return writeNewCronJobErr
	}

	setupCronJob := exec.Command("crontab", cronjobWorkingDir)

	setupCronJobErr = setupCronJob.Run()
	if setupCronJobErr != nil {
		return setupCronJobErr
	}

	return nil
}

func calculateTimeString(timeInMinutes int) string {
	return ""

}

func askForTimeInput() (timeInMinutes int) {
	var time int

	for {
		var inputTime string
		var timeMark string
		var min int
		fmt.Println("Example input: 1 mins, 5 mins, 1 hrs, 15 hrs.")
		fmt.Printf("Your chosen time is: ")

		fmt.Scanf("%s %s", &inputTime, &timeMark)

		if strings.Contains(timeMark, "mins") {
			minString := strings.Split(inputTime, "mins")[0]
			min, _ = strconv.Atoi(minString)
		} else if strings.Contains(timeMark, "hrs") {
			minString := strings.Split(inputTime, "hrs")[0]
			hours, _ := strconv.Atoi(minString)
			min = hours * 60
		} else {
			fmt.Println("It's either minutes or hours")
			continue
		}

		if min >= 1 && min <= 86400 {
			time = min
			break
		} else {
			fmt.Println("Please remember the maximum time")
		}

	}

	return time
}
