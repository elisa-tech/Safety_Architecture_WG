/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigTest", func() {
	_, curreentFilename, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(curreentFilename)
	var testConfig configuration

	BeforeEach(func() {
		testConfig = configuration{
			DBDriver:       "dummy",
			DBDSN:          "dummy",
			Symbol:         "dummy",
			Instance:       1234,
			Mode:           1234,
			ExcludedBefore: []string{"dummy1", "dummy2", "dummy3"},
			ExcludedAfter:  []string{"dummyA", "dummyB", "dummyC"},
			MaxDepth:       1234, //0: no limit
			Jout:           "jsonOutputPlain",
			cmdlineNeeds:   map[string]bool{},
			TargetSubsys:   []string{},
			Graphviz:       oText,
		}
	})

	When("The CLI is invoked with no arguments", func() {
		It("Should fail and inform the user of missing arguments", func() {
			os.Args = []string{"nav"}
			conf, err := argsParse(cmdLineItemInit())
			Expect(conf).To(Equal(defaultConfig))
			Expect(err).ToNot(BeNil())
			Expect(fmt.Sprintf("%s", err)).To(Equal("missing needed arg"))
		})
	})

	When("The CLI is invoked with an invalid switch argument", func() {
		It("Should inform the user about the missing swicth argument", func() {
			os.Args = []string{"nav", "-i", "1", "-s"}
			conf, err := argsParse(cmdLineItemInit())
			Expect(conf).To(Equal(defaultConfig))
			Expect(err).ToNot(BeNil())
			Expect(fmt.Sprintf("%s", err)).To(Equal("missing switch arg"))
		})

		It("Should inform the user about the invalid swicth argument type", func() {
			os.Args = []string{"nav", "-i", "a", "-s", "symb"}
			conf, err := argsParse(cmdLineItemInit())
			Expect(conf).To(Equal(defaultConfig))
			Expect(err).ToNot(BeNil())
			Expect(fmt.Sprintf("%s", err)).To(Equal("strconv.Atoi: parsing \"a\": invalid syntax"))
		})
	})

	When("The CLI is invoked with -f", func() {
		It("Should inform the user that no value was provided", func() {
			os.Args = []string{"nav", "-i", "a", "-s", "symb", "-f"}
			conf, err := argsParse(cmdLineItemInit())
			Expect(conf).To(Equal(defaultConfig))
			Expect(err).ToNot(BeNil())
			Expect(fmt.Sprintf("%s", err)).To(Equal("strconv.Atoi: parsing \"a\": invalid syntax"))
		})

		It("Should inform the user that the file path does not exist", func() {
			os.Args = []string{"nav", "-i", "1", "-s", "symb", "-f", cwd + "/t_files/dummy.json"}
			conf, err := argsParse(cmdLineItemInit())
			Expect(conf).To(Equal(defaultConfig))
			Expect(err).ToNot(BeNil())
			Expect(fmt.Sprintf("%s", err)).To(Equal("open " + cwd + "/t_files/dummy.json: no such file or directory"))
		})

		It("Should parse the config file", func() {
			testConfig.cmdlineNeeds = map[string]bool{"-s": true, "-i": true}
			os.Args = []string{"nav", "-i", "1", "-s", "symb", "-f", cwd + "/t_files/test1.json"}
			conf, err := argsParse(cmdLineItemInit())
			Expect(conf).To(Equal(testConfig))
			Expect(err).To(BeNil())
		})

		It("Should overwrite a config property from the CLI", func() {
			testConfig.cmdlineNeeds = map[string]bool{"-s": true, "-i": true}
			os.Args = []string{
				"nav", "-i", "1",
				"-s", "symb",
				"-f", cwd + "/t_files/test1.json",
				"-u", "new",
			}
			conf, err := argsParse(cmdLineItemInit())
			Expect(conf).To(Equal(testConfig))
			Expect(err).To(BeNil())
		})
	})
})
