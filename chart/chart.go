package chart

import (
	"io"
	"os"

	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/RichardHoa/Network-Latency-Visualizer/network"
	"github.com/RichardHoa/Network-Latency-Visualizer/ping"
	"github.com/RichardHoa/Network-Latency-Visualizer/table"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateLineItems(dataPoints []string) []opts.LineData {
	items := make([]opts.LineData, 0)
	for _, dataPoint := range dataPoints {
		items = append(items, opts.LineData{Value: dataPoint})
	}
	return items
}

func LineLabelPingChart(min []string, avg []string, max []string, stddev []string, timeString []string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Network latency chart",
			Subtitle: "All metrics are measured in ms (milisecond)\nYou can turn on and off the line by clicking on the legend",
			Link:     "https://github.com/go-echarts/go-echarts",
		}),
	)

	line.SetXAxis(timeString).
		AddSeries("Min latency", generateLineItems(min)).
		AddSeries("Avg latency", generateLineItems(avg)).
		AddSeries("Max latency", generateLineItems(max)).
		AddSeries("Standard deviation", generateLineItems(stddev)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				ShowSymbol: opts.Bool(true),
			}),
		)
	return line
}

func LineLabelNetworkPIDChart(TopDesc []string, networkDataMap map[string]network.NetworkData, MBType string) *charts.Line {
	title := fmt.Sprintf("Process %s network chart", MBType)
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: title,
			Link:  "https://github.com/go-echarts/go-echarts",
		}),
	)

	timeString := networkDataMap[TopDesc[0]].Time
	line.SetXAxis(timeString)

	for _, processName := range TopDesc {
		if MBType == "MBIn" {
			line.AddSeries(processName, generateLineItems(networkDataMap[processName].MBIn))
		} else if MBType == "MBOut" {
			line.AddSeries(processName, generateLineItems(networkDataMap[processName].MBOut))
		}
	}

	line.SetSeriesOptions(
		charts.WithLineChartOpts(opts.LineChart{
			ShowSymbol: opts.Bool(true),
		}),
	)
	return line

}

func CreatePingChart() {

	min, avg, max, stddev, timeString, readReportErr := ping.ReadPingReport("report/report.txt")
	if readReportErr != nil {
		log.Fatal(readReportErr)
	}

	page := components.NewPage()
	page.AddCharts(
		LineLabelPingChart(min, avg, max, stddev, timeString),
	)
	err := CreateAndOpenHTML(page, "chart/html/ping.html", "Network latency chart")
	if err != nil {
		log.Fatal(err)
	}

}

func CreateNetworkChart(WORKING_DIR string) error {

	networkDataMap, readNetworkDataErr := network.ReadNetworkData(WORKING_DIR)
	if readNetworkDataErr != nil {
		return readNetworkDataErr
	}

	keysMBInDesc := network.SortNetworkDataMap(networkDataMap, true)
	MBInDescTop := network.GetTopDesc(keysMBInDesc, 3)

	keysMBOutDesc := network.SortNetworkDataMap(networkDataMap, false)
	MBOutDescTop := network.GetTopDesc(keysMBOutDesc, 3)

	MBInPage := components.NewPage()
	MBInPage.AddCharts(
		LineLabelNetworkPIDChart(MBInDescTop, networkDataMap, "MBIn"),
	)

	MBOutPage := components.NewPage()
	MBOutPage.AddCharts(
		LineLabelNetworkPIDChart(MBOutDescTop, networkDataMap, "MBOut"),
	)

	openHTMLMBInErr := CreateAndOpenHTML(MBInPage, "chart/html/networkpid-in.html", "Network in chart")
	if openHTMLMBInErr != nil {
		return openHTMLMBInErr
	}

	openHTMLMBOutErr := CreateAndOpenHTML(MBOutPage, "chart/html/networkpid-out.html", "Network out chart")
	if openHTMLMBOutErr != nil {
		return openHTMLMBOutErr
	}

	table.PrintNetworkingTable(networkDataMap,keysMBInDesc)

	return nil
}

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
