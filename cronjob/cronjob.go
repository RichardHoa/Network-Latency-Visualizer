package cronjob

import (
	"fmt"
	"os/exec"
	"os"
	"strconv"
	"strings"
)

func calculateTimeString(timeInMinutes int) (timeStringInCronjob string) {
	firstPosition := "*"
	secondPosition := "*"

	if timeInMinutes == 1440 {
		firstPosition = "0"
		secondPosition = "0"

	}

	if timeInMinutes == 60 {
		firstPosition = "0"
	}

	if timeInMinutes < 60 {
		firstPosition += "/" + strconv.Itoa(timeInMinutes)
	}

	if timeInMinutes > 60 && timeInMinutes < 1440 {
		secondPosition += "/" + strconv.Itoa(timeInMinutes/60)
	}



	finalString := firstPosition + " " + secondPosition + " * * *"

	return finalString

}

func askForTimeInput() (timeInMinutes int) {
	var time int

	fmt.Println("For the first step, please choose how often you want to check your network latency")
	fmt.Println("Minimum: 1 minutes, maximum: 1 day")

	for {
		var inputTime string
		var timeMark string
		var min int
		fmt.Println("Possible input: 1 - 60 mins, 1 - 24 hrs. We don't support 6.5 hrs or anything like that")
		fmt.Printf("Your chosen time is: ")

		fmt.Scanf("%s %s", &inputTime, &timeMark)

		if strings.Contains(timeMark, "mins") {
			minString := strings.Split(inputTime, "mins")[0]
			min, _ = strconv.Atoi(minString)
			// fmt.Printf("Mins: %d\n", min)
		} else if strings.Contains(timeMark, "hrs") {
			minString := strings.Split(inputTime, "hrs")[0]
			hours, _ := strconv.Atoi(minString)
			min = hours * 60
		} else {
			fmt.Println("It's either minutes or hours")
			continue
		}

		if min >= 1 && min <= 60 || min > 60 && min%60 == 0 {
			time = min
			break
		} else {
			fmt.Println("Please remember the maximum time")
		}

	}

	return time
}


func saveCronJob(timeStringInCronjob string, WORKING_DIR string) error {

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

	cronjob := timeStringInCronjob + " " + WORKING_DIR + "/scanning" + "\n"

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

func SetUpCronJob(WORKING_DIR string) error {
	timeInMinutes := askForTimeInput()

	fmt.Printf("You chose %d minutes\n", timeInMinutes)

	timeStringInCronjob := calculateTimeString(timeInMinutes)

	fmt.Printf("Your cronjob timestring is: %s\n", timeStringInCronjob)

	fmt.Println("Please allow the script to set the cronjob by clicking allow")

	err := saveCronJob(timeStringInCronjob, WORKING_DIR)

	return err 

}



