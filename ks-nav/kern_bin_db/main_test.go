/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"archive/tar"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
)

// Tests Makefile parsing feature
func TestParseFilesMakefile(t *testing.T) {
	var test_makefile = []struct {
		FileName string
		Expected kversion
	}{
		{"t_files/linux-4.9.214_Makefile", kversion{4, 9, 214, ""}},
		{"t_files/linux-5.13.9_Makefile", kversion{5, 13, 9, ""}},
		{"t_files/linux-5.15_Makefile", kversion{5, 15, 0, ""}},
		{"t_files/linux-5.18.4_Makefile", kversion{5, 18, 4, ""}},
		{"t_files/linux-5.4.154_Makefile", kversion{5, 4, 154, ""}},
		{"t_files/linux-6.0-rc2_Makefile", kversion{6, 0, 0, "-rc2"}},
	}

	for _, item := range test_makefile {
		makefile, err := get_FromFile(item.FileName)
		if err != nil {
			s := fmt.Sprintf("Error fetch makefile %s", item.FileName)
			t.Error(s)
		}
		v, err := get_version(makefile)
		if err != nil {
			t.Error("Error parsing makefile")
		}
		if v != item.Expected {
			s := fmt.Sprintf("Error in validating the result: got %d.%d.%d%s, expected %d.%d.%d%s", v.Version, v.Patchlevel, v.Sublevel, v.Extraversion,
				item.Expected.Version, item.Expected.Patchlevel, item.Expected.Sublevel, item.Expected.Extraversion)
			t.Error(s)
		}
	}
}

// Tests kernel build configuration feature
func TestParseFilesConfig(t *testing.T) {
	var test_config = []struct {
		FileName string
		Expected string
	}{
		{"t_files/linux-4.9.214_autoconf.h", "7e3619ddf81d683c15e5cb55c57dd16386b359aa"},
		{"t_files/linux-5.13.9_autoconf.h", "99fd41f9da13c43f880ec71500e56b719db4308f"},
		{"t_files/linux-5.15_autoconf.h", "eaa565eaedbbd1b9aaf7bbceb51804cec3dcca53"},
		{"t_files/linux-5.18.4_autoconf.h", "fec6afca6f92e093433727c3c6d1fd07ffbe5f12"},
		{"t_files/linux-5.4.154_autoconf.h", "d1471ae2dbf261ae65089db3b012676834fceae8"},
		{"t_files/linux-6.0-rc2_autoconf.h", "2a6d6426a81c2f84771c00a286f0a592f4cc6a24"},
	}

	for _, item := range test_config {
		config, err := get_FromFile(item.FileName)
		if err != nil {
			s := fmt.Sprintf("Error fetch config %s", item.FileName)
			t.Error(s)
		}
		kconfig := parse_config(config)
		tconf := ""
		keys := make([]string, 0, len(tconf))
		for k := range kconfig {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			tconf = tconf + fmt.Sprintf("CONFIG_%s=%s\n", k, kconfig[k])
		}
		hasher := sha1.New()
		hasher.Write([]byte(tconf))
		sha := hex.EncodeToString(hasher.Sum(nil))
		if sha != item.Expected {
			s := fmt.Sprintf("Error in validating the result: got %s, expected %s", sha, item.Expected)
			t.Error(s)
		}
	}
}

// Untar utility, used to  present a fake kernel build root.
func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

