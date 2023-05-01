/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/DATA-DOG/go-sqlmock"
)

var _ = Describe("Psql Tests", func() {

	var db *sql.DB
	var mock sqlmock.Sqlmock
	var e entry
	var dok, dko *SqlDB

	BeforeEach(func() {
		e = entry{
			symbol:     "mysymbol",
			sourceRef:  "config.h",
			addressRef: "0",
			subsys:     []string{},
			symId:      1,
		}

		dok = &SqlDB{}
		db, mock, _ = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		dok.db = db
		dok.cache = Cache{}

		dko = &SqlDB{}
		dko.db = nil
		dko.cache = Cache{}
	})

	AfterEach(func() {
		defer GinkgoRecover()
		defer db.Close()
	})

	When("connectDB", func() {
		// TODO: `psql.connectDB` fn refactor needed
	})

	When("getEntryById", func() {
		testQuery := "select symbol_id, symbol_name, subsys_name, file_name from " +
			"(select * from symbols, files where symbols.symbol_file_ref_id=files.file_id and symbols.symbol_instance_id_ref=0) as dummy " +
			"left outer join tags on dummy.symbol_file_ref_id=tags.tag_file_ref_id where symbol_id=0 and symbol_instance_id_ref=0"
		It("Should return a cached result", func() {
			dko.cache.entries = map[int]entry{1: e}
			_entry, err := dko.getEntryById(1, 1)

			Expect(err).To(BeNil())
			Expect(_entry).To(Equal(e))
		})

		It("Should inform the user of an internal database error", func() {
			mock.ExpectExec(testQuery).
				WithArgs(0, 0).
				WillReturnError(fmt.Errorf("myerror"))
			mock.ExpectRollback()

			dok.cache.entries = map[int]entry{}
			_entry, err := dok.getEntryById(0, 0)

			Expect(err).ToNot(BeNil())
			Expect(_entry).To(Equal(entry{}))
		})

		It("Should not return an empty entry without errors", func() {
			rows := sqlmock.NewRows([]string{
				"id",
				"symbol",
				"fn",
				"sourceRef",
				"addressRef",
				"subsys",
				"symId",
			})

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.entries = map[int]entry{}
			_, err := dok.getEntryById(0, 0)

			Expect(err).To(BeNil())
		})

		It("Should fail for a row scan", func() {
			rows := sqlmock.NewRows([]string{
				"symbol",
				"fn",
				"subsys",
				"symId",
			})
			rows.AddRow(nil, nil, nil, nil)

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.entries = map[int]entry{}
			_, err := dok.getEntryById(0, 0)

			Expect(err).ToNot(BeNil())
		})

		It("Should fail for a row error", func() {
			rows := sqlmock.NewRows([]string{
				"symbol",
				"fn",
				"subsys",
				"symId",
			})
			rows.AddRow("0", "1", "2", 3)
			rows.RowError(0, fmt.Errorf("row error"))

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.entries = map[int]entry{}
			_, err := dok.getEntryById(0, 0)

			Expect(err).ToNot(BeNil())
		})

		It("Should find and return an entry", func() {
			rows := sqlmock.NewRows([]string{
				"symbol",
				"fn",
				"subsys",
				"symId",
			})
			rows.AddRow("0", "1", "2", 3)

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.entries = map[int]entry{}
			_entry, err := dok.getEntryById(0, 0)

			expectedEntry := entry{
				symbol:     "1",
				fn:         "3",
				sourceRef:  "",
				addressRef: "",
				subsys:     []string{"2"},
				symId:      0,
			}
			Expect(err).To(BeNil())
			Expect(_entry).To(Equal(expectedEntry))
		})
	})

	When("getSuccessorsById", func() {

		testQuery := "select caller, callee, source_line, ref_addr from xrefs where caller = 0 and xref_instance_id_ref = 0"

		It("Should return cached sucessors", func() {

			dko.cache.successors = map[int][]entry{1: {e}}
			entries, err := dko.getSuccessorsById(1, 1)

			Expect(err).To(BeNil())
			Expect(entries).To(Equal([]entry{e}))
		})

		It("Should panic because of a query error", func() {
			mock.ExpectQuery(testQuery).
				WithArgs(0, 0).
				WillReturnError(fmt.Errorf("myerror"))
			mock.ExpectRollback()

			dok.cache = Cache{}
			Expect(func() { dok.getSuccessorsById(0, 0) }).To(Panic())
		})

		It("Should not fail in case no records are found", func() {
			rows := sqlmock.NewRows([]string{
				"symbol",
				"fn",
				"subsys",
				"symId",
			})

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.successors = map[int][]entry{}
			entries, err := dok.getSuccessorsById(0, 0)

			Expect(err).To(BeNil())
			Expect(entries).To(BeNil())
		})

		It("Should fail for a row scan", func() {
			rows := sqlmock.NewRows([]string{
				"symbol",
				"fn",
				"subsys",
				"symId",
			})
			rows.AddRow(nil, nil, nil, nil)

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.successors = map[int][]entry{}
			_, err := dok.getSuccessorsById(0, 0)

			Expect(err).ToNot(BeNil())
		})

		It("Should fail for a row error", func() {
			rows := sqlmock.NewRows([]string{
				"symbol",
				"fn",
				"subsys",
				"symId",
			})
			rows.AddRow("0", "1", "2", 3)
			rows.RowError(0, fmt.Errorf("row error"))

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.successors = map[int][]entry{}
			entries, err := dok.getSuccessorsById(0, 0)

			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(fmt.Errorf("row error")))
			Expect(entries).To(BeNil())
		})

		It("Should return a list of entries", func() {
			rows := sqlmock.NewRows([]string{
				"symbol",
				"fn",
				"subsys",
				"symId",
			})
			rows.AddRow("0", "1", "2", 3)

			mock.
				ExpectQuery(testQuery).
				//WithArgs(0, 0).
				WillReturnRows(rows)
			mock.ExpectCommit()

			expected := entry{symbol: "", fn: "", sourceRef: "2", addressRef: "3", subsys: nil, symId: 0}

			dok.cache.successors = map[int][]entry{}
			entries, err := dok.getSuccessorsById(0, 0)

			Expect(err).To(BeNil())
			Expect(entries).To(Equal([]entry{expected}))
		})
	})

	When("getSubsysFromSymbolName", func() {
		testQuery := "select (select symbol_type from symbols where symbol_name='mysym_key' and symbol_instance_id_ref=0) as type, subsys_name from " +
			"(select count(*) as cnt, subsys_name from tags where subsys_name in (select subsys_name from symbols, " +
			"tags where symbols.symbol_file_ref_id=tags.tag_file_ref_id and symbols.symbol_name='mysym_key' and symbols.symbol_instance_id_ref=0) " +
			"group by subsys_name order by cnt desc) as tbl"

		It("Should return a cached symbol", func() {
			dko.cache.subSys = map[string]string{"mysym_key": "mysym_val"}
			sym, err := dko.getSubsysFromSymbolName("mysym_key", 0)

			Expect(err).To(BeNil())
			Expect(sym).To(Equal("mysym_val"))
		})

		It("Should panic because of a query error", func() {
			mock.ExpectQuery(testQuery).
				WithArgs("mysym_key", 0).
				WillReturnError(fmt.Errorf("myerror"))
			mock.ExpectRollback()
			dok.cache.subSys = map[string]string{}

			Expect(func() { dok.getSubsysFromSymbolName("mysym_key", 0) }).To(Panic())
		})

		It("Should return nil in case of no rows", func() {
			rows := sqlmock.NewRows([]string{
				"ty",
				"sub",
			})

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.subSys = map[string]string{}
			sym, err := dok.getSubsysFromSymbolName("mysym_key", 0)

			Expect(err).To(BeNil())
			Expect(sym).To(Equal(""))

		})

		It("Should return a row scan error", func() {
			rows := sqlmock.NewRows([]string{
				"ty",
				"sub",
			})
			rows.AddRow("direct", nil)
			mock.
				ExpectQuery(testQuery).
				//WithArgs("subsys", 0).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.subSys = map[string]string{}
			sym, err := dok.getSubsysFromSymbolName("mysym_key", 0)

			Expect(err).ToNot(BeNil())
			Expect(sym).To(Equal(""))
		})

		It("Should return a row error", func() {
			rows := sqlmock.NewRows([]string{
				"ty",
				"sub",
			})
			rows.AddRow("direct", "subsys")
			rows.RowError(0, fmt.Errorf("row error"))

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.subSys = map[string]string{}
			sym, err := dok.getSubsysFromSymbolName("mysym_key", 0)

			Expect(err).ToNot(BeNil())
			Expect(sym).To(Equal(""))
		})

		It("Should find and return a subsystem name", func() {
			rows := sqlmock.NewRows([]string{
				"ty",
				"sub",
			})
			rows.AddRow("direct", "subsys")

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.subSys = map[string]string{}
			sym, err := dok.getSubsysFromSymbolName("mysym_key", 0)

			Expect(err).To(BeNil())
			Expect(sym).To(Equal("subsys"))
			Expect(len(dok.cache.subSys)).To(Equal(1))
			Expect(dok.cache.subSys["mysym_key"]).To(Equal("subsys"))
		})

		It("Should find and return a subsystem name with indirect type", func() {
			rows := sqlmock.NewRows([]string{
				"ty",
				"sub",
			})
			rows.AddRow("indirect", "subsys")

			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			dok.cache.subSys = map[string]string{}
			sym, err := dok.getSubsysFromSymbolName("mysym_key", 0)

			Expect(err).To(BeNil())
			Expect(sym).To(Equal("indirect"))
			Expect(len(dok.cache.subSys)).To(Equal(1))
			Expect(dok.cache.subSys["mysym_key"]).To(Equal("indirect"))
		})
	})

	When("sym2num", func() {
		testQuery := "select symbol_id from symbols where symbols.symbol_name='mysym_key' and symbols.symbol_instance_id_ref=0"

		It("Should panic because of a query error", func() {
			mock.ExpectQuery(testQuery).
				WithArgs("mysim", 0).
				WillReturnError(fmt.Errorf("myerror"))
			mock.ExpectRollback()

			Expect(func() { dok.sym2num("mysym_key", 0) }).To(Panic())
		})

		It("Should return a row scan error", func() {
			rows := sqlmock.NewRows([]string{
				"res",
			})
			rows.AddRow(nil)
			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			symid, err := dok.sym2num("mysym_key", 0)

			Expect(err).ToNot(BeNil())
			Expect(symid).To(Equal(-1))
		})

		It("Should return a row error", func() {
			rows := sqlmock.NewRows([]string{
				"res",
			})
			rows.AddRow(1)
			rows.RowError(0, fmt.Errorf("sym2num row error"))
			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			symid, err := dok.sym2num("mysym_key", 0)

			Expect(err).ToNot(BeNil())
			Expect(symid).To(Equal(-1))
		})

		It("Should not fail with no rows", func() {
			rows := sqlmock.NewRows([]string{
				"res",
			})
			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			symid, err := dok.sym2num("mysym_key", 0)

			Expect(err).ToNot(BeNil())
			Expect(symid).To(Equal(-1))
		})

		It("Should fail if more than one record using the same id", func() {
			rows := sqlmock.NewRows([]string{
				"res",
			})
			rows.AddRow(1)
			rows.AddRow(1)
			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			symid, err := dok.sym2num("mysym_key", 0)

			Expect(err).ToNot(BeNil())
			Expect(symid).To(Equal(1))
		})

		It("Should return the symbol id", func() {
			rows := sqlmock.NewRows([]string{
				"res",
			})
			rows.AddRow(42)
			mock.
				ExpectQuery(testQuery).
				WillReturnRows(rows)
			mock.ExpectCommit()

			symid, err := dok.sym2num("mysym_key", 0)

			Expect(err).To(BeNil())
			Expect(symid).To(Equal(42))
		})
	})

	Describe("symbSubsys", func() {
		var symList []int
		var instance int
		entryIdTestQuery := "select symbol_id, symbol_name, subsys_name, file_name from " +
			"(select * from symbols, files where symbols.symbol_file_ref_id=files.file_id and symbols.symbol_instance_id_ref=0) as dummy " +
			"left outer join tags on dummy.symbol_file_ref_id=tags.tag_file_ref_id where symbol_id=0 and symbol_instance_id_ref=0"
		testQuery := "select subsys_name from tags where tag_file_ref_id= (select symbol_file_ref_id from symbols where symbol_id=0)"
		_getEntryById := func(commit bool) {
			entryRows := sqlmock.NewRows([]string{
				"symbol",
				"fn",
				"subsys",
				"symId",
			})
			entryRows.AddRow("0", "1", "2", 0)

			mock.ExpectQuery(entryIdTestQuery).
				WillReturnRows(entryRows)
			if commit {
				mock.ExpectCommit()
			}
		}

		BeforeEach(func() {
			symList = []int{0}
			instance = 0
			dok.cache.entries = map[int]entry{}
		})

		When("Using an empty symbol list", func() {
			It("Should neither fail or return data", func() {
				out, err := dok.symbSubsys([]int{}, instance)

				Expect(err).To(BeNil())
				Expect(out).To(Equal(""))
			})
		})

		When("A getEntryById error happens", func() {
			It("Should return an error if getEntryById fails", func() {
				mock.ExpectExec(entryIdTestQuery).
					//WithArgs(0, 0).
					WillReturnError(fmt.Errorf("getEntryById query error"))
				mock.ExpectRollback()

				out, err := dok.symbSubsys(symList, instance)
				isErrMatched := strings.Contains(
					fmt.Sprint(err),
					"symbSubsys::getEntryById error")

				Expect(err).ToNot(BeNil())
				Expect(isErrMatched).To(BeTrue())
				Expect(out).To(Equal(""))
			})
		})

		When("A database error happens", func() {
			It("Should fail in case of a db.Query error", func() {
				_getEntryById(true)

				mock.ExpectQuery(testQuery).
					//WithArgs(0).
					WillReturnError(errors.New("symbSubsys db query error"))
				mock.ExpectRollback()

				out, err := dok.symbSubsys(symList, instance)

				Expect(err).ToNot(BeNil())
				Expect(fmt.Sprint(err)).To(Equal("symbSubsys: query failed"))
				Expect(out).To(Equal(""))
			})

			It("Should fail in case of a db.Scan error", func() {
				_getEntryById(false)

				rows := sqlmock.NewRows([]string{
					"subsys",
				})
				rows.AddRow(nil)

				mock.
					ExpectQuery(testQuery).
					WillReturnRows(rows)
				mock.ExpectCommit()

				out, err := dok.symbSubsys(symList, instance)

				Expect(err).ToNot(BeNil())
				Expect(fmt.Sprint(err)).To(Equal("symbSubsys: error while scan query rows"))
				Expect(out).To(Equal(""))
			})
		})

		When("It succeeds", func() {
			It("Should return the function name and its subsystems", func() {
				_getEntryById(false)

				rows := sqlmock.NewRows([]string{
					"subsys",
				})
				rows.AddRow("mock")

				mock.
					ExpectQuery(testQuery).
					WillReturnRows(rows)
				mock.ExpectCommit()

				out, err := dok.symbSubsys(symList, instance)

				Expect(err).To(BeNil())
				Expect(out).To(Equal("{\"FuncName\":\"1\", \"subsystems\":[\"mock\"]}"))
			})
		})
	})
})
