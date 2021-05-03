# Copyright (C) 2021 Intel Corporation
# SPDX-License-Identifier: GPL-2.0-only

'''
Main function of Call Trees Tool
Parses arguments and executes the main flows
'''

import os
import sys
import logging
import argparse
import string

from datetime import datetime
from dbapi import DB
from edgeprocessor import EdgeProcessor
from graphdraw import GraphDrawFile
from functionprocessor import FunctionProcessor
from headerprocessor import HeaderProcessor
from complexityprocessor import ComplexityProcessor
from helperfunctions import Helper

CALLTREEDIR = "calltreelog"
LOGFILE = "calltrees.log"

__version__ = "1.0.0"
__copyright__ = "Copyright (C) 2021 Intel Corporation"
__license__ = "Licensed under GPL-2.0-or-later; see the source for copying conditions."

python_major = sys.version_info[0]
python_minor = sys.version_info[1]
if python_major < 3 or (python_major == 3 and python_minor < 9) :
    print("Python 3.9 or a more recent version is required.")
    sys.exit(os.EX_USAGE)

class ReadablePath(argparse.Action):
    ''' Verifies that the input directory or file exists and is readable.'''
    def __call__(self, parser, namespace, values, option_string=None):
        prospective_path=values
        if not all(c in string.printable for c in prospective_path):
            print("input path contains unprintable characters")
            sys.exit(os.EX_DATAERR)
        if any((c in ' ?*<>[],;:%"') for c in prospective_path):
            print("input path contains invalid characters")
            sys.exit(os.EX_DATAERR)
        if len(os.path.basename(prospective_path)) > os.pathconf('/', 'PC_NAME_MAX'):
            print("file name too long")
            sys.exit(os.EX_DATAERR)
        if len(prospective_path) > os.pathconf('/', 'PC_PATH_MAX'):
            print("path is too long")
            sys.exit(os.EX_DATAERR)
        if self.dest != "dbfile":
            if not os.path.exists(prospective_path):
                print("{} is not a valid path".format(prospective_path))
                sys.exit(os.EX_DATAERR)
            if not os.access(prospective_path, os.R_OK):
                print("{} is not a readable path".format(prospective_path))
                sys.exit(os.EX_DATAERR)

        setattr(namespace, self.dest, prospective_path)


def index_files(args, database):
    '''Enter files into the database.'''
    logging.info("Indexing files")

    with open(args.linuxlog) as buildfile:
        for line in buildfile:
            if line is None:
                continue
            line = line.rstrip()
            fullpath = os.path.join(args.sourcepath, line)
            if os.path.exists(fullpath) is False:
                logging.warning("File %s does not exist", line)
                continue
            name = line.split("/")[-1]
            database.add_file(fullpath, line, name)


def setup_logging(args):
    '''Set up logging.'''
    levels = {
        "critical": logging.CRITICAL,
        "error": logging.ERROR,
        "warning": logging.WARNING,
        "info": logging.INFO,
        "debug": logging.DEBUG,
        "none": logging.NOTSET
    }
    timestamp = datetime.now().strftime("%d%b%Y-%H%M%S")
    root_logger = logging.getLogger()

    file_handler = logging.FileHandler(LOGFILE)
    file_handler.setLevel(levels[args.loglevel])
    root_logger.addHandler(file_handler)

    console_handler = logging.StreamHandler()
    console_handler.setLevel(levels[args.loglevel])
    root_logger.addHandler(console_handler)
    logging.info("Calltrees started at %s", timestamp)

def parsearguments():
    '''Parse command line arguments.'''
    version_string = "Call-Tree Tool {}\n{}\n{}".format(
            __version__, __copyright__, __license__)

    parser = argparse.ArgumentParser(prog="Call-Tree Tool",
            formatter_class=argparse.RawTextHelpFormatter)
    parser.add_argument("dbfile", help="the sqlite3 database", action=ReadablePath)
    parser.add_argument("--index-files", dest="index_files",
                        help="index files", action="store_true")
    parser.add_argument("--index-functions", dest="index_functions",
                        help="index functions", action="store_true")
    parser.add_argument("--draw-function", dest="draw_file",
                        help="draw a calltree", action="store")
    parser.add_argument("--draw-depth", dest="draw_depth", type=int,
                        help="depth of the call tree", action="store")
    parser.add_argument("--index-edges", dest="index_edges",
                        help="index edges - function calls",
                        action="store_true")
    parser.add_argument("--linux-log", dest="linuxlog",
                        help="file containing the list of files to scan",
                        action=ReadablePath)
    parser.add_argument("--source-path", dest="sourcepath",
                        help="path to the source code", action=ReadablePath)
    parser.add_argument("--index-headers", dest="headers_path",
                        help="index functions from header files",
                        action=ReadablePath)
    parser.add_argument("--complexity", dest="complexity",
                        help="perform a complexity analysis",
                        action="store_true")
    parser.add_argument("--loglevel", dest="loglevel",
                        help="Logging level: <critical, error, "
                        "warning, info, debug, none",
                        choices=['critical', 'error', 'warning', 'info',
                            'debug', 'none'],
                        action="store")
    parser.add_argument("--version", action="version", version=version_string)
    args = parser.parse_args()
    return args

def main():
    ''' Parse arguments and execute.'''
    args = parsearguments()
    setup_logging(args)
    database = DB(args.dbfile)
    database.createdb()
    helper = Helper()

    helper.create_log_dirs()

    if args.draw_file is not None:
        origin = database.find_origin(args.draw_file)
        if origin is None:
            print("function not found in database")
            sys.exit(os.EX_DATAERR)

        node = (origin[0], origin[1])
        calltreelogfile = "./{}/{}-{}".format(CALLTREEDIR,
                                              args.draw_file, args.draw_depth)
        drawer = GraphDrawFile(database, node,
                               calltreelogfile, int(args.draw_depth))
        drawer.draw()

    if args.index_files and args.linuxlog is not None \
            and args.sourcepath is not None:
        index_files(args, database)

    if args.index_functions and args.linuxlog is not None \
            and args.sourcepath is not None:
        func_processor = FunctionProcessor(database, args)
        func_processor.process()

    if args.index_edges and args.linuxlog is not None \
            and args.sourcepath is not None:
        edge_processor = EdgeProcessor(database, args)
        edge_processor.process()

    if args.complexity and args.linuxlog is not None \
            and args.sourcepath is not None:
        complexity_processor = ComplexityProcessor(database, args)
        complexity_processor.process()

    if args.headers_path is not None:
        header_processor = HeaderProcessor(database, args)
        header_processor.process()


if __name__ == "__main__":
    main()
