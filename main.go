package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Welcome to the Network-Latency-Visualizer!")

	crontabJobs, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		log.Fatal(err)
	}
	crontabArray := string(crontabJobs)
	if strings.Contains(crontabArray, "networking") {
		fmt.Println("Crontab has been established")
	}

	file, openFileErr := os.OpenFile("report.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if openFileErr != nil {
		log.Fatal(openFileErr)
	}

	defer file.Close()


	data := []byte("1 * * * * /Users/hoathaidang/Documents/bootdev/go-networking/networking\n")

	_, writeFileErr := file.Write(data)
	if writeFileErr != nil {
		log.Fatal(writeFileErr)
	}



	// var lengthInMinutes int
	// fmt.Println("For the first step, please choose how often you want to check your network latency")
	// fmt.Println("Minimum: 1 minutes, maximum: 1 day (86400 minutes)")
	// fmt.Printf("Please type the number in minutes: ")
	// fmt.Scan(&lengthInMinutes)

	// output, err := exec.Command("ping", "google.com", "-c", "10").Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// outputString := string(output)
	// outputArray := strings.Split(outputString, "\n")
	// fmt.Println(outputArray[len(outputArray)-2])
	// fmt.Printf("The length is %d minutes\n", lengthInMinutes)
}
