# Copyright (C) 2021 Intel Corporation
# SPDX-License-Identifier: GPL-2.0-only

'''
Header Processor
Process CFlow output to find functions in header files
'''

import logging
import pathlib
from functionprocessor import FunctionProcessor
from complexityprocessor import ComplexityProcessor


class HeaderProcessor:
    '''Process header files.'''
    def __init__(self, database, args):
        self.headers_path = args.headers_path
        self.function_processor = FunctionProcessor(database, args)
        self.complexity_processor = ComplexityProcessor(database, args)

    def process_file(self, fullpath, file):
        '''Process a single header file.'''
        logging.info("Processing header file: %s", file)
        self.function_processor.process_file(fullpath, file)
        self.complexity_processor.process_file(fullpath, file)

    def process(self):
        '''Find header files and process them.'''
        logging.info("Header processor starting")
        for filepath in pathlib.Path(self.headers_path).glob('**/*.h'):
            strpath = str(filepath)
            shortpath = strpath[len(self.headers_path) + 1:]
            # the following block is an example of skipping irrelevant files
            if "arch" in shortpath:
                if "x86" not in shortpath:
                    logging.info("skipping arch file not relatd "
                          "to x86: %s", shortpath)
                    continue
            self.process_file(strpath, shortpath)
