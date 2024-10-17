package table

import (
    "os"

    "github.com/jedib0t/go-pretty/table"
)

func PrintTable() {
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
	
    t.AppendHeader(table.Row{"#", "First Name", "Last Name", "Salary"})
    t.AppendRows([]table.Row{
        {1, "Arya", "Stark", 3000},
        {20, "Jon", "Snow", 2000, "You know nothing, Jon Snow!"},
    })
    t.AppendRow([]interface{}{300, "Tyrion", "Lannister", 5000})
	t.SetAutoIndex(true)
	t.SetStyle(table.StyleColoredBlackOnBlueWhite)
    t.Render()
}