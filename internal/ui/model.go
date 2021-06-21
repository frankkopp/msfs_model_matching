/*
 * MatchMaker - create model matching files for VATSIM vPilot
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
	"github.com/frankkopp/MatchMaker/internal/rules"
	"github.com/lxn/walk"
)

// LiveryModel is a data structure for the ui to use ase a model. It bridges the
// underlying data structures to the model. It also controls a much of the ui
// elements based on data changes.
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
	tabBarWidget.SetEnabled(false)
	scanButton.SetEnabled(false)
	liveryTableView.SetEnabled(false)
	StatusBar1.SetText(fmt.Sprintf("Scanning %s ...", config.Configuration.Ini.Section("paths").Key("liveryDir").Value()))
	StatusBar5.SetText(fmt.Sprint("Rules not copied or saved yet."))
	// use parallel execution to allow the ui to be responsive
	go m.scanLiveries()
}

func (m *LiveryModel) scanLiveries() {
	liveries, err := livery.ScanLiveryFolder(config.Configuration.Ini.Section("paths").Key("liveryDir").Value())
	if err != nil {
		StatusBar1.SetText(fmt.Sprintf("Scanning failed: %s ...", err))
		m.onUpdateList()
		return
	}
	m.items = liveries
	m.Sort(m.sortColumn, m.sortOrder)
	m.onUpdateList()
}

// called every time when there is a change in the list of liveries
func (m *LiveryModel) onUpdateList() {
	tabBarWidget.SetEnabled(false)
	scanButton.SetEnabled(false)
	liveryTableView.SetEnabled(false)
	m.PublishRowsReset()
	StatusBar1.SetText(fmt.Sprintf("Number of liveries found: %d", m.RowCount()))
	StatusBar2.SetText(fmt.Sprintf("Number of liveries queued: %d", m.QueuedCount()))
	rules.CalculateRules(m.items)
	StatusBar3.SetText(fmt.Sprintf("Generating %d mappings...", rules.Counter))
	StatusBar4.SetText(fmt.Sprint("Generating XML lines..."))
	if config.Configuration.Dirty {
		StatusBar6.SetText(fmt.Sprint("Configuration not saved yet."))
		configTabPage.SetName(configTabPage.Name() + " (changed)")
		configTabPage.SizeChanged()
	}
	// use parallel execution to allow the ui to be responsive
	go m.buildXML()
}

// builds the actual XML from the calculated rules
func (m *LiveryModel) buildXML() {
	rulesText.SetText("")

	output, numberOfLines := rules.GenerateXML()

	// show in view
	rulesText.SetText(output.String())
	StatusBar3.SetText(fmt.Sprintf("Generated %d mappings.", rules.Counter))
	StatusBar4.SetText(fmt.Sprintf("Generated %d rule lines.", numberOfLines))
	StatusBar5.SetText(fmt.Sprint("Rules not copied or saved yet."))

	liveryTableView.SetEnabled(true)
	scanButton.SetEnabled(true)
	tabBarWidget.SetEnabled(true)

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
		return item.Custom
	case 2:
		return item.Icao
	case 3:
		return item.Title
	case 4:
		return item.BaseContainer
	case 5:
		return item.AircraftCfgFile
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
			return compare(a.Custom)
		case 2:
			return compare(a.Icao < b.Icao)
		case 3:
			return compare(a.Title < b.Title)
		case 4:
			return compare(a.BaseContainer < b.BaseContainer)
		case 5:
			return compare(a.AircraftCfgFile < b.AircraftCfgFile)
		}
		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

func (m *LiveryModel) Clear() {
	m.items = []*livery.Livery{}
	m.onUpdateList()
}
