/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"database/sql"
	"errors"
	"fmt"
	addr2line "github.com/elazarl/addr2line"
	"path/filepath"
	"strings"
	"sync"
)

var Query_fmts = [...][DBTYPE_Last]string{
	{
	"insert into instances (version_string, note) values ('%d.%d.%d%s', '%s'); SELECT currval('instances_instance_id_seq');",
	"insert into instances (version_string, note) values ('%d.%d.%d%s', '%s'); select LAST_INSERT_ID();",
	"insert into instances (version_string, note) values ('%d.%d.%d%s', '%s') returning instance_id; ",
	"insert into instances (version_string, note) values ('%d.%d.%d%s', '%s') returning instance_id; ",
	},
	{
	"insert into configs (config_symbol, config_value, config_instance_id_ref) values ('%s', '%s', %d);",
	"insert into configs (config_symbol, config_value, config_instance_id_ref) values ('%s', '%s', %d);",
	"insert into configs (config_symbol, config_value, config_instance_id_ref) values ('%s', '%s', %d);",
	"insert into configs (config_symbol, config_value, config_instance_id_ref) values ('%s', '%s', %d);",
	},
	{
	"insert into files (file_name, file_instance_id_ref) select 'NoFile',%d;",
	"insert into files (file_name, file_instance_id_ref) select 'NoFile',%d;",
	"insert into files (file_name, file_instance_id_ref) select 'NoFile',%d;",
	"insert into files (file_name, file_instance_id_ref) select 'NoFile',%d;",
	},
	{
	"insert into symbols (symbol_name,symbol_address,symbol_type,symbol_file_ref_id,symbol_instance_id_ref) " + "select (select 'Indirect call'), '0x00000000', 'indirect', (select file_id from files where file_name ='NoFile' and file_instance_id_ref=%[1]d), %[1]d;",
	"insert into symbols (symbol_name,symbol_address,symbol_type,symbol_file_ref_id,symbol_instance_id_ref) " + "select (select 'Indirect call'), '0x00000000', 'indirect', (select file_id from files where file_name ='NoFile' and file_instance_id_ref=%[1]d), %[1]d;",
	"insert into symbols (symbol_name,symbol_address,symbol_type,symbol_file_ref_id,symbol_instance_id_ref) " + "select (select 'Indirect call'), '0x00000000', 'indirect', (select file_id from files where file_name ='NoFile' and file_instance_id_ref=%[1]d), %[1]d;",
	"insert into symbols (symbol_name,symbol_address,symbol_type,symbol_file_ref_id,symbol_instance_id_ref) " + "select (select 'Indirect call'), '0x00000000', 'indirect', (select file_id from files where file_name ='NoFile' and file_instance_id_ref=%[1]d), %[1]d;",
	},
	{
	"insert into tags (subsys_name, tag_file_ref_id, tag_instance_id_ref) select (select 'Indirect Calls'), " + "(select file_id from files where file_name='NoFile' and file_instance_id_ref=%[1]d), %[1]d;",
	"insert into tags (subsys_name, tag_file_ref_id, tag_instance_id_ref) select (select 'Indirect Calls'), " + "(select file_id from files where file_name='NoFile' and file_instance_id_ref=%[1]d), %[1]d;",
	"insert into tags (subsys_name, tag_file_ref_id, tag_instance_id_ref) select (select 'Indirect Calls'), " + "(select file_id from files where file_name='NoFile' and file_instance_id_ref=%[1]d), %[1]d;",
	"insert into tags (subsys_name, tag_file_ref_id, tag_instance_id_ref) select (select 'Indirect Calls'), " + "(select file_id from files where file_name='NoFile' and file_instance_id_ref=%[1]d), %[1]d;",
	},
	{
	"insert into files (file_name, file_instance_id_ref) Select '%%[1]s', %[1]d Where not exists " + "(select * from files where file_name='%%[1]s' and file_instance_id_ref=%[1]d);" + "insert into symbols (symbol_name, symbol_address, symbol_type, symbol_file_ref_id, symbol_instance_id_ref) " + "select '%[2]s', '%[3]s', '%[4]s', (select file_id from files where file_name='%%[1]s' and file_instance_id_ref=%[1]d), %[1]d;",
	"insert into files (file_name, file_instance_id_ref) Select '%%[1]s', %[1]d Where not exists " + "(select * from files where file_name='%%[1]s' and file_instance_id_ref=%[1]d);" + "insert into symbols (symbol_name, symbol_address, symbol_type, symbol_file_ref_id, symbol_instance_id_ref) " + "select '%[2]s', '%[3]s', '%[4]s', (select file_id from files where file_name='%%[1]s' and file_instance_id_ref=%[1]d), %[1]d;",
	"insert into files (file_name, file_instance_id_ref) Select '%%[1]s', %[1]d Where not exists " + "(select * from files where file_name='%%[1]s' and file_instance_id_ref=%[1]d);" + "insert into symbols (symbol_name, symbol_address, symbol_type, symbol_file_ref_id, symbol_instance_id_ref) " + "select '%[2]s', '%[3]s', '%[4]s', (select file_id from files where file_name='%%[1]s' and file_instance_id_ref=%[1]d), %[1]d;",
	"insert into files (file_name, file_instance_id_ref) Select '%%[1]s', %[1]d Where not exists " + "(select * from files where file_name='%%[1]s' and file_instance_id_ref=%[1]d);" + "insert into symbols (symbol_name, symbol_address, symbol_type, symbol_file_ref_id, symbol_instance_id_ref) " + "select '%[2]s', '%[3]s', '%[4]s', (select file_id from files where file_name='%%[1]s' and file_instance_id_ref=%[1]d), %[1]d;",
	},
	{
	"insert into xrefs (caller, callee, ref_addr, source_line, xref_instance_id_ref) " + "select (Select symbol_id from symbols where symbol_address ='0x%08[1]x' and symbol_instance_id_ref=%[3]d), " + "(Select symbol_id from symbols where symbol_address ='0x%08[2]x' and symbol_instance_id_ref=%[3]d limit 1), " + "'0x%08[5]x', " + "'%[4]s', " + "%[3]d;",
	"insert into xrefs (caller, callee, ref_addr, source_line, xref_instance_id_ref) " + "select (Select symbol_id from symbols where symbol_address ='0x%08[1]x' and symbol_instance_id_ref=%[3]d), " + "(Select symbol_id from symbols where symbol_address ='0x%08[2]x' and symbol_instance_id_ref=%[3]d limit 1), " + "'0x%08[5]x', " + "'%[4]s', " + "%[3]d;",
	"insert into xrefs (caller, callee, ref_addr, source_line, xref_instance_id_ref) " + "select (Select symbol_id from symbols where symbol_address ='0x%08[1]x' and symbol_instance_id_ref=%[3]d), " + "(Select symbol_id from symbols where symbol_address ='0x%08[2]x' and symbol_instance_id_ref=%[3]d limit 1), " + "'0x%08[5]x', " + "'%[4]s', " + "%[3]d;",
	"insert into xrefs (caller, callee, ref_addr, source_line, xref_instance_id_ref) " + "select (Select symbol_id from symbols where symbol_address ='0x%08[1]x' and symbol_instance_id_ref=%[3]d), " + "(Select symbol_id from symbols where symbol_address ='0x%08[2]x' and symbol_instance_id_ref=%[3]d limit 1), " + "'0x%08[5]x', " + "'%[4]s', " + "%[3]d;",
	},
	{
	"insert into tags (subsys_name, tag_file_ref_id, tag_instance_id_ref) select '%%[1]s', " + "(select file_id from files where file_name='%[1]s%%[2]s' and file_instance_id_ref=%%[3]d) as fn_id, %%[3]d " + "WHERE EXISTS ( select file_id from files where file_name='%[1]s%%[2]s' and file_instance_id_ref=%%[3]d);",
	"insert into tags (subsys_name, tag_file_ref_id, tag_instance_id_ref) select '%%[1]s', " + "(select file_id from files where file_name='%[1]s%%[2]s' and file_instance_id_ref=%%[3]d) as fn_id, %%[3]d " + "WHERE EXISTS ( select file_id from files where file_name='%[1]s%%[2]s' and file_instance_id_ref=%%[3]d);",
	"insert into tags (subsys_name, tag_file_ref_id, tag_instance_id_ref) select '%%[1]s', " + "(select file_id from files where file_name='%[1]s%%[2]s' and file_instance_id_ref=%%[3]d) as fn_id, %%[3]d " + "WHERE EXISTS ( select file_id from files where file_name='%[1]s%%[2]s' and file_instance_id_ref=%%[3]d);",
	"insert into tags (subsys_name, tag_file_ref_id, tag_instance_id_ref) select '%%[1]s', " + "(select file_id from files where file_name='%[1]s%%[2]s' and file_instance_id_ref=%%[3]d) as fn_id, %%[3]d " + "WHERE EXISTS ( select file_id from files where file_name='%[1]s%%[2]s' and file_instance_id_ref=%%[3]d);",
	},
}

