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
	"github.com/RichardHoa/Network-Latency-Visualizer/speedtest"
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

func LineLabelNetworkPIDChart(TopDesc []string, networkDataMap map[string]network.NetworkData, MBType string) *charts.Line {
	var title string
	if MBType == "MBIn" {
		title = "Incoming Network Data"
	} else {
		title = "Outgoing Network Data"
	}
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
			line.AddSeries(processName, generateLineItems(networkDataMap[processName].ReceivedMB))
		} else if MBType == "MBOut" {
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

	openHTMLMBInErr := CreateAndOpenHTML(MBInPage, "chart/html/networkpid-in.html", "Incoming network data")
	if openHTMLMBInErr != nil {
		return openHTMLMBInErr
	}

	openHTMLMBOutErr := CreateAndOpenHTML(MBOutPage, "chart/html/networkpid-out.html", "Outgoing network data")
	if openHTMLMBOutErr != nil {
		return openHTMLMBOutErr
	}

	table.PrintNetworkingTable(networkDataMap, keysMBInDesc)

	return nil
}

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
