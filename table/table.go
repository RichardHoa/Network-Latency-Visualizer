package table

import (
	"os"

	"github.com/RichardHoa/Network-Latency-Visualizer/network"
	"github.com/jedib0t/go-pretty/table"
)

// Function to print out beautiful table along with network consumption chart
func PrintNetworkingTable(networkDataMap map[string]network.NetworkData, keyDesc []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	// Set the header
	t.AppendHeader(table.Row{"Process Name", "Incoming data (MB)", "Outgoing data (MB)", "Time"})

	// Key is sorted by MBIn (incoming network)
	for _, processName := range keyDesc {
		dataLength := len(networkDataMap[processName].Time)

		// Get the latest network in data
		MBIn := networkDataMap[processName].ReceivedMB[dataLength-1]
		// Get the latest network out data
		MBOut := networkDataMap[processName].SentMB[dataLength-1]
		// Get the latest time recorded
		Time := networkDataMap[processName].Time[dataLength-1]
		// Append them all to the row
		t.AppendRow(table.Row{processName, MBIn, MBOut, Time})
	}
	t.AppendFooter(table.Row{"Table is sorted by Incoming network"})
	// Set auto index
	t.SetAutoIndex(true)
	// Set style
	t.SetStyle(table.StyleColoredBlackOnMagentaWhite)
	// Render the table to the terminal
	t.Render()
}
