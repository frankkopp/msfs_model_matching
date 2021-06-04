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
	"strings"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/frankkopp/MatchMaker/internal/util"

	// "github.com/lxn/walk"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	RulesTabHandle *walk.TabPage
	RulesText      *walk.TextEdit
)

func rulesTab() TabPage {
	// monoFont, _ := walk.NewFont("Lucida Sans Typewriter", 9, walk.FontBold)

	return TabPage{
		AssignTo: &RulesTabHandle,
		Title:    "Generated Rules",
		Layout:   VBox{},
		Children: []Widget{
			TextEdit{
				AssignTo: &RulesText,
				Text:     "No rules generated yet.",
				ReadOnly: true,
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
						Text: "Copy All",
						OnClicked: func() {
							walk.Clipboard().SetText(RulesText.Text())
							StatusBar5.SetText(fmt.Sprint("Rules copied to Clipboard."))
							StatusBar5.SetToolTipText(fmt.Sprint("Rules copied to Clipboard."))
						},
					},
					PushButton{
						Text: "Save",
						OnClicked: func() {
							var output = strings.Builder{}
							output.WriteString(RulesText.Text())
							util.SaveToFile(config.Configuration.Ini.Section("paths").Key("outputFile").Value(), output)
							StatusBar5.SetText(fmt.Sprintf("Rules saved to file: %s", config.Configuration.Ini.Section("paths").Key("outputFile").Value()))
						},
					},
				},
			},
		},
	}
}
