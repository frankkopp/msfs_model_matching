package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/karrick/godirwalk"
)

var (
	// TODO: make this commandline args or base config
	dirname = "D:\\Games\\MSFS2020\\Community"
	// dirname = "D:\\Games\\MSFS2020\\ADD_ONS\\Liveries"
	defaultTypesFile   = "D:\\_DEV\\msfs_model_matching\\config\\default_types.txt"
	typeVariationsFile = "D:\\_DEV\\msfs_model_matching\\config\\typeVariations.txt"
	icaoVariationsFile = "D:\\_DEV\\msfs_model_matching\\config\\icaoVariations.txt"

	coreRules = map[string]map[string][]string{}
)

func main() {

	// build config data structure
	defaultTypes := readDefaultTypes(defaultTypesFile)
	typeVariations := readTypeVariations(typeVariationsFile)
	icaoVariations := readIcaoVariations(icaoVariationsFile)

	// find all aircraft.cfg files
	aircraftCfgPaths := findAllAircraftCfg(dirname)
	fmt.Printf("%d entries\n", len(aircraftCfgPaths))

	// process all found aircraft.cfg file and fill the coreRules data structure
	for _, item := range aircraftCfgPaths {
		processAircraftCfg(item, icaoVariations)
	}

	// output of the XML rules
	var output strings.Builder
	output.Grow(100_000)

	output.WriteString(xml.Header)
	output.WriteString("<ModelMatchRuleSet>\n")

	// DEFAULTS
	fmt.Fprintf(&output, "\n")
	fmt.Fprintf(&output, "<!-- DEFAULTS -->\n")
	// MODEL
	sortedBaseKeys := sortBaseKeys(defaultTypes)
	for _, baseKey := range sortedBaseKeys {
		fmt.Fprintf(&output, "<!-- BASE TYPE: %s -->\n", baseKey)
		if len(defaultTypes[baseKey]) == 0 || defaultTypes[baseKey][0] == "" {
			fmt.Fprintf(&output, "<!-- NO DEFAULTS -->\n")
			continue
		}
		for _, typeKey := range typeVariations[baseKey] {
			fmt.Fprintf(&output, "<ModelMatchRule TypeCode=\"%s\" ModelName=\"", typeKey)
			// VARIATIONS
			for i, n := range defaultTypes[baseKey] {
				if i != 0 {
					fmt.Fprint(&output, "//")
				}
				fmt.Fprintf(&output, "%s", n)
			}
			fmt.Fprintf(&output, "\" />\n")
		}
	}

	// ICAO
	sortedIcaoKeys := sortIcaoKeys(coreRules)
	for _, icaoKey := range sortedIcaoKeys {
		fmt.Fprintf(&output, "\n")
		fmt.Fprintf(&output, "<!-- AIRLINE ICAO: %s -->\n", icaoKey)
		// MODEL
		sortedBaseKeys := sortBaseKeys(coreRules[icaoKey])
		for _, baseKey := range sortedBaseKeys {
			fmt.Fprintf(&output, "<!-- BASE TYPE: %s -->\n", baseKey)
			for _, typeKey := range typeVariations[baseKey] {
				fmt.Fprintf(&output, "<ModelMatchRule CallsignPrefix=\"%s\" TypeCode=\"%s\" ModelName=\"", icaoKey, typeKey)
				// VARIATIONS
				for i, n := range coreRules[icaoKey][baseKey] {
					if i != 0 {
						fmt.Fprint(&output, "//")
					}
					fmt.Fprintf(&output, "%s", n)
				}
				fmt.Fprintf(&output, "\" />\n")
			}
		}
	}

	output.WriteString("\n</ModelMatchRuleSet>\n")

	// TODO: Write configurable file
	fmt.Println(output.String())

}

// TODO clean redundant code
func readIcaoVariations(file string) map[string][]string {
	icaoVariations := map[string][]string{}
	lines, _ := readFile(file)
	for _, line := range *lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		tokens := strings.Split(line, ":")
		modelType := strings.TrimSpace(tokens[0])
		for _, t := range tokens[1:] {
			t = strings.TrimSpace(t)
			icaoVariations[modelType] = append(icaoVariations[modelType], t)
		}
	}
	return icaoVariations
}

// TODO clean redundant code
func readDefaultTypes(file string) map[string][]string {
	defaultTypes := map[string][]string{}
	lines, _ := readFile(file)
	for _, line := range *lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		tokens := strings.Split(line, ":")
		modelType := strings.TrimSpace(tokens[0])
		for _, t := range tokens[1:] {
			t = strings.TrimSpace(t)
			defaultTypes[modelType] = append(defaultTypes[modelType], t)
		}
	}
	return defaultTypes
}

// TODO clean redundant code
func readTypeVariations(file string) map[string][]string {
	typeVariations := map[string][]string{}
	lines, _ := readFile(file)
	for _, line := range *lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		tokens := strings.Split(line, ":")
		modelType := strings.TrimSpace(tokens[0])
		for _, t := range tokens[1:] {
			t = strings.TrimSpace(t)
			typeVariations[modelType] = append(typeVariations[modelType], t)
		}
	}
	return typeVariations
}

// as go does not support a sorted map iteration we use a slice of all keys and sort it
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

// get all aircraft.cfg files with in the start directory
func findAllAircraftCfg(filePath string) []string {
	var aircraftCfgPaths []string
	err := godirwalk.Walk(filePath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsRegular() {
				if de.Name() != "aircraft.cfg" {
					return godirwalk.SkipThis
				}
				aircraftCfgPaths = append(aircraftCfgPaths, osPathname)
			}
			return nil
		},
		Unsorted:            false, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
		FollowSymbolicLinks: true,
	})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	return aircraftCfgPaths
}

// parse the file and try to find the three relevant data points:
// ICAO = airline code
// base = base plane model
// name = title of the variation
func processAircraftCfg(path string, icaoVariations map[string][]string) {

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
		if name == "" && strings.Contains(line, "title") {
			name = getValue(line)
			continue
		}
	}

	// if we have all value in the file we save into data structure - otherwise this file will be ignored
	if base != "" && icao != "" && name != "" {

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
var regexValue = regexp.MustCompile(`"\.?\.?\\?(.*)"`)

func getValue(line string) string {
	match := regexValue.FindStringSubmatch(line)
	if len(match) <= 1 {
		return ""
	}
	return match[1]
}

// reads a complete file into a slice of strings
func readFile(path string) (*[]string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("File \"%s\" could not be read; %s\n", path, err)
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")
	return &lines, nil
}
