	/*
	 * ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	 *
	 *   Name: nav - Kernel source code analysis tool
	 *   Description: Extract call trees for kernel API
	 *
	 *   Author: Alessandro Carminati <acarmina@redhat.com>
	 *   Author: Maurizio Papini <mpapini@redhat.com>
	 *
	 * ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	 *
	 *   Copyright (c) 2022 Red Hat, Inc. All rights reserved.
	 *
	 *   This copyrighted material is made available to anyone wishing
	 *   to use, modify, copy, or redistribute it subject to the terms
	 *   and conditions of the GNU General Public License version 2.
	 *
	 *   This program is distributed in the hope that it will be
	 *   useful, but WITHOUT ANY WARRANTY; without even the implied
	 *   warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
	 *   PURPOSE. See the GNU General Public License for more details.
	 *
	 *   You should have received a copy of the GNU General Public
	 *   License along with this program; if not, write to the Free
	 *   Software Foundation, Inc., 51 Franklin Street, Fifth Floor,
	 *   Boston, MA 02110-1301, USA.
	 *
	 * ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	 */

package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"database/sql"
	"encoding/base64"
	)

const (
	GraphOnly int		= 1
	JsonOutputPlain int	= 2
	JsonOutputB64 int	= 3
	JsonOutputGZB64 int	= 4
	)

const JsonOutputFMT string = "{\"graph\": \"%s\",\"graph_type\":\"%s\",\"symbols\": [%s]}"

var fmt_dot = []string {
		"",
		"\"%s\"->\"%s\" \n",
		"\\\"%s\\\"->\\\"%s\\\" \\\\\\n",
		"\"%s\"->\"%s\" \n",
		"\"%s\"->\"%s\" \n",
		}

var fmt_dot_header = []string {
		"",
		"digraph G {\n",
		"digraph G {\\\\\\n",
		"digraph G {\n",
		"digraph G {\n",
		}

func opt2num(s string) int{
        var opt = map[string]int{
                "GraphOnly":1,
                "JsonOutputPlain":2,
                "JsonOutputB64":3,
                "JsonOutputGZB64":4,
                }
        val, ok := opt[s]
        if !ok {
                return 0
        }
        return val
}

func generate_output(db *sql.DB, conf *configuration) (string, error){
	var	GraphOutput	string
	var 	JsonOutput	string
	var	prod =		map[string]int{}
	var	visited		[]int
	var	entry_name	string
	var	output		string

	cache := make(map[int][]Entry)
	cache2 := make(map[int]Entry)
	cache3 := make(map[string]string)

	start, err:=sym2num(db, (*conf).Symbol, (*conf).Instance)
	if err!=nil{
		fmt.Println("symbol not found")
		return "", err
		}

	GraphOutput=fmt_dot_header[opt2num((*conf).Jout)]
	entry, err := get_entry_by_id(db, start, (*conf).Instance, cache2)
		if err!=nil {
			entry_name="Unknown";
			return "",err
			} else {
				entry_name=entry.Symbol
				}

	Navigate(db, start, entry_name, &visited, prod, (*conf).Instance, Cache{cache, cache2, cache3}, (*conf).Mode, (*conf).Excluded, 0, (*conf).MaxDepth, fmt_dot[opt2num((*conf).Jout)], &output)
	GraphOutput=GraphOutput+output
	GraphOutput=GraphOutput+"}"

	symbdata, err := symbSubsys(db, visited, (*conf).Instance, Cache{cache, cache2, cache3})
	if err != nil{
		return "",err
		}

	switch opt2num((*conf).Jout) {
		case GraphOnly:
			JsonOutput=GraphOutput
		case JsonOutputPlain:
			JsonOutput=fmt.Sprintf(JsonOutputFMT, GraphOutput, (*conf).Jout, symbdata)
		case JsonOutputB64:
			b64dot:=base64.StdEncoding.EncodeToString([]byte(GraphOutput))
			JsonOutput=fmt.Sprintf(JsonOutputFMT, b64dot, (*conf).Jout, symbdata)

		case JsonOutputGZB64:
			var b bytes.Buffer
			gz := gzip.NewWriter(&b)
			if _, err := gz.Write([]byte(GraphOutput)); err != nil {
				return "", errors.New("gzip failed")
				}
			if err := gz.Close(); err != nil {
				return "", errors.New("gzip failed")
				}
			b64dot:=base64.StdEncoding.EncodeToString(b.Bytes())
			JsonOutput=fmt.Sprintf(JsonOutputFMT, b64dot, (*conf).Jout, symbdata)

		default:
			return "", errors.New("Unknown output mode")
	}
	return JsonOutput, nil
}

func main() {

        conf, err := args_parse(cmd_line_item_init())
        if err!=nil {
		if err.Error() != "dummy"{
			fmt.Println(err.Error())
			}
                print_help(cmd_line_item_init());
                os.Exit(-1)
                }
	if opt2num(conf.Jout)==0 {
		fmt.Printf("unknown mode %s\n", conf.Jout)
		os.Exit(-2)
		}
	t:=Connect_token{ conf.DBURL, conf.DBPort,  conf.DBUser, conf.DBPassword, conf.DBTargetDB}
	db:=Connect_db(&t)

	output, err := generate_output(db, &conf)
	if err!=nil{
		fmt.Println("internal error", err)
		os.Exit(-3)
		}
	fmt.Println(output)

}
