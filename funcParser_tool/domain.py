"""
SPDX-License-Identifier: GPL-2.0-only
SPDX-FileCopyrightText: Copyright (C) 2021 Intel Corporation

This program is free software; you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation; version 2.
This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.
You should have received a copy of the GNU General Public License along with
this program; if not, write to the Free Software Foundation, Inc., 51 Franklin
Street, Fifth Floor, Boston, MA 02110-1301, USA.

Author: Stefano Dell'Osa <stefano.dellosa@intel.com>
"""

import re
from sqlalchemy import Column, String, ForeignKey
from sqlalchemy.orm import registry
from sqlalchemy.orm import declarative_base
from setup_logger import logger
import sys

# Database config
Base = declarative_base()
mapper_registry = registry()


# Classes
class SwModule(Base):
    """This class represents a SW Module, intended as a collection of source
     files."""
    __tablename__ = 'Modules'

    name = Column('Name', String, primary_key=True, nullable=False)

    def __repr__(self):
        return "<SwModule(name='%s')>" & self.name


class SourceFile(Base):
    """This class represents a source file."""
    __tablename__ = 'Source_Files'

    name = Column('Name', String, nullable=False)
    path = Column('Path', String, primary_key=True, nullable=False,
                  index=True)
    module = Column('SW_Module', String, ForeignKey("Modules.Name",
                                                    ondelete="CASCADE"))

    def __repr__(self):
        return "<SourceFile(name='%s', path='%s', module='%s')>" &\
               (self.path, self.name, self.module)

    def find_function_definitions(self):
        """this method finds all functions defined within a given source
        file."""
        try:
            with open(self.path, "r") as src_file:
                text = src_file.read()
        except (IOError, OSError) as e:
            logger.error("Error while opening %s source file. %s" % (self.path,
                                                                     e))
            return None

        func_definitions = list(filter(None, re.findall(
            r"^[a-z0-9_ \n*]*[*&]{0,2}?[a-z0-9_]+\([\w\s,*]*?\)"
            r"[\s]*[\r\n]*{.*[\r\n]",
            text, re.M)))

        logger.debug("Defined functions for file %s: %s" % (self.name,
                                                            func_definitions))

        return self.__format_function_definitions(func_definitions)

    @staticmethod
    def __format_function_definitions(func_definitions):
        """This method formats function attributes out of a raw output."""

        formatted_functions = []
        for line in func_definitions:
            line = line.replace('\n', ' ')
            line = line.replace('\t', '')

            # Sample format:
            # 'long vfs_ioctl(struct file *filp, unsigned int cmd, long arg) {'
            matched_type = re.findall(r'^.*(?=\s[a-z0-9_*]*\()', line)
            matched_name = re.findall(r'(?<=\s)[a-z0-9_*]*(?=\()', line)
            matched_params = re.findall(r'\(.*\)', line)

            if len(matched_name) > 1 or len(matched_type) > 1 or \
                    len(matched_params) > 1:
                logger.error("Multiple results for in regex in line: %s."
                             " Skipping.", line)
                continue

            func_name = matched_name.pop()
            func_type = matched_type.pop()
            # It might happen that * preceds the func name. In this case we move
            # the * to the type string
            if '*' in func_name:
                func_name = func_name.replace('*', '')
                func_type+=str(' *')

            cur_func = {
                "type": func_type,
                "name": func_name,
                "params": matched_params.pop()
            }
            formatted_functions.append(cur_func)

        logger.debug("Formatted functions: %s", formatted_functions)

        return formatted_functions

    @staticmethod
    def __format_macros(macros):
        """This method formats macros attributes out of a raw output."""

        formatted_macros = []
        for line in macros:
            line = line.replace('\n', ' ')
            line = line.replace('\t', '')

            # Sample format: 'SYSCALL_DEFINE3(ioctl, unsigned int, fd,
            # unsigned int, cmd, unsigned long, arg) {'
            cur_func = {
                "name": list(filter(None, re.findall(r'[A-Z0-9_]*(?=\()',
                                                     line))),
                "params": re.findall(r'\(.*\)', line)
            }
            formatted_macros.append(cur_func)
        logger.debug("Formatted macros: %s", formatted_macros)

        return formatted_macros

    def find_function_calls(self):
        """This method search for function calls inside the source file."""

        try:
            with open(self.path, "r") as src_file:
                text = src_file.read()
        except (IOError, OSError) as e:
            logger.error("Error while opening %s source file. %s" % (self.path,
                                                                     e))
            return None

        # Regex to search function calls defined as: funcname(param, &param,
        # param->value, "value");
        func_calls = list(filter(None,
                                 re.findall(
                                     r'(?<=[\s!])[a-z0-9_]*(?=\('
                                     r'[*,&]{0,2}\w[,\s*\w\-\>\_\&\*\\(\)'
                                     r'"\/%]*\)[ ;,]{1})',
                                     text)))

        # Regex to find function calls in parenthesys, e.g.: (func(), func2());
        func_betw_parentheses = list(filter(
            None,
            re.findall(r'(?<=\()[a-z0-9_!]*(?=\('
                       r'[*&!]{0,2}[\w]{0,1}[,\s*\w\-\>\_\&\*\(\)"|=%+-]*\)'
                       r'[ ;{\r\n]{1})',
                       text)))
        # Remove ! char from func name (if present)
        func_pars_purified = [func.replace('!','') for func in
                              func_betw_parentheses]
        func_calls = self.__merge_lists(func_calls, func_pars_purified)

        # Regex to find function calls with no parameters. E.g.: ktime_get()
        func_no_params = list(filter(
            None,
            re.findall(r'(?<=[\s!])[a-z0-9_]*(?=\(\)[ ;,]{1})',
                       text)))
        func_calls = self.__merge_lists(func_calls, func_no_params)

        # Regex to find function declarations E.g. int func_name(par1, par2);
        func_decl = list(filter(None,
                                re.findall(r"^(\w* )[a-z0-9_]*\(.*\);[\r\n]",
                                           text)))

        logger.debug("Function declarations for file %s: %s" % (self.name,
                                                                func_decl))

        for cur_func_declr in func_decl:
            if not cur_func_declr:
                continue
            func_name = re.findall(r'(?<=\s)\w*(?=\()', cur_func_declr)
            if func_name and func_name in func_calls:
                logger.debug("Removing %s from function call list since is a"
                             " function declaration.")
                func_calls.remove(func_name)
                continue

        # Return a list of unique elements
        func_calls = list(dict.fromkeys(func_calls))
        logger.debug("Function calls for file %s: %s" % (self.name, func_calls))

        return func_calls

    def find_function_pointers(self):
        """This method search for function pointers inside the source file."""

        try:
            with open(self.path, "r") as src_file:
                text = src_file.read()
        except (IOError, OSError) as e:
            logger.error("Error while opening %s source file. %s" % (self.path,
                                                                     e))
            return None

        # Regex to find function calls by means of function pointers:
        # struct->func_ptr(param, param);
        func_ptrs = list(filter(None,
                                re.findall(r"(?<=->)[a-z0-9_]*(?=\(.*\);)",
                                           text)))

        logger.debug("Function pointers for file %s: %s" % (self.name,
                                                            func_ptrs))

        return func_ptrs

    @staticmethod
    def __merge_lists(main_list, subsidiary_list):
        """This method merges item of subsidiary list in main list avoiding to
        add duplicated items.
        :rtype: object
        """
        for cur_item in subsidiary_list:
            if cur_item is not None and cur_item not in main_list:
                main_list.append(cur_item)

        return main_list

    def find_macros(self):
        """this method finds all macros defined within a given source file."""

        try:
            with open(self.path, "r") as src_file:
                text = src_file.read()
        except (IOError, OSError) as e:
            logger.error("Error while opening %s source file. %s" % (self.path,
                                                                     e))
            return None

        macro_definitions = re.findall(r"^[A-Z0-9_\s]*[*&]{0,2}?[A-Z0-9_]\("
                                       r"\w[,\s*\w]*\)[\s]*[\r\n]*{.*[\r\n]",
                                       text, re.M)

        logger.debug("Defined macros for file %s: %s" % (self.name,
                                                         macro_definitions))

        return self.__format_macros(macro_definitions)


