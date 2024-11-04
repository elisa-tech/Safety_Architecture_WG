#!/bin/bash
#
# SPDX-License-Identifier: GPL-2.0-only
#
# WHAT DO I DO?
# This script detects source code exporting symbols and yet not cross-referenced from
# the Linux Documentation files
#
# RATIONALE
# In https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html it is stated
# that "Every function that is exported to loadable modules using EXPORT_SYMBOL or
# EXPORT_SYMBOL_GPL should have a kernel-doc comment".
# In order to import such kernel-doc comments into the Linux Kernel Documentation
# builds, the respective source files, where such symbols are defined and exported,
# must be cross-referenced (according to
# "https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html#including-kernel-doc-comments")
# So if such source code is not cross-referenced, the respective exported symbols
# go undocumented
#
# HOWTO:
# 1) cd into a local directory
# 2) copy this file in such a local directory 
# 3) in such a local directory run "git clone https://github.com/torvalds/linux.git"
# 4) execute the script as in "./find_unreferenced_sources linux/<search_path>"
#    where <search_path> is the linux directory where we are looking for unreferenced
#    source code.
#    Example: "./find_unreferenced_sources linux/mm"
#
# Copyright (c) 2024 Red Hat, Inc.
#
for i in `grep EXPORT_SYMBOL $1 -R | cut -d: -f 1 | sort -u`;
do
	if [ -z "$(grep "kernel-doc:: ${i##*linux/}" $i "${i%linux*}linux/Documentation" -R)" ]
	then
                echo "${i##*linux/} is undocumented yet contains exported symbols"
	else
		echo "${i##*linux/} is documented"
       fi;
done
