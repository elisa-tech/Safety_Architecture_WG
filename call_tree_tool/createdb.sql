-- Copyright (C) 2021 Intel Corporation
-- SPDX-License-Identifier: GPL-2.0-only

CREATE TABLE IF NOT EXISTS files (
    fullpath     TEXT,
    shortpath    TEXT,
    name	     TEXT,
    line_count   INTEGER,
    scan_date 	 TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_sr     	 TEXT CHECK (is_sr in ('yes', 'no', 'mixed')),
    checksum     TEXT,
    PRIMARY KEY(shortpath)
);

CREATE TABLE IF NOT EXISTS functions (
    fullpath     TEXT,
    shortpath    TEXT,
    name         TEXT,
    line         INTEGER,
    mccabe       INTEGER,
    line_count   INTEGER,
    is_sr     	 TEXT CHECK (is_sr in ('yes', 'no')),
    min_depth    INTEGER,
    checksum     TEXT,
    signature    TEXT,
    PRIMARY KEY(shortpath, name)
);

CREATE INDEX IF NOT EXISTS idx_functions ON functions(name);

CREATE TABLE  IF NOT EXISTS edges (
    frompath	TEXT,
    fromfunc	TEXT,
    linenumber  INTEGER,
    topath	    TEXT,
    tofunc	    TEXT,
    PRIMARY KEY(frompath, fromfunc, topath, tofunc)
);

CREATE INDEX  IF NOT EXISTS idx_edge ON edges (frompath, fromfunc);

