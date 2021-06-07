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

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	parseTabPage    *walk.TabPage
	liveryTableView *walk.TableView
	scanButton      *walk.PushButton

	model = NewLiveryModel()
)

func parseTab() TabPage {

	return TabPage{
		AssignTo: &parseTabPage,
		Title:    "Parse Liveries",
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				AssignTo:  &scanButton,
				Text:      fmt.Sprintf("Scan: %s", config.Configuration.Ini.Section("paths").Key("liveryDir").Value()),
				OnClicked: model.ScanLiveriesAction,
			},
			TableView{
				AssignTo:         &liveryTableView,
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
					{Title: "Custom", Width: 50},
				},
				StyleCell: func(style *walk.CellStyle) {
					item := model.items[style.Row()]
					// style row
					if item.Custom {
						style.TextColor = walk.RGB(0, 130, 40)
					}
					// style individual cell
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
						// mark base containers which are not configured to be mapped
						if !config.Configuration.Ini.Section("defaultTypes").HasKey(item.BaseContainer) {
							style.TextColor = walk.RGB(146, 43, 33)
						}
					case 3:
						// placeholder
					case 4:
						// placeholder
					case 5:
						// placeholder
					}
				},
				Model: model,
				ContextMenuItems: []MenuItem{
					Action{
						Text:        "Edit custom",
						OnTriggered: OnItemEditAction,
					},
					Action{
						Text:        "Remove custom",
						OnTriggered: OnItemRemoveCustomAction,
					},
					Action{
						Text:        "Activate Livery",
						OnTriggered: OnItemActivatedAction,
					},
					Action{
						Text:        "Deactivate Livery",
						OnTriggered: OnItemDeactivatedAction,
					},
				},
				OnItemActivated: func() {
					if model.items[liveryTableView.CurrentIndex()].Process {
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
	if len(liveryTableView.SelectedIndexes()) == 0 {
		return
	}
	// we only allow to edit one item - use first selected and deselect the others
	first := liveryTableView.SelectedIndexes()[0]
	err := liveryTableView.SetSelectedIndexes([]int{first})
	if err != nil {
		return
	}

	// restore original data and remove custom data
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
	model.onUpdateList()
}

func OnItemEditAction() {
	if len(liveryTableView.SelectedIndexes()) == 0 {
		return
	}
	// we only allow to edit one item - use first selected and deselect the others
	first := liveryTableView.SelectedIndexes()[0]
	err := liveryTableView.SetSelectedIndexes([]int{first})
	if err != nil {
		return
	}

	if dlg, err := EditDialog(mainWindow, liveryTableView.SelectedIndexes()); err != nil {
		fmt.Print(err)
	} else if dlg == walk.DlgCmdOK {
		model.onUpdateList()
	}
}

func OnItemDeactivatedAction() {
	customData := config.Configuration.Custom
	// multiple selected items allowed and iterated over
	for _, i := range liveryTableView.SelectedIndexes() {
		item := model.items[i]
		item.Process = false
		item.Custom = true
		// add this entry to custom-data now that is has been altered
		customData.SetProcessFlag(item.AircraftCfgFile, item.Process, item.Icao)
	}
	model.onUpdateList()
}

func OnItemActivatedAction() {
	customData := config.Configuration.Custom
	for _, i := range liveryTableView.SelectedIndexes() {
		item := model.items[i]
		if item.Complete {
			item.Process = true
			item.Custom = true
			customData.SetProcessFlag(item.AircraftCfgFile, item.Process, item.Icao)
		}
	}
	model.onUpdateList()
}
