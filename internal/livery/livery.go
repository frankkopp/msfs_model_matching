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
	"regexp"
	"strings"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/frankkopp/MatchMaker/internal/util"
	"github.com/karrick/godirwalk"
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
	var liveries []*Livery
	err := godirwalk.Walk(filePath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				if de.Name() != config.FileName {
					return godirwalk.SkipThis
				}
				livery := processAircraftCfg(osPathname)
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
// ICAO = airline code
// base = base plane model
// name = title of the variation
// returns nil if file was not a belonging to a livery or new Livery instance otherwise
func processAircraftCfg(path string) *Livery {

	var livery *Livery

	lines, _ := util.ReadFile(path)

	var base, icao, title string
	for _, line := range *lines {
		if strings.Contains(line, "base_container") {
			base = getValue(line)
			continue
		}
		if strings.Contains(line, "icao_airline") {
			icao = strings.ToUpper(getValue(line))
			continue
		}
		// only take the first title - ignore the rest
		// TODO: could be optimized to look for AI variations as they
		//  are most likely less resource hungry for model matching
		if title == "" && strings.Contains(line, "title") {
			title = getValue(line)
			continue
		}
	}

	// Liveries always have a base container. We skip other files
	if base != "" {
		livery = NewLivery(path)
		(*livery).BaseContainer = base
		(*livery).Title = title
		(*livery).Icao = icao
		if title != "" && icao != "" {
			(*livery).Process = true
			(*livery).Complete = true
		}
	}

	// returns nil if file was not a belonging to a livery or new Livery instance otherwise
	return livery
}

// extract the value from the aircraft.cfg line
var regexValue = regexp.MustCompile(`^.*\s*=\s*"?[.\\/]*(.*?)"?\s*(;.*)?$`)

// extract the value from the aircraft.cfg line
func getValue(line string) string {
	match := regexValue.FindStringSubmatch(line)
	if len(match) <= 1 {
		return ""
	}
	return match[1]
}
