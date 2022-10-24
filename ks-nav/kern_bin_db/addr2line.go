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
	"fmt"
	"strings"
	"database/sql"
	"path/filepath"
	addr2line "github.com/elazarl/addr2line"
)

// Represents one task for the addr to line subsystem.
type workloads struct{
	Addr	uint64
	Name	string
	Query	string
	DB	*sql.DB
	}

// Caches item elements
type Addr2line_items struct {
	Addr		uint64
	File_name	string
	}

// Commandline handle functions prototype
type ins_f func(*sql.DB, string, bool)

// The adrr2line item cache
var Addr2line_cache []Addr2line_items;

// Initializes the addr to line subsystem
// prepares a channel for the addr2line resolver communication.
func addr2line_init(fn string) (chan workloads){
	a, err := addr2line.New(fn)
	if err != nil {
		panic( err)
		}
	adresses := make(chan workloads, 16)
	go workload(a, adresses, Insert_data)
	return adresses
}

// Checks the current symbol is in cache, if present returns data from the cache.
func in_cache(Addr uint64, Addr2line_cache []Addr2line_items)(bool, string){
	for _,a := range Addr2line_cache {
		if a.Addr == Addr {
			return true, a.File_name
			}
		}
	return false, ""
}

// The goroutine responsible to resolve queries.
// it reads from a channel for workloads elements.
// because it manages the database it can also receive
// raw queries. Raw queries are workloads whose Name="None"
// if proper workloads, a resolution is triggered.
func workload(a *addr2line.Addr2line, addresses chan workloads, insert_func ins_f){
	var e	workloads
	var qready string

	for {
		e = <-addresses
		switch e.Name {
		case "None":
			insert_func(e.DB, e.Query, false)
			break
		default:
			rs, _ := a.Resolve(e.Addr)
			if len(rs)==0 {
				qready=fmt.Sprintf(e.Query, "NONE")
				}
			for _, a:=range rs{
				qready=fmt.Sprintf(e.Query, filepath.Clean(a.File))
				if a.Function == strings.ReplaceAll(e.Name, "sym.", "") {
					break
					}
				}
			insert_func(e.DB, qready, false)
			break
			}
	}
}

// Sends a workload to the resolver.
func spawn_query(db *sql.DB, addr uint64, name string, addresses chan workloads, query string) {
	addresses <- workloads{addr, name, query, db}
}
