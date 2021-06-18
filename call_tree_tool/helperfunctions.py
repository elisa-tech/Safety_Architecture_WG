# Copyright (C) 2021 Intel Corporation
# SPDX-License-Identifier: GPL-2.0-only

'''
Helper (utility) functions
'''

import os
import subprocess
import sys
import logging


class Helper:
    '''Provide helper functions.'''
    def __init__(self):
        self.logdir = "calltreelog"

    def create_log_dirs(self):
        '''Create log directories.'''
        if not os.path.exists(self.logdir):
            os.mkdir(self.logdir)
            os.mkdir(os.path.join(self.logdir, "edges"))
            os.mkdir(os.path.join(self.logdir, "functions"))
            os.mkdir(os.path.join(self.logdir, "variables"))
            os.mkdir(os.path.join(self.logdir, "complexity"))

    def log_file_name(self, folder, info):
        '''Generate a temp file path for log files.'''
        filename = "{}/{}/{}.log".format(self.logdir, folder, info)
        return filename

    def write_to_log(self, filename, loglevel, text):
        '''Writes data to the appropriate log file.'''
        if loglevel == "debug":
            with open(filename, "w") as logfile:
                logfile.write(text)

    def fix_source_file(self, infile_name, outfile, outfilename):
        '''Fix source files to allow processing by cflow.'''
        # check that the temporary file's permissions are 0600
        file_mode = oct(os.stat(outfilename).st_mode)[-3:]
        if file_mode != "600":
            logging.error("temporary file created with wrong permissions")
            sys.exit(os.EX_DATAERR)

        with open(infile_name, "r") as infile:
            outfile.write('#include "calltrees-config.h"\n')
            for line in infile:
                if "#include" in line:
                    continue
                try:
                    if "SYSCALL_DEFINE" in line and "MAXARGS" not in line:
                        prefix = ''
                        if "COMPAT" in line:
                            prefix = "compat_"
                        parts = line.split("(")[1].split(",")
                        funcname = parts[0]
                        outfile.write("int {}sys_{}({}".
                                      format(prefix, funcname,
                                             ''.join(parts[1:])))
                        continue
                except Exception:
                    logging.error("fix_source_file: problem occured in "
                                  "file: %s line: %s", infile_name, line)
                outfile.write(line)
            outfile.flush()

    def pre_process_file(self, linux_src_path, source_shortpath):
        '''
        Run pre processor on source file.
        :linux_src_path: The path to linux sources
        :source_shortpath: The name of the source file to preprocess.
        '''

        # pre-processed files extension is .i
        preproc_source_shortpath = source_shortpath.replace(".c", ".i")
        preproc_source_fullpath = os.path.join(linux_src_path, preproc_source_shortpath)
        if not os.path.exists(preproc_source_fullpath):
            # TODO: check a better way to compile single source without having to cd into root dir
            cur_dir = os.getcwd()
            command = "make %s" % preproc_source_shortpath
            # Cd Linux root dir to trigger make command
            os.chdir(linux_src_path)
            try:
                output = subprocess.run(command,
                                        capture_output=True, text=True,
                                        check=True, shell=True)
            except subprocess.CalledProcessError as e:
                os.chdir(cur_dir)
                logging.error("Failed executing preprocessor on file %s: %s" % (source_shortpath, e))
                return None
            # Cd to previous path
            os.chdir(cur_dir)
            with open(preproc_source_fullpath, 'r') as f_readonly:
                lines = f_readonly.readlines()

            with open(preproc_source_fullpath, 'w') as f:
                for line in lines:
                    if not line.startswith('#'):
                        f.write(line)

        return preproc_source_fullpath