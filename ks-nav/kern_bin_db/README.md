# kern_bin_db - kernel source symbols extractor and DB builder


## Motivation
kern_bin_db has been developed to produce a structured kernel symbol and 
xrefs database (SQL).
Such database can help manual source code analysis, easing the source code, 
and allowing automated checks such as the scan for recursive functions and 
similar queries.

## Prerequisites
In order to develop this project you need to install the radare2 dev package:

#### Fedora
```azure
$ sudo dnf install radare2-devel
```
#### Arch Linux 
```azure
$ sudo pacman -S radare2
```
## Build

kern_bin_db is implemented in Golang. Golang applications are usually easy 
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
        -f      <v>     specifies json configuration file
        -s      <v>     Forces use specified strip binary
        -e      <v>     Forces to use a specified DB Driver (i.e. postgres, mysql or sqlite3)
        -d      <v>     Forces to use a specified DB DSN
        -n      <v>     Forecs use specified note (default 'upstream')
        -c              Checks dependencies
        -h              This Help

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
"DBDriver":             "postgres",
"DBDSN":                "host=dbs.hqhome163.com port=5432 user=alessandro password=<password> dbname=kernel_bin sslmode=disable",
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
|LinuxWODebug  |File created after the strip operation, and on which the R" tool operates on        |string  |vmlinux.work                |
|StripBin      |Executable that performs the strip operation to the selected architecture.          |string  |/usr/bin/strip              |
|DBDriver      |Name of DB engine driver, i.e. postgres, mysql or sqlite3                           |string  |postgres                    |
|DBDSN         |DSN in the engine specific format                                                   |string  |See Note                    |
|Maintainers_fn|The path to MAINTAINERS file, typically in the kernel source tree                   |string  |MAINTAINERS                 |
|KConfig_fn    |The path to autoconf file containing the current build configuration                |string  |include/generated/autoconf.h|
|KMakefile     |The path to main kernel sourcecode Makefile, typically sitting on the kernel tree / |string  |Makefile                    |
|Mode          |Mode of operation, use only for debug purpose. Defaults to 15                       |integer |15                          |
|Note          |The string gets copied to the database. Consider a sort of tag for the data set     |string  |upstream                    |

**NOTE:** Defaults are designed to make the tool work out of the box, if 
the executable is placed in the Linux kernel source code root directory. 

Currently the default DBDSN value is set to: 
host=dbs.hqhome163.com port=5432 user=alessandro password=<password> dbname=kernel_bin sslmode=disable
to be consistent with the previous default configuration.

For sqlite the DBDSN can be as simple as just the filename of the database
file: **kernel-symbol.db**.

In any case, before starting kern_bin_db, the DB schema needs to be created
manually with one of the provided sql-files.

# DSN examples

| DBMS          | Example                                                                                                |
|---------------|--------------------------------------------------------------------------------------------------------|
| MySQL/MariaDB | alessandro:<password>@tcp(dbs.hqhome163.com:3306)/kernel_bin?multiStatements=true                      |
| Postgresql    | host=dbs.hqhome163.com port=5432 user=alessandro password=<password> dbname=kernel_bin sslmode=disable |

# TODO
* currently, kern_bin_db scans only the kernel binary image. Loadable 
modules are not considered, which means that symbols defined in modules 
are not considered. A useful extension would be to scan modules in a 
given kernel tree and add symbols to the database.
* in relation to the previous, it would be useful to track the symbol
source (kernel image/module). The database does not currently support 
this feature, which means that the database schema needs to be adapted.
* Indirect calls ate difficult to follow, but the place in the source 
where they are called can be easily spotted. It makes sense to note 
this place for future automated/manual investigations. In addition, 
it makes also sense to reserve a place in the same row of the same 
table for a field that links another table, where possible references 
are kept.
In this way, the future nav version can expand indirect calls as a 
finite number of calls.
