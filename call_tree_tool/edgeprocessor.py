# Copyright (C) 2021 Intel Corporation
# SPDX-License-Identifier: GPL-2.0-only

'''
Edge Processor
Processes CFlow output to find edges
'''

import logging
import os
import sys
import subprocess
import tempfile
from helperfunctions import Helper

INDENTATION = 4
CFLOW = "/usr/bin/cflow"

class EdgeProcessor:
    '''Process files to detect calling relations between functions.'''
    def __init__(self, database, args):
        self.args = args
        self.database = database
        self.helper = Helper()
        self.current_function = []
        self.current_indent = 0

        if not os.path.exists(CFLOW):
            print("GNU cflow tool not found. Install GNU cflow")
            print("and change the CFLOW variable in edgeprocessor.py to")
            print("point to it.")
            sys.exit(1)

        logging.info("starting edge processor")

    def log_file_name(self, shortpath):
        '''Generate the correct path for logfiles.'''
        info = shortpath.replace("/", "-")
        name = self.helper.log_file_name("edges", info)
        return name

    def process_line(self, line):
        '''Process a line of cflow output.'''
        indentation = len(line) - len(line.lstrip(' '))
        funcname = line.split()[0]

        # macro names don't have () at the end
        if funcname[-2:] == "()":
            funcname = funcname[:-2]

        funcrecord = self.database.get_function_by_name(funcname)
        if not funcrecord:
            # check if function is in filter list (E.g. Compiler extensions to C language)
            with open("func_filter", "r") as filtered_funcs:
                 if funcname in filtered_funcs.read():
                     logging.debug("Filtering function %s...", funcname)
                 else:
                    funcdata = ("unknown", funcname)
            return
        else:
            funcdata = (funcrecord[0][1], funcname)

        if indentation > self.current_indent:
            self.process_child(funcdata)
            return
        if indentation < self.current_indent:
            self.process_parent(funcdata, indentation)
            return
        if indentation == self.current_indent:
            self.process_sibling(funcdata)
            return

    def process_sibling(self, funcdata):
        '''Process lines at the same level of indentation.'''
        if self.current_indent > 0:
            self.current_function.pop()
            parent = self.current_function[-1]
            self.database.add_edge(parent, funcdata)
        self.current_function.append(funcdata)

    def process_child(self, funcdata):
        '''Process a line that is forward indented.'''
        parent = self.current_function[-1]
        self.database.add_edge(parent, funcdata)
        self.current_function.append(funcdata)
        self.current_indent += INDENTATION

    def process_parent(self, funcdata, indentation):
        '''Process a line that is backwards indented.'''
        self.current_function.pop()
        for i in range(indentation, self.current_indent, INDENTATION):
            self.current_function.pop()
        if indentation > 0:
            parent = self.current_function[-1]
            self.database.add_edge(parent, funcdata)
        self.current_function.append(funcdata)
        self.current_indent = indentation

    def process_file(self, shortpath, fullpath):
        '''Process a file.'''
        logging.info("index edges for %s", shortpath)
        pre_processed_file = self.helper.pre_process_file(self.args.sourcepath, shortpath)
        if pre_processed_file is None:
            logging.warning("Skipping edges for %s since preprocessor run failed..." % shortpath)
            return

        try:
            output = subprocess.run([CFLOW, '-i', '-s', pre_processed_file],
                                    capture_output=True, text=True,
                                    check=True, shell=False)
        except subprocess.CalledProcessError:
            logging.error("Failed executing cflow")
            return

        self.helper.write_to_log(self.log_file_name(shortpath),
                                 self.args.loglevel, output.stdout)

        # the [:-1] at the end is because the last line also
        # contains \n and creates an additional empty line
        lines = output.stdout.split('\n')[:-1]

        self.current_function = []
        self.current_indent = 0
        for line in lines:
            self.process_line(line)

    def process(self):
        '''Perform edge processing.'''
        with open(self.args.linuxlog) as buildfile:
            for shortpath in buildfile:
                if shortpath is None:
                    continue
                shortpath = shortpath.rstrip()
                fullpath = os.path.join(self.args.sourcepath, shortpath)
                if not os.path.exists(fullpath):
                    logging.error("file not found: %s", shortpath)
                    continue
                self.process_file(shortpath, fullpath)