type Workload_Type int64

// Const values for configuration mode field.
const (
	_ Workload_Type = iota
	GENERATE_QUERY
	EXECUTE_QUERY_ONLY
	GENERATE_QUERY_AND_EXECUTE
	GENERATE_QUERY_AND_EXECUTE_W_A2L
	Workload_Type_Last
)

type Workload struct {
	Addr2ln_offset uint64
	Addr2ln_name   string
	Query_str      string
	Query_args     interface{}
	Workload_type  Workload_Type
}

// Context type
type Context struct {
	a2l         *addr2line.Addr2line
	ch_workload chan Workload
	mu          sync.Mutex
	DB          *sql.DB
	DBT         string
}

// Caches item elements
type Addr2line_items struct {
	Addr      uint64
	File_name string
}

// Commandline handle functions prototype
type ins_f func(*Context, string)

var Test_result []string

func Fake_Insert_data(context *Context, query string) {
	Test_result = append(Test_result)
}

func A2L_resolver__init(fn string, DB_inst *sql.DB, DBT string, test bool) *Context {
	a, err := addr2line.New(fn)
	if err != nil {
		panic(err)
	}
	addresses := make(chan Workload, 16)
	context := &Context{a2l: a, ch_workload: addresses, DB: DB_inst, DBT: DBT}

	if !test {
		go workload(context, Insert_data)
	} else {
		go workload(context, Fake_Insert_data)
	}
	return context
}

