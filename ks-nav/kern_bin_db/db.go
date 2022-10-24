/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"database/sql"
	"fmt"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)


type DBtype int64

// Const values for DB type.
const (
        PSQL DBtype = iota
        MYSQL
        MARIADB
        SQLITE3
        DBTYPE_Last
)

// Sql connection configuration
type Connect_token struct {
	DBDriver string
	DBDSN    string
}

type Insert_Instance_Args struct {
	Version      int64
	Patchlevel   int64
	Sublevel     int64
	Extraversion string
	Note         string
}

type Insert_Config_Args struct {
	Config_key  string
	Config_val  string
	Instance_no int
}
type Insert_Files_Ind_Args struct {
	Id int
}
type Insert_Symbols_Ind_Args struct {
	Id int
}
type Insert_Tags_Ind_Args struct {
	Id int
}

type Insert_Symbols_Files_Args struct {
	Id            int
	Symbol_Name   string
	Symbol_Offset string
	Symbol_Type   string
}

type Insert_Xrefs_Args struct {
	Caller_Offset  uint64
	Calling_Offset uint64
	Callee_Offset  uint64
	Id             int
	Source_line    string
}

type Insert_Tags_Args struct {
	addr2line_prefix string
}

// Connects the target db and returns the handle
func Connect_db(t *Connect_token) *sql.DB {
	db, err := sql.Open((*t).DBDriver, (*t).DBDSN)
	fmt.Println((*t).DBDriver, (*t).DBDSN)
	if err != nil {
		panic(err)
	}
	return db
}

// Executes insert queries
func Insert_data(context *Context, query string) {

	_, err := (*context).DB.Exec(query)
	if err != nil {
		fmt.Println("##################################################")
		fmt.Println(query)
		fmt.Println("##################################################")
		panic(err)
	}
}

// Executes insert query for instance table and returns the id allocated
func Insert_datawID(context *Context, query string) int {
	var res int

	rows, err := (*context).DB.Query(query)
	if err != nil {
		fmt.Println("##################################################")
		fmt.Println(query)
		fmt.Println("##################################################")
		panic(err)
	}
	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&res); err != nil {
		panic(err)
	}

	return res
}
func DBSN2DBtype(DBT string) (DBtype, error) {

	switch DBT {
		case "psql", "PSQL", "postgres":
			return PSQL, nil
		case "mysql", "MySQL", "MYSQL":
			return MYSQL, nil
		case "mariadb", "MariaDB", "MARIADB":
			return MARIADB, nil
		case "sqlite3", "SQLITE3":
			return SQLITE3, nil
		}
	return DBTYPE_Last, errors.New("Unknown Database type")
}
