# Copyright (C) 2021 Intel Corporation
# SPDX-License-Identifier: GPL-2.0-only

'''
DB API
Interface to the SQLite database
'''

import os
import sys
import sqlite3
import logging

CREATEDB = "./createdb.sql"


class DB:
    '''Provide an interface to the database.'''
    def __init__(self, dbfile):
        try:
            self.conn = sqlite3.connect(dbfile, isolation_level=None)
            self.cur = self.conn.cursor()
        except sqlite3.Error as err:
            logging.error("Failed openning to the DB. Error: %s", err)
            sys.exit(os.EX_DATAERR)

    def createdb(self):
        '''Create the SQLite database.'''
        with open(CREATEDB, "r") as scriptfile:
            script = scriptfile.read()
        try:
            self.cur.executescript(script)
        except sqlite3.Error as err:
            logging.warning("error executing the createdb.sql script. %s",
                             err.args)

    def add_file(self, fullpath, shortpath, name):
        '''Store a record into the files table.'''
        with open(fullpath) as infile:
            linecount = sum(1 for line in infile)
        command = "INSERT INTO files(fullpath, shortpath, line_count, "
        command += "name) VALUES(?, ?, ?, ?)"
        logging.debug(command)
        try:
            self.cur.execute(command, (fullpath, shortpath, linecount, name))
        except sqlite3.Error:
            logging.warning("Error adding file %s to DB", shortpath)

    def add_function(self, fullpath, shortpath, name, line, signature):
        '''Store an entry into the functions table'''
        command = "INSERT INTO functions (fullpath, shortpath, "
        command += "name, line, signature) VALUES(?, ?, ?, ?, ?)"
        logging.debug(command)
        logging.info(command)
        try:
            self.cur.execute(command,
                             (fullpath, shortpath, name, line, signature))
        except sqlite3.Error:
            logging.warning("function %s:%s already exists "
                            "in DB", shortpath, name)

    def get_function_by_name(self, funcname):
        '''Retrieve a function record.'''
        command = "SELECT * FROM functions WHERE name = ?"
        logging.debug(command)
        self.cur.execute(command, [funcname])
        data = self.cur.fetchall()
        return data

    def add_edge(self, origin, to_func):
        '''Store a record into the edges table.'''
        command = "INSERT INTO edges (frompath, fromfunc, topath, tofunc) "
        command += "VALUES(?, ?, ?, ?)"
        logging.debug(command, (origin[0], origin[1], to_func[0], to_func[1]))
        try:
            self.cur.execute(command, (origin[0], origin[1], to_func[0], to_func[1]))
        except sqlite3.Error:
            logging.warning("edge is already in the DB")

    def get_edges(self, origin):
        '''Retrieve records from the edges table'''
        '''
        command = "SELECT a.topath, a.tofunc, b.line_coverage, "
        command += "b.branch_coverage, b.line_count, b.branch_count "
        command += "FROM edges a LEFT JOIN coverage "
        command += "b ON a.tofunc = b.func WHERE a.frompath = ? "
        command += "AND a.fromfunc = ?"
        '''
        command =  "SELECT topath, tofunc FROM edges "
        command += "WHERE frompath = ? and fromfunc = ?"
        logging.debug(command)
        self.cur.execute(command, (origin[0], origin[1]))
        results = self.cur.fetchall()
        return results

    def find_origin(self, origin):
        '''Retrieve the shortpath and name for a function.'''
        command = "SELECT shortpath, name FROM functions WHERE name = ?"
        logging.debug(command)
        self.cur.execute(command, (origin,))
        result = self.cur.fetchone()
        return result

    def update_complexity(self, func, ccn, line_count):
        '''Update a function record with complexity data.'''
        query = "UPDATE functions SET mccabe = ?,  line_count = ? "
        query += "WHERE name = ?"
        logging.debug(query)
        try:
            self.cur.execute(query, (ccn, line_count, func))
        except sqlite3.Error:
            logging.warning("error updating complexity "
                            "for function %s", func)
