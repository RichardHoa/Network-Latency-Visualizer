package table

import (
	"os"

	"github.com/RichardHoa/Network-Latency-Visualizer/network"
	"github.com/jedib0t/go-pretty/table"
)

func PrintNetworkingTable(networkDataMap map[string]network.NetworkData, keyDesc []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Process Name", "Incoming usage (MB)", "Outgoing usage (MB)", "Time"})

	for _, processName := range keyDesc {

		dataLength := len(networkDataMap[processName].Time)

		MBIn := networkDataMap[processName].MBIn[dataLength-1]

		MBOut := networkDataMap[processName].MBOut[dataLength-1]

		Time := networkDataMap[processName].Time[dataLength-1]

		t.AppendRow(table.Row{processName, MBIn, MBOut, Time})
	}
    t.AppendFooter(table.Row{"Table is sorted by MB In"})

	t.SetAutoIndex(true)
	t.SetStyle(table.StyleColoredBlackOnMagentaWhite)
	t.Render()
}
