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

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Sql connection configuration.
type connectToken struct {
	DBDriver string
	DBDSN    string
}

type Cache struct {
	successors map[int][]entry
	entries    map[int]entry
	subSys     map[string]string
}

type SqlDB struct {
	db    *sql.DB
	cache Cache
}

// Connects the target db and returns the handle.
func (d *SqlDB) init(arg interface{}) (err error) {
	t, ok := arg.(*connectToken)
	if !ok {
		var ok1 bool
		d.db, ok1 = arg.(*sql.DB)
		if !ok1 {
			return errors.New("invalid type")
		}
	}
	if ok {
		d.db, err = sql.Open(t.DBDriver, t.DBDSN)
	}
	if err == nil {
		d.cache.successors = make(map[int][]entry)
		d.cache.entries = make(map[int]entry)
		d.cache.subSys = make(map[string]string)
	}
	return err
}

func (d *SqlDB) GetExploredSubsystemByName(subs string) string {
	debugIOPrintln("input subs=", subs)
	debugIOPrintln("output =", subs)
	return d.cache.subSys[subs]
}

// Returns function details from a given id.
func (d *SqlDB) getEntryById(symbolId int, instance int) (entry, error) {
	var e entry
	var s sql.NullString

	debugIOPrintf("input symbolId=%d, instance=%d\n", symbolId, instance)
	if e, ok := d.cache.entries[symbolId]; ok {
		debugIOPrintf("output entry=%+v, error=%s\n", e, "nil")
		return e, nil
	}

	query := "select symbol_id, symbol_name, subsys_name, file_name from " +
		"(select * from symbols, files where symbols.symbol_file_ref_id=files.file_id and symbols.symbol_instance_id_ref=%[2]d) as dummy " +
		"left outer join tags on dummy.symbol_file_ref_id=tags.tag_file_ref_id where symbol_id=%[1]d and symbol_instance_id_ref=%[2]d"
	query = fmt.Sprintf(query, symbolId, instance)
	debugQueryPrintln(query)
	rows, err := d.db.Query(query)
	if err != nil {
		debugIOPrintf("output entry=%+v, error=%s\n", entry{}, err)
		return entry{}, err
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&e.symId, &e.symbol, &s, &e.fn); err != nil {
			debugIOPrintf("output entry=%+v, error=%s\n", entry{}, err)
			return e, err
		}
		if s.Valid {
			e.subsys = append(e.subsys, s.String)
		}
	}
	if err = rows.Err(); err != nil {
		debugIOPrintf("output entry=%+v, error=%s\n", entry{}, err)
		return e, err
	}
	d.cache.entries[symbolId] = e
	debugIOPrintf("output entry=%+v, error=%s\n", e, "nil")
	return e, nil
}

// Returns the list of successors (called function) for a given function.
func (d *SqlDB) getSuccessorsById(symbolId int, instance int) ([]entry, error) {
	var e edge
	var res []entry

	debugIOPrintf("input symbolId=%d, instance=%d\n", symbolId, instance)
	if res, ok := d.cache.successors[symbolId]; ok {
		debugIOPrintf("output []entry=%+v, error=%s\n", res, "nil")
		return res, nil
	}

	query := "select caller, callee, source_line, ref_addr from xrefs where caller = %[1]d and xref_instance_id_ref = %[2]d"
	query = fmt.Sprintf(query, symbolId, instance)
	debugQueryPrintln(query)
	rows, err := d.db.Query(query)
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&e.caller, &e.callee, &e.sourceRef, &e.addressRef); err != nil {
			debugIOPrintf("output []entry=%+v, error=%s\n", nil, err)
			return nil, err
		}
		successor, _ := d.getEntryById(e.callee, instance)
		successor.sourceRef = e.sourceRef
		successor.addressRef = e.addressRef
		res = append(res, successor)
	}
	if err = rows.Err(); err != nil {
		debugIOPrintf("output []entry=%+v, error=%s\n", nil, err)
		return nil, err
	}
	d.cache.successors[symbolId] = res
	debugIOPrintf("output []entry=%+v, error=%s\n", res, "nil")
	return res, nil
}

