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

// Package main
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	. "github.com/frankkopp/MatchMaker/internal/config"
	"github.com/frankkopp/MatchMaker/internal/ui"
)

const (
	Version = "dev_v0.3"
	IniFile = "matchmaker.ini"
)

func main() {

	// take care of command line argument
	Configuration.IniFileName = flag.String("ini", IniFile, "path to ini file")
	liveryDirectory := flag.String("dir", "", "path where liveries are searched recursively")
	outputFile := flag.String("outputFile", "", "path and filename to output file")
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

	// create the main window and open it.
	mainWindow := ui.NewMainWindow()
	_, err := mainWindow.Run()
	if err != nil {
		log.Fatal("Could not create main window: " + err.Error())
	}
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
