/// +build ignore

package main

import (
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type nodeValue string

func (nv nodeValue) String() string {
	return string(nv)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	nodes := []*widgets.TreeNode{
		{
			Value: nodeValue("Key 1"),
			Nodes: []*widgets.TreeNode{
				{
					Value: nodeValue("Key 1.1"),
					Nodes: []*widgets.TreeNode{
						{
							Value: nodeValue("Key 1.1.1"),
							Nodes: nil,
						},
						{
							Value: nodeValue("Key 1.1.2"),
							Nodes: nil,
						},
					},
				},
				{
					Value: nodeValue("Key 1.2"),
					Nodes: nil,
				},
			},
		},
		{
			Value: nodeValue("Key 2"),
			Nodes: []*widgets.TreeNode{
				{
					Value: nodeValue("Key 2.1"),
					Nodes: nil,
				},
				{
					Value: nodeValue("Key 2.2"),
					Nodes: nil,
				},
				{
					Value: nodeValue("Key 2.3"),
					Nodes: nil,
				},
			},
		},
		{
			Value: nodeValue("Key 3"),
			Nodes: nil,
		},
	}

	l := widgets.NewTree()
	l.ActiveBorderStyle = ui.NewStyle(ui.ColorRed)
	l.Active = true
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetNodes(nodes)

	l2 := widgets.NewTree()
	l2.ActiveBorderStyle = ui.NewStyle(ui.ColorRed)
	l2.TextStyle = ui.NewStyle(ui.ColorYellow)
	l2.WrapText = false
	l2.SetNodes(nodes)

	x, y := ui.TerminalDimensions()

	l.SetRect(0, 0, x/2, y)
	l2.SetRect(x/2+1, 0, x, y)

	uiWidgets := []*widgets.Tree{l, l2}

	for _, w := range uiWidgets {
		ui.Render(w)
	}

	var activeWidget *widgets.Tree

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
                for _, w := range(uiWidgets) {
                        if w.Active {
                                activeWidget = w
                                break
                        }
                }

		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			activeWidget.ScrollDown()
		case "k", "<Up>":
			activeWidget.ScrollUp()
		case "<C-d>":
			activeWidget.ScrollHalfPageDown()
		case "<C-u>":
			activeWidget.ScrollHalfPageUp()
		case "<C-f>":
			activeWidget.ScrollPageDown()
		case "<C-b>":
			activeWidget.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				activeWidget.ScrollTop()
			}
		case "<Home>":
			activeWidget.ScrollTop()
		case "<Enter>":
			activeWidget.ToggleExpand()
		case "G", "<End>":
			activeWidget.ScrollBottom()
		case "E":
			activeWidget.ExpandAll()
		case "C":
			activeWidget.CollapseAll()
		case "<Resize>":
			x, y := ui.TerminalDimensions()
			l.SetRect(0, 0, x/2, y)
			l2.SetRect(0, 0, x/2+1, y)
		case "<Tab>":
			l.Active = !l.Active
			l2.Active = !l2.Active
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		for _, w := range uiWidgets {
			ui.Render(w)
		}
	}
}
