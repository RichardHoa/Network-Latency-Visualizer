package main

import (
	"fmt"
	"github.com/RichardHoa/Network-Latency-Visualizer/chart"
	"github.com/RichardHoa/Network-Latency-Visualizer/cronjob"
	"github.com/RichardHoa/Network-Latency-Visualizer/speedtest"
	"github.com/nexidian/gocliselect"
	"log"
	"os"
	"os/exec"
)

func RunTerminal(WORKING_DIR string) {

	fmt.Println("advanced mode")

	// Create a terminal menu for the user
	menu := gocliselect.NewMenu("What do you want to do?")

	// Create option for the user
	menu.AddItem("Cronjob options", "cronjob options")
	menu.AddItem("Show chart of each network", "network pid")
	menu.AddItem("Show network latency chart", "chart")
	menu.AddItem("Speed testing", "speed testing")
	menu.AddItem("Quit", "quit")

	for {
		clearTerminal()
		// Get the choice from the user
		choice := menu.Display()

		switch choice {
		case "network pid":
			err := chart.CreateNetworkChart(WORKING_DIR)
			if err != nil {
				log.Fatal(err)
			}
			// table.PrintTable()
			os.Exit(1)

		case "cronjob options":
			cronJobOPtions(WORKING_DIR)

		case "chart":
			chart.CreatePingChart()
			os.Exit(1)

		case "speed testing":
			fmt.Println("We are running speed testing, please wait....")
			speedtest.SpeedTesting()
			os.Exit(1)

		case "edit cronjob":
			cronjob.SaveCronJob("", WORKING_DIR, "remove")
			err := cronjob.SetUpCronJob(WORKING_DIR)
			if err != nil {
				log.Fatal(err)
			}

		case "quit":
			fmt.Println("Goodbye! See you later")
			os.Exit(1)

		}

	}

}

func cronJobOPtions(WORKING_DIR string) {
	menu := gocliselect.NewMenu("Cronjob options")

	menu.AddItem("Edit cronjob", "edit cronjob")
	menu.AddItem("Remove cronjob", "remove cronjob")

	clearTerminal()
	choice := menu.Display()

	switch choice {

	case "edit cronjob":
		cronjob.SaveCronJob("", WORKING_DIR, "remove")
		err := cronjob.SetUpCronJob(WORKING_DIR)
		if err != nil {
			log.Fatal(err)
		}
	case "remove cronjob":
		cronjob.SaveCronJob("", WORKING_DIR, "remove")
	}
}

func clearTerminal() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()

}
