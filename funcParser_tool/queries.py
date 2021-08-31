"""
SPDX-License-Identifier: GPL-2.0-only
SPDX-FileCopyrightText: Copyright (C) 2021 Intel Corporation

This program is free software; you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation; version 2.
This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.
You should have received a copy of the GNU General Public License along with
this program; if not, write to the Free Software Foundation, Inc., 51 Franklin
Street, Fifth Floor, Boston, MA 02110-1301, USA.

Author: Stefano Dell'Osa <stefano.dellosa@intel.com>

This program is a POC that parses C source files and headers looking for
 functions definitions and function calls.
They are then mapped into logical containers called SW Modules (or SW Blocks)
for further processing.
All data is then stored in SQLite database.

"""

import os
from setup_logger import logger
from domain import SwFunctionDefinition
from domain import SwFunctionCall
from domain import SwfunctionPointers
from domain import SwMacro
from domain import SwModule
from domain import SourceFile

def db_add_module(session, module_name):
    """This function executes queries to insert a function definition into the
        db is not already present."""
    # Chek if module already exists in db
    if session.query(SwModule).filter(
            SwModule.name == module_name).one_or_none() is None:
        # Insert SW module into the db
        sw_module = SwModule(name=module_name)
        session.add(sw_module)
    else:
        logger.debug("Skipping module %s since already present in db.",
                     module_name)


def db_add_source_file(session, src_path, module_name, linux_src_path):
    """This function executes queries to insert a source file into the
            db is not already present."""
    # Check if source files already exist in db
    if session.query(SourceFile).filter(
            SourceFile.path == src_path).one_or_none() is None:
        # Insert source file into the db
        source_file = SourceFile(name=os.path.basename(src_path),
                                 path=linux_src_path,
                                 module=module_name)
        session.add(source_file)
        session.commit()
    else:
        source_file = None
        logger.debug("Skipping source file %s since already"
                     " present in db.", src_path)

    return source_file


def db_add_func_def(session, func_name, source_file, func_type, func_params):
    """This function executes queries to insert a function definition into the
    db is not already present."""
    # Check if function definition is already in db
    if session.query(SwFunctionDefinition).filter(
            SwFunctionDefinition.name == func_name,
            SwFunctionDefinition.source_file_path ==
            source_file.path).one_or_none() is None:
        # Insert function definition into the db
        function = SwFunctionDefinition(type=func_type,
                                        name=func_name,
                                        params=func_params,
                                        source_file_path=source_file.path)
        session.add(function)
        session.commit()
    else:
        logger.debug(
            "Skipping %s function from %s since already in db" %
            (func_type+' '+func_name+func_params, source_file.path))


def db_add_func_call(session, func_call, source_file):
    """This function executes queries to insert a function call into the
        db is not already present."""
    # Check if function is already in db
    if session.query(SwFunctionCall).filter(
            SwFunctionCall.name == func_call,
            SwFunctionCall.source_file_path ==
            source_file.path).one_or_none() is None:
        # Insert function into the db
        func_call = SwFunctionCall(name=func_call,
                                   source_file_path=source_file.path)
        session.add(func_call)
        session.commit()
    else:
        logger.debug("Skipping %s function from %s since already in db" %
                     (func_call, source_file.path))


def db_add_func_pointer(session, func_ptr, source_file):
    """This function executes queries to insert a function call into the
            db is not already present."""
    # Check if function pointer is already in db
    if session.query(SwfunctionPointers).filter(
            SwfunctionPointers.name == func_ptr,
            SwfunctionPointers.source_file_path ==
            source_file.path).one_or_none() is None:
        # Insert function pointer into the db
        func_ptr = SwfunctionPointers(name=func_ptr,
                                      source_file_path=source_file.path)
        session.add(func_ptr)
        session.commit()
    else:
        logger.debug("Skipping %s function pointer from %s since already in db"
                     % (func_ptr, source_file.path))


def db_add_macro(session, macro_name, macro_params, source_file):
    """This function executes queries to insert a function call into the
            db is not already present."""
    # Check if macro is already in db
    if session.query(SwMacro).filter(
            SwMacro.name == macro_name,
            SwMacro.source_file_path == source_file.path).one_or_none() is None:
        # Insert function pointer into the db
        macro = SwMacro(name=macro_name, params=macro_params,
                        source_file_path=source_file.path)
        session.add(macro)
        session.commit()
    else:
        logger.debug("Skipping %s macro from %s since already in db" %
                     (macro_name, source_file.path))

