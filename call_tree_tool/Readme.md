# Call-Tree Tool User's Guide

## Pre-requisite tools
Call-Tree Tool uses several external tools. Download and install these tools before using Call-Tree Tool. 
Call-Tree Tool requires Python 3.9 or above. See the instructions below for setting up the required virtual environment.

### Set up the Python virtual environment
Follow the instructions in https://docs.python.org/3/library/venv.html to set up and activate a virtual environment for Python 3.9.

From within the virtual environment install the "lizard" package by executing: pip instlal lizard.

### GNU cflow
Call-Tree Tool uses GNU cflow as the main analysis tool to identify functions and edges (function calls) between functions.

Call-Tree Tool uses cflow version 1.6.
Download cflow from: https://www.gnu.org/software/cflow/
Build it and make accessible to Call-Tree Tool.

cflow is used in two python files:
* edgeprocessor.py
* functionprocessor.py

In both files change the value of the CFLOW variable to point to the cflow executable.

## sqlite3
Call-Tree Tool does not use the sqlite3 command line utility, but since all the data is stored in an sqlite3 database, it is convenient to have the sqlite3 utility, or an sqlite3 GUI browser.

Download and build sqlite3 from: https://www.sqlite.org/download.html or install a package for your system (e.g.: "apt install sqlite3" ).

## Using Call-Tree Tool
The following sections explain how to use Call-Tree Tool to scan source files and how to access the data it collects.

### Database schema
Call-Tree Tool stores all information in a SQLite3 database. The schema of the DB is very simple and self-explanatory and can be viewed by the ".schema" command of the sqlite3 tool, or with any other sqlite3 browser.

Some less self-explanatory fields are listed below:
* is_sr - means "is safety-related?". This field indicates whether a function is safety related. The script to classify functions as SR or NSR (non safety-related) will be provided at a later stage.

### Files table 
Use this command to fill the Files table in the DB:

python ./calltrees.py --index-files --source-path \<source path> --linux-log \<linux build log> --loglevel \<loglevel> \<sqlite3 DB file>

The arguments to this command are:
* --index-files - the action to be done
* --source-path - full path to the linux source directory (the "linux" directory")
* --linux-log - a file that lists the files to scan. This file contains records, one on each line. Each record is the relative path to the source path. For example, the file .../linux/kernel/fork.c should be listed as "kernel/fork.c"
* --log-level- the logging level of the Call-Tree Tool script. The levels are:
  * none - writes no information
  * critical
  * error
  * warning
  * info
  * debug - writes most information
* path the the sqlite3 DB file into which data will be written

Note that these command line arguments are used in the other commands as well.

### Functions table 
Use the command:
python ./calltrees.py --index-functions --source-path \<source path> --linux-log \<linux build log> --loglevel \<loglevel> \<sqlite3 DB file>

The command line arguments are the same as for the index-files command, except for the --index-functions command

### Index functions from header files
Use this command:
python ./calltrees.py --index-headers \<source path> --loglevel \<loglevel> \<sqlite3 DB file>

This command is similar but the path to the linux sources (including the "linux" directory) is given as an argument to the --index-header option.

### Index edges - function calls between functions
Use the command:
python ./calltrees.py --index-edges --source-path \<source path> --linux-log \<linux build log> --loglevel \<loglevel> \<sqlite3 DB file>

The arguments are the same as for the other commands.

### Note
The shell script: calltrees-elisa.sh performs all the aforementioned commands in sequence. It takes 3 arguments:
* Path to the root of the Linux source tree
* List of files to scan (the linuxbuild.log file).
* Path to the sqlite3 DB file

### Functions filter
In order to handle some special cases detected as functions but considered false positives (e.g. Compiler extensions to C language), caltrees exposes a filter file where the user
can list the cases he does not want to consider as function. File is present in calltrees root folder and provded empty by default. Edit the file insterting patterns to filter (if any) followed by a newline (list separator).

### Print a call tree
To print a calltree to stdout use this command:
python ./calltrees.py --draw-function <function> --draw-depth <depth> --loglevel <log level> <database file>

For example:
python ./calltrees.py --draw-function _do_fork --draw-depth 10 --loglevel info ./calltrees.db
 
The arguments are:
* --draw-function - the name of the function to draw. Note that you may have to look up the exact name of the function in the DB. For example, system calls are prefixed with sys_ (e.g. sys_mmap).
* --draw-depth - the depth of the tree to draw
* --loglevel - the log level for the command, as listed above
* the path to the sqlite3 database file



