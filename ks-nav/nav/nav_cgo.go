//go:build CGO
// +build CGO

/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

/*
#cgo LDFLAGS: -ldotparser -Ldot_parser/
#include <stdio.h>
#include "dot_parser/dot.tab.h"
extern int dotparse(void);
extern void set_input_string(const char* in);
extern void end_lexical_scan(void);
int parse_string(const char* in) {
	set_input_string(in);
	int rv = dotparse();
	end_lexical_scan();
	return rv;
}
*/
import "C"

func valid_dot(dot string) bool {
	if C.parse_string(C.CString(dot)) == 0 {
		return true
	}
	return false
}
