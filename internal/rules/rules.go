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

// Package rules manages rules as a map of map of a list of liveries
// This gives convenient access to the liveries for each icao and type-code
// while implicitly avoiding doublets
package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/frankkopp/MatchMaker/internal/livery"
	"github.com/frankkopp/MatchMaker/internal/util"
	"gopkg.in/ini.v1"
)

var (

	// Rules map[ICAO][TypeCode][]liveries
	Rules = map[string]map[string][]string{}

	Counter = 0
	Dirty   = false

	DefaultTypes   map[string][]string
	TypeVariations map[string][]string
	IcaoVariations map[string][]string
)

// CalculateRules (re-)calculates the rules based on current configuration and livery data.
// Will be stored in rules.Rules
func CalculateRules(liveries []*livery.Livery) {
	// clear previous calculation
	Counter = 0
	Rules = map[string]map[string][]string{}

	// read and build config data structures
	// do this for every run to use the latest configuration
	DefaultTypes = readConfig(config.Configuration.Ini.Section("defaultTypes"))
	TypeVariations = readConfig(config.Configuration.Ini.Section("typeVariations"))
	IcaoVariations = readConfig(config.Configuration.Ini.Section("icaoVariations"))

	// create default rules for each type variation
	// add ICAO and all variations to coreRules
	// uses "default" as a special icao
	Rules["default"] = map[string][]string{}
	for baseContainer := range DefaultTypes {
		for _, typeVariation := range TypeVariations[baseContainer] {
			Rules["default"][typeVariation] = append(Rules["default"][typeVariation], DefaultTypes[baseContainer]...)
			Counter++
		}
	}

	// create an entry for each icao and type-code
	// <ModelMatchRule CallsignPrefix="DLH" TypeCode="A380" ModelName="Boeing 747-8i Lufthansa" />
	for _, cLivery := range liveries {
		// skip all which are not to process or invalid
		if !(cLivery.Process && cLivery.Complete) {
			continue
		}
		// only process liveries with configured base containers
		// check if this base container is part of the configuration otherwise skip
		if !config.Configuration.Ini.Section("defaultTypes").HasKey(cLivery.BaseContainer) {
			continue
		}
		// get alternate icao variations from config if available
		icaoList := findIcaoVariations(cLivery, IcaoVariations)
		// add icao and all variations to rules
		for _, icao := range icaoList {
			// init map for new entries
			if _, ok := Rules[icao]; !ok {
				Rules[icao] = map[string][]string{}
			}
			for _, typeVariation := range TypeVariations[cLivery.BaseContainer] {
				// add the cLivery to the rule
				Rules[icao][typeVariation] = append(Rules[icao][typeVariation], cLivery.Title)
				// this counts each cLivery for the same icao and type
				Counter++
			}
		}
	}
	Dirty = true
}

// check if the ICAO has alternative ICAOs which should use the same livery
func findIcaoVariations(l *livery.Livery, icaoVariations map[string][]string) []string {
	var icaoList []string
	for k := range icaoVariations {
		// is this ICAO part of the variations?
		for _, v := range icaoVariations[k] {
			if v == l.Icao {
				icaoList = []string{}
				icaoList = append(icaoList, icaoVariations[k]...)
				return icaoList
			}
		}
	}
	return []string{l.Icao}
}

// Reads lines of strings separated by"," from a section of the ini.
// Each line will be mapped in a map with the ini-key as map key and the rest entries as list of strings
func readConfig(section *ini.Section) map[string][]string {
	data := map[string][]string{}
	for _, s := range section.Keys() {
		tokens := s.Strings(",") // it is possible to use the key object directly to get to the value
		data[s.Name()] = append(data[s.Name()], tokens...)
	}
	return data
}

// SortIcaoKeys as go does not support a sorted map iteration we use a slice of all keys and sort it
// Could probably be done with generics
func SortIcaoKeys(m map[string]map[string][]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortBaseKeys as go does not support a sorted map iteration we use a slice of all keys and sort it
// Could probably be done with generics
func SortBaseKeys(m map[string][]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// GenerateXML generates a string with the XML representation of all matching rules.
// Also returns the number of rules generated.
func GenerateXML() (strings.Builder, int) {
	numberOfLines := 0

	var output strings.Builder
	output.Grow(100_000)

	// Header
	output.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n")
	output.WriteString("<ModelMatchRuleSet>\r\n\r\n")

	// default rules
	fmt.Fprintf(&output, "<!-- DEFAULTS -->\r\n")
	for _, icaoKey := range SortIcaoKeys(Rules) {
		if icaoKey != "default" {
			continue
		}
		if len(Rules[icaoKey]) == 0 {
			break
		}
		for _, baseKey := range SortBaseKeys(DefaultTypes) { // only iterate over types with default livery
			fmt.Fprintf(&output, "<!-- BASE: %s -->\r\n", baseKey)
			for _, typeKey := range TypeVariations[baseKey] {
				if len(Rules[icaoKey][typeKey]) == 0 {
					continue
				}
				fmt.Fprintf(&output, "<ModelMatchRule TypeCode=\"%s\" ModelName=\"", typeKey)
				for i, livery := range Rules[icaoKey][typeKey] {
					if i != 0 {
						fmt.Fprint(&output, "//")
					}
					fmt.Fprintf(&output, "%s", livery)
				}
				fmt.Fprintf(&output, "\" />\r\n")
			}
		}
	}
	fmt.Fprintf(&output, "\r\n")

	// ICAO based rules
	fmt.Fprintf(&output, "<!-- PER ICAO RULES -->\r\n")
	for _, icaoKey := range SortIcaoKeys(Rules) {
		if icaoKey == "default" || len(Rules[icaoKey]) == 0 {
			continue
		}
		fmt.Fprintf(&output, "<!-- ICAO: %s -->\r\n", icaoKey)
		for _, baseKey := range SortBaseKeys(DefaultTypes) { // only iterate over types with default livery
			fmt.Fprintf(&output, "<!-- BASE: %s -->\r\n", baseKey)
			for _, typeKey := range TypeVariations[baseKey] {
				if len(Rules[icaoKey][typeKey]) == 0 {
					continue
				}
				fmt.Fprintf(&output, "<ModelMatchRule CallsignPrefix=\"%s\" TypeCode=\"%s\" ModelName=\"", icaoKey, typeKey)
				for i, livery := range Rules[icaoKey][typeKey] {
					if i != 0 {
						fmt.Fprint(&output, "//")
					}
					fmt.Fprintf(&output, "%s", livery)
				}
				fmt.Fprintf(&output, "\" />\r\n")
				numberOfLines++
			}
		}
		fmt.Fprintf(&output, "\r\n")
	}

	// Footer
	output.WriteString("\r\n</ModelMatchRuleSet>\r\n")
	return output, numberOfLines
}

func SaveRulesToFile() error {
	outPutFile := config.Configuration.Ini.Section("paths").Key("outputFile").String()
	err := util.CreateBackup(outPutFile)
	if err != nil {
		return err
	}
	var output, _ = GenerateXML()
	err = util.SaveToFile(outPutFile, output)
	if err != nil {
		return err
	}
	Dirty = false
	return nil
}
