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
	model           = NewLiveryModel()
)

func parseTab() TabPage {

	return TabPage{
		Title:  "Parse Liveries",
		Layout: VBox{},
		Children: []Widget{
			PushButton{
				AssignTo:  &ScanButton,
				Text:      fmt.Sprintf("Scan: %s", config.Configuration.Ini.Section("paths").Key("liveryDir").Value()),
				OnClicked: model.ScanLiveriesAction,
			},
			TableView{
				AssignTo:         &LiveryTableView,
				AlternatingRowBG: true,
				CheckBoxes:       false,
				ColumnsOrderable: true,
				ColumnsSizable:   true,
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
					if item.Custom {
						style.TextColor = walk.RGB(0, 130, 40)
					}
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
				ContextMenuItems: []MenuItem{
					Action{
						Text:        "Activate Livery",
						OnTriggered: OnItemActivatedAction,
					},
					Action{
						Text:        "Deactivate Livery",
						OnTriggered: OnItemDeactivatedAction,
					},
					Action{
						Text:        "Edit custom",
						OnTriggered: OnItemEditAction,
					},
					Action{
						Text:        "Remove custom",
						OnTriggered: OnItemRemoveCustomAction,
					},
				},
				OnItemActivated: func() {
					if model.items[LiveryTableView.CurrentIndex()].Process {
						OnItemDeactivatedAction()
					} else {
						OnItemActivatedAction()
					}
				},
			},
		},
	}
}

func OnItemRemoveCustomAction() {
	if len(LiveryTableView.SelectedIndexes()) == 0 {
		return
	}
	// we only allow to edit one item
	first := LiveryTableView.SelectedIndexes()[0]
	LiveryTableView.SetSelectedIndexes([]int{first})

	// restore list data and remove custom data
	selectedItem := model.items[first]
	aircraftCfgFile := selectedItem.AircraftCfgFile
	if !config.Configuration.Custom.HasEntry(aircraftCfgFile) {
		return
	}
	selectedItem.Icao = config.Configuration.Custom.GetEntry(aircraftCfgFile).OriginalIcao
	selectedItem.Custom = false
	if err := config.Configuration.Custom.RemoveEntry(aircraftCfgFile); err != nil {
		fmt.Printf("Could not remove entry for: %s\n", aircraftCfgFile)
	}
	config.Configuration.UpdateIniCustomData()
}

func OnItemEditAction() {
	if len(LiveryTableView.SelectedIndexes()) == 0 {
		return
	}
	// we only allow to edit one item
	first := LiveryTableView.SelectedIndexes()[0]
	LiveryTableView.SetSelectedIndexes([]int{first})

	if dlg, err := EditDialog(MainWindowHandle, LiveryTableView.SelectedIndexes()); err != nil {
		fmt.Print(err)
	} else if dlg == walk.DlgCmdOK {
		model.handleUpdate()
	}
}

func OnItemDeactivatedAction() {
	customData := config.Configuration.Custom

	for _, i := range LiveryTableView.SelectedIndexes() {
		item := model.items[i]
		item.Process = false
		item.Custom = true
		if customData.HasEntry(item.AircraftCfgFile) {
			customData.AddOrChangeEntry(
				customData.GetEntry(item.AircraftCfgFile).AircraftCfgFile,
				customData.GetEntry(item.AircraftCfgFile).Process,
				customData.GetEntry(item.AircraftCfgFile).OriginalIcao,
				customData.GetEntry(item.AircraftCfgFile).CustomIcao)
		} else {
			customData.AddOrChangeEntry(
				item.AircraftCfgFile,
				item.Process,
				item.Icao,
				"")
		}
	}
	config.Configuration.UpdateIniCustomData()
	model.handleUpdate()
}

func OnItemActivatedAction() {
	customData := config.Configuration.Custom
	for _, i := range LiveryTableView.SelectedIndexes() {
		item := model.items[i]
		if item.Complete {
			item.Process = true
			item.Custom = true
			if customData.HasEntry(item.AircraftCfgFile) {
				customData.AddOrChangeEntry(
					customData.GetEntry(item.AircraftCfgFile).AircraftCfgFile,
					item.Process,
					customData.GetEntry(item.AircraftCfgFile).OriginalIcao,
					customData.GetEntry(item.AircraftCfgFile).CustomIcao)
			} else {
				customData.AddOrChangeEntry(
					item.AircraftCfgFile,
					item.Process,
					item.Icao,
					"")
			}
			config.Configuration.UpdateIniCustomData()
		}
	}
	model.handleUpdate()
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

func (m *LiveryModel) ScanLiveriesAction() {
	ScanButton.SetEnabled(false)
	LiveryTableView.SetEnabled(false)
	StatusBar1.SetText(fmt.Sprintf("Scanning %s ...", config.Configuration.Ini.Section("paths").Key("liveryDir").Value()))
	StatusBar5.SetText(fmt.Sprint(""))
	go m.scanLiveries()
}

func (m *LiveryModel) scanLiveries() {
	liveries, err := livery.ScanLiveryFolder(config.Configuration.Ini.Section("paths").Key("liveryDir").Value())
	if err != nil {
		// TODO: error handling
		m.handleUpdate()
		return
	}
	m.items = liveries
	m.Sort(m.sortColumn, m.sortOrder)
	m.handleUpdate()
}

// called every time when there is a change in the list of liveries
func (m *LiveryModel) handleUpdate() {
	ScanButton.SetEnabled(false)
	LiveryTableView.SetEnabled(false)
	m.PublishRowsReset()
	StatusBar1.SetText(fmt.Sprintf("Number of liveries found: %d", m.RowCount()))
	StatusBar2.SetText(fmt.Sprintf("Number of liveries queued: %d", m.QueuedCount()))
	rules.CalculateRules(m.items)
	StatusBar3.SetText(fmt.Sprintf("Generating %d mappings...", rules.Counter))
	StatusBar4.SetText(fmt.Sprint("Generating XML lines..."))
	go m.buildXML()
}

// builds the actual XML from the calculated rules
// not sure this belongs to this view class?
func (m *LiveryModel) buildXML() {
	RulesText.SetText("")
	numberOfLines := 0

	var output strings.Builder
	output.Grow(100_000)

	// Header
	output.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n")
	output.WriteString("<ModelMatchRuleSet>\r\n\r\n")

	// default rules
	fmt.Fprintf(&output, "<!-- DEFAULTS -->\r\n")
	for _, icaoKey := range rules.SortIcaoKeys(rules.Rules) {
		if icaoKey != "default" {
			continue
		}
		for _, baseKey := range rules.SortBaseKeys(rules.TypeVariations) {
			fmt.Fprintf(&output, "<!-- BASE: %s -->\r\n", baseKey)
			for _, typeKey := range rules.TypeVariations[baseKey] {
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
				numberOfLines++
			}
		}
		fmt.Fprintf(&output, "\r\n")
	}

	// Footer
	output.WriteString("\r\n</ModelMatchRuleSet>\r\n")
	RulesText.SetText(output.String())
	StatusBar3.SetText(fmt.Sprintf("Generated %d mappings.", rules.Counter))
	StatusBar4.SetText(fmt.Sprintf("Generated %d rule lines.", numberOfLines))
	LiveryTableView.SetEnabled(true)
	ScanButton.SetEnabled(true)
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
