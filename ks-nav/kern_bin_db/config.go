/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// Default config file
var conf_fn = "conf.json"

type Arg_func func(*configuration, []string) error

type Secret struct {
	Secret  string `json:"secret"`
	Fix     string `json:"fix"`
	Counter uint64 `json:"counter"`
	Len     int    `json:"len"`
}

// Command line switch elements
type cmd_line_items struct {
	id       int
	Switch   string
	Help_srt string
	Has_arg  bool
	Func     Arg_func
}

// Represents the application configuration
type configuration struct {
	LinuxWDebug    string
	LinuxWODebug   string
	StripBin       string
	DBDriver       string
	DBDSN          string
	Maintainers_fn string
	KConfig_fn     string
	KMakefile      string
	Mode           int
	Note           string
}

// Instance of default configuration values
var Default_config configuration = configuration{
	LinuxWDebug:    "vmlinux",
	LinuxWODebug:   "vmlinux.work",
	StripBin:       "/usr/bin/strip",
	DBDriver:       "postgres",
	DBDSN:          "host=dbs.hqhome163.com port=5432 user=alessandro password=<password> dbname=kernel_bin sslmode=disable",
	Maintainers_fn: "MAINTAINERS",
	KConfig_fn:     "include/generated/autoconf.h",
	KMakefile:      "Makefile",
	Mode:           ENABLE_SYBOLSNFILES | ENABLE_XREFS | ENABLE_MAINTAINERS | ENABLE_VERSION_CONFIG,
	Note:           "upstream",
}

// Inserts a commandline item item, which is composed by:
// * switch string
// * switch descriptio
// * if the switch requires an additiona argument
// * a pointer to the function that manages the switch
// * the configuration that gets updated
func push_cmd_line_item(Switch string, Help_str string, Has_arg bool, Func Arg_func, cmd_line *[]cmd_line_items) {
	*cmd_line = append(*cmd_line, cmd_line_items{id: len(*cmd_line) + 1, Switch: Switch, Help_srt: Help_str, Has_arg: Has_arg, Func: Func})
}

// This function initializes configuration parser subsystem
// Inserts all the commandline switches suppported by the application
func cmd_line_item_init() []cmd_line_items {
	var res []cmd_line_items

	push_cmd_line_item("-f", "specifies json configuration file", true, func_jconf, &res)
	push_cmd_line_item("-s", "Forces use specified strip binary", true, func_forceStrip, &res)
	push_cmd_line_item("-e", "Forces to use a specified DB Driver (i.e. postgres, mysql or sqlite3)", true, func_DBDriver, &res)
	push_cmd_line_item("-d", "Forces to use a specified DB DSN", true, func_DBDSN, &res)
	push_cmd_line_item("-n", "Forecs use specified note (default 'upstream')", true, func_Note, &res)
	push_cmd_line_item("-c", "Checks dependencies", false, func_check, &res)
	push_cmd_line_item("-h", "This Help", false, func_help, &res)

	return res
}

func func_help(conf *configuration, fn []string) error {
	return errors.New("Dummy")
}

func func_jconf(conf *configuration, fn []string) error {
	jsonFile, err := os.Open(fn[0])
	if err != nil {
		return err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	err = json.Unmarshal(byteValue, conf)
	if err != nil {
		return err
	}
	return nil
}

func func_forceStrip(conf *configuration, fn []string) error {
	(*conf).StripBin = fn[0]
	return nil
}

func func_DBDriver(conf *configuration, db_driver []string) error {
	(*conf).DBDriver = db_driver[0]
	return nil
}

func func_DBDSN(conf *configuration, db_dsn []string) error {
	(*conf).DBDSN = db_dsn[0]
	return nil
}

func func_Note(conf *configuration, note []string) error {
	(*conf).Note = note[0]
	return nil
}

func func_check(conf *configuration, dummy []string) error {
	return nil
}

// Uses commandline args to generate the help string
func print_help(lines []cmd_line_items) {

	for _, item := range lines {
		fmt.Printf(
			"\t%s\t%s\t%s\n",
			item.Switch,
			func(a bool) string {
				if a {
					return "<v>"
				}
				return ""
			}(item.Has_arg),
			item.Help_srt,
		)
	}
}

// Used to parse the command line and generate the command line
func args_parse(lines []cmd_line_items) (configuration, error) {
	var extra bool = false
	var conf configuration = Default_config
	var f Arg_func

	for _, os_arg := range os.Args[1:] {
		if !extra {
			for _, arg := range lines {
				if arg.Switch == os_arg {
					if arg.Has_arg {
						f = arg.Func
						extra = true
						break
					}
					err := arg.Func(&conf, []string{})
					if err != nil {
						return Default_config, err
					}
				}
			}
			continue
		}
		if extra {
			err := f(&conf, []string{os_arg})
			if err != nil {
				return Default_config, err
			}
			extra = false
		}
	}
	if extra {
		return conf, errors.New("extra arg needed but none")
	} else {
		return conf, nil
	}
}