func resolve_addr(context *Context, address uint64) string {
	var res string = ""
	context.mu.Lock()
	rs, _ := context.a2l.Resolve(address)
	context.mu.Unlock()
	if len(rs) == 0 {
		res = "NONE"
	}
	for _, a := range rs {
		res = fmt.Sprintf("%s:%d", filepath.Clean(a.File), a.Line)
	}
	return res
}

func workload(context *Context, insert_func ins_f) {
	var e Workload
	var qready string

	for {
		e = <-context.ch_workload
		switch e.Workload_type {
		case GENERATE_QUERY_AND_EXECUTE, EXECUTE_QUERY_ONLY:
			insert_func(context, e.Query_str)
			break
		case GENERATE_QUERY_AND_EXECUTE_W_A2L:
			context.mu.Lock()
			rs, _ := context.a2l.Resolve(e.Addr2ln_offset)
			context.mu.Unlock()
			if len(rs) == 0 {
				qready = fmt.Sprintf(e.Query_str, "NONE")
			}
			for _, a := range rs {
				qready = fmt.Sprintf(e.Query_str, filepath.Clean(a.File))
				if a.Function == strings.ReplaceAll(e.Addr2ln_name, "sym.", "") {
					break
				}
			}
			insert_func(context, qready)
			break
		default:
		}
	}
}

func Generate_Query_Str(Q_WL *Workload, DBT DBtype) error {
	var err error = nil

	switch arg := (*Q_WL).Query_args.(type) {
	case Insert_Instance_Args:
		(*Q_WL).Query_str = fmt.Sprintf(Query_fmts[0][DBT], arg.Version, arg.Patchlevel, arg.Sublevel, arg.Extraversion, arg.Note)
	case Insert_Config_Args:
		(*Q_WL).Query_str = fmt.Sprintf(Query_fmts[1][DBT], arg.Config_key, arg.Config_val, arg.Instance_no)
	case Insert_Files_Ind_Args:
		(*Q_WL).Query_str = fmt.Sprintf(Query_fmts[2][DBT], arg.Id)
	case Insert_Symbols_Ind_Args:
		(*Q_WL).Query_str = fmt.Sprintf(Query_fmts[3][DBT], arg.Id)
	case Insert_Tags_Ind_Args:
		(*Q_WL).Query_str = fmt.Sprintf(Query_fmts[4][DBT], arg.Id)
	case Insert_Symbols_Files_Args:
		(*Q_WL).Query_str = fmt.Sprintf(Query_fmts[5][DBT], arg.Id, arg.Symbol_Name, arg.Symbol_Offset, arg.Symbol_Type)
	case Insert_Xrefs_Args:
		(*Q_WL).Query_str = fmt.Sprintf(Query_fmts[6][DBT], arg.Caller_Offset, arg.Callee_Offset, arg.Id, arg.Source_line, arg.Calling_Offset)
	case Insert_Tags_Args:
		(*Q_WL).Query_str = fmt.Sprintf(Query_fmts[7][DBT], arg.addr2line_prefix)
	default:
		err = errors.New("GENERATE_QUERY: Unknown workload argument")
	}
	return err
}

func query_mgmt(ctx *Context, Q_WL *Workload) error {
	var err error

	switch (*Q_WL).Workload_type {
	case GENERATE_QUERY:
		DBT, err := DBSN2DBtype((*ctx).DBT)
		if err != nil {
			panic(err)
		}
		err = Generate_Query_Str(Q_WL, DBT)
	case EXECUTE_QUERY_ONLY:
		(*ctx).ch_workload <- *Q_WL
	case GENERATE_QUERY_AND_EXECUTE, GENERATE_QUERY_AND_EXECUTE_W_A2L:
		DBT, err := DBSN2DBtype((*ctx).DBT)
		if err != nil {
			panic(err)
		}
		err = Generate_Query_Str(Q_WL, DBT)
		if err == nil {
			(*ctx).ch_workload <- *Q_WL
		}
	default:
		err = errors.New("Unknown workload type")
	}
	return err
}
