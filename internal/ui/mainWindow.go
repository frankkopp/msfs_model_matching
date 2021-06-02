package ui

import (
	"log"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// ///////////////////////////////////////////////////////////
// Public
// ///////////////////////////////////////////////////////////

func NewMainWindow(version string, configuration config.Config) *MainWindow {

	// boldFont, _ := walk.NewFont("Segoe UI", 9, walk.FontBold)
	// goodIcon, _ := walk.Resources.Icon("../img/check.ico")
	// badIcon, _ := walk.Resources.Icon("../img/stop.ico")

	// var openAction *walk.Action

	icon1, err := walk.NewIconFromFile("../img/check.ico")
	if err != nil {
		log.Fatal(err)
	}

	mw := MainWindow{
		Title: "vPilot MatchMaker " + version,

		MenuItems: []MenuItem{
			Menu{
				Text: "&File",
				Items: []MenuItem{
					Action{
						// AssignTo:    &openAction,
						Text: "&Open",
						// Image:       "img/open.png",
						// OnTriggered: mw.openAction_Triggered,
					},
					Separator{},
					Action{
						Text: "Exit",
						// OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&Help",
				Items: []MenuItem{
					Action{
						Text: "About",
						// OnTriggered: mw.aboutAction_Triggered,
					},
				},
			},
		},

		ToolBar: ToolBar{
			ButtonStyle: ToolBarButtonImageBeforeText,
			Items: []MenuItem{
				Action{
					Text: "Special",
					// Image:       "img/system-shutdown.png",
					// Enabled:     Bind("isSpecialMode && enabledCB.Checked"),
					// OnTriggered: mw.specialAction_Triggered,
				},
			},
		},

		Size:   Size{800, 600},
		Layout: VBox{},

		Children: []Widget{
			TabWidget{
				Pages: []TabPage{
					TabPage{
						Title:  "Configuration",
						Layout: VBox{},
						Children: []Widget{
							Composite{Layout: HBox{},
								Children: []Widget{
									TextLabel{Text: "Livery Folder: ", MinSize: Size{Width: 150, Height: 15}},
									TextLabel{Text: configuration.LiveryDirectory},
									HSpacer{},
								},
							},
							Composite{Layout: HBox{},
								Children: []Widget{
									TextLabel{Text: "Default Types Config: ", MinSize: Size{Width: 150, Height: 15}},
									TextLabel{Text: configuration.DefaultTypesFile},
									HSpacer{},
								},
							},
							Composite{Layout: HBox{},
								Children: []Widget{
									TextLabel{Text: "Type Variation Config: ", MinSize: Size{Width: 150, Height: 15}},
									TextLabel{Text: configuration.TypeVariationsFile},
									HSpacer{},
								},
							},
							Composite{Layout: HBox{},
								Children: []Widget{
									TextLabel{Text: "ICAO Variation Config: ", MinSize: Size{Width: 150, Height: 15}},
									TextLabel{Text: configuration.IcaoVariationsFile},
									HSpacer{},
								},
							},
							Composite{Layout: HBox{},
								Children: []Widget{
									TextLabel{Text: "Custom Data Config: ", MinSize: Size{Width: 150, Height: 15}},
									TextLabel{Text: configuration.CustomDataFile},
									HSpacer{},
								},
							},
							Composite{Layout: HBox{},
								Children: []Widget{
									TextLabel{Text: "Output File: ", MinSize: Size{Width: 150, Height: 15}},
									TextLabel{Text: configuration.OutputFile},
									HSpacer{},
								},
							},
							VSpacer{},
						},
					},
					TabPage{
						Title: "TAB2",
					},
				},
			},
		},

		StatusBarItems: []StatusBarItem{
			StatusBarItem{
				// AssignTo: &sbi,
				Icon:  icon1,
				Text:  "click",
				Width: 80,
				// OnClicked: func() {
				// 	if sbi.Text() == "click" {
				// 		sbi.SetText("again")
				// 		sbi.SetIcon(icon2)
				// 	} else {
				// 		sbi.SetText("click")
				// 		sbi.SetIcon(icon1)
				// 	}
				// },
			},
			StatusBarItem{
				Text:        "left",
				ToolTipText: "no tooltip for me",
			},
			StatusBarItem{
				Text: "\tcenter",
			},
			StatusBarItem{
				Text: "\t\tright",
			},
			StatusBarItem{
				Icon:        icon1,
				ToolTipText: "An icon with a tooltip",
			},
		},
	}

	return &mw
}
