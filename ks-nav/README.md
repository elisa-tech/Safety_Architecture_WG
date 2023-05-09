# ks-nav

ks-nav is a set of two tools developed to support ELISA activities.
The first tool, **kern_bin_db** parses a Linux kernel image and produces data 
that feeds a backend Postgresql database. The data contains symbols, source 
files, build time configuration, and subsytems. The data gets organized into 
a database with the following schema:

```
+---------+       +---------+       +-----+
|instances|--<^>--| symbols |--<^>--|files|
+---------+       +---------+       +-----+
    |              |       |           |
   <^>             |       |          <^>
    |              |   ^   |           |
+-------+          | /   \ |        +----+
|configs|          +<xrefs>+        |tags|
+-------+            \   /          +----+
                       V
```
The second tool, **nav** uses the database to produce calltree graphs images.

# Dependencies

Required dependencies:
* Postgresql
* Radare 2 v5.7.8
* addr2line 2.38
* binutil strip for the selected target architecture
* golang 1.18
* optional UPX for compress binary

Golang used packages:

* github.com/VividCortex/ewma v1.1.1 
* github.com/cheggaaa/pb/v3 v3.1.0 
* github.com/elazarl/addr2line v0.0.0-20160815095215-325e5a3858c1 
* github.com/fatih/color v1.10.0 
* github.com/lib/pq v1.10.6
* github.com/mattn/go-colorable v0.1.8 
* github.com/mattn/go-isatty v0.0.12 
* github.com/mattn/go-runewidth v0.0.12 
* github.com/radareorg/r2pipe-go v0.2.1 
* github.com/rivo/uniseg v0.2.0 
* golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6 


# tools usage

For specific directions on tool usage and build, refer to specific README file

# Authors
* Alessandro Carminati <acarmina@redhat.com>
* Maurizio Papini <mpapini@redhat.com>
