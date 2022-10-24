/*
 * Copyright (c) 2022 Red Hat, Inc.
 * SPDX-License-Identifier: GPL-2.0-or-later
 */

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/goccy/go-graphviz"
	"os"
	"strings"
)

const (
	invalidOutput int = iota
	graphOnly
	jsonOutputPlain
	jsonOutputB64
	jsonOutputGZB64
)

const jsonOutputFMT string = "{\"graph\": \"%s\",\"graph_type\":\"%s\",\"symbols\": [%s]}"

var fmtDot = []string{
	"",
	"\"%s\"->\"%s\" \n",
	"\\\"%s\\\"->\\\"%s\\\" \\\\\\n",
	"\"%s\"->\"%s\" \n",
	"\"%s\"->\"%s\" \n",
}

var fmtDotHeader = []string{
	"",
	"digraph G {\nrankdir=\"LR\"\n",
	"digraph G {\\\\\\nrankdir=\"LR\"\\\\\\n",
	"digraph G {\nrankdir=\"LR\"\n",
	"digraph G {\nrankdir=\"LR\"\n",
}

var fmtDotNodeHighlightWSymb = "\"%[1]s\" [shape=record style=\"rounded,filled,bold\" fillcolor=yellow label=\"%[1]s|%[2]s\"]\n"
var fmtDotNodeHighlightWoSymb = "\"%[1]s\" [shape=record style=\"rounded,filled,bold\" fillcolor=yellow label=\"%[1]s\"]\n"

func opt2num(s string) int {
	var opt = map[string]int{
		"graphOnly":       graphOnly,
		"jsonOutputPlain": jsonOutputPlain,
		"jsonOutputB64":   jsonOutputB64,
		"jsonOutputGZB64": jsonOutputGZB64,
	}
	val, ok := opt[s]
	if !ok {
		return 0
	}
	return val
}

func decorateLine(l string, r string, adjm []adjM) string {
	var res = " [label=\""

	for _, item := range adjm {
		if (item.l.subsys == l) && (item.r.subsys == r) {
			tmp := fmt.Sprintf("%s([%s]%s),\\n", item.r.symbol, item.r.addressRef, item.r.sourceRef)
			if !strings.Contains(res, tmp) {
				res += fmt.Sprintf("%s([%s]%s),\\n", item.r.symbol, item.r.addressRef, item.r.sourceRef)
			}
		}
	}
	res += "\"]"
	return res
}

func decorate(dotStr string, adjm []adjM) string {
	var res string

	dotBody := strings.Split(dotStr, "\n")
	for i, line := range dotBody {
		split := strings.Split(line, "->")
		if len(split) == 2 {
			res = res + dotBody[i] + decorateLine(strings.TrimSpace(strings.ReplaceAll(split[0], "\"", "")), strings.TrimSpace(strings.ReplaceAll(split[1], "\"", "")), adjm) + "\n"
		}
	}
	return res
}

func do_graphviz(dot string, output_type outIMode) error {
	var buf bytes.Buffer
	var format graphviz.Format

	switch output_type {
	case oPNG:
		format = graphviz.PNG
	case oJPG:
		format = graphviz.JPG
	case oSVG:
		format = graphviz.SVG
	default:
		return errors.New("Unknown format")
	}

	graph, _ := graphviz.ParseBytes([]byte(dot))
	g := graphviz.New()
	defer func() {
		if err := graph.Close(); err != nil {
			panic(err)
		}
		g.Close()
	}()
	g.Render(graph, format, &buf)
	binary.Write(os.Stdout, binary.LittleEndian, buf.Bytes())
	return nil
}

func generateOutput(d Datasource, conf *configuration) (string, error) {
	var graphOutput string
	var jsonOutput string
	var prod = map[string]int{}
	var visited []int
	var entryName string
	var output string
	var adjm []adjM

	start, err := d.sym2num(conf.Symbol, conf.Instance)
	if err != nil {
		fmt.Println("Symbol not found")
		return "", err
	}

	graphOutput = fmtDotHeader[opt2num(conf.Jout)]
	entry, err := d.getEntryById(start, conf.Instance)
	if err != nil {
		return "", err
	} else {
		entryName = entry.symbol
	}
	startSubsys, _ := d.getSubsysFromSymbolName(entryName, conf.Instance)
	if startSubsys == "" {
		startSubsys = SUBSYS_UNDEF
	}

	if (conf.Mode == printTargeted) && len(conf.TargetSubsys) == 0 {
		targSubsysTmp, err := d.getSubsysFromSymbolName(conf.Symbol, conf.Instance)
		if err != nil {
			panic(err)
		}
		conf.TargetSubsys = append(conf.TargetSubsys, targSubsysTmp)
	}

	navigate(d, start, node{startSubsys, entryName, "entry point", "0x0"}, conf.TargetSubsys, &visited, &adjm, prod, conf.Instance, conf.Mode, conf.ExcludedAfter, conf.ExcludedBefore, 0, conf.MaxDepth, fmtDot[opt2num(conf.Jout)], &output)

	if (conf.Mode == printSubsysWs) || (conf.Mode == printTargeted) {
		output = decorate(output, adjm)
	}

	graphOutput += output
	if conf.Mode == printTargeted {
		for _, i := range conf.TargetSubsys {
			if d.GetExploredSubsystemByName(conf.Symbol) == i {
				graphOutput += fmt.Sprintf(fmtDotNodeHighlightWSymb, i, conf.Symbol)
			} else {
				graphOutput += fmt.Sprintf(fmtDotNodeHighlightWoSymb, i)
			}
		}
	}
	graphOutput += "}"

	symbdata, err := d.symbSubsys(visited, conf.Instance)
	if err != nil {
		return "", err
	}

	switch opt2num(conf.Jout) {
	case graphOnly:
		jsonOutput = graphOutput
	case jsonOutputPlain:
		jsonOutput = fmt.Sprintf(jsonOutputFMT, graphOutput, conf.Jout, symbdata)
	case jsonOutputB64:
		b64dot := base64.StdEncoding.EncodeToString([]byte(graphOutput))
		jsonOutput = fmt.Sprintf(jsonOutputFMT, b64dot, conf.Jout, symbdata)

	case jsonOutputGZB64:
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write([]byte(graphOutput)); err != nil {
			return "", errors.New("gzip failed")
		}
		if err := gz.Close(); err != nil {
			return "", errors.New("gzip failed")
		}
		b64dot := base64.StdEncoding.EncodeToString(b.Bytes())
		jsonOutput = fmt.Sprintf(jsonOutputFMT, b64dot, conf.Jout, symbdata)

	default:
		return "", errors.New("unknown output mode")
	}
	return jsonOutput, nil
}

func main() {
	conf, err := argsParse(cmdLineItemInit())
	if err != nil {
		if err.Error() != "dummy" {
			fmt.Println(err.Error())
		}
		printHelp(cmdLineItemInit())
		os.Exit(-1)
	}
	if opt2num(conf.Jout) == 0 {
		fmt.Printf("Unknown mode %s\n", conf.Jout)
		os.Exit(-2)
	}
	t := connectToken{conf.DBDriver, conf.DBDSN}
	d := &SqlDB{}
	err = d.init(&t)
	if err != nil {
		panic(err)
	}

	output, err := generateOutput(d, &conf)
	if err != nil {
		fmt.Println("Internal error", err)
		os.Exit(-3)
	}
	if conf.Graphviz != oText {
		err = do_graphviz(output, conf.Graphviz)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(output)
	}
}
