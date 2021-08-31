# funcParser

funcParser tool is a POC developed to support ELISA FFI activities. It parses source files in order to parse functions logically grouping them a given SW module/Block.
It is intended to be part of a tool suite supporting the process of proving FFI on Linux kernel SW modules/block.
It is released to ELISA project with the aim of refining it with contribution of the whole communuty.

# Dependencies
Required host dependencies:
- python 3.9
- SQLite 2.8.17

Required python packages:
- sqlalchemy 1.4

# Usage:
	Run 'main.py -h' for help and instructions

# Key concepts:
- **Freedom From Interference (FFI)**: absence of cascading failures between components with No/different ASIL/SIL levels.
- **SW Module/Block**: logical container gouping a set of source files.
	 **NOTE:** source files belong to a specified module and cannot be shared between different modules.
- **Incoming Functions**: functions defined within a given SW Module/Block
- **Outgoing Functions**: functions called from a given SW Module/Block (Query will filter out functions calls referring to functions defined within a SW Module/Block)

# Sample queries:

## Create view listing all functions defined in a given module (e.g. fs)
	CREATE VIEW fs_incoming_functions AS
	SELECT Function_Definitions.Name, Function_Definitions.Source_file, Source_Files.SW_Module
	FROM Function_Definitions, Source_Files
	WHERE Function_Definitions.Source_file == Source_Files.Path AND Source_Files.SW_Module like 'fs'

## Query all outgoing function calls from a given module
	SELECT DISTINCT Function_Calls.Name, Function_Calls.Source_File
	FROM Function_Calls, fs_incoming_functions, Source_Files
	WHERE Function_Calls.Source_File == Source_Files.Path AND fs_incoming_functions.SW_Module == Source_Files.SW_Module AND
	Function_Calls.Name NOT IN (
	SELECT Name FROM fs_incoming_functions)

# Author
- Stefano Dell'Osa <stefano.dellosa@intel.com>

# Contributor(s)
