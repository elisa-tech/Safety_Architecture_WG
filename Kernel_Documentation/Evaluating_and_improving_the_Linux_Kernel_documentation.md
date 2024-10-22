
# Evaluating and improving the Linux Kernel documentation


# Goal and Rationale

The goal of this document is to analyze the current templates and guidelines that are available in the [Linux Kernel documentation](https://docs.kernel.org/), evaluate if and how they fulfill architecture and design aspects required by functional safety, define improvements also in consideration of maintenance challenges deriving from a continuously evolving code baseline.

The main reason for such an analysis is to improve the baseline of available documentation evidence for SW Architectures and Design aspects, starting from the process of documenting these, in a way that can simplify the adoption of Linux in a functional safety system.

There are however other reasons for documenting SW Architectures and Design aspects that go beyond functional safety; such as:



* Facilitating developers in understanding the code and in contributing;
* Facilitating the development of test cases and verify their correctness and completeness;
* Facilitating integrators to understand the Linux expected behavior and assess the adequacy of it with respect to their specific context of use.


# ISO-26262 architecture and unit design objectives vs  Linux Documentation repositories

With respect to ISO-26262 If we consider the “objectives” associated with “software architectural design” (part 6.7) and with “software unit design and implementation” (part 6.8), we respectively have the set of objectives defined below. For each of them, let's analyze the available Linux documentation and the role in achieving the objective.



* “Software architectural design” objectives:
    1.  “*to develop a software architectural design that satisfies the software safety requirements and the other software requirements;*”. <BR /><BR />While Linux is not developed according to a predefined set of requirements, an integrator wishing to use Linux in a specific use case, including safety use cases, must be able to assess the capabilities of Linux to meet his requirements allocated to the operating system. In this regard it is important to have comprehensive documentation that explains the Linux behavior from an integrator point of view.  
        To this extend the available Linux documentation is:
        - [The Linux Manpage project](https://git.kernel.org/pub/scm/docs/man-pages/man-pages.git/);
        - [The Linux kernel user-space API guide](https://docs.kernel.org/userspace-api/index.html#the-linux-kernel-user-space-api-guide);
        - [The Linux kernel user’s and administrator’s guide](https://docs.kernel.org/admin-guide/index.html);
        - [The GNU C Library Reference Manual](https://sourceware.org/glibc/manual/latest/html_mono/libc.html);
    
    2.  “to verify that the software architectural design is suitable to satisfy the software safety requirements with the required ASIL; and”<BR /><BR />While from an integrator point of view the Documentation covering the Linux expected behavior at the user interface is needed to assess the adequacy of Linux to meet the integrator’s allocated requirements; in order to verify the Linux internal architecture and design against the above mentioned expected behavior we need documentation covering such internal elements. To this extend the available Linux documentation is:
        *   [Internal API manuals](https://docs.kernel.org/index.html#internal-api-manuals)
        *   Firmware interfaces documentation:
            *    [Devicetree](https://www.kernel.org/doc/html/latest/devicetree/index.html)
            *    [ACPI](https://www.kernel.org/doc/html/latest/firmware-guide/index.html)

        Another place where the expected behavior of the code can be looked up is the Kernel mailing list. In fact this is the place where developers informally discuss the code and the reason behind a specific implementation proposal. Especially patchsets’ cover letters usually provide an overview of what the code does and can be helpful in understanding it.<BR />
[cregit](https://cregit.linuxsources.org/) is a web based tool that can be used to retrieve commits associated with a set of code lines and their respective mailing list discussions. By entering the cregit website and choosing a specific Kernel version, it is then possible to navigate the source tree, select the file of interest, select then the source code line of interest, click on these, and look at the respective mailing list discussions (that also include the cover letter, if there is one).

    3.  “to support the implementation and verification of the software.”
        *   All documentation mentioned above is required to support the software implementation and verification. In fact, documentation mapping to goal a. Is needed to verify Linux at the user space interfaces, while documentation mapping to goal b. is needed to verify Linux at the level of its internal components. Accordingly tests can be written at different hierarchical levels.
        *   To support the software implementation the section “[Working with the development community](https://docs.kernel.org/index.html#working-with-the-development-community)” describes the Kernel development process and the guidelines for submitting patches
        *   Finally the section “[Development tools for the kernel](https://docs.kernel.org/dev-tools/index.html#development-tools-for-the-kernel)” describes a list of tools that can support the implementation and verification of the Kernel code.<BR />

* “Software unit design and implementation” objectives
    1.  “to develop a software unit design in accordance with the software architectural design, the design criteria and the allocated software requirements which supports the implementation and verification of the software unit; and”<BR /><BR />With respect to this goal the overall behavior of subsystems and drivers can be documented in the [Internal API manuals](https://docs.kernel.org/index.html#internal-api-manuals). Such information can be documented either natively in the RST files (e.g. “Documentation/watchdog/watchdog-api.rst”) or referenced from the source code by using the [Overview](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html#overview-documentation-comments) documentation's comments. The last ones seem to be the most appropriate ones to be used to explain the behavior of a specific subsystem or driver and how it contributes to achieve the overall behavior at the user-space API (since they allow to track which source code is undocumented)<BR /><BR />
    2.  “to implement the software units as specified."<BR /><BR />With respect to this objective the [Internal API manuals](https://docs.kernel.org/index.html#internal-api-manuals) includes kernel-doc headers that allow documenting single [functions](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html#function-documentation) as well as [members](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html#members) of relevant data structures.<BR />On top of this, all considerations made under goal c. of part 6.7 are also valid for this goal


# Architecture and design aspects currently covered by the Linux documentation templates

In this chapter we’ll evaluate the static and dynamic aspects of SW Architecture and Design required by Functional Safety standards (including ISO-26262) and if/how the information associated with such aspects can be mapped to the existing Kernel documentation.


* **Static Design Aspects**:
    *  **main components**.<BR />
This aspect must be considered from different perspectives:
        *   The integrator’s perspective.<BR />
The main components can be considered the APIs that are available to the integrator.
            *    User-space APIs. As mentioned above these are documented mainly in:
                *     [The Linux Manpage project](https://git.kernel.org/pub/scm/docs/man-pages/man-pages.git/):
                *     [The Linux kernel user-space API guide](https://docs.kernel.org/userspace-api/index.html#the-linux-kernel-user-space-api-guide)
                *     [The Linux kernel user’s and administrator’s guide](https://docs.kernel.org/admin-guide/index.html)
                *     [The GNU C Library Reference Manual](https://sourceware.org/glibc/manual/latest/html_mono/libc.html)<BR />
            *    Kernel-space APIs for loadable modules: these are represented by the symbols exported by EXPORT_SYMBOL or EXPORT_SYMBOL_GPL and shall be documented following the guidelines in “[Writing kernel-doc comments](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html)”.

        *   The developer’s perspective.<BR />
For the Linux Kernel there are different aspects that should be considered.
            *    Drivers/Subsystems: the Linux Kernel is already partitioned into Drivers and Subsystems in the [MAINTAINERs](https://github.com/torvalds/linux/blob/master/MAINTAINERS) file. It seems natural to follow such partitioning that also rules the hierarchy of maintainers and the Kernel development process. As mentioned above the overall behavior of such components can be documented by following the [Overview](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html#overview-documentation-comments) documentation guidelines.
            *    Single functions’: each driver or subsystem is composed of multiple functions; these can be documented following the guidelines in “[Writing kernel-doc comments](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html)” (they use the same formalism used by exported symbols)

    *  **SW/HW interfaces**.<BR />
Also this aspect shall be considered from different perspectives:
        *   Integrators perspective: the relevant interfaces and associated behavior are defined in the bullet above “main components” and the respective documentation is hyperlinked.
        *   Developer’s perspective: from a developer’s perspective this aspect is tricky. While the SW, HW, FW interfaces exposed by the single drivers and subsystems can theoretically be documented in the respective [Overview](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html#overview-documentation-comments) section, as of today there is no guideline asking to do so.<BR />So in order to analyze dependencies the main source of information is the code itself and hence a possibility to effectively parse and visualize dependencies is the [ks-nav](https://github.com/elisa-tech/ks-nav) tool.
        *   Firmware: the Linux Kernel uses firmware interfaces to parse HW information used to configure the single drivers and hence modifying the overall Kernel behavior. The kernel provides two types of Firmware interfaces:
            *    [Devicetree](https://www.kernel.org/doc/html/latest/devicetree/index.html)
            *    [ACPI](https://www.kernel.org/doc/html/latest/firmware-guide/index.html)
            
            For Devicetree HW specific information, the binding files under “Documentation/devicetree/bindings” contain a description how the provided HW parameters according to the [schema](https://www.kernel.org/doc/html/latest/devicetree/bindings/writing-schema.html) template; however today such binding files are not compiled as part of the Linux Kernel documentation.<BR />For ACPI HW specific information the reference tables and their meaning is described in the ACPI specification documents ([ACPI Specification 6.5](https://uefi.org/specs/ACPI/6.5/) at the time of writing)

    *  **SW/HW resources**.<BR />
Internal SW/HW resources relevant for the integrator are exposed through the user-space APIs mentioned above with respective documentation hyperlinks.<BR /><BR />
From a developer perspective relevant SW/HW resources can be documented following the [members](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html#members) documentation template, however such template does not enforce specifying which subsystem or driver uses them. In order to fill this gap the [ks-nav](https://github.com/elisa-tech/ks-nav) tool has the capability to highlight global and static data that is used by and eventually shared between subsystems and drivers.

* **Dynamic Aspects**.<BR />
Functional Safety standards usually demand the following aspects to be documented 
    *  chain of events/behaviour
    *  the logical sequence of data processing
    *  control flow and data flow
    *  temporal constraints

    Unfortunately today all aspects above are not enforced from a Documentation template point of view. These aspects can be documented, especially in the [Overview](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html#overview-documentation-comments) sections or in the drivers’ or subsystems’ specific RST files, however there are no specific fields mapping to them.<BR /><BR />

# Assessment of the currently available documentation

An initial assessment of the currently available documentation was performed to verify the following statement from the section “[Writing kernel-doc comments](https://www.kernel.org/doc/html/latest/doc-guide/kernel-doc.html)” in the Linux Kernel Documentation: “Every function that is exported to loadable modules using EXPORT_SYMBOL or EXPORT_SYMBOL_GPL should have a kernel-doc comment”.<BR />
In order to have a kernel-doc header associated with a symbol come up in the official Kernel Documentation the following should happen:

1. The corresponding file exporting the symbol shall be cross-referenced from the rst file in the Kernel documentation folder;
2. The symbol shall be documented with a Kernel doc header.

The initial assessment focused on point 1. and consisted in running the “[find_unreferenced_sources.sh](https://github.com/elisa-tech/Safety_Architecture_WG/blob/main/Kernel_Documentation_tools/find_unreferenced_sources.sh)” script on the Linux kernel source tree (see the output [here](./exported_symbols_files_cross_reference_report)). From the report we have 6069 source files exporting symbols and only 547 of them are cross-referenced from the Linux Kernel documentation; therefore today more than 90% of files to be cross-referenced for symbols documentation are not.<BR />
As a final note, the initial assessment did not cover eventually duplicated information between different repos (e.g. symbol specific info already documented in the md files); this is something that will be covered later after a suitable documentation process has been defined.<BR />
Given the significant gap already found, a [session](https://lpc.events/event/18/contributions/1894/) was held at Linux Plumbers 2024 to explain the software design aspects that are expected to be documented by quality and safety international standards, where the Linux documentation process stands in that regard and what should be done in order to fill the gaps. The agreed next steps are as follows:

1. Defining a template (or templates) that can be used to document Kernel code / User APIs: sysfs APIs have been proposed as a good starting point to elaborate on the aspects that can be covered by such template;
2. Getting such a template (or templates) reviewed internally in ELISA first and then externally by a restricted group of maintainers (Jonathan Corbet, Thomas Gleixner, Steve Rostedt, Bjorn Andersson);
3. Starting to document a specific subsystem using such templates (possibly mm) and getting the documentation upstream;
4. Having LF to kick-off a new project for extending such documentation and for defining the process to enforce it to the rest of the Kernel.<BR /><BR />

Note: with respect to point 2. above it is crucial to have the documentation as close as possible to the code and to have a tool integrated in the Kernel CI/CD process that is able to detect non-compliant development (e.g. documentation not matching the function prototypes or critical design elements present in the code but not documented).
