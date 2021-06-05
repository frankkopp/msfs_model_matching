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

// Package ui implements a user interface based on the walk library (https://github.com/lxn/walk)
package ui

import (
	"fmt"
	"strconv"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	mainWindow   *walk.MainWindow
	tabBarWidget *walk.TabWidget
)

func NewMainWindow() (*walk.MainWindow, error) {

	parseTabPage := parseTab()
	rulesTabPage := rulesTab()
	configTabPage := configTab()

	err := MainWindow{
		AssignTo: &mainWindow,
		Title:    "vPilot MatchMaker " + config.Configuration.Version,
		// MenuItems: NewMenuItem(),
		// ToolBar: toolbar(),
		Layout: VBox{},
		Children: []Widget{
			TabWidget{
				AssignTo: &tabBarWidget,
				Pages: []TabPage{
					parseTabPage,
					rulesTabPage,
					configTabPage,
				},
				// update data when viw (tab) is changed
				OnCurrentIndexChanged: func() {
					switch tabBarWidget.CurrentIndex() {
					case 0:
						scanButton.SetText(fmt.Sprintf("Scan: %s", config.Configuration.Ini.Section("paths").Key("liveryDir").Value()))
					case 1:
						// ignore
					case 2:
						LoadToView()
						checkConfigChange()
					}
				},
			},
		},
		StatusBarItems: statusbar(),
	}.Create()
	if err != nil {
		return nil, err
	}

	// restore previous window state from ini
	x, err := strconv.Atoi(config.Configuration.Ini.Section("application").Key("PosX").Value())
	if err != nil {
		x = 20
	}
	y, err := strconv.Atoi(config.Configuration.Ini.Section("application").Key("PosY").Value())
	if err != nil {
		y = 20
	}
	w, _ := strconv.Atoi(config.Configuration.Ini.Section("application").Key("Width").Value())
	if err != nil {
		w = 1400
	}
	h, _ := strconv.Atoi(config.Configuration.Ini.Section("application").Key("Height").Value())
	if err != nil {
		h = 600
	}
	mainWindow.SetX(x)
	mainWindow.SetY(y)
	mainWindow.SetWidth(w)
	mainWindow.SetHeight(h)

	// store window state to ini when closing window
	mainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		fmt.Printf("Saving window state...\n")
		err := config.Configuration.Ini.Reload()
		if err != nil {
			fmt.Printf("Could not reload ini: %s\n", err)
			return
		}
		config.Configuration.Ini.Section("application").Key("PosX").SetValue(strconv.Itoa(mainWindow.X()))
		config.Configuration.Ini.Section("application").Key("PosY").SetValue(strconv.Itoa(mainWindow.Y()))
		config.Configuration.Ini.Section("application").Key("Width").SetValue(strconv.Itoa(mainWindow.Width()))
		config.Configuration.Ini.Section("application").Key("Height").SetValue(strconv.Itoa(mainWindow.Height()))
		err = config.Configuration.SaveIni()
		if err != nil {
			fmt.Printf("Could not save ini: %s\n", err)
			return
		}
		fmt.Println("Done. Exiting.")
	})

	mainWindow.Run()

	return mainWindow, nil
}
