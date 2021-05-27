//
// VATSIM vPilot MatchMaker
//
// MIT License
//
// Copyright (c) 2021 Frank Kopp
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

// Package main
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/karrick/godirwalk"
)

const (
	fileName = "aircraft.cfg"
)

var (
	version    = "0.2a"
	coreRules  = map[string]map[string][]string{}
	showCustom = false
	counter    = 0
)

func main() {

	// take care of command line argument
	versionInfo := flag.Bool("version", false, "prints version and exits")
	showCustomFlag := flag.Bool("showCustom", false, "shows liveries which are custom by configuration")
	liveryDirectory := flag.String("dir", ".", "path where "+fileName+" are searched recursively")
	defaultTypesFile := flag.String("defaultTypesFile", "..\\config\\defaultTypes.txt", "path and filename to default types config file")
	typeVariationsFile := flag.String("typeVariationsFile", "..\\config\\typeVariations.txt", "path and filename to type variations config file")
	icaoVariationsFile := flag.String("icaoVariationsFile", "..\\config\\icaoVariations.txt", "path and filename to icao variations config file")
	fixLiveriesFile := flag.String("fixLiveriesFile", "..\\config\\customData.txt", "path and filename to fix liveries config file")
	outPutFile := flag.String("outPutFile", ".\\MatchMakingRules.vmr", "path and filename to output file")

	flag.Parse()

	showCustom = *showCustomFlag

	// print version info and exit
	if *versionInfo {
		printVersionInfo()
		return
	}

	startReading := time.Now()

	// build config data structures
	fmt.Printf("Reading configuration\n")
	defaultTypes := readConfig(*defaultTypesFile)
	typeVariations := readConfig(*typeVariationsFile)
	icaoVariations := readConfig(*icaoVariationsFile)
	fixLiveries := readConfig(*fixLiveriesFile)

	// find all aircraft.cfg files
	fmt.Printf("Searching for %s in %s\n", fileName, *liveryDirectory)
	aircraftCfgPaths := findAllAircraftCfg(*liveryDirectory)
	fmt.Printf("Found %d %s files in %d ms\n", len(aircraftCfgPaths), fileName, time.Since(startReading).Milliseconds())

	// process all found aircraft.cfg file and fill the coreRules data structure
	fmt.Printf("Processing all %s\n", fileName)
	for _, item := range aircraftCfgPaths {
		processAircraftCfg(item, icaoVariations, fixLiveries)
	}

	// create output of the XML rules
	fmt.Printf("Creating XML rules\n")
	output := createXMLRules(defaultTypes, typeVariations)

	// Save output inti vmr file
	fmt.Printf("Saving VMR file to %s\n", *outPutFile)
	saveToFile(outPutFile, output)

	elapsedReading := time.Since(startReading)
	fmt.Printf("Created %d rules in %d ms\n", counter, elapsedReading.Milliseconds())
}

// use the data structures to create the XML output and rules
func createXMLRules(defaultTypes map[string][]string, typeVariations map[string][]string) strings.Builder {
	var output strings.Builder
	output.Grow(100_000)

	// Header
	output.WriteString(xml.Header)
	output.WriteString("<ModelMatchRuleSet>\n\n")

	// Default rules
	fmt.Fprintf(&output, "<!-- DEFAULTS -->\n")

	// Create a section for each base plane model we have available
	sortedBaseKeys := sortBaseKeys(defaultTypes)
	for _, baseKey := range sortedBaseKeys {
		fmt.Fprintf(&output, "<!-- BASE TYPE: %s -->\n", baseKey)

		// only continue if a default is defined in config
		if len(defaultTypes[baseKey]) == 0 || defaultTypes[baseKey][0] == "" {
			fmt.Fprintf(&output, "<!-- NO DEFAULTS -->\n")
			continue
		}

		// create a rule for each plane type variation
		for _, typeKey := range typeVariations[baseKey] {
			fmt.Fprintf(&output, "<ModelMatchRule TypeCode=\"%s\" ModelName=\"", typeKey)
			// define default liveries for this plane type - might be multiple liveries - vPilot will choose randomly
			for i, n := range defaultTypes[baseKey] {
				if i != 0 {
					fmt.Fprint(&output, "//")
				}
				fmt.Fprintf(&output, "%s", n)
			}
			fmt.Fprintf(&output, "\" />\n")
			counter++
		}
	}

	// Create a section for each airline ICAO we have found in the aircraft.cfg files
	sortedIcaoKeys := sortIcaoKeys(coreRules)
	for _, icaoKey := range sortedIcaoKeys {
		fmt.Fprintf(&output, "\n")
		fmt.Fprintf(&output, "<!-- AIRLINE ICAO: %s -->\n", icaoKey)
		// Create a section for each base plane model we have available
		sortedBaseKeys := sortBaseKeys(coreRules[icaoKey])
		for _, baseKey := range sortedBaseKeys {
			fmt.Fprintf(&output, "<!-- BASE TYPE: %s -->\n", baseKey)
			// create a rule for each plane type variation
			for _, typeKey := range typeVariations[baseKey] {
				fmt.Fprintf(&output, "<ModelMatchRule CallsignPrefix=\"%s\" TypeCode=\"%s\" ModelName=\"", icaoKey, typeKey)
				// define liveries for this airline and plane type - might be multiple liveries - vPilot will choose randomly
				for i, n := range coreRules[icaoKey][baseKey] {
					if i != 0 {
						fmt.Fprint(&output, "//")
					}
					fmt.Fprintf(&output, "%s", n)
				}
				fmt.Fprintf(&output, "\" />\n")
				counter++
			}
		}
	}

	// Footer
	output.WriteString("\n</ModelMatchRuleSet>\n")

	return output
}

