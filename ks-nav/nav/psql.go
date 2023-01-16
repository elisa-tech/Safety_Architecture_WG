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
	"database/sql"
	"fmt"
	"strings"
	"regexp"
	"errors"
	"sort"
	_ "github.com/lib/pq"
)

type OutMode int64

// Const values for configuration mode field.
const (
	_		OutMode		= iota
	PRINT_ALL
	PRINT_SUBSYS
	PRINT_SUBSYS_WS
	PRINT_TARGETED
	OutModeLast
)
const 	SUBSYS_UNDEF	= "The REST"

// parent node
type Node struct {
	Subsys		string
	Symbol		string
	Source_ref	string
	Address_ref	string
}
type AdjM struct {
	l		Node
	r		Node
}

// Sql connection configuration
type Connect_token struct{
	Host	string
	Port	int
	User	string
	Pass	string
	Dbname	string
}

type Entry struct {
	Sym_id		int
	Symbol		string
	Exported	bool
	Type		string
	Subsys		[]string
	Fn		string
	Source_ref	string
	Address_ref	string
}

type Edge struct {
	Caller		int
	Callee		int
	Source_ref	string
	Address_ref	string
}

type Cache struct {
	Successors	map[int][]Entry
	Entries		map[int]Entry
	SubSys		map[string]string
}

var check int = 0

var chached int = 0

// Connects the target db and returns the handle
func Connect_db(t *Connect_token) (*sql.DB){
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", (*t).Host, (*t).Port, (*t).User, (*t).Pass, (*t).Dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err!= nil {
		panic(err)
		}
	return db
}

// Returns function details from a given id
func get_entry_by_id(db *sql.DB, symbol_id int, instance int,cache map[int]Entry)(Entry, error){
	var e		Entry
	var s		sql.NullString

	if e, ok := cache[symbol_id]; ok {
		return e, nil
		}

	query:="select symbol_id, symbol_name, subsys_name, file_name from "+
		"(select * from symbols, files where symbols.symbol_file_ref_id=files.file_id and symbols.symbol_instance_id_ref=$2) as dummy "+
		"left outer join tags on dummy.symbol_file_ref_id=tags.tag_file_ref_id where symbol_id=$1 and symbol_instance_id_ref=$2"
	rows, err := db.Query(query, symbol_id, instance)
	if err!= nil {
		panic(err)
		}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&e.Sym_id, &e.Symbol, &s, &e.Fn); err != nil {
			fmt.Println("this error hit3")
			fmt.Println(err)
			return e, err
			}
		if s.Valid {
			e.Subsys=append(e.Subsys, s.String)
			}
		}
	if err = rows.Err(); err != nil {
		fmt.Println("this error hit")
		return e, err
		}
	cache[symbol_id]=e
	return e, nil
}

// Returns the list of successors (called function) for a given function
func get_successors_by_id(db *sql.DB, symbol_id int, instance int, cache Cache )([]Entry, error){
	var e		Edge
	var res		[]Entry


	if res, ok := cache.Successors[symbol_id]; ok {
		return res, nil
		}

	query:="select caller, callee, source_line, ref_addr from xrefs where caller =$1 and xref_instance_id_ref=$2"
	rows, err := db.Query(query, symbol_id, instance)
	if err!= nil {
		panic(err)
		}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&e.Caller, &e.Callee, &e.Source_ref, &e.Address_ref); err != nil {
			fmt.Println("get_successors_by_id: this error hit1 ", err)
			return nil, err
			}
		successor,_ := get_entry_by_id(db, e.Callee, instance, cache.Entries)
		successor.Source_ref = e.Source_ref
		successor.Address_ref = e.Address_ref
		res=append(res, successor )
		}
	if err = rows.Err(); err != nil {
		fmt.Println("this error hit2 ")
		return nil, err
		}
	cache.Successors[symbol_id]=res
	return res, nil
}

// Return id an item is already in the list
func Not_in(list []int, v int) bool {

	for _, a := range list {
		if a == v {
			return false
			}
		}
	return true
}

// Removes duplicates resulting by the exploration of a call tree
func removeDuplicate(list []Entry) []Entry {

	sort.SliceStable(list, func(i, j int) bool { return list[i].Sym_id < list[j].Sym_id })
	allKeys := make(map[int]bool)
	res := []Entry{}
	for _, item := range list {
		if _, value := allKeys[item.Sym_id]; !value {
			allKeys[item.Sym_id] = true
			res = append(res, item)
			}
		}
	return res
}

// Given a function returns the lager subsystem it belongs
func get_subsys_from_symbol_name(db *sql.DB, symbol string, instance int, subsytems_cache map[string]string)(string, error){
	var ty,sub string

	if res, ok := subsytems_cache[symbol]; ok {
		return res, nil
		}
	query:="select (select symbol_type from symbols where symbol_name=$1 and symbol_instance_id_ref=$2) as type, subsys_name from "+
		"(select count(*) as cnt, subsys_name from tags where subsys_name in (select subsys_name from symbols, "+
		"tags where symbols.symbol_file_ref_id=tags.tag_file_ref_id and symbols.symbol_name=$1 and symbols.symbol_instance_id_ref=$2) "+
		"group by subsys_name order by cnt desc) as tbl;"

	rows, err := db.Query(query, symbol, instance)
	if err!= nil {
		fmt.Println(query)
		fmt.Println(symbol, instance)
		panic(err)
		}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ty,&sub); err != nil {
			fmt.Println("get_subsys_from_symbol_name: this error hit1 ")
			return "", err
			}
		}

	if ty== "indirect" {
		sub=ty
		}
	subsytems_cache[symbol]=sub
	return sub, nil
}

