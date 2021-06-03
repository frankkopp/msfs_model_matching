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

package rules

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/frankkopp/MatchMaker/internal/config"
	"github.com/frankkopp/MatchMaker/internal/livery"
	"github.com/frankkopp/MatchMaker/internal/util"
)

var (
	Counter = 0

	// Rules map[ICAO][TypeCode][]liveries
	Rules = map[string]map[string][]string{}

	DefaultTypes   map[string][]string
	TypeVariations map[string][]string
	IcaoVariations map[string][]string
)

func CalculateRules(liveries []*livery.Livery) error {
	// clear previous calculation
	Counter = 0
	Rules = map[string]map[string][]string{}

	// read and build config data structures
	DefaultTypes = readConfig(*config.Configuration.DefaultTypesFile)
	TypeVariations = readConfig(*config.Configuration.TypeVariationsFile)
	IcaoVariations = readConfig(*config.Configuration.IcaoVariationsFile)
	// customData := readConfig(*config.Configuration.CustomDataFile)

	// create default rules for each type variation
	// add ICAO and all variations to coreRules
	Rules["default"] = map[string][]string{}
	for baseContainer := range DefaultTypes {
		for _, typeVariation := range TypeVariations[baseContainer] {
			Rules["default"][typeVariation] = append(Rules["default"][typeVariation], DefaultTypes[baseContainer]...)
			Counter++
		}
	}

	// create an entry for each icao based rule
	// <ModelMatchRule CallsignPrefix="DLH" TypeCode="A380" ModelName="Boeing 747-8i Lufthansa" />
	for _, l := range liveries {
		// skip all which are not to process or invalid
		if !l.Process || !l.Complete {
			continue
		}
		// lookup if configuration has this base_container,
		icaoList := findIcaoVariations(l, IcaoVariations)
		// add ICAO and all variations to rules
		for _, icao := range icaoList {
			if _, ok := Rules[icao]; !ok {
				Rules[icao] = map[string][]string{}
			}
			for _, typeVariation := range TypeVariations[l.BaseContainer] {
				// add the livery to the rule
				Rules[icao][typeVariation] = append(Rules[icao][typeVariation], l.Title)
				// this counts each livery for the same icao and type
				Counter++
			}
		}
	}
	return nil
}

// check if the ICAO has alternative ICAOs which should use the same livery
func findIcaoVariations(l *livery.Livery, icaoVariations map[string][]string) []string {
	icaoList := []string{l.Icao}
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
	return icaoList
}

// Reads a config file with the format:
// lines of strings separated by";"
// each line will be mapped in a map with the first entry as map key and the rest entries as list of strings
func readConfig(file string) map[string][]string {
	data := map[string][]string{}

	lines, err := util.ReadFile(file)
	if err != nil {
		fmt.Printf("Could not read config file %s\n", file)
		fmt.Println("Exiting")
		os.Exit(1)
	}

	for _, line := range *lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		tokens := strings.Split(line, ";")
		modelType := strings.TrimSpace(tokens[0])
		for _, t := range tokens[1:] {
			t = strings.TrimSpace(t)
			data[modelType] = append(data[modelType], t)
		}
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

// SortIcaoKeys as go does not support a sorted map iteration we use a slice of all keys and sort it
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
