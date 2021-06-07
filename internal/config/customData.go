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

// Package config customData represents a data structure which allows to overwrite
// livery data to either complete the data (often ICAO airline code is missing) or
// to disable liveries for the matching rules calculation
package config

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Entry represents custom data for one particular livery
type Entry struct {
	AircraftCfgFile string
	Process         bool
	CustomIcao      string
	OriginalIcao    string
}

// CustomData holds a map of all entries mapped against their aircraft.cfg file-path
type CustomData struct {
	data map[string]*Entry
}

// newCustomData creates an instance of CustomData from a given string holding the
// custom-data data structure
// The custom-data data structure is a 1 or more lines containing a ,-separated list of
// <aircraft.cfg-Filepath>,<process [true|false]>,<original icao code>, <custom icao code>
func newCustomData(body string) *CustomData {
	newData := &CustomData{}
	newData.data = map[string]*Entry{}

	// read in the customData body block and split at newline.
	// then split each line via separator "," and add it to the data structure
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		tokens := strings.Split(line, ",")
		if len(tokens) < 4 {
			fmt.Sprintf("Custom data entry invalid: %s", line)
			continue
		}
		process, err := strconv.ParseBool(tokens[1])
		if err != nil {
			process = false
		}
		entry := &Entry{
			AircraftCfgFile: strings.TrimSpace(tokens[0]),
			Process:         process,
			OriginalIcao:    strings.TrimSpace(tokens[2]),
			CustomIcao:      strings.TrimSpace(tokens[3]),
		}
		newData.data[entry.AircraftCfgFile] = entry
	}

	return newData
}

// NumberOfEntries  returns the number of custom data entries
func (d CustomData) NumberOfEntries() int {
	return len(d.data)
}

// AddOrChangeEntry adds a new entry or changes an existing entry to the custom-data data structure.
// When entry exists overwrites the complete entry with the given parameters.
// Updates the ini data structure (section "customData"
func (d CustomData) AddOrChangeEntry(aircraftCfg string, process bool, originalIcao string, customIcao string) {
	if aircraftCfg == "" {
		return
	}
	d.data[aircraftCfg] = &Entry{
		AircraftCfgFile: aircraftCfg,
		Process:         process,
		CustomIcao:      customIcao,
		OriginalIcao:    originalIcao,
	}
	Configuration.UpdateIniCustomData()
}

// SetProcessFlag creates a new custom-data entry or sets the process flag of an existing entry
func (d CustomData) SetProcessFlag(aircraftCfg string, process bool, icao string) {
	if d.HasEntry(aircraftCfg) {
		d.GetEntry(aircraftCfg).Process = process
	} else {
		d.AddOrChangeEntry(
			aircraftCfg,
			process,
			icao,
			"")
	}
	Configuration.UpdateIniCustomData()
}

// RemoveEntry removes an entry from the custom-data data structure.
// Returns error if entry not found.
// Updates the ini data structure (section "customData"
func (d CustomData) RemoveEntry(aircraftCfg string) error {
	if _, ok := d.data[aircraftCfg]; !ok {
		return errors.New("entry not found")
	}
	delete(d.data, aircraftCfg)
	Configuration.UpdateIniCustomData()
	return nil
}

// HasEntry checks if an entry exists.
func (d CustomData) HasEntry(aircraftCfg string) bool {
	if _, ok := d.data[aircraftCfg]; !ok {
		return false
	}
	return true
}

// GetEntry returns a pointer to the entry if found, nil otherwise.
func (d CustomData) GetEntry(aircraftCfg string) *Entry {
	if _, ok := d.data[aircraftCfg]; !ok {
		return nil
	}
	return d.data[aircraftCfg]
}

// GetDataBody returns a string with all custom-data as it would be stored into the ini file.
func (d CustomData) GetDataBody() string {
	body := strings.Builder{}

	for _, key := range sortKeys(d.data) {
		fmt.Fprintf(&body, "%s,%t,%s,%s\r\n", key, d.data[key].Process, d.data[key].OriginalIcao, d.data[key].CustomIcao)
	}
	return body.String()
}

// sortKeys as go does not support a sorted map iteration we use a slice of all keys and sort it
// Could probably be done with generics
func sortKeys(m map[string]*Entry) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
