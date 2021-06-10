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

// Package config takes care of all configuration data and is accessible as "Singleton"
// from the whole application. It contains a the ini-file configuration data and
// the custom-data once it is extracted from the ini structure.
// It handles loading and saving of the configuration and also provides synchronization
// functions to keep custom-data and ini-data in sync.
package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/frankkopp/MatchMaker/internal/util"
	"gopkg.in/ini.v1"
)

const (
	FileName = "aircraft.cfg"
)

var (
	// Configuration is a pseudo Singleton for Config
	Configuration = Config{}
)

type Config struct {
	Version     string
	IniFileName *string
	Ini         *ini.File
	Custom      *CustomData
	Dirty       bool
	Verbose     *bool
}

// LoadIni loads configuration from the configured ini file and applies it
// to the given configuration
func (c *Config) LoadIni() {
	tmpIni, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{"customData"},
	}, *c.IniFileName)
	if err != nil {
		log.Printf("No ini file found. Using default configuration: %v", err)
		tmpIni = loadDefaults()
	}
	// TODO: validate ini configuration
	c.Ini = tmpIni
	c.ExtractCustomDataFromIni()
	c.Dirty = false
}

// LoadFromString loads configuration from the configuration view in the UI and applies it
func (c *Config) LoadFromString(iniString string) error {
	tmpIni, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{"customData"},
	}, []byte(iniString))
	if err != nil {
		return err
	}
	// TODO: validate ini configuration
	c.Ini = tmpIni
	c.ExtractCustomDataFromIni()
	c.Dirty = true
	return nil
}

// SaveIni save the current configuration to the configured ini file
func (c *Config) SaveIni() error {
	if *c.IniFileName == "" {
		return errors.New("no ini file path given")
	}
	err := util.CreateBackup(*c.IniFileName)
	if err != nil {
		return err
	}
	c.UpdateIniCustomData()
	err = c.Ini.SaveTo(*c.IniFileName)
	if err != nil {
		return err
	}
	c.Dirty = false
	return nil
}

// ExtractCustomDataFromIni reads and parses the ini section "customData" and
// stores it in a separate data structure.
// It is important to keep this in sync.
func (c *Config) ExtractCustomDataFromIni() {
	c.Custom = newCustomData(c.Ini.Section("customData").Body())
}

// UpdateIniCustomData writes the separate custom-data data structure in the
// "customData" section of the ini data structure.
// It is important to keep this in sync.
func (c *Config) UpdateIniCustomData() {
	dataBody := c.Custom.GetDataBody()
	c.Ini.Section("customData").SetBody(dataBody + "-- end of customData - do not delete --\r\n")
	c.Dirty = true
}

// SetLiveryDirectory sets the liveryDir value in the paths sections of the ini
func (c *Config) SetLiveryDirectory(s string) {
	c.Ini.Section("paths").Key("liveryDir").SetValue(s)
	c.Dirty = true
}

// SetOutputFile sets the outputFile value in the paths sections of the ini
func (c *Config) SetOutputFile(s string) {
	c.Ini.Section("paths").Key("outputFile").SetValue(s)
	c.Dirty = true
}

// loads default configuration from a hard coded string containing a
// default ini file structure
func loadDefaults() *ini.File {
	tmpIni, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{"customData"},
	}, []byte(defaultIni))
	if err != nil {
		fmt.Printf("Error in default ini: %v", err)
		panic(err)
	}
	return tmpIni
}
