# Copyright (C) 2021 Intel Corporation
# SPDX-License-Identifier: GPL-2.0-only

'''
Edge Processor
Processes CFlow output to find edges
'''

import logging
import os
import subprocess
import tempfile
from helperfunctions import Helper

INDENTATION = 4
CFLOW = "/usr/local/bin/cflow"

class EdgeProcessor:
    '''Process files to detect calling relations between functions.'''
    def __init__(self, database, args):
        self.args = args
        self.database = database
        self.helper = Helper()
        self.current_function = []
        self.current_indent = 0
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
            funcdata = ("unknown", funcname)
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
        with tempfile.NamedTemporaryFile(mode='w+t', dir=os.getcwd()) as temp:
            # remove all #include statements from the source file
            self.helper.fix_source_file(fullpath, temp, temp.name)
            try:
                output = subprocess.run([CFLOW, temp.name],
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
