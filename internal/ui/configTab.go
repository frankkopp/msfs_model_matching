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
	"bytes"
	"fmt"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	ConfigIniText *walk.TextEdit
)

func configTab() TabPage {
	return TabPage{
		Title:  "Configuration",
		Layout: VBox{},
		Children: []Widget{
			TextLabel{
				Text: "The below is a provisional way to edit the configuration. A more structured approach will be implemented in the future.",
			},
			TextLabel{
				Text: "Be very careful as it is not very robust and might destroy your configuration.",
			},
			// Composite{Layout: HBox{},
			// 	Children: []Widget{
			// 		TextLabel{Text: "Livery Folder: ", MinSize: Size{Width: 150, Height: 15}},
			// 		TextLabel{Text: config.Configuration.Ini.Section("paths").Key("liveryDir").Value()},
			// 		HSpacer{},
			// 	},
			// },
			// Composite{Layout: HBox{},
			// 	Children: []Widget{
			// 		TextLabel{Text: "Output File: ", MinSize: Size{Width: 150, Height: 15}},
			// 		TextLabel{Text: config.Configuration.Ini.Section("paths").Key("outputFile").Value()},
			// 		HSpacer{},
			// 	},
			// },
			TextEdit{
				AssignTo: &ConfigIniText,
				Text:     "",
				ReadOnly: false,
				VScroll:  true,
				Font: Font{
					Family:    "Lucida Sans Typewriter",
					PointSize: 8,
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					PushButton{
						Text: "Load to view",
						OnClicked: func() {
							var tmp bytes.Buffer
							config.Configuration.Ini.WriteTo(&tmp)
							ConfigIniText.SetText(tmp.String())
							StatusBar6.SetText(fmt.Sprintf("Configuration loaded."))
						},
					},
					PushButton{
						Text: "Use",
						OnClicked: func() {
							// TODO: catch error
							config.Configuration.LoadFromView(ConfigIniText.Text())
							StatusBar6.SetText(fmt.Sprintf("Applied configuration."))
						},
					},
					PushButton{
						Text: "Save to File",
						OnClicked: func() {
							// TODO: catch error
							config.Configuration.SaveIni()
							StatusBar6.SetText(fmt.Sprintf("Configuration saved to %s", *config.Configuration.IniFileName))
						},
					},
				},
			},
		},
	}
}
