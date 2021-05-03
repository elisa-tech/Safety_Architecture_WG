# Copyright (C) 2021 Intel Corporation
# SPDX-License-Identifier: GPL-2.0-only

'''
Complexity Processor
Adds complexity info (McCabe score and statement count) for functions
'''

import logging
import os
import tempfile
import lizard
from helperfunctions import Helper

CALLTREEDIR = "calltreelog"


class ComplexityProcessor:
    '''Analyze complexity of functions and populate the database.'''
    def __init__(self, database, args):
        self.database = database
        self.args = args
        self.helper = Helper()

    def log_file_name(self, shortpath):
        '''Generate the correct log file path.'''
        info = shortpath.replace("/", "-")
        name = self.helper.log_file_name("complexity", info)
        return name

    def process_file(self, fullpath, shortpath):
        '''Process a file.'''
        logging.info("processing complexity for file %s", shortpath)
        with tempfile.NamedTemporaryFile(mode='w+t') as temp:
            # remove all #include statements from the source file
            self.helper.fix_source_file(fullpath, temp, temp.name)
            complexity_info = lizard.analyze_file(temp.name)
            for info in complexity_info.function_list:
                self.database.update_complexity(
                    info.name, info.cyclomatic_complexity, info.nloc)

    def process(self):
        '''Select files to processs.'''
        with open(self.args.linuxlog) as buildfile:
            for shortpath in buildfile:
                if shortpath is None:
                    continue
                shortpath = shortpath.rstrip()
                fullpath = os.path.join(self.args.sourcepath, shortpath)
                if not os.path.exists(fullpath):
                    logging.error("File not found: %s", shortpath)
                    continue
                self.process_file(fullpath, shortpath)
