/*
 *  VATSIM vPilot MatchMaker
 *
 *  MIT License
 *
 *  Copyright (c) 2021 Frank Kopp
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 *
 */

package ui

import (
	"fmt"
	"sort"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/frankkopp/MatchMaker/internal/livery"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	configuration config.Config

	GenerateButton *walk.PushButton
)

func parseTab(conf config.Config) TabPage {
	configuration = conf

	// boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)
	// goodIcon, _ := walk.Resources.Icon("../img/check.ico")
	// badIcon, _ := walk.Resources.Icon("../img/stop.ico")

	model := NewLiveryModel()
	var tv *walk.TableView

	return TabPage{
		Title:  "Parse Liveries",
		Layout: VBox{},
		Children: []Widget{
			PushButton{
				Text:      "Scan",
				OnClicked: model.ScanLiveries,
			},
			TableView{
				AssignTo:         &tv,
				AlternatingRowBG: true,
				CheckBoxes:       false,
				ColumnsOrderable: true,
				MultiSelection:   true,
				Columns: []TableViewColumn{
					{Title: "Process", Width: 50, Alignment: AlignCenter},
					{Title: "Livery Configuration File", Width: 650},
					{Title: "Base Container", Width: 200},
					{Title: "Title", Width: 200},
					{Title: "ICAO", Width: 50},
					{Title: "Has Custom Config", Width: 50},
				},
				StyleCell: func(style *walk.CellStyle) {
					item := model.items[style.Row()]

					switch style.Col() {
					case 0:
						if item.Complete {
							style.BackgroundColor = walk.RGB(162, 202, 112)
						} else {
							style.BackgroundColor = walk.RGB(204, 102, 102)
						}
					case 1:
						// placeholder
					case 2:
						// placeholder
					case 3:
						// placeholder
					case 4:
						// placeholder
					}

				},
				Model: model,
				OnSelectedIndexesChanged: func() {
					// fmt.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
				},
				OnItemActivated: func() {
					model.items[tv.CurrentIndex()].Process = !model.items[tv.CurrentIndex()].Process
					model.handleUpdate()
				},
			},
			PushButton{
				AssignTo:  &GenerateButton,
				Text:      "Generate Rules File",
				Enabled:   false,
				OnClicked: model.GenerateRules,
			},
		},
	}
}

type LiveryModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*livery.Livery
}

func NewLiveryModel() *LiveryModel {
	m := new(LiveryModel)
	return m
}

func (m *LiveryModel) GenerateRules() {
	StatusBar3.SetText(fmt.Sprintf("Generated rules: %d", m.QueuedCount()))
	RulesText.SetText(fmt.Sprintf("We would generate at least %d rules here!", m.QueuedCount()))
	RulesTabHandle.SetFocus()
}

func (m *LiveryModel) ScanLiveries() {
	liveries, err := livery.ScanLiveryFolder(*configuration.LiveryDirectory)
	if err != nil {
		return
	}
	// Notify TableView and other interested parties about the reset.
	m.items = liveries
	m.Sort(m.sortColumn, m.sortOrder)
	m.handleUpdate()
}

func (m *LiveryModel) handleUpdate() {
	m.PublishRowsReset()
	StatusBar1.SetText(fmt.Sprintf("Number of liveries found: %d", m.RowCount()))
	StatusBar2.SetText(fmt.Sprintf("Number of liveries queued: %d", m.QueuedCount()))
	GenerateButton.SetEnabled(m.QueuedCount() > 0)
}

func (m *LiveryModel) QueuedCount() int {
	i := 0
	for _, item := range m.items {
		if item.Process && item.Complete {
			i++
		}
	}
	return i
}

// RowCount Called by the TableView from SetModel and every time the model publishes a RowsReset event.
func (m *LiveryModel) RowCount() int {
	return len(m.items)
}

// Value Called by the TableView when it needs the text to display for a given cell.
func (m *LiveryModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Process

	case 1:
		return item.AircraftCfgFile

	case 2:
		return item.BaseContainer

	case 3:
		return item.Title

	case 4:
		return item.Icao

	case 5:
		return item.Custom
	}

	panic("unexpected col")
}

// Sort Called by the TableView to sort the model.
func (m *LiveryModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		compare := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}

		// sort for Process first the Complete
		switch m.sortColumn {
		case 0:
			aS := 0
			bS := 0
			if a.Process {
				aS += 2
			}
			if a.Complete {
				aS++
			}
			if b.Process {
				bS += 2
			}
			if b.Complete {
				bS++
			}
			return compare(aS > bS)

		case 1:
			return compare(a.AircraftCfgFile < b.AircraftCfgFile)

		case 2:
			return compare(a.BaseContainer < b.BaseContainer)

		case 3:
			return compare(a.Title < b.Title)

		case 4:
			return compare(a.Icao < b.Icao)

		case 5:
			return compare(a.Custom)

		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}
