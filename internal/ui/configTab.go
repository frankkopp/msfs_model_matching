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
	"bytes"
	"fmt"
	"time"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	configTabPage  *walk.TabPage
	configIniText  *walk.TextEdit
	apply, discard *walk.PushButton
)

func configTab() TabPage {
	return TabPage{
		AssignTo: &configTabPage,
		Title:    "Configuration",
		Layout:   VBox{},
		Children: []Widget{
			TextLabel{
				Text: "The below is a provisional way to edit the configuration. A more structured approach will be implemented in the future.",
			},
			TextLabel{
				Text: "Be very careful as it is not very robust and might destroy your configuration.",
			},
			TextLabel{
				Text: "If the text is red the configuration is not valid (minimum required configuration).",
			},
			TextEdit{
				AssignTo: &configIniText,
				Text:     "",
				ReadOnly: false,
				VScroll:  true,
				Font: Font{
					Family:    "Lucida Sans Typewriter",
					PointSize: 8,
				},
				OnTextChanged: func() {
					onTextChanged()
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					PushButton{
						AssignTo:  &discard,
						Text:      "Discard changes",
						OnClicked: LoadToView,
					},
					PushButton{
						AssignTo:  &apply,
						Text:      "Apply changes",
						OnClicked: applyConfig,
					},
					PushButton{
						Text:      "Save to File",
						OnClicked: saveToFile,
					},
				},
			},
		},
	}
}

func saveToFile() {
	err := config.Configuration.SaveIni()
	if err != nil {
		StatusBar6.SetText(fmt.Sprintf("Faild to save configuration to %s: %s", *config.Configuration.IniFileName, err))
		return
	}
	StatusBar6.SetText(fmt.Sprintf("Configuration saved to %s", *config.Configuration.IniFileName))
}

func applyConfig() {
	err := config.Configuration.LoadFromString(configIniText.Text())
	if err != nil {
		StatusBar6.SetText(fmt.Sprintf("Failed tp ally configuration: %s", err))
		return
	}
	// when applying a new configuration we need to rescan to recalculate the rules
	model.Clear()
	StatusBar1.SetText(fmt.Sprintf("New configuration. Please Rescan"))
	StatusBar6.SetText(fmt.Sprintf("Applied configuration. Please Rescan"))
	apply.SetEnabled(false)
	discard.SetEnabled(false)
	go func() {
		time.Sleep(2 * time.Second)
		checkConfigChange()
	}()
}

// onTextChanged enables the Discard and Apply buttons and checks to update the status bar
func onTextChanged() {
	apply.SetEnabled(true)
	discard.SetEnabled(true)
	checkConfigChange()
}

// checkConfigChange checks if the current configuration has been changed (is dirty) and updates
// status bar accordingly.
func checkConfigChange() {
	if config.Configuration.Dirty {
		StatusBar6.SetText(fmt.Sprint("Configuration not saved yet."))
	} else {
		StatusBar6.SetText(fmt.Sprint("Configuration loaded."))
	}
	if config.Configuration.Valid {
		configIniText.SetTextColor(walk.RGB(0, 0, 0))
	} else {
		configIniText.SetTextColor(walk.RGB(255, 0, 0))
	}

}

// LoadToView loads the current configuration as ini-file text into the text edit view.
// Deactivates the Discard and Apply buttons as they only get activated on manual changes to the text edit .
func LoadToView() {
	var tmp bytes.Buffer
	tmpDirty := config.Configuration.Dirty // save the status of dirty as SetText would mark it dirty
	config.Configuration.Ini.WriteTo(&tmp)
	configIniText.SetText(tmp.String())
	configIniText.SetTextSelection(1, 1)
	config.Configuration.Dirty = tmpDirty // restore dirty flag
	apply.SetEnabled(false)
	discard.SetEnabled(false)
	checkConfigChange()
}
