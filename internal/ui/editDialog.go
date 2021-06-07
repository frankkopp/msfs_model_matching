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
	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func EditDialog(owner walk.Form, indexes []int) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	if len(indexes) == 0 {
		return 0, nil
	}

	// initially we only take the first selected item
	item := model.items[indexes[0]]

	var (
		processCheck *walk.CheckBox
		customIcao   *walk.LineEdit
	)

	_ = Dialog{
		AssignTo:      &dlg,
		Title:         "Edit custom rule",
		MinSize:       Size{Width: 800, Height: 300},
		DefaultButton: &cancelPB,
		CancelButton:  &cancelPB,
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Process:",
					},
					CheckBox{
						AssignTo: &processCheck,
						Checked:  item.Process,
					},
					Label{
						Text: "Livery:",
					},
					LineEdit{
						Text:     item.AircraftCfgFile,
						ReadOnly: true,
					},
					Label{
						Text: "Base Container:",
					},
					LineEdit{
						Text:     item.BaseContainer,
						ReadOnly: true,
					},
					Label{
						Text: "Title:",
					},
					LineEdit{
						Text:     item.Title,
						ReadOnly: true,
					},
					Label{
						Text: "ICAO:",
					},
					LineEdit{
						AssignTo: &customIcao,
						Text:     item.Icao,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if processCheck.Checked() == item.Process && customIcao.Text() == item.Icao {
								// no changes
								return
							}

							if customIcao.Text() != "" {
								item.Custom = true
								config.Configuration.Custom.AddOrChangeEntry(item.AircraftCfgFile, processCheck.Checked(), item.Icao, customIcao.Text())
								item.Complete = true
								item.Process = processCheck.Checked()
								item.Icao = customIcao.Text()
							} else {
								item.Complete = false
								item.Process = false
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Create(owner)

	customIcao.SetFocus()
	return dlg.Run(), nil
}
