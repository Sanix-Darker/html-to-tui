package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/rivo/tview"
	"golang.org/x/net/html"
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
		for j, cellData := range rowData {
			if j%2 == 0 {
				t.SetCell(i, j/2, tview.NewTableCell(cellData.Text).
					SetAlign(tview.AlignCenter))
			}
		}
		app.SetRoot(t, true).Run()
	}
}

func extractTableData(doc *html.Node) ([][]*tview.TableCell, error) {
	var tableData [][]*tview.TableCell
	tableElements := extractTableElements(doc)
	for _, table := range tableElements {
		rows := extractTableRows(table)
		for _, row := range rows {
			cells := extractTableCells(row)
			var rowData []*tview.TableCell
			for _, cell := range cells {
				text := extractText(cell)
				rowData = append(rowData, tview.NewTableCell(text))
			}
			tableData = append(tableData, rowData)
		}
	}
	return tableData, nil
}

func extractTableCells(row *html.Node) []*html.Node {
	var cells []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "td" || n.Data == "th") {
			cells = append(cells, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(row)
	return cells
}

func extractTableRows(table *html.Node) []*html.Node {
	var rows []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			rows = append(rows, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(table)
	return rows
}

func extractTableElements(doc *html.Node) []*html.Node {
	var tableElements []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			tableElements = append(tableElements, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return tableElements
}

func parseHTML(htmlData []byte) (*html.Node, error) {
	return html.Parse(strings.NewReader(string(htmlData)))
}

func extractText(n *html.Node) string {
	var text string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return text
}
