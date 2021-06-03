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

package config

import (
	"errors"

	"gopkg.in/ini.v1"
)

const (
	FileName = "aircraft.cfg"
)

var (
	Configuration = Config{}
)

type Config struct {
	Version     string
	IniFileName string

	// cmd line options
	VersionInfo *bool

	// paths
	LiveryDirectory    *string
	DefaultTypesFile   *string
	TypeVariationsFile *string
	IcaoVariationsFile *string
	CustomDataFile     *string
	OutputFile         *string
}

func (c *Config) SaveIni() error {
	if c.IniFileName == "" {
		return errors.New("no ini file path given")
	}
	iniFile := ini.Empty()

	section, err := iniFile.NewSection("paths")
	if err != nil {
		return err
	}
	section.Key("liveryDir").SetValue(*c.LiveryDirectory)
	section.Key("defaultTypesFile").SetValue(*c.DefaultTypesFile)
	section.Key("typeVariationsFile").SetValue(*c.TypeVariationsFile)
	section.Key("icaoVariationsFile").SetValue(*c.IcaoVariationsFile)
	section.Key("customDataFile").SetValue(*c.CustomDataFile)
	section.Key("outputFile").SetValue(*c.OutputFile)

	err = iniFile.SaveTo(c.IniFileName)
	if err != nil {
		return err
	}
	return nil
}