// Given a function returns the lager subsystem it belongs.
func (d *SqlDB) getSubsysFromSymbolName(symbol string, instance int) (string, error) {
	var ty, sub string

	debugIOPrintf("input symbol=%s, instance=%d\n", symbol, instance)
	if res, ok := d.cache.subSys[symbol]; ok {
		debugIOPrintf("output  string=%s, error=%s\n", res, "nil")
		return res, nil
	}
	query := "select (select symbol_type from symbols where symbol_name='%[1]s' and symbol_instance_id_ref=%[2]d) as type, subsys_name from " +
		"(select count(*) as cnt, subsys_name from tags where subsys_name in (select subsys_name from symbols, " +
		"tags where symbols.symbol_file_ref_id=tags.tag_file_ref_id and symbols.symbol_name='%[1]s' and symbols.symbol_instance_id_ref=%[2]d) " +
		"group by subsys_name order by cnt desc) as tbl"

	query = fmt.Sprintf(query, symbol, instance)
	debugQueryPrintln(query)
	rows, err := d.db.Query(query)
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&ty, &sub); err != nil {
			debugIOPrintf("output  string=%s, error=%s\n", "", err)
			return "", err
		}
	}

	if err = rows.Err(); err != nil {
		debugIOPrintf("output  string=%s, error=%s\n", "", err)
		return "", err
	}

	if sub == "" {
		debugIOPrintf("output  string=%s, error=%s\n", "", "nil")
		return "", nil
	}

	if ty == "indirect" {
		sub = ty
	}
	d.cache.subSys[symbol] = sub
	debugIOPrintf("output  string=%s, error=%s\n", sub, "nil")
	return sub, nil
}

// Returns the id of a given function name.
func (d *SqlDB) sym2num(symb string, instance int) (int, error) {
	var res = -1
	var cnt = 0

	debugIOPrintf("input symbol=%s, instance=%d\n", symb, instance)
	query := "select symbol_id from symbols where symbols.symbol_name='%[1]s' and symbols.symbol_instance_id_ref=%[2]d"
	query = fmt.Sprintf(query, symb, instance)
	debugQueryPrintln(query)
	rows, err := d.db.Query(query)
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := rows.Close()
		if err == nil {
			err = closeErr
		}
	}()

	for rows.Next() {
		cnt++
		if err := rows.Scan(&res); err != nil {
			debugIOPrintf("output int=%d, error=%s\n", res, err)
			return res, err
		}
	}

	if err = rows.Err(); err != nil {
		debugIOPrintf("output int=%d, error=%s\n", -1, err)
		return -1, err
	}

	if cnt != 1 {
		return res, errors.New("duplicate ID in the DB")
	}
	debugIOPrintf("output int=%d, error=%s\n", res, "nil")
	return res, nil
}

// Returns the subsystem list associated with a given function name.
func (d *SqlDB) symbSubsys(symblist []int, instance int) (string, error) {
	var out string
	var res string

	debugIOPrintf("input symblist=%+v, instance=%d\n", symblist, instance)
	for _, symbid := range symblist {
		// Resolve symb.
		symb, err := d.getEntryById(symbid, instance)
		if err != nil {
			return "", fmt.Errorf("symbSubsys::getEntryById error: %s", err)
		}
		out += fmt.Sprintf("{\"FuncName\":\"%s\", \"subsystems\":[", symb.symbol)
		query := fmt.Sprintf("select subsys_name from tags where tag_file_ref_id= (select symbol_file_ref_id from symbols where symbol_id=%d)", symbid)
		debugQueryPrintln(query)
		rows, err := d.db.Query(query)
		if err != nil {
			err = errors.New("symbSubsys: query failed")
			debugIOPrintf("output string=%s, error=%s\n", "", err)
			return "", err
		}

		defer func() {
			closeErr := rows.Close()
			if err == nil {
				err = closeErr
			}
		}()

		for rows.Next() {
			if err := rows.Scan(&res); err != nil {
				err = errors.New("symbSubsys: error while scan query rows")
				debugIOPrintf("output string=%s, error=%s\n", "", err)
				return "", err
			}
			out += fmt.Sprintf("\"%s\",", res)
		}
		out = strings.TrimSuffix(out, ",") + "]},"
	}
	out = strings.TrimSuffix(out, ",")
	debugIOPrintf("output string=%s, error=%s\n", out, "nil")
	return out, nil
}
