package chart

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/RichardHoa/Network-Latency-Visualizer/network"
	"github.com/RichardHoa/Network-Latency-Visualizer/ping"
	"github.com/RichardHoa/Network-Latency-Visualizer/speedtest"
	"github.com/RichardHoa/Network-Latency-Visualizer/table"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Function to generate line items for the chart
func generateLineItems(dataPoints []string) []opts.LineData {
	items := make([]opts.LineData, 0)
	for _, dataPoint := range dataPoints {
		items = append(items, opts.LineData{Value: dataPoint})
	}
	return items
}

// Function to set up the ping chart
func LineLabelPingChart(pingStats ping.PingStats) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Network latency chart",
			Subtitle: "All metrics are measured in ms (milisecond)\nYou can turn on and off the line by clicking on the legend",
			Link:     "https://github.com/go-echarts/go-echarts",
		}),
	)

	line.SetXAxis(pingStats.TimeString).
		AddSeries("Min latency", generateLineItems(pingStats.Min)).
		AddSeries("Avg latency", generateLineItems(pingStats.Avg)).
		AddSeries("Max latency", generateLineItems(pingStats.Max)).
		AddSeries("Standard deviation", generateLineItems(pingStats.Sttdev)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				ShowSymbol: opts.Bool(true),
			}),
		)
	return line
}

// Function to set up the network usage chart
func LineLabelProcessNetworkUsageChart(TopDesc []string, networkDataMap map[string]network.NetworkData, typeOfNetwork string) *charts.Line {
	var title string
	if typeOfNetwork == "received" {
		title = "Received Network Data"
	}

	if typeOfNetwork == "sent" {
		title = "Sent Network Data"
	}
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: title,
			Link:  "https://github.com/go-echarts/go-echarts",
		}),
	)

	// Get the time string of the process that either received or sent the most data
	
	processNameLongestTime := network.FindLongestTime(TopDesc, networkDataMap)
	networkDataMap = network.EqualizeTopKey(networkDataMap, TopDesc, processNameLongestTime)

	timeString := networkDataMap[processNameLongestTime].Time

	line.SetXAxis(timeString)

	for _, processName := range TopDesc {
		if typeOfNetwork == "received" {
			line.AddSeries(processName, generateLineItems(networkDataMap[processName].ReceivedMB))
		} else if typeOfNetwork == "sent" {
			fmt.Printf("process name: %s\n", processName)
			fmt.Printf("Time: %s\n", networkDataMap[processName].Time)
			line.AddSeries(processName, generateLineItems(networkDataMap[processName].SentMB))
		}
	}

	line.SetSeriesOptions(
		charts.WithLineChartOpts(opts.LineChart{
			ShowSymbol: opts.Bool(true),
		}),
	)
	return line
}

// Function to set up the speedtest chart
func LineLabelSpeedtestChart(DLSpeed []string, ULSpeed []string, timeString []string) *charts.Line {

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Speedtest chart",
			Subtitle: "All metrics are measured in MB (megabyte)",
			Link:     "https://github.com/go-echarts/go-echarts",
		}),
	)
	line.SetXAxis(timeString).
		AddSeries("Download Speed", generateLineItems(DLSpeed)).
		AddSeries("Upload Speed", generateLineItems(ULSpeed)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				ShowSymbol: opts.Bool(true),
			}),
		)
	return line
}

// Function to create the network latency chart
func CreatePingChart() {

	pingStats, readReportErr := ping.ReadPingReport("ping/ping.txt")
	if readReportErr != nil {
		log.Fatal(readReportErr)
	}

	page := components.NewPage()
	page.AddCharts(
		LineLabelPingChart(pingStats),
	)
	err := CreateAndOpenHTML(page, "chart/html/ping.html", "Network latency chart")
	if err != nil {
		log.Fatal(err)
	}

}

// Function to create the process network usage chart
func CreateNetworkChart(WORKING_DIR string) error {

	// Get the network data Map
	networkDataMap, readNetworkDataErr := network.ReadNetworkData(WORKING_DIR)
	if readNetworkDataErr != nil {
		return readNetworkDataErr
	}

	// Sort the map in descending order for received data
	receivedKeysDesc := network.SortNetworkDataMap(networkDataMap, true)
	// Get the top 3 keys with the most received data
	receivedKeysTop := network.GetTopDesc(receivedKeysDesc, 3)

	// Sort the map in descending order for sent data
	sentKeysDesc := network.SortNetworkDataMap(networkDataMap, false)
	// Get the top 3 keys with the most sent data
	sentKeysTop := network.GetTopDesc(sentKeysDesc, 3)

	// Create 2 charts
	receivedNetworkpage := components.NewPage()
	receivedNetworkpage.AddCharts(
		LineLabelProcessNetworkUsageChart(receivedKeysTop, networkDataMap, "received"),
	)

	sentNetworkPage := components.NewPage()
	sentNetworkPage.AddCharts(
		LineLabelProcessNetworkUsageChart(sentKeysTop, networkDataMap, "sent"),
	)

	receivedHTMLOpenErr := CreateAndOpenHTML(receivedNetworkpage, "chart/html/networkpid-in.html", "Received network data")
	if receivedHTMLOpenErr != nil {
		return receivedHTMLOpenErr
	}

	sentHTMLOpenErr := CreateAndOpenHTML(sentNetworkPage, "chart/html/networkpid-out.html", "Sent network data")
	if sentHTMLOpenErr != nil {
		return sentHTMLOpenErr
	}

	table.PrintNetworkingTable(networkDataMap, receivedKeysDesc)

	return nil
}

// Function to create the speedtest chart
func CreateSpeedtestChart() error {

	DLSpeed, UPSpeed, timeString, readReportErr := speedtest.ReadSpeedTestReport("speedtest/speedtest.txt")
	if readReportErr != nil {
		return readReportErr
	}

	page := components.NewPage()
	page.AddCharts(
		LineLabelSpeedtestChart(DLSpeed, UPSpeed, timeString),
	)
	OpenHTMLErr := CreateAndOpenHTML(page, "chart/html/speedtest.html", "Speedtest chart")
	if OpenHTMLErr != nil {
		return OpenHTMLErr
	}

	return nil
}

// Helper functions
func CreateAndOpenHTML(page *components.Page, filePath string, title string) error {

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	page.Render(io.MultiWriter(file))

	htmlContent, _ := os.ReadFile(filePath)

	htmlTitle := fmt.Sprintf("<title>%s</title>", title)

	updatedContent := strings.Replace(string(htmlContent), "<title>Awesome go-echarts</title>", htmlTitle, 1)

	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	openHTML := exec.Command("open", filePath)
	err = openHTML.Run()
	if err != nil {
		log.Fatal(err)
	}
	return nil

}

