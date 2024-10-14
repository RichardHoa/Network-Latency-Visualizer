package chart

import (
	"io"
	// "math/rand"
	"os"

	// "fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"log"
	"os/exec"
	"strings"
)

// var (
// 	itemCntLine = 6
// 	dataPoints  = []string{"Apple", "Banana", "Peach ", "Lemon", "Pear", "Cherry", "Something", "something", "something", "something", "something"}
// )

func generateLineItems(dataPoints []string) []opts.LineData {
	items := make([]opts.LineData, 0)
	for _, dataPoint := range dataPoints {
		items = append(items, opts.LineData{Value: dataPoint})
	}
	return items
}

func lineShowLabel(min []string, avg []string, max []string, stddev []string, timeString []string) *charts.Line {
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

type LineExamples struct{}

func (LineExamples) Examples() {

	report, openFileErr := os.ReadFile("report/report.txt")
	if openFileErr != nil {
		log.Fatal(openFileErr)
	}
	lines := strings.Split(string(report), "\n")
	lines = lines[:len(lines)-1]

	min := make([]string, 0)
	avg := make([]string, 0)
	max := make([]string, 0)
	stddev := make([]string, 0)
	timeString := make([]string, 0)

	for _, line := range lines {
		lineSlice := strings.Split(line, "|")
		time := lineSlice[1]
		data := lineSlice[0]
		stats := strings.Split(strings.Split(data, "=")[1], "/")

		min = append(min, stats[0])
		avg = append(avg, stats[1])
		max = append(max, stats[2])
		stddev = append(stddev, strings.Split(stats[3], " ")[0])
		timeString = append(timeString, time)
	}

	page := components.NewPage()
	page.AddCharts(
		lineShowLabel(min, avg, max, stddev, timeString),
	)
	f, err := os.Create("chart/html/line.html")
	if err != nil {
		log.Fatal(err)
	}
	page.Render(io.MultiWriter(f))

	openHTML := exec.Command("open", "./chart/html/line.html")
	err = openHTML.Run()
	if err != nil {
		log.Fatal(err)
	}
}
