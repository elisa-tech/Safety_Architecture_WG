/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	appName  string = "App Name: nav"
	appDescr string = "Descr: kernel symbol navigator"
)

type argFunc func(*configuration, []string) error

// Command line switch elements.
type cmdLineItems struct {
	function  argFunc
	switchStr string
	helpStr   string
	id        int
	hasArg    bool
	needed    bool
}

// Represents the application configuration.
type configuration struct {
	cmdlineNeeds   map[string]bool
	DBDriver       string
	DBDSN          string
	Symbol         string
	Jout           string
	ExcludedBefore []string
	ExcludedAfter  []string
	TargetSubsys   []string
	Instance       int
	MaxDepth       int
	Graphviz       outIMode
	Mode           outMode
}

// Instance of default configuration values.
var defaultConfig = configuration{
	DBDriver:       "postgres",
	DBDSN:          "host=dbs.hqhome163.com port=5432 user=alessandro password=<password> dbname=kernel_bin sslmode=disable",
	Symbol:         "",
	Instance:       0,
	Mode:           printSubsys,
	ExcludedBefore: []string{},
	ExcludedAfter:  []string{},
	TargetSubsys:   []string{},
	MaxDepth:       0, //0: no limit
	Jout:           "graphOnly",
	Graphviz:       oText,
	cmdlineNeeds:   map[string]bool{},
}

// Inserts a commandline item, which is composed by:
// * switch string
// * switch description
// * if the switch requires an additional argument
// * a pointer to the function that manages the switch
// * the configuration that gets updated.
func pushCmdLineItem(switchStr string, helpStr string, hasArg bool, needed bool, function argFunc, cmdLine *[]cmdLineItems) {
	*cmdLine = append(*cmdLine, cmdLineItems{id: len(*cmdLine) + 1, switchStr: switchStr, helpStr: helpStr, hasArg: hasArg, needed: needed, function: function})
}

// This function initializes configuration parser subsystem
// Inserts all the commandline switches supported by the application.
func cmdLineItemInit() []cmdLineItems {
	var res []cmdLineItems

	pushCmdLineItem("-j", "Force Json output with subsystems data", true, false, funcOutType, &res)
	pushCmdLineItem("-s", "Specifies symbol", true, true, funcSymbol, &res)
	pushCmdLineItem("-i", "Specifies instance", true, true, funcInstance, &res)
	pushCmdLineItem("-f", "Specifies config file", true, false, funcJconf, &res)
	pushCmdLineItem("-e", "Forces to use a specified DB Driver (i.e. postgres, mysql or sqlite3)", true, false, funcDBDriver, &res)
	pushCmdLineItem("-d", "Forces to use a specified DB DSN", true, false, funcDBDSN, &res)
	pushCmdLineItem("-m", "Sets display mode 2=subsystems,1=all", true, false, funcMode, &res)
	pushCmdLineItem("-x", "Specify Max depth in call flow exploration", true, false, funcDepth, &res)
	pushCmdLineItem("-g", "if -j option is graphOnly (default) output PNG in place of dot", true, false, funcGraphviz, &res)
	pushCmdLineItem("-h", "This help", false, false, funcHelp, &res)

	return res
}

func funcHelp(conf *configuration, fn []string) error {
	return errors.New("command help")
}

func funcGraphviz(conf *configuration, imode []string) error {
	s, err := strconv.Atoi(imode[0])
	if err != nil {
		return err
	}
	if outIMode(s) < oText || outIMode(s) >= OutIModeLast {
		return errors.New("unsupported mode")
	}
	conf.Graphviz = outIMode(s)
	return nil
}

func funcOutType(conf *configuration, jout []string) error {
	conf.Jout = jout[0]
	return nil
}

func funcJconf(conf *configuration, fn []string) error {
	jsonFile, err := os.Open(fn[0])
	if err != nil {
		return err
	}
	defer func() {
		closeErr := jsonFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, conf)
	if err != nil {
		return err
	}
	return nil
}

func funcSymbol(conf *configuration, fn []string) error {
	conf.Symbol = fn[0]
	return nil
}

func funcDBDriver(conf *configuration, driver []string) error {
	conf.DBDriver = driver[0]
	return nil
}

func funcDBDSN(conf *configuration, dsn []string) error {
	conf.DBDSN = dsn[0]
	return nil
}

func funcDepth(conf *configuration, depth []string) error {
	s, err := strconv.Atoi(depth[0])
	if err != nil {
		return err
	}
	if s < 0 {
		return errors.New("depth must be >= 0")
	}
	conf.MaxDepth = s
	return nil
}

func funcInstance(conf *configuration, instance []string) error {
	s, err := strconv.Atoi(instance[0])
	if err != nil {
		return err
	}
	conf.Instance = s
	return nil
}

func funcMode(conf *configuration, mode []string) error {
	s, err := strconv.Atoi(mode[0])
	if err != nil {
		return err
	}
	if outMode(s) < printAll || outMode(s) >= OutModeLast {
		return errors.New("unsupported mode")
	}
	conf.Mode = outMode(s)
	return nil
}

// Uses commandline args to generate the help string.
func printHelp(lines []cmdLineItems) {

	fmt.Println(appName)
	fmt.Println(appDescr)
	for _, item := range lines {
		fmt.Printf(
			"\t%s\t%s\t%s\n",
			item.switchStr,
			func(a bool) string {
				if a {
					return "<v>"
				}
				return ""
			}(item.hasArg),
			item.helpStr,
		)
	}
}

// Used to parse the command line and generate the command line.
func argsParse(lines []cmdLineItems) (configuration, error) {
	var extra = false
	var conf = defaultConfig
	var f argFunc

	for _, item := range lines {
		if item.needed {
			conf.cmdlineNeeds[item.switchStr] = false
		}
	}

	for _, osArg := range os.Args[1:] {
		if !extra {
			for _, arg := range lines {
				if arg.switchStr == osArg {
					if arg.needed {
						conf.cmdlineNeeds[arg.switchStr] = true
					}
					if arg.hasArg {
						f = arg.function
						extra = true
						break
					}
					err := arg.function(&conf, []string{})
					if err != nil {
						return defaultConfig, err
					}
				}
			}
			continue
		}
		if extra {
			err := f(&conf, []string{osArg})
			if err != nil {
				return defaultConfig, err
			}
			extra = false
		}

	}
	if extra {
		return defaultConfig, errors.New("missing switch arg")
	}

	res := true
	for _, element := range conf.cmdlineNeeds {
		res = res && element
	}
	if res {
		return conf, nil
	}
	return defaultConfig, errors.New("missing needed arg")
}
