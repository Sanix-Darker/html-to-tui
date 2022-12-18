package main

import (
	"fmt"
	"io/ioutil"

	"github.com/rivo/tview"
)

func main() {
	htmlData, err := ioutil.ReadFile("table.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	doc, err := parseHTML(htmlData)
	if err != nil {
		fmt.Println(err)
		return
	}

	app := tview.NewApplication()

	tableData, err := extractTableData(doc)
	if err != nil {
		fmt.Println(err)
		return
	}

	t := tview.NewTable().SetBorders(true)
	for i, rowData := range tableData {
		for j, cell := range rowData {
			t.SetCell(i, j, cell)
		}
	}
	app.SetRoot(t, true).Run()
}