class SwFunctionDefinition(Base):
    """This class represents a function defined in a source file."""
    __tablename__ = 'Function_Definitions'

    name = Column('Name', String, primary_key=True, nullable=False, index=True)
    type = Column('Type', String, nullable=False)
    params = Column('Params', String, nullable=False)
    source_file_path = Column('Source_file',
                              String,
                              ForeignKey('Source_Files.Path',
                                         ondelete="CASCADE"),
                              primary_key=True, index=True)

    def __repr__(self):
        return "<SwFunctionDefinition(type='%s', name='%s', params='%s'," \
               " source_file_path='%s')>" & \
               (self.type, self.name, self.params, self.source_file_path)


class SwFunctionCall(Base):
    """This class represents a function call."""
    __tablename__ = 'Function_Calls'

    name = Column('Name', String, primary_key=True, nullable=False, index=True)
    source_file_path = Column('Source_File', String, primary_key=True,
                              nullable=False, index=True)

    def __repr__(self):
        return "<SwFunctionCall(name='%s', source_file_path='%s')>" & \
               (self.name, self.source_file_path)


# TODO:best would be to use preprocessor output, in that case Macros would be
#  resolved and this class will be non needed any longer.
class SwMacro(Base):
    """This class represents a Macro"""
    __tablename__ = 'Function_Macros'

    name = Column('Name', String, primary_key=True, nullable=False, index=True)
    params = Column('Params', String, nullable=False)
    source_file_path = Column('Source_File', String, primary_key=True,
                              nullable=False)

    def __repr__(self):
        return "<SwMacro(name='%s', params='%s', source_file_path='%s')>" & \
               (self.name, self.params,
                self.source_file_path)

class SwfunctionPointers(Base):
    """This class represents a Function Pointer"""
    __tablename__ = 'Function_Pointers'

    name = Column('Name', String, primary_key=True, nullable=False, index=True)
    source_file_path = Column('Source_File', String, primary_key=True,
                              nullable=False, index=True)


