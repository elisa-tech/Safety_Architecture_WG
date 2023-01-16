/*
 * ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *   Name: kern_bin_db - Kernel source code analysis tool database creator
 *   Description: Parses kernel source tree and binary images and builds the DB
 *
 *   Author: Alessandro Carminati <acarmina@redhat.com>
 *   Author: Maurizio Papini <mpapini@redhat.com>
 *
 * ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *
 *   Copyright (c) 2022 Red Hat, Inc. All rights reserved.
 *
 *   This copyrighted material is made available to anyone wishing
 *   to use, modify, copy, or redistribute it subject to the terms
 *   and conditions of the GNU General Public License version 2.
 *
 *   This program is distributed in the hope that it will be
 *   useful, but WITHOUT ANY WARRANTY; without even the implied
 *   warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
 *   PURPOSE. See the GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public
 *   License along with this program; if not, write to the Free
 *   Software Foundation, Inc., 51 Franklin Street, Fifth Floor,
 *   Boston, MA 02110-1301, USA.
 *
 * ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
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
		DBURL:          "dbs.hqhome163.com",
		DBPort:         5432,
		DBUser:         "alessandro",
		DBPassword:     "<password>",
		DBTargetDB:     "kernel_bin",
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
	os.Args = []string{"kern_bin_db", "-u", "None2"}
	conf, err = args_parse(cmd_line_item_init())
	if conf.DBUser != "None2" {
		t.Error("Error parsing database userid arg")
	}
	os.Args = []string{"kern_bin_db", "-p", "None3"}
	conf, err = args_parse(cmd_line_item_init())
	if conf.DBPassword != "None3" {
		t.Error("Error parsing password arg")
	}
	os.Args = []string{"kern_bin_db", "-d", "None4"}
	conf, err = args_parse(cmd_line_item_init())
	if conf.DBURL != "None4" {
		t.Error("Error parsing db url arg")
	}
	os.Args = []string{"kern_bin_db", "-o", "1234"}
	conf, err = args_parse(cmd_line_item_init())
	if conf.DBPort != 1234 {
		t.Error("Error parsing db port arg")
	}
}
