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

// Package livery contains a data structure and functions for representing liveries meta data as
// read from an "aircraft.cfg" file. It also keep additional data for managing the
// the rules generation process. E.g. skipping, custom icao, etc.
package livery

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/karrick/godirwalk"
	"gopkg.in/ini.v1"
)

// Livery is a data structure for representing liveries meta data as
// read from an "aircraft.cfg" file. It also keep additional data for managing the
// the rules generation process. E.g. skipping, custom icao, etc.
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
	// scan liveries
	var liveries []*Livery
	err := godirwalk.Walk(filePath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				if de.Name() != config.FileName {
					return godirwalk.SkipThis
				}
				liveries = append(liveries, processAircraftCfg(osPathname, config.Configuration.Custom)...)
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
func processAircraftCfg(path string, custom *config.CustomData) []*Livery {

	// this is to catch bad aircraft.cfg files which cause the ini library to throw a panic
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	// load the aircraft.cfg file as if it were a ini file
	cfg, err := ini.InsensitiveLoad(path)
	if err != nil {
		return nil
	}

	isPlane := cfg.Section("GENERAL").Key("Category").String() == "airplane"
	isVariation := cfg.Section("VARIATION").HasKey("base_container")

	// Liveries always have a base container. We skip other files
	if !isPlane && !isVariation {
		if *config.Configuration.Verbose {
			fmt.Printf("Not a plane or livery: %s\n", path)
		}
		return nil
	}

	var baseContainer string
	if isVariation {
		baseContainer = getBaseName(cfg.Section("VARIATION").Key("base_container").Value())
		// fmt.Printf("BASE %s\n", baseContainer)
	} else { // is plane
		baseContainer = filepath.Base(filepath.Dir(path))
		// fmt.Printf("PLANE %s\n", baseContainer)
	}

	var liveries []*Livery
	for index := 0; cfg.Section("FLTSIM." + strconv.Itoa(index)).HasKey("title"); index++ {
		// fmt.Printf("FLTSIM %d: %s\n", index, cfg.Section("FLTSIM."+strconv.Itoa(index)).Key("title").String())

		// generate a unique key from path and index
		variationKey := getVariationKey(path, index)

		// create the Livery instance from the aircraft.cfg data
		livery := NewLivery(variationKey)
		livery.BaseContainer = baseContainer
		livery.Title = cleanUp(cfg.Section("FLTSIM." + strconv.Itoa(index)).Key("title").String())
		livery.Icao = cleanUp(cfg.Section("FLTSIM." + strconv.Itoa(index)).Key("icao_airline").String())
		livery.Complete = livery.Title != "" && livery.Icao != ""
		livery.Process = livery.Complete && config.Configuration.Ini.Section("defaultTypes").HasKey(livery.BaseContainer)

		// check for custom data and overwrite livery data if necessary
		if custom.HasEntry(variationKey) {
			if *config.Configuration.Verbose {
				fmt.Println("Custom data applied for ", variationKey)
			}
			entry := custom.GetEntry(variationKey)
			if entry.CustomIcao != "" {
				livery.Icao = entry.CustomIcao
				livery.Complete = true
				livery.Custom = true
			}
			// to be able to be processed the livery data has to be complete e.g. has an ICAO
			livery.Process = entry.Process && livery.Complete && config.Configuration.Ini.Section("defaultTypes").HasKey(livery.BaseContainer)
		}

		// add to list
		liveries = append(liveries, livery)
	}

	// returns nil if file was not a belonging to a livery or new Livery instance otherwise
	return liveries
}

func getVariationKey(path string, index int) string {
	return path + ":" + strconv.Itoa(index)
}

// some liveries unfortunately have incorrect structures - mostly non closing "
// these issues will be corrected here
func cleanUp(value string) string {
	reg, err := regexp.Compile("[\"]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(value, "")
}

// extract the base container name from the aircraft.cfg base_container value as this is often a relative path
func getBaseName(line string) string {
	filePath := filepath.Base(line)
	if len(filePath) <= 1 {
		return ""
	}
	return filePath
}
