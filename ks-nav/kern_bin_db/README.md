# kern_bin_db - kernel source symbols extractor and DB builder


## Motivation
kern_bin_db has been developed to produce a structured kernel symbol and 
xrefs database (SQL).
Such database can help manual source code analysis, easing the source code, 
and allowing automated checks such as the scan for recursive function and 
similar queries.

## Build

kern_bin_db is implemented in Golang. Golang application are usually easy 
to build. In addition to this, it has very few dependencies other than the 
standard Golang packages.
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

The following is the command line switches list from the kern_bin_db help.
```
$ ./kern_bin_db -h
Kernel symbol fetcher
	-f	<v>	specifies json configuration file
	-s	<v>	Forces use specified strip binary
	-u	<v>	Forces use specified database userid
	-p	<v>	Forecs use specified password
	-d	<v>	Forecs use specified DBhost
	-o	<v>	Forecs use specified DBPort
	-n	<v>	Forecs use specified note (default 'upstream')
	-c		Checks dependencies
	-h		This Help

```
After having compiled a good default by specifying the Postgres database backend, start 
collecting symbols is as easy as issuing the following command:

```
$ ./kern_bin_db -f conf.json -n "Custom kernel from NXP bsp"
```

## Sample configuration:

```
{
"LinuxWDebug":          "vmlinux",
"LinuxWODebug":         "vmlinux.work",
"StripBin":             "/usr/bin/aarch64-linux-gnu-strip",
"DBURL":                "dbs.hqhome163.com",
"DBPort":               5432,
"DBUser":               "alessandro",
"DBPassword":           "<password>",
"DBTargetDB":           "kernel_bin",
"Maintainers_fn":       "MAINTAINERS",
"KConfig_fn":           "include/generated/autoconf.h",
"KMakefile":            "Makefile",
"Mode":                 15,
"Note":                 "upstream"
}
```

Configuration is a file containing a JSON serialized conf object

|Field         |description                                                                         |type    |Default value               |
|--------------|------------------------------------------------------------------------------------|--------|----------------------------|
|LinuxWDebug   |Linux image built with the debug symbols, input for the operation                   |string  |vmlinux                     |
|LinuxWODebug  |File created after the strip operation,and on which the R" tool operates on         |string  |vmlinux.work                |
|StripBin      |Executable that performs the strip operation to the selected architecture.          |string  |/usr/bin/strip              |
|DBURL         |Host name ot ip address of the psql instance                                        |string  |dbs.hqhome163.com           |
|DBPort        |tcp port where psql instance is listening                                           |integer |5432                        |
|DBUser        |Valid username on the psql instance                                                 |string  |alessandro                  |
|DBPassword    |Valid password on the psql instance                                                 |string  |<password>                  |
|DBTargetDB    |The identifier for the DB containing symbols                                        |string  |kernel_bin                  |
|Maintainers_fn|The path to MAINTAINERS file, typically in the kernel source tree                   |string  |MAINTAINERS                 |
|KConfig_fn    |The path to autoconf file containing the current build configuration                |string  |include/generated/autoconf.h|
|KMakefile     |The path to main kernel sourcecode Makefile, tipically sitting on the kernel tree / |string  |Makefile                    |
|Mode          |Mode of operation, use only for debug purpose. Defaluts at 15                       |integer |15                          |
|Note          |The string gets copied to the database. Consider a sort of tag for the data set     |string  |upstream                    |

**NOTE:** Defaults  are designed to make the tool work out of the box, if 
the executable is placed in the Linux kernel source code root directory.