// Reads a config file with the format:
// lines of strings seperated by":"
// each line will be mapped in a map with the first entry as map key and the rest entries as list of strings
func readConfig(file string) map[string][]string {
	data := map[string][]string{}
	lines, err := readFile(file)

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

// as go does not support a sorted map iteration we use a slice of all keys and sort it
// Could probably be done with generics
func sortBaseKeys(m map[string][]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// as go does not support a sorted map iteration we use a slice of all keys and sort it
// Could probably be done with generics
func sortIcaoKeys(m map[string]map[string][]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// find all aircraft.cfg files as paths with in the start directory
func findAllAircraftCfg(filePath string) []string {
	var aircraftCfgPaths []string
	err := godirwalk.Walk(filePath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {

				if de.Name() != fileName {
					return godirwalk.SkipThis
				}
				aircraftCfgPaths = append(aircraftCfgPaths, osPathname)
			}
			return nil
		},
		Unsorted:            true,
		FollowSymbolicLinks: true,
	})
	if err != nil {
		fmt.Printf("Error during search for %s files: %v\n", fileName, err)
	}
	return aircraftCfgPaths
}

// parse the file and try to find the three relevant data points:
// ICAO = airline code
// base = base plane model
// name = title of the variation
func processAircraftCfg(path string, icaoVariations map[string][]string, customLiveries map[string][]string) {

	lines, _ := readFile(path)

	var base, icao, name string
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
		// TODO: could be optimized to look for AI variations as they are most likely less resource hungry
		//  for model matching
		if name == "" && strings.Contains(line, "title") {
			name = getValue(line)
			continue
		}
	}

	// handle custom and faulty liveries
	if _, ok := customLiveries[path]; ok {
		name = customLiveries[path][0]
		base = customLiveries[path][1]
		icao = customLiveries[path][2]
		if showCustom {
			fmt.Printf("CUSTOM : %s;%s;%s;%s\n", path, name, base, icao)
		}
		if name == "skip" {
			return
		}
	}

	// if we have all value in the file we save into data structure - otherwise this file will be skipped
	if !(base != "" && icao != "" && name != "") {
		fmt.Printf("SKIPPED: %s;%s;%s;%s\n", path, name, base, icao)
	} else {
		// check if the ICAO has alternative ICAOs which should use the same livery
		icaoList := []string{icao}
		for k := range icaoVariations {
			// do we match the main ICAO?
			if k == icao {
				icaoList = append(icaoList, icaoVariations[icao]...)
				break
			}
			// is this ICAO part of the variations?
			for _, v := range icaoVariations[icao] {
				if v == icao {
					icaoList = append(icaoList, icaoVariations[icao]...)
					break
				}
			}
		}

		// add ICAO and all variations to coreRules
		for _, i := range icaoList {
			if _, ok := coreRules[i]; !ok {
				coreRules[i] = map[string][]string{}
			}
			coreRules[i][base] = append(coreRules[i][base], name)
		}
	}
}

// extract the value from the aircraft.cfg line
var regexValue = regexp.MustCompile(`^.*\s*=\s*"?\.*\\*(.*?)"?\s*(;.*)?$`)

// extract the value from the aircraft.cfg line
func getValue(line string) string {
	match := regexValue.FindStringSubmatch(line)
	if len(match) <= 1 {
		return ""
	}
	return match[1]
}

// save a strings.Builder instance to a file
func saveToFile(outPutFile *string, output strings.Builder) {
	outFile, err := os.Create(*outPutFile)
	if err != nil {
		fmt.Println("Error while creating VMR file to %S\n", *outPutFile)
	}
	_, err = outFile.WriteString(output.String())
	if err != nil {
		fmt.Println("Error while saving to VMR file to %S\n", *outPutFile)
	}
	err = outFile.Close()
	if err != nil {
		fmt.Println("Error while closing VMR file to %S\n", *outPutFile)
	}
}

// reads a complete file into a slice of strings
func readFile(path string) (*[]string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")
	return &lines, nil
}

func printVersionInfo() {
	fmt.Printf("MatchMaker %s\n", version)
	fmt.Println("Environment:")
	fmt.Printf("  Using GO version %s\n", runtime.Version())
	fmt.Printf("  Running %s using %s as a compiler\n", runtime.GOARCH, runtime.Compiler)
	fmt.Printf("  Number of CPU: %d\n", runtime.NumCPU())
	fmt.Printf("  Number of Goroutines: %d\n", runtime.NumGoroutine())
	cwd, _ := os.Getwd()
	fmt.Printf("  Working directory: %s\n", cwd)
}
