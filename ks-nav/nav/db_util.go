/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"fmt"
	"regexp"
	"sort"
)

type Datasource interface {
	init(arg interface{}) (err error)
	GetExploredSubsystemByName(subs string) string
	getSuccessorsById(symbolId int, instance int) ([]entry, error)
	getSubsysFromSymbolName(symbol string, instance int) (string, error)
	sym2num(symb string, instance int) (int, error)
	symbSubsys(symblist []int, instance int) (string, error)
	getEntryById(symbolId int, instance int) (entry, error)
}

type outMode int64
type outIMode int64

// Const values for configuration mode field.
const (
	_ outMode = iota
	printAll
	printSubsys
	printSubsysWs
	printTargeted
	OutModeLast
)

const (
	_ outIMode = iota
	oText
	oPNG
	oJPG
	oSVG
	OutIModeLast
)

const SUBSYS_UNDEF = "The REST"

// Parent node.
type node struct {
	subsys     string
	symbol     string
	sourceRef  string
	addressRef string
}
type adjM struct {
	l node
	r node
}

type entry struct {
	symbol     string
	fn         string
	sourceRef  string
	addressRef string
	subsys     []string
	symId      int
}

type edge struct {
	sourceRef  string
	addressRef string
	caller     int
	callee     int
}

// Return id an item is already in the list.
func notIn(list []int, v int) bool {

	for _, a := range list {
		if a == v {
			return false
		}
	}
	return true
}

// Removes duplicates resulting by the exploration of a call tree.
func removeDuplicate(list []entry) []entry {

	sort.SliceStable(list, func(i, j int) bool { return list[i].symId < list[j].symId })
	allKeys := make(map[int]bool)
	var res []entry
	for _, item := range list {
		if _, value := allKeys[item.symId]; !value {
			allKeys[item.symId] = true
			res = append(res, item)
		}
	}
	return res
}

// Checks if a given function needs to be explored.
func notExcluded(symbol string, excluded []string) bool {
	for _, s := range excluded {
		if match, _ := regexp.MatchString(s, symbol); match {
			return false
		}
	}
	return true
}

// Computes the call tree of a given function name
// TODO: refactory needed:
// What is the problem: too many args.
// suggestion: New version with input and output structs.
func navigate(d Datasource, symbolId int, parentDispaly node, targets []string, visited *[]int, AdjMap *[]adjM, prod map[string]int, instance int, mode outMode, excludedAfter []string, excludedBefore []string, depth int, maxdepth int, dotFmt string, output *string) {
	var tmp, s string
	var l, r, ll node

	*visited = append(*visited, symbolId)
	l = parentDispaly
	successors, err := d.getSuccessorsById(symbolId, instance)
	if mode == printAll {
		successors = removeDuplicate(successors)
	}
	if err == nil {
		for _, curr := range successors {
			depthInc := 0
			if notExcluded(curr.symbol, excludedBefore) {
				r.symbol = curr.symbol
				r.sourceRef = curr.sourceRef
				r.addressRef = curr.addressRef
				tmp, _ = d.getSubsysFromSymbolName(r.symbol, instance)
				if tmp == "" {
					r.subsys = SUBSYS_UNDEF
				}

				switch mode {
				case printAll:
					s = fmt.Sprintf(dotFmt, l.symbol, r.symbol)
					ll = r
					depthInc = 1
				case printSubsys, printSubsysWs, printTargeted:
					if tmp, _ = d.getSubsysFromSymbolName(r.symbol, instance); r.subsys != tmp {
						if tmp != "" {
							r.subsys = tmp
						} else {
							r.subsys = SUBSYS_UNDEF
						}
					}

					if l.subsys != r.subsys {
						s = fmt.Sprintf(dotFmt, l.subsys, r.subsys)
						*AdjMap = append(*AdjMap, adjM{l, r})
						depthInc = 1
					} else {
						s = ""
					}
					ll = r
				default:
					panic(mode)
				}
				if _, ok := prod[s]; ok {
					prod[s]++
				} else {
					prod[s] = 1
					if s != "" {
						if (mode != printTargeted) || (intargets(targets, l.subsys, r.subsys)) {
							*output = (*output) + s
						}
					}
				}
				if notIn(*visited, curr.symId) {
					if (notExcluded(curr.symbol, excludedAfter) && notExcluded(curr.symbol, excludedBefore)) && (maxdepth == 0 || ((maxdepth > 0) && (depth+depthInc < maxdepth))) {
						navigate(d, curr.symId, ll, targets, visited, AdjMap, prod, instance, mode, excludedAfter, excludedBefore, depth+depthInc, maxdepth, dotFmt, output)
					}
				}
			}
		}
	}
}

// returns true if one of the nodes n1, n2 is a target node.
func intargets(targets []string, n1 string, n2 string) bool {

	for _, t := range targets {
		if (t == n1) || (t == n2) {
			return true
		}
	}
	return false
}