// Utility to copy test files.
func cp(source string, destination string) error {

	bytesRead, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destination, bytesRead, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Tests the maintainer file paths extraction
func TestMaintainer(t *testing.T) {
	var FakeLinuxTreeTest string = "t_files/linux-fake.tar"
	var Fakedir string = "/tmp/linux-fake"
	var testData = []struct {
		filename string
		subs     int
		files    int
	}{
		{"t_files/linux-4.9.214_MAINTAINERS", 1627, 113029},
		{"t_files/linux-5.10.57_MAINTAINERS", 2251, 135588},
		{"t_files/linux-5.13.9_MAINTAINERS", 2336, 136637},
		{"t_files/linux-5.8.2_MAINTAINERS", 2208, 133981},
		{"t_files/linux-6.0_MAINTAINERS", 2585, 142382},
	}
	defer os.RemoveAll(Fakedir)

	err := Untar(FakeLinuxTreeTest, Fakedir)
	if err != nil {
		t.Error("Error cant initialize fake linux directory", err)
	}
	_, filename, _, _ := runtime.Caller(0)
	current := filepath.Dir(filename)
	err = os.Chdir(Fakedir)
	if err != nil {
		t.Error("Error cant initialize fake linux directory", err)
	}
	for _, f := range testData {
		err = cp(current+"/"+f.filename, "MAINTAINERS")
		if err != nil {
			t.Error("Error cant use maintainer file", f.filename)
		}
		s, err := get_FromFile("MAINTAINERS")
		if err != nil {
			t.Error("Error cant read maintainers", err)
		}
		ss := s[seek2data(s):]
		items := parse_maintainers(ss)
		queries := generate_queries("MAINTAINERS", items, "Subsystem: %[1]s' FileName='%[2]s", 0)
		if (f.subs != len(items)) && (f.files != len(queries)) {
			t.Error("Error validating number of files and subsystems in ", f.filename)
		}
	}

}

// Tests the ability to extract the configuration from command line arguments
func TestConfig(t *testing.T) {
	var Default_config configuration = configuration{
		LinuxWDebug:    "vmlinux",
		LinuxWODebug:   "vmlinux.work",
		StripBin:       "/usr/bin/strip",
		DBDriver:       "postgres",
		DBDSN:          "host=dbs.hqhome163.com port=5432 user=alessandro password=<password> dbname=kernel_bin sslmode=disable",
		Maintainers_fn: "MAINTAINERS",
		KConfig_fn:     "include/generated/autoconf.h",
		KMakefile:      "Makefile",
		Mode:           15,
		Note:           "upstream",
	}

	os.Args = []string{"kern_bin_db"}
	conf, err := args_parse(cmd_line_item_init())
	if err != nil {
		t.Error("Error validating empty command line input")
	}
	if conf != Default_config {
		t.Error("Error parsing empty command line input")
	}
	os.Args = []string{"kern_bin_db", "-f"}
	conf, err = args_parse(cmd_line_item_init())
	if err == nil {
		t.Error("error cmd line not detected", conf)
	}
	_, filename, _, _ := runtime.Caller(0)
	current := filepath.Dir(filename)
	os.Args = []string{"kern_bin_db", "-f", current + "/t_files/test1.json"}
	conf, err = args_parse(cmd_line_item_init())
	if err != nil {
		t.Error("error loading sample test configuration 1", err)
	}
	if conf == Default_config {
		t.Error("Error parsing sample test configuration 1", conf)
	}
	if conf.LinuxWDebug != "dummy" {
		t.Error("Error parsing sample test configuration 1", conf)
	}
	os.Args = []string{"kern_bin_db", "-s", "None1"}
	conf, err = args_parse(cmd_line_item_init())
	if conf.StripBin != "None1" {
		t.Error("Error parsing strip binary arg")
	}
	os.Args = []string{"kern_bin_db", "-e", "None2"}
	conf, err = args_parse(cmd_line_item_init())
	if conf.DBDriver != "None2" {
		t.Error("Error parsing database driver arg")
	}
	os.Args = []string{"kern_bin_db", "-d", "None3"}
	conf, err = args_parse(cmd_line_item_init())
	if conf.DBDSN != "None3" {
		t.Error("Error parsing database DSN arg")
	}
	os.Args = []string{"kern_bin_db", "-n", "None4"}
	conf, err = args_parse(cmd_line_item_init())
	if conf.Note != "None4" {
		t.Error("Error parsing note")
	}
}

// Tests the ability to remove duplicates xrefs from the xref list
func TestRemoveDuplicate(t *testing.T) {
	xref_test1 := []xref{ //should find no duplicates
		xref{"indirect", 0xffffffff826a98a7, 0}, xref{"indirect", 0xffffffff826a98ef, 0}, xref{"indirect", 0xffffffff826ab514, 0}, xref{"indirect", 0xffffffff826ab5f3, 0},
		xref{"direct", 0xffffffff82681b9a, 0xffffffff826831c4}, xref{"direct", 0xffffffff8268214a, 0xffffffff82685b2a}, xref{"direct", 0xffffffff8268215a, 0xffffffff826878ab}, xref{"direct", 0xffffffff8268226a, 0xffffffff82687591},
		xref{"direct", 0xffffffff826822aa, 0xffffffff82685dc1}, xref{"direct", 0xffffffff8268245b, 0xffffffff826826aa}, xref{"direct", 0xffffffff8268255b, 0xffffffff826833a5}, xref{"direct", 0xffffffff8268259b, 0xffffffff82680329},
		xref{"direct", 0xffffffff8268269b, 0xffffffff82687778}, xref{"direct", 0xffffffff8268275b, 0xffffffff82681f54}, xref{"direct", 0xffffffff8268295b, 0xffffffff82680948}, xref{"direct", 0xffffffff82682d7b, 0xffffffff8268062e},
		xref{"direct", 0xffffffff82682dcb, 0xffffffff82680a21}, xref{"direct", 0xffffffff82682f2b, 0xffffffff82680e85}, xref{"direct", 0xffffffff8268305b, 0xffffffff82682533}, xref{"direct", 0xffffffff8268315b, 0xffffffff82681357},
		xref{"direct", 0xffffffff8268335b, 0xffffffff82684240}, xref{"direct", 0xffffffff8268377b, 0xffffffff826801af}, xref{"indirect", 0xffffffff826ab60b, 0}, xref{"direct", 0xffffffff8268385b, 0xffffffff826838f8},
		xref{"direct", 0xffffffff82683b5b, 0xffffffff826873c6}, xref{"indirect", 0xffffffff826ab62f, 0}, xref{"indirect", 0xffffffff826ac49f, 0}, xref{"indirect", 0xffffffff826acecb, 0},
		xref{"indirect", 0xffffffff826ad6f3, 0}, xref{"indirect", 0xffffffff826ad993, 0}, xref{"indirect", 0xffffffff826af05b, 0}, xref{"indirect", 0xffffffff826aff13, 0},
		xref{"indirect", 0xffffffff826b091b, 0}, xref{"indirect", 0xffffffff826b0ab3, 0}, xref{"indirect", 0xffffffff826b0c57, 0}, xref{"indirect", 0xffffffff826b180f, 0},
		xref{"indirect", 0xffffffff826b2306, 0}, xref{"indirect", 0xffffffff826b30db, 0}, xref{"indirect", 0xffffffff826b34a7, 0}, xref{"direct", 0xffffffff82683c5b, 0xffffffff82682750},
		xref{"direct", 0xffffffff82683d3b, 0xffffffff82681e45}, xref{"direct", 0xffffffff82683e3b, 0xffffffff82683ce2}, xref{"direct", 0xffffffff82692001, 0xffffffff82683dc1}, xref{"direct", 0xffffffff82692243, 0xffffffff82684737},
		xref{"direct", 0xffffffff8269264a, 0xffffffff82684816}, xref{"direct", 0xffffffff82692a47, 0xffffffff82681f18}, xref{"indirect", 0xffffffff826b3b43, 0}, xref{"indirect", 0xffffffff826b43a7, 0},
		xref{"indirect", 0xffffffff826b448a, 0}, xref{"indirect", 0xffffffff826b4593, 0}, xref{"indirect", 0xffffffff826b47a2, 0}, xref{"indirect", 0xffffffff826ee65f, 0}, xref{"direct", 0xffffffff8269551e, 0xffffffff82686005},
		xref{"direct", 0xffffffff8269564a, 0xffffffff82682056}, xref{"direct", 0xffffffff826956ff, 0xffffffff82684fb0}, xref{"indirect", 0xffffffff826eeabf, 0}, xref{"indirect", 0xffffffff826ef91f, 0},
		xref{"indirect", 0xffffffff826efa87, 0}, xref{"indirect", 0xffffffff826efc17, 0}, xref{"direct", 0xffffffff8269aa6f, 0xffffffff826847a9}, xref{"indirect", 0xffffffff826efd7f, 0},
		xref{"direct", 0xffffffff8269ac97, 0xffffffff82687579}, xref{"indirect", 0xffffffff826f04af, 0}, xref{"direct", 0xffffffff8269aecb, 0xffffffff82686991}, xref{"indirect", 0xffffffff826f05c7, 0},
		xref{"indirect", 0xffffffff826f1157, 0}, xref{"indirect", 0xffffffff826f1427, 0}, xref{"indirect", 0xffffffff826f17e7, 0}, xref{"indirect", 0xffffffff826f18ff, 0},
		xref{"indirect", 0xffffffff826f1c47, 0}, xref{"indirect", 0xffffffff826f1cbf, 0}, xref{"indirect", 0xffffffff826f1e27, 0}, xref{"indirect", 0xffffffff826f27d7, 0},
		xref{"indirect", 0xffffffff826f2a2f, 0},
	}
	xref_test2 := []xref{ //should find 2 duplicates
		xref{"indirect", 0xffffffff826a98a7, 0}, xref{"indirect", 0xffffffff826a98ef, 0}, xref{"indirect", 0xffffffff826ab514, 0}, xref{"indirect", 0xffffffff826ab5f3, 0},
		xref{"direct", 0xffffffff82681b9a, 0xffffffff826831c4}, xref{"direct", 0xffffffff8268214a, 0xffffffff82685b2a}, xref{"direct", 0xffffffff8268215a, 0xffffffff826878ab}, xref{"direct", 0xffffffff8268226a, 0xffffffff82687591},
		xref{"direct", 0xffffffff826822aa, 0xffffffff82685dc1}, xref{"direct", 0xffffffff8268245b, 0xffffffff826826aa}, xref{"direct", 0xffffffff8268255b, 0xffffffff826833a5}, xref{"direct", 0xffffffff8268259b, 0xffffffff82680329},
		xref{"direct", 0xffffffff8268269b, 0xffffffff82687778}, xref{"direct", 0xffffffff8268275b, 0xffffffff82681f54}, xref{"direct", 0xffffffff8268295b, 0xffffffff82680948}, xref{"direct", 0xffffffff82682d7b, 0xffffffff8268062e},
		xref{"direct", 0xffffffff82682dcb, 0xffffffff82680a21}, xref{"direct", 0xffffffff82682f2b, 0xffffffff82680e85}, xref{"direct", 0xffffffff8268305b, 0xffffffff82682533}, xref{"direct", 0xffffffff8268315b, 0xffffffff82681357},
		xref{"direct", 0xffffffff8268335b, 0xffffffff82684240}, xref{"direct", 0xffffffff8268377b, 0xffffffff826801af}, xref{"indirect", 0xffffffff826ab60b, 0}, xref{"direct", 0xffffffff8268385b, 0xffffffff826838f8},
		xref{"direct", 0xffffffff82683b5b, 0xffffffff826873c6}, xref{"indirect", 0xffffffff826ab62f, 0}, xref{"indirect", 0xffffffff826ac49f, 0}, xref{"indirect", 0xffffffff826acecb, 0},
		xref{"indirect", 0xffffffff826ad6f3, 0}, xref{"indirect", 0xffffffff826ad993, 0}, xref{"indirect", 0xffffffff826af05b, 0}, xref{"indirect", 0xffffffff826aff13, 0},
		xref{"indirect", 0xffffffff826b091b, 0}, xref{"indirect", 0xffffffff826b0ab3, 0}, xref{"indirect", 0xffffffff826b0c57, 0}, xref{"indirect", 0xffffffff826b180f, 0},
		xref{"indirect", 0xffffffff826b2306, 0}, xref{"indirect", 0xffffffff826b30db, 0}, xref{"indirect", 0xffffffff826b34a7, 0}, xref{"direct", 0xffffffff82683c5b, 0xffffffff82682750},
		xref{"direct", 0xffffffff82683d3b, 0xffffffff82681e45}, xref{"direct", 0xffffffff82683e3b, 0xffffffff82683ce2}, xref{"direct", 0xffffffff82692001, 0xffffffff82683dc1}, xref{"direct", 0xffffffff82692243, 0xffffffff82684737},
		xref{"direct", 0xffffffff8269264a, 0xffffffff82684816}, xref{"direct", 0xffffffff82692a47, 0xffffffff82681f18}, xref{"indirect", 0xffffffff826b3b43, 0}, xref{"indirect", 0xffffffff826b43a7, 0},
		xref{"indirect", 0xffffffff826b448a, 0}, xref{"indirect", 0xffffffff826b4593, 0}, xref{"indirect", 0xffffffff826b47a2, 0}, xref{"indirect", 0xffffffff826ee65f, 0}, xref{"direct", 0xffffffff8269551e, 0xffffffff82686005},
		xref{"direct", 0xffffffff8269564a, 0xffffffff82682056}, xref{"direct", 0xffffffff826956ff, 0xffffffff82684fb0}, xref{"indirect", 0xffffffff826eeabf, 0}, xref{"indirect", 0xffffffff826ef91f, 0},
		xref{"indirect", 0xffffffff826efa87, 0}, xref{"indirect", 0xffffffff826efc17, 0}, xref{"direct", 0xffffffff8269aa6f, 0xffffffff826847a9}, xref{"indirect", 0xffffffff826efd7f, 0},
		xref{"direct", 0xffffffff8269ac97, 0xffffffff82687579}, xref{"indirect", 0xffffffff826f04af, 0}, xref{"direct", 0xffffffff8269aecb, 0xffffffff82686991}, xref{"indirect", 0xffffffff826f05c7, 0},
		xref{"indirect", 0xffffffff826f1157, 0}, xref{"indirect", 0xffffffff826f1427, 0}, xref{"indirect", 0xffffffff826f17e7, 0}, xref{"indirect", 0xffffffff826f18ff, 0},
		xref{"indirect", 0xffffffff826f1c47, 0}, xref{"indirect", 0xffffffff826f1cbf, 0}, xref{"indirect", 0xffffffff826f1e27, 0}, xref{"indirect", 0xffffffff826f27d7, 0},
		xref{"indirect", 0xffffffff826f2a2f, 0}, xref{"indirect", 0xffffffff826a98a7, 0}, xref{"direct", 0xffffffff82682766, 0xffffffff82681f54},
	}

	res := removeDuplicate(xref_test1)
	if len(xref_test1) != len(res) {
		t.Error("Expect to find no duplicates. Test failed.")
	}

	res = removeDuplicate(xref_test2)
	if len(xref_test2)-2 != len(res) {
		t.Error("Expect to find duplicates. Test failed.")
	}
}
