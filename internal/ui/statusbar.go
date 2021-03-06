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
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	StatusBar1 *walk.StatusBarItem
	StatusBar2 *walk.StatusBarItem
	StatusBar3 *walk.StatusBarItem
	StatusBar4 *walk.StatusBarItem
	StatusBar5 *walk.StatusBarItem
	StatusBar6 *walk.StatusBarItem
)

func statusbar() []StatusBarItem {

	return []StatusBarItem{
		StatusBarItem{
			AssignTo: &StatusBar1,
			Text:     "No liveries scanned yet.",
			Width:    240,
		},
		StatusBarItem{
			AssignTo: &StatusBar2,
			Text:     "No liveries queued yet.",
			Width:    180,
		},
		StatusBarItem{
			AssignTo: &StatusBar3,
			Text:     "No mapping calculated yet.",
			Width:    180,
		},
		StatusBarItem{
			AssignTo: &StatusBar4,
			Text:     "No rules generated yet.",
			Width:    180,
		},
		StatusBarItem{
			AssignTo: &StatusBar5,
			Text:     "No rules yet",
			Width:    300,
		},
		StatusBarItem{
			AssignTo: &StatusBar6,
			Text:     "Configuration loaded.",
			Width:    220,
		},
	}
}
