package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/net/html"
)

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
				color, err := extractColor(cell)
				if err != nil {
					return nil, err
				}
				bgColor, err := extractBackgroundColor(cell)
				if err != nil {
					bgColor = tcell.ColorDefault
				}
				rowData = append(rowData, tview.NewTableCell(text).
					SetTextColor(color).
					SetBackgroundColor(bgColor).
					SetStyle(tcell.StyleDefault))
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

func extractColor(n *html.Node) (tcell.Color, error) {
	for _, a := range n.Attr {
		if a.Key == "style" {
			for _, s := range strings.Split(a.Val, ";") {
				if strings.HasPrefix(s, "color:") {
					color := strings.TrimSpace(strings.TrimPrefix(s, "color:"))
					return tcell.GetColor(color), nil
				}
			}
		}
	}
	return tcell.ColorWhite, nil
}

func extractBackgroundColor(n *html.Node) (tcell.Color, error) {
	for _, a := range n.Attr {
		if a.Key == "style" {
			for _, s := range strings.Split(a.Val, ";") {
				if strings.HasPrefix(s, "background-color:") {
					color := strings.TrimSpace(strings.TrimPrefix(s, "background-color:"))
					return tcell.GetColor(color), nil
				}
			}
		}
	}
	return tcell.ColorWhite, nil
}
