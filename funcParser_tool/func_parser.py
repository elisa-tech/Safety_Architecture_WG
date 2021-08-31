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

# This is a POC that shows how to classify func_definitions of a given SW module
import argparse
import sys
import os
import json
import glob

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from setup_logger import logger
from domain import Base
import queries

version = '1.0'

def parse_arguments():
    """This function defiles the arguments to be parsed at program startup"""
    parser = argparse.ArgumentParser()
    parser = argparse.ArgumentParser(prog='funcParser',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog='''In order to setup funcParser to scan source file(s) make sure\
that the JSON configuration file properly reflects your module(s) and file(s)\
 list.''')
    parser.add_argument("-src", "--linux_sources_path", required=True,
                        help="Absolute path to Linux sources")
    parser.add_argument("-mdl", "--modules_file",
                        help="JSON configuration file name. It contains the list"
                             " of SW module(s) and related source file(s)"
                             "to scan. "
                             "Configuration file must be located in funcParser"
                             " folder.")
    parser.add_argument("-dbg", "--debug", action="store_false", default=False,
                        help="Enable debug mode loading SQLite database in "
                             "memory (no persistence). Default: False")
    return parser

def init_db(args):
    """This function initializes the SQLite database"""

    db_path = ""
    if args.debug:
        db_path = ":memory:"
        logger.info("Database is not persistent since debug mode is selected.")
    else:
        db_path = os.path.join(os.getcwd(), "funcParser.db")
        logger.info("Database path: %s", db_path)
    db = 'sqlite:///%s' % db_path

    # Check if file exists and delete it
    if not args.debug and os.path.isfile(db_path):
        try:
            logger.info("Deleting pre-existing db file %s ...", db_path)
            os.remove(db_path)
        except OSError:
            logger.error("Unable to delete pre existing db file.")
            sys.exit()

    engine = create_engine(db)
    Base.metadata.create_all(engine)
    Session = sessionmaker(bind=engine)
    session = Session()
    return session


def load_modules_from_json(args):
    """This function loads SW Modules/Blocks info from JSON file"""
    cur_dir = os.getcwd()
    json_cfg_file = os.path.join(cur_dir, args.modules_file)
    if not os.path.isfile(json_cfg_file):
        logger.error("Error in loading JSON modules file.")
        sys.exit(-2)

    try:
        with open(json_cfg_file) as json_file:
            data = json.load(json_file)
    except (IOError, OSError) as e:
        logger.error("Error while opening JSON module file. %s" % e)
        sys.exit(-1)

    return data


def process_function(func, source_file, session):
    """This function processes function definitions and inserts them into the
     db."""
    if not func['type'] or not func['name'] or not func['params']:
        logger.warning("Find Defined Functions: skipping %s since required"
                       " field is missing in %s."
                       "Skipping" % (func, source_file.path))
        return

    func_name = func.get('name')
    func_type = func.get('type')
    func_params = func.get('params')
    # Add function definition into db if not already present
    queries.db_add_func_def(session, func_name, source_file, func_type,
                            func_params)


def process_function_call(func_call, source_file, session):
    """This function processes function calls and inserts them into the db."""
    if not func_call:
        logger.warning("Find Function call: missing required field 'function"
                       " call name' in %s. "
                       " Skipping.", source_file.path)
        return
    # Add function call into db if not already present
    queries.db_add_func_call(session, func_call, source_file)


def process_function_ptr(func_ptr, source_file, session):
    """This function processes function pointers and inserts them into the
    db."""
    if not func_ptr:
        logger.warning("Find Function pointer: missing required field 'function"
                       " pointer name' in %s. "
                       " Skipping.", source_file.path)
        return
    # Add function pointer into db if not already present
    queries.db_add_func_pointer(session, func_ptr, source_file)


def process_macro(macro, source_file, session):
    """This function processes function macros and inserts them into the db."""
    if not macro['name'] or not macro['params']:
        logger.warning("Find Macros: skipping %s since required field is "
                       "missing.", macro)
        return

    macro_name = macro.get('name').pop()
    macro_params = macro.get('params').pop()

    # Add macro into db if not already present
    queries.db_add_macro(session, macro_name, macro_params, source_file)


def main():
    """This is the Main function"""
    logger.info("Welcome to funcParser %s!", version)

    args = parse_arguments().parse_args()
    session = init_db(args)
    modules_dict = load_modules_from_json(args)
    logger.debug("Modules: %s", modules_dict)

    # Iterate SW modules list and derive list of files, func_definitions
    # definitions and func_definitions calls
    for module in modules_dict['modules']:
        module_name = module['name']
        logger.info("Parsing module %s", module_name)

        # Add SW module in db if not already present
        queries.db_add_module(session, module_name)

        # Iterate over source files list
        # TODO: implement recursion into subdirs
        source_files = module['source_files']
        for mdl_src in source_files:
            path = mdl_src['path']
            src_paths = []

            # Check if path include wildcard
            for file_path in glob.glob(os.path.join(args.linux_sources_path,
                                                    path)):
                if os.path.isfile(file_path) and file_path not in src_paths and\
                        (file_path.endswith('.c') or file_path.endswith('.h')):
                    src_paths.append(file_path)
            for src_path in src_paths:
                logger.info("Parsing source file %s ...", src_path)

                # Insert source file into the db if not already present
                source_file = queries.db_add_source_file(
                    session, src_path, module_name, os.path.join(
                        args.linux_sources_path, src_path))

                if source_file is not None:
                    # Browse source file looking for func_definitions
                    for func in source_file.find_function_definitions():
                        process_function(func, source_file, session)

                    # Browse source file looking for function calls
                    for cur_func_call in source_file.find_function_calls():
                        process_function_call(cur_func_call, source_file,
                                              session)

                    # Browse source file looking for function pointers
                    for cur_func_ptr in source_file.find_function_pointers():
                        process_function_ptr(cur_func_ptr, source_file, session)

                    # Browse source file looking for Macros
                    for cur_macro in source_file.find_macros():
                        process_macro(cur_macro, source_file, session)

        # Write to db
        logger.info("Completed data insertion for module '%s' into db...",
                    module_name)

    logger.info("All done. Data stored and available in db.")


if __name__ == "__main__":
    main()
