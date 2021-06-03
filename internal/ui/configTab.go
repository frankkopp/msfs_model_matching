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

	"github.com/frankkopp/MatchMaker/internal/config"
	. "github.com/lxn/walk/declarative"
)

func configTab() TabPage {
	return TabPage{
		Title:  "Configuration",
		Layout: VBox{},
		Children: []Widget{
			Composite{Layout: HBox{},
				Children: []Widget{
					TextLabel{Text: "Livery Folder: ", MinSize: Size{Width: 150, Height: 15}},
					TextLabel{Text: *config.Configuration.LiveryDirectory},
					HSpacer{},
				},
			},
			Composite{Layout: HBox{},
				Children: []Widget{
					TextLabel{Text: "Default Types Config: ", MinSize: Size{Width: 150, Height: 15}},
					TextLabel{Text: *config.Configuration.DefaultTypesFile},
					HSpacer{},
				},
			},
			Composite{Layout: HBox{},
				Children: []Widget{
					TextLabel{Text: "Type Variation Config: ", MinSize: Size{Width: 150, Height: 15}},
					TextLabel{Text: *config.Configuration.TypeVariationsFile},
					HSpacer{},
				},
			},
			Composite{Layout: HBox{},
				Children: []Widget{
					TextLabel{Text: "ICAO Variation Config: ", MinSize: Size{Width: 150, Height: 15}},
					TextLabel{Text: *config.Configuration.IcaoVariationsFile},
					HSpacer{},
				},
			},
			Composite{Layout: HBox{},
				Children: []Widget{
					TextLabel{Text: "Custom Data Config: ", MinSize: Size{Width: 150, Height: 15}},
					TextLabel{Text: *config.Configuration.CustomDataFile},
					HSpacer{},
				},
			},
			Composite{Layout: HBox{},
				Children: []Widget{
					TextLabel{Text: "Output File: ", MinSize: Size{Width: 150, Height: 15}},
					TextLabel{Text: *config.Configuration.OutputFile},
					HSpacer{},
				},
			},
			VSpacer{},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					PushButton{
						Text: "Load",
						OnClicked: func() {

						},
					},
					PushButton{
						Text: "Save",
						OnClicked: func() {
							// TODO: catch error
							config.Configuration.SaveIni()
							StatusBar6.SetText(fmt.Sprintf("Configuration saved to %s", config.Configuration.IniFileName))
						},
					},
				},
			},
		},
	}
}
