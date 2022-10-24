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

// Const values for configuration mode field.
const (
	PRINT_ALL int	= 1
	PRINT_SUBSYS	= 2
)

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
}

type Edge struct {
	Caller	int
	Callee	int
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

	query:="select caller, callee from xrefs where caller =$1 and xref_instance_id_ref=$2"
	rows, err := db.Query(query, symbol_id, instance)
	if err!= nil {
		panic(err)
		}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&e.Caller, &e.Callee); err != nil {
			fmt.Println("this error hit1 ")
			return nil, err
			}
		successors,_ := get_entry_by_id(db, e.Callee, instance, cache.Entries)
		res=append(res, successors )
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
	var res string

	if res, ok := subsytems_cache[symbol]; ok {
		return res, nil
		}
	query:="select subsys_name from (select count(*)as cnt, subsys_name from tags where subsys_name in (select subsys_name from symbols, "+
		"tags where symbols.symbol_file_ref_id=tags.tag_file_ref_id and symbols.symbol_name=$1 and symbols.symbol_instance_id_ref=$2) "+
		"group by subsys_name order by cnt desc) as tbl;"

	rows, err := db.Query(query, symbol, instance)
	if err!= nil {
		panic(err)
		}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&res); err != nil {
			fmt.Println("this error hit1 ")
			return "", err
			}
		}
	subsytems_cache[symbol]=res
	return res, nil
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
func Navigate(db *sql.DB, symbol_id int, parent_dispaly string, visited *[]int, prod map[string]int, instance int, cache Cache, mode int, excluded []string, depth uint, maxdepth uint, dot_fmt string, output *string) {
	var tmp,s,l,ll,r	string
	var depthInc		uint	= 0

	*visited=append(*visited, symbol_id)
	l=parent_dispaly
	successors, err:=get_successors_by_id(db, symbol_id, instance, cache);
	successors=removeDuplicate(successors)
	if err==nil {
		for _, curr := range successors{
			entry, err := get_entry_by_id(db, curr.Sym_id, instance, cache.Entries)
			if err!=nil {
				r="Unknown";
				} else {
					r=entry.Symbol
					}
			switch mode {
				case PRINT_ALL:
					s=fmt.Sprintf(dot_fmt, l, r)
					ll=r
					depthInc = 1
					break
				case PRINT_SUBSYS:
					if tmp, err=get_subsys_from_symbol_name(db,r, instance, cache.SubSys); r!=tmp {
						if tmp != "" {
							r=tmp
							} else {
								r="UNDEFINED SUBSYSTEM"
								}
						}
					if l!=r {
						s=fmt.Sprintf(dot_fmt, l, r)
						depthInc = 1
						} else {
							s="";
							}
					ll=r
					break

				}
			if _, ok := prod[s]; ok {
				prod[s]++
				} else {
					prod[s]=1
					if s!="" {
						(*output)=(*output)+s
						}
					}

			if Not_in(*visited, curr.Sym_id){
				if not_exluded(entry.Symbol, excluded) && (maxdepth == 0 || (maxdepth > 0 && depth < maxdepth)){
					Navigate(db, curr.Sym_id, ll, visited, prod, instance, cache, mode, excluded, depth+depthInc, maxdepth, dot_fmt, output)
					}
				}
			}
		}
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
