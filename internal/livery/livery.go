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

package livery

import (
	"fmt"
	"regexp"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/frankkopp/MatchMaker/internal/customData"
	"github.com/karrick/godirwalk"
	"gopkg.in/ini.v1"
)

type Livery struct {
	AircraftCfgFile string
	BaseContainer   string
	Title           string
	Icao            string
	Custom          bool // has custom config
	Process         bool // rules should be created
	Complete        bool // rules should be created
}

// NewLivery creates a new instance of a Livery
func NewLivery(aircraftCfgFile string) *Livery {
	return &Livery{
		AircraftCfgFile: aircraftCfgFile,
	}
}

// ScanLiveryFolder find all aircraft.cfg files as paths with in the start directory
func ScanLiveryFolder(filePath string) ([]*Livery, error) {
	// load custom rules
	custom := customData.NewCustomData(config.Configuration.Ini.Section("customData").Body())

	// scan liveries
	var liveries []*Livery
	err := godirwalk.Walk(filePath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				if de.Name() != config.FileName {
					return godirwalk.SkipThis
				}
				livery := processAircraftCfg(osPathname, custom)
				if livery != nil {
					liveries = append(liveries, livery)
				}
			}
			return nil
		},
		Unsorted:            true,
		FollowSymbolicLinks: true,
	})
	return liveries, err
}

// parse the file and try to find the three relevant data points:
// base = base plane model
// icao = airline code
// name = title of the variation
// returns nil if file was invalid or not a livery aircraft.cfg
func processAircraftCfg(path string, custom *customData.CustomData) *Livery {

	// this is to catch bad *aircraft.cfg files which cause the ini library to throw a panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	cfg, err := ini.Load(path)
	if err != nil {
		return nil
	}

	// Liveries always have a base container. We skip other files
	if !cfg.Section("VARIATION").HasKey("base_container") {
		fmt.Printf("No VARIATION section: %s\n", path)
		return nil
	}

	baseContainer := getValue(cfg.Section("VARIATION").Key("base_container").Value())

	// check if this base container s part of the configuration
	if !config.Configuration.Ini.Section("defaultTypes").HasKey(baseContainer) {
		fmt.Printf("Not part of default types: %s %s\n", baseContainer, path)
		return nil
	}

	livery := NewLivery(path)
	(*livery).BaseContainer = baseContainer
	(*livery).Icao = cfg.Section("FLTSIM.0").Key("icao_airline").Value()
	(*livery).Title = cfg.Section("FLTSIM.0").Key("title").Value()
	if (*livery).Title != "" && (*livery).Icao != "" {
		(*livery).Process = true
		(*livery).Complete = true
	}

	// check for custom data and overwrite livery data if necessary
	if custom.HasEntry(path) {
		fmt.Println("CUSTOM for ", path)
		entry := custom.GetEntry(path)
		if entry.CustomIcao != "" {
			(*livery).Icao = entry.CustomIcao
			(*livery).Complete = true
			(*livery).Custom = true
		}
		(*livery).Process = entry.Process && (*livery).Complete
	}

	// returns nil if file was not a belonging to a livery or new Livery instance otherwise
	return livery
}

// extract the value from the aircraft.cfg line
var regexValue = regexp.MustCompile(`^[.\\/]*(.*?)$`)

// extract the value from the aircraft.cfg line
func getValue(line string) string {
	match := regexValue.FindStringSubmatch(line)
	if len(match) <= 1 {
		return ""
	}
	return match[1]
}
