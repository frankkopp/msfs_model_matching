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
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/frankkopp/MatchMaker/internal/livery"
	"github.com/frankkopp/MatchMaker/internal/rules"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	LiveryTableView *walk.TableView
	ScanButton      *walk.PushButton
	GenerateButton  *walk.PushButton
)

func parseTab() TabPage {

	// boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)
	// goodIcon, _ := walk.Resources.Icon("../img/check.ico")
	// badIcon, _ := walk.Resources.Icon("../img/stop.ico")

	model := NewLiveryModel()

	return TabPage{
		Title:  "Parse Liveries",
		Layout: VBox{},
		Children: []Widget{
			PushButton{
				AssignTo:  &ScanButton,
				Text:      "Scan",
				OnClicked: model.ScanLiveries,
			},
			TableView{
				AssignTo:         &LiveryTableView,
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
					// fmt.Printf("SelectedIndexes: %v\n", LiveryTableView.SelectedIndexes())
				},
				OnItemActivated: func() {
					if model.items[LiveryTableView.CurrentIndex()].Complete {
						model.items[LiveryTableView.CurrentIndex()].Process = !model.items[LiveryTableView.CurrentIndex()].Process
						model.handleUpdate()
					}
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
	TabBarHandle.SetCurrentIndex(1)
	StatusBar3.SetText(fmt.Sprintf("Generating %d rules...", rules.Counter))
	RulesText.SetText("")
	go m.buildXML()
}

func (m *LiveryModel) buildXML() {
	var output strings.Builder
	output.Grow(100_000)

	// Header
	output.WriteString(xml.Header)
	output.WriteString("\r\n\r\n<ModelMatchRuleSet>\r\n\r\n")

	// default rules
	fmt.Fprintf(&output, "<!-- DEFAULTS -->\r\n")
	for _, icaoKey := range rules.SortIcaoKeys(rules.Rules) {
		if icaoKey != "default" {
			continue
		}
		for typeKey := range rules.Rules[icaoKey] {
			fmt.Fprintf(&output, "<ModelMatchRule TypeCode=\"%s\" ModelName=\"", typeKey)
			for i, livery := range rules.Rules[icaoKey][typeKey] {
				if i != 0 {
					fmt.Fprint(&output, "//")
				}
				fmt.Fprintf(&output, "%s", livery)
			}
			fmt.Fprintf(&output, "\" />\r\n")
		}
	}
	fmt.Fprintf(&output, "\r\n")

	// ICAO based rules
	fmt.Fprintf(&output, "<!-- PER ICAO RULES -->\r\n")
	for _, icaoKey := range rules.SortIcaoKeys(rules.Rules) {
		if icaoKey == "default" {
			continue
		}
		fmt.Fprintf(&output, "<!-- ICAO:  %s -->\r\n", icaoKey)
		for _, baseKey := range rules.SortBaseKeys(rules.TypeVariations) {
			fmt.Fprintf(&output, "<!-- BASE: %s -->\r\n", baseKey)
			for _, typeKey := range rules.TypeVariations[baseKey] {
				if len(rules.Rules[icaoKey][typeKey]) == 0 {
					continue
				}
				fmt.Fprintf(&output, "<ModelMatchRule CallsignPrefix=\"%s\" TypeCode=\"%s\" ModelName=\"", icaoKey, typeKey)
				for i, livery := range rules.Rules[icaoKey][typeKey] {
					if i != 0 {
						fmt.Fprint(&output, "//")
					}
					fmt.Fprintf(&output, "%s", livery)
				}
				fmt.Fprintf(&output, "\" />\r\n")
			}
		}
		fmt.Fprintf(&output, "\r\n")
	}

	// Footer
	output.WriteString("\r\n</ModelMatchRuleSet>\r\n")

	RulesText.SetText(output.String())
	StatusBar3.SetText(fmt.Sprintf("Generated %d rules.", rules.Counter))
}

func (m *LiveryModel) ScanLiveries() {
	ScanButton.SetEnabled(false)
	GenerateButton.SetEnabled(false)
	LiveryTableView.SetEnabled(false)
	StatusBar1.SetText(fmt.Sprintf("Scanning %s ...", *config.Configuration.LiveryDirectory))

	go m.scanLiveries()
}

func (m *LiveryModel) scanLiveries() {
	liveries, err := livery.ScanLiveryFolder(*config.Configuration.LiveryDirectory)
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
	rules.CalculateRules(m.items)
	StatusBar3.SetText(fmt.Sprintf("Number of rules to be generated: %d", rules.Counter))
	LiveryTableView.SetEnabled(true)
	ScanButton.SetEnabled(true)
	GenerateButton.SetEnabled(rules.Counter > 0)
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
