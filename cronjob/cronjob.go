package cronjob

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Calculate timestring in cronjob format, ex: * * * * *
func calculateTimeString(timeInMinutes int) (timeStringInCronjob string) {
	// Initialize the first and second position
	firstPosition := "*"
	secondPosition := "*"

	// If 24 hours
	if timeInMinutes == 1440 {
		firstPosition = "0"
		secondPosition = "0"

	}
	// If 60  minute
	if timeInMinutes == 60 {
		firstPosition = "0"
	}

	// If less than 60 minutes and not 1 minutes
	if timeInMinutes < 60 && timeInMinutes != 1 {
		firstPosition += "/" + strconv.Itoa(timeInMinutes)
	}

	// If more than 60 minutes and below 24 hours
	if timeInMinutes > 60 && timeInMinutes < 1440 {
		secondPosition += "/" + strconv.Itoa(timeInMinutes/60)
	}

	finalString := firstPosition + " " + secondPosition + " * * *"

	return finalString

}

// Ask for user time input, return time in minutes
func askForTimeInput() (timeInMinutes int) {
	var time int

	fmt.Println("Please choose how often you want to check your network latency")
	fmt.Println("Minimum: 1 minutes, maximum: 1 day")

	for {
		var inputTime string
		var timeMark string
		var min int
		fmt.Println("Possible input: 1 - 60 mins, 1 - 24 hrs. We don't support decimal hrs and mins")
		fmt.Printf("Your chosen time is: ")

		// Get user input
		fmt.Scanf("%s %s", &inputTime, &timeMark)

		// If user input in mins, convert that to min
		if strings.Contains(timeMark, "mins") {
			minString := strings.Split(inputTime, "mins")[0]
			min, _ = strconv.Atoi(minString)
			// If user input in hrs, convert that to min
		} else if strings.Contains(timeMark, "hrs") {
			minString := strings.Split(inputTime, "hrs")[0]
			hours, _ := strconv.Atoi(minString)
			min = hours * 60
			// If user input is not hrs or mins, prompt again
		} else {
			fmt.Println("It's either minutes or hours")
			continue
		}
		//  If min is between 1 and 60 or min is between 1 hr and 24 hrs and the number is whole
		if min >= 1 && min <= 60 || min > 60 && min%60 == 0 && min <= 1440 {
			time = min
			break
			// If the time is not within the limit, prompt user again
		} else {
			fmt.Println("Please remember the maximum time")
		}

	}

	return time
}

// Save the cronjob to the system, has 2 mode: add and remove
// add mode adds the cronjob to the system
// remove mode removes the cronjob from the system
func SaveCronJob(timeStringInCronjob string, WORKING_DIR string, mode string) error {
	// Initialize the path of txt.file
	cronTXTPath := WORKING_DIR + "/cronjob/cron.txt"

	// IF the file exist, then delete it
	if _, err := os.Stat(cronTXTPath); err == nil {
		e := os.Remove(cronTXTPath)
		if e != nil {
			return e
		}
	}

	// Create a new file
	file, openFileErr := os.OpenFile(cronTXTPath, os.O_RDWR|os.O_CREATE, 0777)
	if openFileErr != nil {
		return openFileErr
	}
	defer file.Close()

	// List existing cronjobs
	crontabJobs, setupCronJobErr := exec.Command("crontab", "-l").Output()
	if setupCronJobErr != nil {
		return setupCronJobErr
	}

	if mode == "remove" {
		// Convert []byte to string array
		cronjobArray := strings.Split(string(crontabJobs), "\n")
		// Scanning target to remove
		scanningCronjob := WORKING_DIR + "/scanning"
		envEnvironment := "WORKING_DIR=" + WORKING_DIR
		// Remove the line that contain the scanning target
		for index, cronjob := range cronjobArray {
			if strings.Contains(cronjob, scanningCronjob) {
				cronjobArray = append(cronjobArray[:index], cronjobArray[index+1:]...)
			}

			if strings.Contains(cronjob, envEnvironment) {
				cronjobArray = append(cronjobArray[:index], cronjobArray[index+1:]...)
			}
		}
		// Convert string array to []byte
		joinedString := strings.Join(cronjobArray, "\n")
		// Assign the []byte to the existing cronjobs
		crontabJobs = []byte(joinedString)

	}

	if mode == "add" {
		_, writeENV := file.WriteString("WORKING_DIR=" + WORKING_DIR + "\n")
		if writeENV != nil {
			return writeENV
		}

	}

	// Write existing cronjobs to the file
	_, writeExistingCronJobErr := file.Write(crontabJobs)
	if writeExistingCronJobErr != nil {
		return writeExistingCronJobErr
	}

	if mode == "add" {
		// Create a new cronjob string
		cronjob := timeStringInCronjob + " " + WORKING_DIR + "/scanning >> /tmp/scanning.out 2>> /tmp/scanning.err" + "\n"
		// Write the new cronjob to the file
		_, writeNewCronJobErr := file.WriteString(cronjob)
		if writeNewCronJobErr != nil {
			return writeNewCronJobErr
		}

	}

	// Set up the cronjob to the system
	setupCronJob := exec.Command("crontab", cronTXTPath)
	setupCronJobErr = setupCronJob.Run()
	if setupCronJobErr != nil {
		return setupCronJobErr
	}

	// If no error then return nil
	return nil
}

// Set up cronjob
func SetUpCronJob(WORKING_DIR string) error {
	// Ask for time in minutes
	timeInMinutes := askForTimeInput()
	fmt.Printf("You chose %d minutes\n", timeInMinutes)

	// Calculate time string in cronjob
	timeStringInCronjob := calculateTimeString(timeInMinutes)
	fmt.Printf("Your cronjob timestring is: %s\n", timeStringInCronjob)
	fmt.Println("Please allow the script to set the cronjob by clicking allow")

	// Save the cronjob
	err := SaveCronJob(timeStringInCronjob, WORKING_DIR, "add")

	return err

}
