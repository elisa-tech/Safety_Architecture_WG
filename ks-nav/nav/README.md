# Nav - kernel source code navigator

Nav is a tool that uses a pre-constituted database to emit call trees graphs that can be used stand alone or feed into a graph display system to help engineers do static analysis.

## Motivation
Although other similar tool do exist, the motivation for this tool is to be developed, is to solve a specific need: do a kernel source code analysis aimed at feature level and  shows it as functions call tree or subsystems call tree. 

## Build
Nav is implemented in Golang. Golang application are usually easy to build. In addition to this, it has very few dependencies other than the standard Golang packages.
A `Makefile` is provided to ease the build process. 
Just typing 
```
$ make
```
the native build is made.
In addition to the default target,  `amd64`, `arm64`, and `upx`, targets exists.
|target |function                                                                   |
|-------|---------------------------------------------------------------------------|
|amd64  |forces the build to amd64 aka x86_64 regardless the underlying architecture|
|arm64  |forces the build to arm64 aka aarc64 regardless the underlying architecture|
|upx    |triggers compress previously generated executable using UPX                |
As example, this builds aarch64 upx compressed executable:
```
$ make arm64
$ make upx
```
## Usage example
As the nav compiled executable is available, it is essential to provide the configuration to query the backend database. The easiest way to provide the configuration to nav is to specify a configuration file.
Although the nav tool has an internal default for all the configuration parameters, that are used if not otherwise specified, this default can be overridden by both configuration file or command line switches.
The configuration file is a plain json object, and it can be passed by using the command line switch `-f`.
The order on which the  configuration is evaluated is as depicted here:
```
+-------------------+    +--------------+   +----------+
|Nav builtin default|--->|conf json file|-->|CLI switch|
+-------------------+    +--------------+   +----------+
```
So the following example the built-in default is overridden with the conf.json and the arguments in the command line in the end override the final configuration.

```
$ ./nav -f conf.json -s kernel_init
```
The following is the command line switches list from the nav help.
```
$ ./nav -h
Command Help
App Name: nav
Descr: kernel symbol navigator
	-j	<v>	Force Json output with subsystems data
	-s	<v>	Specifies symbol
	-i	<v>	Specifies instance
	-f	<v>	Specifies config file
	-u	<v>	Forces use specified database userid
	-p	<v>	Forecs use specified password
	-d	<v>	Forecs use specified DBhost
	-p	<v>	Forecs use specified DBPort
	-m	<v>	Sets display mode 2=subsystems,1=all
	-h		This Help
```

## Sample configuration:
```
{
"DBURL":"dbs.hqhome163.com",
"DBPort":5432,
"DBUser":"alessandro",
"DBPassword":"<password>",
"DBTargetDB":"kernel_bin",
"Symbol":"__arm64_sys_getppid",
"Instance":1,
"Mode":1,
"Excluded": ["rcu_.*", "kmalloc", "kfree"],
"MaxDepth":0,
"Jout": "JsonOutputPlain"
}
```
Configuration is a file containing a JSON serialized conf object

|Field        |description                                                                                                |type    |Default value      |
|-------------|-----------------------------------------------------------------------------------------------------------|--------|-------------------|
|DBURL        |Host name ot ip address of the psql instance                                                               |string  |dbs.hqhome163.com  |
|DBPort       |tcp port where psql instance is listening                                                                  |integer |5432               |
|DBUser       |Valid username on the psql instance                                                                        |string  |alessandro         |
|DBPassword   |Valid password on the psql instance                                                                        |string  |<password>         |
|DBTargetDB   |The identifier for the DB containing symbols                                                               |string  |kernel_bin         |
|Symbol       |The symbol where start the navigation                                                                      |string  |NULL               |
|Instance     |The interesting symbols instance identifier                                                                |integer |1                  |
|Mode         |Mode of plotting: 1 symbols, 2 subsystems, 3 subsystems with labels,4 target subsystem isolation           |integer |2                  |
|Excluded     |List of symbols/subsystem not to be expanded                                                               |string[]|["rcu_.*"]         |
|MaxDepth     |Max number of levels to explore 0 no limit                                                                 |integer |0                  |
|Jout         |Type of output: GraphOnly, JsonOutputPlain, JsonOutputB64, JsonOutputGZB64                                 |enum    |GraphOnly          |
|Target_sybsys|List of subsys that need to be highlighted. if empty, only the subs that contain the start is highlighted  |string  |[]                 | 
