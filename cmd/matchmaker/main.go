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

// Package main
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	. "github.com/frankkopp/MatchMaker/internal/config"
	"github.com/frankkopp/MatchMaker/internal/livery"
	"github.com/frankkopp/MatchMaker/internal/rules"
	"github.com/frankkopp/MatchMaker/internal/ui"
	"github.com/frankkopp/MatchMaker/internal/util"
)

const (
	Version = "v1.0"
	IniFile = "matchmaker.ini"
)

func main() {

	// take care of command line argument
	Configuration.IniFileName = flag.String("ini", IniFile, "path to ini file")
	liveryDirectory := flag.String("dir", "", "path where liveries are searched recursively")
	outputFile := flag.String("outputFile", "", "path and filename to output file")
	noUI := flag.Bool("noUI", false, "does not use ui and starts directly with given configuration")
	Configuration.Verbose = flag.Bool("verbose", false, "prints additional information to console")
	versionInfo := flag.Bool("version", false, "prints version and exits")

	flag.Parse()

	// print version info and exit
	if *versionInfo {
		printVersionInfo()
		return
	}

	// store version in configuration
	Configuration.Version = Version

	// load ini file - it loads a default configuration if the ini file is not provided
	Configuration.LoadIni()

	// overwrite the ini configuration with command line options
	if *liveryDirectory != "" {
		Configuration.SetLiveryDirectory(*liveryDirectory)
	}
	if *outputFile != "" {
		Configuration.SetOutputFile(*outputFile)
	}

	// Command line processing without any UI
	if *noUI {
		if err := commandLineProcessing(); err != nil {
			log.Print(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// create the main window and open it.
	_, err := ui.NewMainWindow()
	if err != nil {
		log.Fatal("Could not create main window: " + err.Error())
	}
}

func commandLineProcessing() error {
	fmt.Printf("vPilot MatchMaker by Frank Kopp %s\n", Version)
	fmt.Println("======================================================================================")

	// Step 1: search for liveries
	fmt.Printf("Searching for liveries in folder %s...\n",
		Configuration.Ini.Section("paths").Key("liveryDir").Value())
	liveries, err := livery.ScanLiveryFolder(
		Configuration.Ini.Section("paths").Key("liveryDir").Value())
	if err != nil {
		return err
	}
	fmt.Printf("Found %d liveries.\n", len(liveries))

	// Step 2: calculate rules
	fmt.Printf("Calculating rules...\n")
	rules.CalculateRules(liveries)
	fmt.Printf("Calculated %d rules.\n", rules.Counter)

	// Step 3: write rules to file as XML
	outputFile := Configuration.Ini.Section("paths").Key("outputFile").Value()
	fmt.Printf("Saving vmr file to %s...\n", outputFile)
	var output = strings.Builder{}
	xml, number := rules.GenerateXML()
	output.WriteString(xml)
	err = util.SaveToFile(outputFile, output)
	if err != nil {
		return err
	}
	fmt.Printf("Rules file written to %s (%d XML rules).\n", outputFile, number)

	fmt.Printf("DONE\n")
	return nil
}

func printVersionInfo() {
	fmt.Printf("MatchMaker %s\n", Version)
	fmt.Println("Environment:")
	fmt.Printf("  Using GO version %s\n", runtime.Version())
	fmt.Printf("  Running %s using %s as a compiler\n", runtime.GOARCH, runtime.Compiler)
	fmt.Printf("  Number of CPU: %d\n", runtime.NumCPU())
	fmt.Printf("  Number of Goroutines: %d\n", runtime.NumGoroutine())
	cwd, _ := os.Getwd()
	fmt.Printf("  Working directory: %s\n", cwd)
}
