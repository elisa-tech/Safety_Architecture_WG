//go:build !CGO
// +build !CGO

/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nav Tests", func() {
	Describe("opt2num", func() {
		When("Using a valid options key", func() {
			It("Should return the correct value for graphOnly", func() {
				Expect(opt2num("graphOnly")).To(Equal(graphOnly))
			})

			It("Should return the correct value for jsonOutputPlain", func() {
				Expect(opt2num("jsonOutputPlain")).To(Equal(jsonOutputPlain))
			})

			It("Should return the correct value for jsonOutputB64", func() {
				Expect(opt2num("jsonOutputB64")).To(Equal(jsonOutputB64))
			})

			It("Should return the correct value for jsonOutputGZB64", func() {
				Expect(opt2num("jsonOutputGZB64")).To(Equal(jsonOutputGZB64))
			})
		})

		When("Using an invalid options key", func() {
			It("Should return 0", func() {
				Expect(opt2num("invalidKey")).To(Equal(invalidOutput))
			})
		})
	})

	Describe("decorateLine", func() {
		var l string
		var r string
		var adjm []adjM
		BeforeEach(func() {
			l = "lsys"
			r = "rsys"
			adjm = []adjM{
				{
					l: node{
						subsys:     "lsys",
						symbol:     "lsym",
						sourceRef:  "lsource",
						addressRef: "laddr",
					},
					r: node{
						subsys:     "rsys",
						symbol:     "rsym",
						sourceRef:  "rsource",
						addressRef: "raddr",
					},
				},
			}
		})
		When("Subsystems match", func() {
			It("Should return a list of all matched subsystems", func() {
				actual := decorateLine(l, r, adjm)
				expected := " [label=\"rsym([raddr]rsource),\\n\"]"

				Expect(actual).To(Equal(expected))
			})

			It("Should ignore duplicated entries", func() {
				duplicated := adjM{
					l: node{
						subsys:     "lsys",
						symbol:     "lsym",
						sourceRef:  "lsource",
						addressRef: "laddr",
					},
					r: node{
						subsys:     "rsys",
						symbol:     "rsym",
						sourceRef:  "rsource",
						addressRef: "raddr",
					},
				}
				adjm = append(adjm, duplicated)

				actual := decorateLine(l, r, adjm)
				expected := " [label=\"rsym([raddr]rsource),\\n\"]"

				Expect(actual).To(Equal(expected))
			})

			It("Should return more than one entry", func() {
				extra := adjM{
					l: node{
						subsys:     "lsys",
						symbol:     "lsym2",
						sourceRef:  "lsource2",
						addressRef: "laddr2",
					},
					r: node{
						subsys:     "rsys",
						symbol:     "rsym2",
						sourceRef:  "rsource2",
						addressRef: "raddr2",
					},
				}
				adjm = append(adjm, extra)

				actual := decorateLine(l, r, adjm)
				expected := " [label=\"rsym([raddr]rsource),\\nrsym2([raddr2]rsource2),\\n\"]"

				Expect(actual).To(Equal(expected))
			})
		})

		When("Subsystems do not match", func() {
			It("Should return an empty list if using an empty slice", func() {
				actual := decorateLine(l, r, []adjM{})
				expected := " [label=\"\"]"

				Expect(actual).To(Equal(expected))

			})

			It("Should return an empty list if nodes do not match", func() {
				actual := decorateLine(l, "asym", adjm)
				expected := " [label=\"\"]"

				Expect(actual).To(Equal(expected))

			})
		})
	})

	Describe("generateOutput using sqlmock", func() {
		var d *sqlMock
		expectedDot := `digraph G {
rankdir="LR"
"__x64_sys_getpid"->"__task_pid_nr_ns" 
"__task_pid_nr_ns"->"__rcu_read_lock" 
"__task_pid_nr_ns"->"__rcu_read_unlock" 
}`
		d = &sqlMock{}
		d.init(nil)
		d.LOADsym2numValues("__x64_sys_getpid", 16, 472055, nil)
		d.LOADgetEntryByIdValues(472055, 16, entry{symbol: "__x64_sys_getpid", fn: "kernel/sys.c", sourceRef: "", addressRef: "", subsys: []string{}, symId: 472055}, nil)
		d.LOADgetSubsysFromSymbolNameValues("__x64_sys_getpid", 16, "", nil)
		d.LOADgetSuccessorsByIdValues(472055, 16, []entry{
			entry{symbol: "__fentry__", fn: "arch/x86/kernel/ftrace_64.S", sourceRef: "kernel/sys.c:892", addressRef: "0xffffffff81077570", subsys: []string{"X86 ARCHITECTURE (32-BIT AND 64-BIT)"}, symId: 501994},
			entry{symbol: "__task_pid_nr_ns", fn: "kernel/pid.c", sourceRef: "kernel/sys.c:893", addressRef: "0xffffffff81077589", subsys: []string{}, symId: 472243},
		}, nil)
		d.LOADgetEntryByIdValues(501994, 16, entry{symbol: "__fentry__", fn: "arch/x86/kernel/ftrace_64.S", sourceRef: "", addressRef: "", subsys: []string{"X86 ARCHITECTURE (32-BIT AND 64-BIT)"}, symId: 501994}, nil)
		d.LOADgetEntryByIdValues(472243, 16, entry{symbol: "__task_pid_nr_ns", fn: "kernel/pid.c", sourceRef: "", addressRef: "", subsys: []string{}, symId: 472243}, nil)
		d.LOADgetSubsysFromSymbolNameValues("symbol=__task_pid_nr_ns", 16, "", nil)
		d.LOADgetSuccessorsByIdValues(472243, 16, []entry{
			entry{symbol: "__fentry__", fn: "arch/x86/kernel/ftrace_64.S", sourceRef: "kernel/pid.c:427", addressRef: "0xffffffff810824e0", subsys: []string{"X86 ARCHITECTURE (32-BIT AND 64-BIT)"}, symId: 501994},
			entry{symbol: "__rcu_read_lock", fn: "kernel/rcu/tree_plugin.h", sourceRef: "kernel/pid.c:430", addressRef: "0xffffffff810824f7", subsys: []string{"READ-COPY UPDATE (RCU)"}, symId: 473674},
			entry{symbol: "__rcu_read_unlock", fn: "kernel/rcu/tree_plugin.h", sourceRef: "kernel/pid.c:435", addressRef: "0xffffffff81082540", subsys: []string{"READ-COPY UPDATE (RCU)"}, symId: 473716},
			entry{symbol: "__rcu_read_unlock", fn: "kernel/rcu/tree_plugin.h", sourceRef: "kernel/pid.c:435", addressRef: "0xffffffff81082584", subsys: []string{"READ-COPY UPDATE (RCU)"}, symId: 473716},
		}, nil)
		d.LOADgetEntryByIdValues(501994, 16, entry{symbol: "__fentry__", fn: "arch/x86/kernel/ftrace_64.S", sourceRef: "", addressRef: "", subsys: []string{"X86 ARCHITECTURE (32-BIT AND 64-BIT)"}, symId: 501994}, nil)
		d.LOADgetEntryByIdValues(473674, 16, entry{symbol: "__rcu_read_lock", fn: "kernel/rcu/tree_plugin.h", sourceRef: "", addressRef: "", subsys: []string{"READ-COPY UPDATE (RCU)"}, symId: 473674}, nil)
		d.LOADgetEntryByIdValues(473716, 16, entry{symbol: "__rcu_read_unlock", fn: "kernel/rcu/tree_plugin.h", sourceRef: "", addressRef: "", subsys: []string{"READ-COPY UPDATE (RCU)"}, symId: 473716}, nil)
		d.LOADgetEntryByIdValues(473716, 16, entry{symbol: "__rcu_read_unlock", fn: "kernel/rcu/tree_plugin.h", sourceRef: "", addressRef: "", subsys: []string{"READ-COPY UPDATE (RCU)"}, symId: 473716}, nil)
		d.LOADgetSubsysFromSymbolNameValues("_rcu_read_lock", 16, "READ-COPY UPDATE (RCU)", nil)
		d.LOADgetSubsysFromSymbolNameValues("__rcu_read_unlock", 16, "READ-COPY UPDATE (RCU)", nil)
		d.LOADsymbSubsysValues([]int{472055, 472243}, 16, "{\"FuncName\":\"__x64_sys_getpid\", \"subsystems\":[]},{\"FuncName\":\"__task_pid_nr_ns\", \"subsystems\":[]}", nil)
		d.LOADgetEntryByIdValues(472055, 16, entry{symbol: "__x64_sys_getpid", fn: "kernel/sys.c", sourceRef: "", addressRef: "", subsys: []string{}, symId: 472055}, nil)
		d.LOADgetEntryByIdValues(472243, 16, entry{symbol: "__task_pid_nr_ns", fn: "kernel/pid.c", sourceRef: "", addressRef: "", subsys: []string{}, symId: 472243}, nil)
		testConfig := configuration{
			DBDriver:       "postgres",
			DBDSN:          "host=dbs.hqhome163.com port=5432 user=alessandro password=<password> dbname=kernel_bin sslmode=disable",
			Symbol:         "__x64_sys_getpid",
			Instance:       16,
			Mode:           printAll,
			ExcludedBefore: []string{"__fentry__", "__stack_chk_fail"},
			ExcludedAfter:  []string{"^kfree$", "^_raw_spin_lock$", "^_raw_spin_unlock$", "^panic$", "^call_rcu$", "^__call_rcu$", "__rcu_read_unlock", "__rcu_read_lock", "path_openat"},
			TargetSubsys:   []string{},
			MaxDepth:       0, //0: no limit
			Jout:           "graphOnly",
			Graphviz:       oText,
			cmdlineNeeds:   map[string]bool{},
		}
		dot, err := generateOutput(d, &testConfig)
		It("Should return syntax correct json with no error", func() {
			Expect(err).To(BeNil())
			Expect(dot).To(Equal(expectedDot))
		})

	})

	Describe("generateOutput using go-sqlmock", func() {
		type mockQueries struct {
			querySTR     string
			resultHead   []string
			resultValues [][]driver.Value
		}
		var queryTestSerie []mockQueries
		var db *sql.DB
		var mock sqlmock.Sqlmock
		var dok *SqlDB
		expectedDot := `digraph G {
rankdir="LR"
"__x64_sys_getpid"->"__task_pid_nr_ns" 
"__task_pid_nr_ns"->"__rcu_read_lock" 
"__task_pid_nr_ns"->"__rcu_read_unlock" 
}`
		dok = &SqlDB{}
		db, mock, _ = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		dok.db = db
		dok.cache = Cache{}

		queryTestSerie = []mockQueries{}
		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR:     "select symbol_id from symbols where symbols.symbol_name='__x64_sys_getpid' and symbols.symbol_instance_id_ref=16",
			resultHead:   []string{"symbol_id"},
			resultValues: [][]driver.Value{{"472055"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR: "select symbol_id, symbol_name, subsys_name, file_name from (select * from symbols, files where symbols.symbol_file_ref_id=files.file_id and " +
				"symbols.symbol_instance_id_ref=16) as dummy left outer join tags on dummy.symbol_file_ref_id=tags.tag_file_ref_id where symbol_id=472055 " +
				"and symbol_instance_id_ref=16",
			resultHead:   []string{"symbol_id", "symbol_name", "subsys_name", "file_name"},
			resultValues: [][]driver.Value{{"472055", "__x64_sys_getpid", "", "kernel/sys.c"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR: "select (select symbol_type from symbols where symbol_name='__x64_sys_getpid' and symbol_instance_id_ref=16) as type, subsys_name from " +
				"(select count(*) as cnt, subsys_name from tags where subsys_name in (select subsys_name from symbols, tags where " +
				"symbols.symbol_file_ref_id=tags.tag_file_ref_id and symbols.symbol_name='__x64_sys_getpid' and symbols.symbol_instance_id_ref=16) " +
				"group by subsys_name order by cnt desc) as tbl",
			resultHead:   []string{"type", "subsys_name"},
			resultValues: nil,
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR:   "select caller, callee, source_line, ref_addr from xrefs where caller = 472055 and xref_instance_id_ref = 16",
			resultHead: []string{"caller", "callee", "source_line", "ref_addr"},
			resultValues: [][]driver.Value{{"472055", "501994", "kernel/sys.c:892", "0xffffffff81077570"},
				{"472055", "472243", "kernel/sys.c:893", "0xffffffff81077589"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR: "select symbol_id, symbol_name, subsys_name, file_name from (select * from symbols, files where symbols.symbol_file_ref_id=files.file_id " +
				"and symbols.symbol_instance_id_ref=16) as dummy left outer join tags on dummy.symbol_file_ref_id=tags.tag_file_ref_id where " +
				"symbol_id=501994 and symbol_instance_id_ref=16",
			resultHead:   []string{"symbol_id", "symbol_name", "subsys_name", "file_name"},
			resultValues: [][]driver.Value{{"501994", "__fentry__", "X86 ARCHITECTURE (32-BIT AND 64-BIT)", "arch/x86/kernel/ftrace_64.S"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR: "select symbol_id, symbol_name, subsys_name, file_name from (select * from symbols, files where symbols.symbol_file_ref_id=files.file_id " +
				"and symbols.symbol_instance_id_ref=16) as dummy left outer join tags on dummy.symbol_file_ref_id=tags.tag_file_ref_id where " +
				"symbol_id=472243 and symbol_instance_id_ref=16",
			resultHead:   []string{"symbol_id", "symbol_name", "subsys_name", "file_name"},
			resultValues: [][]driver.Value{{"472243", "__task_pid_nr_ns", "", "kernel/pid.c"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR: "select (select symbol_type from symbols where symbol_name='__task_pid_nr_ns' and symbol_instance_id_ref=16) as type, subsys_name from " +
				"(select count(*) as cnt, subsys_name from tags where subsys_name in (select subsys_name from symbols, tags where " +
				"symbols.symbol_file_ref_id=tags.tag_file_ref_id and symbols.symbol_name='__task_pid_nr_ns' and symbols.symbol_instance_id_ref=16) " +
				"group by subsys_name order by cnt desc) as tbl",
			resultHead:   []string{"type", "subsys_name"},
			resultValues: nil,
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR:   "select caller, callee, source_line, ref_addr from xrefs where caller = 472243 and xref_instance_id_ref = 16",
			resultHead: []string{"caller", "callee", "source_line", "ref_addr"},
			resultValues: [][]driver.Value{{"472243", "501994", "kernel/pid.c:427", "0xffffffff810824e0"},
				{"472243", "473674", "kernel/pid.c:430", "0xffffffff810824f7"},
				{"472243", "473716", "kernel/pid.c:435", "0xffffffff81082540"},
				{"472243", "473716", "kernel/pid.c:435", "0xffffffff81082584"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR: "select symbol_id, symbol_name, subsys_name, file_name from (select * from symbols, files where symbols.symbol_file_ref_id=files.file_id " +
				"and symbols.symbol_instance_id_ref=16) as dummy left outer join tags on dummy.symbol_file_ref_id=tags.tag_file_ref_id where " +
				"symbol_id=473674 and symbol_instance_id_ref=16",
			resultHead:   []string{"symbol_id", "symbol_name", "subsys_name", "file_name"},
			resultValues: [][]driver.Value{{"473674", "__rcu_read_lock", "READ-COPY UPDATE (RCU)", "kernel/rcu/tree_plugin.h"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR: "select symbol_id, symbol_name, subsys_name, file_name from (select * from symbols, files where symbols.symbol_file_ref_id=files.file_id " +
				"and symbols.symbol_instance_id_ref=16) as dummy left outer join tags on dummy.symbol_file_ref_id=tags.tag_file_ref_id where " +
				"symbol_id=473716 and symbol_instance_id_ref=16",
			resultHead:   []string{"symbol_id", "symbol_name", "subsys_name", "file_name"},
			resultValues: [][]driver.Value{{"473716", "__rcu_read_unlock", "READ-COPY UPDATE (RCU)", "kernel/rcu/tree_plugin.h"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR: "select (select symbol_type from symbols where symbol_name='__rcu_read_lock' and symbol_instance_id_ref=16) as type, subsys_name from " +
				"(select count(*) as cnt, subsys_name from tags where subsys_name in (select subsys_name from symbols, tags where " +
				"symbols.symbol_file_ref_id=tags.tag_file_ref_id and symbols.symbol_name='__rcu_read_lock' and symbols.symbol_instance_id_ref=16) " +
				"group by subsys_name order by cnt desc) as tbl",
			resultHead:   []string{"type", "subsys_name"},
			resultValues: [][]driver.Value{{"direct", "READ-COPY UPDATE (RCU)"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR: "select (select symbol_type from symbols where symbol_name='__rcu_read_unlock' and symbol_instance_id_ref=16) as type, subsys_name from " +
				"(select count(*) as cnt, subsys_name from tags where subsys_name in (select subsys_name from symbols, tags where " +
				"symbols.symbol_file_ref_id=tags.tag_file_ref_id and symbols.symbol_name='__rcu_read_unlock' and symbols.symbol_instance_id_ref=16) " +
				"group by subsys_name order by cnt desc) as tbl",
			resultHead:   []string{"type", "subsys_name"},
			resultValues: [][]driver.Value{{"direct", "READ-COPY UPDATE (RCU)"}},
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR:     "select subsys_name from tags where tag_file_ref_id= (select symbol_file_ref_id from symbols where symbol_id=472055)",
			resultHead:   []string{"subsys_name"},
			resultValues: nil,
		})

		queryTestSerie = append(queryTestSerie, mockQueries{
			querySTR:     "select subsys_name from tags where tag_file_ref_id= (select symbol_file_ref_id from symbols where symbol_id=472243)",
			resultHead:   []string{"subsys_name"},
			resultValues: nil,
		})

		for _, a := range queryTestSerie {
			rows := sqlmock.NewRows(a.resultHead)
			for _, v := range a.resultValues {
				rows.AddRow(v...)
			}
			mock.ExpectQuery(a.querySTR).WillReturnRows(rows)
		}
		mock.ExpectCommit()
		dok.cache.entries = map[int]entry{}
		dok.cache.successors = map[int][]entry{}
		dok.cache.subSys = map[string]string{}
		testConfig := configuration{
			DBDriver:       "postgres",
			DBDSN:          "host=dbs.hqhome163.com port=5432 user=alessandro password=<password> dbname=kernel_bin sslmode=disable",
			Symbol:         "__x64_sys_getpid",
			Instance:       16,
			Mode:           printAll,
			ExcludedBefore: []string{"__fentry__", "__stack_chk_fail"},
			ExcludedAfter:  []string{"^kfree$", "^_raw_spin_lock$", "^_raw_spin_unlock$", "^panic$", "^call_rcu$", "^__call_rcu$", "__rcu_read_unlock", "__rcu_read_lock", "path_openat"},
			TargetSubsys:   []string{},
			MaxDepth:       0, //0: no limit
			Jout:           "graphOnly",
			Graphviz:       oText,
			cmdlineNeeds:   map[string]bool{},
		}

		dot, err := generateOutput(dok, &testConfig)
		It("Should return syntax correct json with no error", func() {
			Expect(err).To(BeNil())
			Expect(dot).To(Equal(expectedDot))
		})
	})

	Describe("main", func() {
		// TODO: `nav.main` refactoring needed
	})
})
