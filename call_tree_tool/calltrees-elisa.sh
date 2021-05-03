#!/bin/sh

if [ "$#" -ne 3 ]; then
  echo "Usage: $0 <root of linux sources> <list of source files to scan> <db file name>"
  exit 1
fi

#
# clean old version of generated files
#
rm -f $3
rm -rf calltreelog
rm -f calltrees.log
rm -f calltrees-config.h
#
# the following lines generate the #define statements according to the configuration
#
sed -n 's/\(CONFIG.*\)=y/#define \1/p' $1/prod_output/.config > calltrees-config.h
sed -n 's/.*\(CONFIG.*\) is not set/#undef \1/p' $1/prod_output/.config >> calltrees-config.h

#
# analyze file, functions, edges and complexity
# log-levels are: critical, error, warning, info, debug, nono
#
python ./calltrees.py --index-files --source-path $1 --linux-log $2 --loglevel error $3 
python ./calltrees.py --index-functions --source-path $1  --linux-log $2  --loglevel error $3
python ./calltrees.py --index-headers $1  --loglevel error $3
python ./calltrees.py --complexity --source-path $1 --linux-log $2 --loglevel error $3
python ./calltrees.py --index-edges --source-path $1 --linux-log $2  --loglevel error $3

#
# remove the calltrees-config.h file
#
rm -f calltrees-config.h

