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
	"github.com/frankkopp/MatchMaker/internal/rules"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func SaveDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var (
		configSave *walk.CheckBox
		rulesSave  *walk.CheckBox
		errorMsg   *walk.TextLabel
	)

	_ = Dialog{
		AssignTo:      &dlg,
		Title:         "Saving configuration and rules",
		MinSize:       Size{Width: 400, Height: 200},
		DefaultButton: &cancelPB,
		CancelButton:  &cancelPB,
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Save configuration:",
					},
					CheckBox{
						AssignTo: &configSave,
						Checked:  config.Configuration.Dirty,
					},
					Label{
						Text: "Save Rules:",
					},
					CheckBox{
						AssignTo: &rulesSave,
						Checked:  rules.Dirty,
					},
					TextLabel{
						AssignTo:      &errorMsg,
						Text:          "",
						TextColor:     walk.RGB(255, 255, 255),
						ColumnSpan:    2,
						TextAlignment: AlignHNearVNear,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "Save",
						OnClicked: func() {
							if configSave.Checked() {
								err := config.Configuration.SaveIni()
								if err != nil {
									errMsg := fmt.Sprintf("Faild to save configuration to %s: %s", *config.Configuration.IniFileName, err)
									StatusBar6.SetText(errMsg)
									brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
									errorMsg.SetBackground(brush)
									// errorMsg.SetTextColor(walk.RGB(255,0,0))
									errorMsg.SetText(errMsg)
									return
								}
							}
							if rulesSave.Checked() {
								err := rules.SaveRulesToFile()
								if err != nil {
									errMsg := fmt.Sprintf("Saving rules to \"%s\" failed: %s", config.Configuration.Ini.Section("paths").Key("outputFile").Value(), err)
									StatusBar5.SetText(errMsg)
									brush, _ := walk.NewSolidColorBrush(walk.RGB(255, 0, 0))
									errorMsg.SetBackground(brush)
									// errorMsg.SetTextColor(walk.RGB(255,0,0))
									errorMsg.SetText(errMsg)
									return
								}
							}
							dlg.Close(walk.DlgCmdYes)
						},
					},
					PushButton{
						Text: "Don't Save",
						OnClicked: func() {
							dlg.Close(walk.DlgCmdNo)
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Close(walk.DlgCmdCancel)
						},
					},
				},
			},
		},
	}.Create(owner)

	return dlg.Run(), nil
}