// Returns the id of a given function name
func sym2num(db *sql.DB, symb string, instance int)(int, error){
	var 	res	int=0
	var	cnt 	int=0
	query:="select symbol_id from symbols where symbols.symbol_name=$1 and symbols.symbol_instance_id_ref=$2"
	rows, err := db.Query(query, symb, instance)
	if err!= nil {
		panic(err)
		}
	defer rows.Close()

	for rows.Next() {
		cnt++
		if err := rows.Scan(&res,); err != nil {
			fmt.Println("this error hit7")
			fmt.Println(err)
			return res, err
			}
		}
	if cnt!=1 {
		return res, errors.New("id is not unique")
		}
	return res, nil
}

// Checks if a given function needs to be explored
func not_exluded(symbol string, excluded []string)bool{

	for _,s:=range excluded{
		if match,_ :=regexp.MatchString(s, symbol); match{
			return false
			}
		}
	return true
}

// Computes the call tree of a given function name
// TODO: refactory needed:
// What is the problem: too many args.
// suggestion: New version with input and output structs.
func Navigate(db *sql.DB, symbol_id int, parent_dispaly Node, targets []string, visited *[]int, AdjMap *[]AdjM, prod map[string]int, instance int, cache Cache, mode OutMode, excluded_after []string, excluded_before []string, depth int, maxdepth int, dot_fmt string, output *string) {
	var tmp, s		string
	var l, r, ll		Node
	var depthInc		int	= 0

	*visited=append(*visited, symbol_id)
	l=parent_dispaly
	successors, err:=get_successors_by_id(db, symbol_id, instance, cache);
	if mode == PRINT_ALL {
		successors=removeDuplicate(successors)
		}
	if err==nil {
		for _, curr := range successors{
			if not_exluded(curr.Symbol, excluded_before) {
				r.Symbol=curr.Symbol
				r.Source_ref = curr.Source_ref
				r.Address_ref = curr.Address_ref
				tmp, _ =get_subsys_from_symbol_name(db,r.Symbol, instance, cache.SubSys)
				if tmp=="" {
					r.Subsys=SUBSYS_UNDEF
					}

				switch mode {
					case PRINT_ALL:
						s=fmt.Sprintf(dot_fmt, l.Symbol, r.Symbol)
						ll=r
						depthInc = 1
						break
					case PRINT_SUBSYS, PRINT_SUBSYS_WS, PRINT_TARGETED:
						if tmp, err=get_subsys_from_symbol_name(db,r.Symbol, instance, cache.SubSys); r.Subsys!=tmp {
							if tmp != "" {
								r.Subsys=tmp
								} else {
									r.Subsys=SUBSYS_UNDEF
									}
							}

						if l.Subsys!=r.Subsys {
							s=fmt.Sprintf(dot_fmt, l.Subsys, r.Subsys)
							*AdjMap=append(*AdjMap, AdjM{l,r})
							depthInc = 1
							} else {
								s="";
								}
						ll=r
						break
					default:
						panic(mode)
					}
				if _, ok := prod[s]; ok {
					prod[s]++
					} else {
						prod[s]=1
						if s!="" {
							if (mode != PRINT_TARGETED) || (intargets(targets, l.Subsys,r.Subsys)) {
								(*output)=(*output)+s
								}
							}
						}

				if Not_in(*visited, curr.Sym_id){
					if (not_exluded(curr.Symbol, excluded_after) || not_exluded(curr.Symbol, excluded_before)) && (  maxdepth == 0  ||  (  (maxdepth > 0)   &&   (depth < maxdepth) ) ){
						Navigate(db, curr.Sym_id, ll, targets, visited, AdjMap, prod, instance, cache, mode, excluded_before, excluded_before, depth+depthInc, maxdepth, dot_fmt, output)
						}
					}
				}
			}
		}
}

//returns true if one of the nodes n1, n2 is a target node
func intargets(targets []string, n1 string, n2 string) bool {

	for _, t := range targets {
		if (t == n1) || (t == n2) {
			return true
			}
		}
	return false
}

// Returns the subsystem list associated with a given function name
func symbSubsys(db *sql.DB, symblist []int, instance int, cache Cache,)(string, error){
	var out	string
	var res	string

	for _, symbid := range symblist {
		//resolve sybm
		symb, _ := get_entry_by_id(db, symbid, instance, cache.Entries)
		out=out+fmt.Sprintf("{\"FuncName\":\"%s\", \"subsystems\":[", symb.Symbol)
		query:=fmt.Sprintf("select subsys_name from tags where tag_file_ref_id= (select symbol_file_ref_id from symbols where symbol_id=%d);", symbid)
		rows, err := db.Query(query)
			if err!= nil {
				return "", errors.New("symbSubsys query failed")
				}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&res,); err != nil {
				return "", errors.New("symbSubsys query browsing failed")
				}
			out=out+fmt.Sprintf("\"%s\",", res)
			}
		out=strings.TrimSuffix(out, ",")+"]},"
		}
	out=strings.TrimSuffix(out, ",")
	return out, nil
}
