# Copyright (C) 2021 Intel Corporation
# SPDX-License-Identifier: GPL-2.0-only

'''
Function Processor
Parses CFlow output to find functions in code
'''

import logging
import os
import subprocess
import tempfile
from helperfunctions import Helper


CFLOW = "/usr/local/bin/cflow"


class FunctionProcessor:
    '''Process functions data from cflow and populate the database.'''
    def __init__(self, database, args):
        self.database = database
        self.args = args
        self.helper = Helper()

    def log_file_name(self, shortpath):
        '''Generate the correct log file name.'''
        info = shortpath.replace("/", "-")
        name = self.helper.log_file_name("functions", info)
        return name

    def process_file(self, fullpath, shortpath):
        '''Process cflow output for a single file.'''
        # pylint: disable=unused-variable
        logging.info("indexing functions for %s", shortpath)
        with tempfile.NamedTemporaryFile(mode='w+t', dir=os.getcwd()) as temp:
            # remove all #include files and fix SYSCALL_DEFINE
            self.helper.fix_source_file(fullpath, temp, temp.name)
            try:
                output = subprocess.run([CFLOW, "-x", "-i", "_st", temp.name],
                                        capture_output=True, text=True,
                                        check=True, shell=False)
            except subprocess.CalledProcessError as err:
                logging.error("cflow failed for function %s %s", shortpath, err)
                return

        self.helper.write_to_log(self.log_file_name(shortpath),
                                 self.args.loglevel, output.stdout)

        # the [:-1] at the end is because the last line also
        # contains \n and creates an additional empty line
        lines = output.stdout.split('\n')[:-1]
        for line in lines:
            components = line.split()
            if components[1] != "*":
                continue
            funcname = components[0]
            filelocation = components[2]
            filename, linenumber = filelocation.split(":")
            signature = " ".join(components[3:])
            self.database.add_function(fullpath, shortpath,
                                       funcname, linenumber, signature)

    def process(self):
        '''Process all files.'''
        with open(self.args.linuxlog) as buildfile:
            for shortpath in buildfile:
                if not shortpath:
                    continue
                shortpath = shortpath.rstrip()
                fullpath = os.path.join(self.args.sourcepath, shortpath)
                if not os.path.exists(fullpath):
                    logging.warning("file not found: %s", shortpath)
                    continue
                self.process_file(fullpath, shortpath)
