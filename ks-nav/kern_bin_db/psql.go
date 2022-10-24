	/*
	 * ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	 *
	 *   Name: kern_bin_db - Kernel source code analysis tool database creator
	 *   Description: Parses kernel source tree and binary images and builds the DB
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
	_ "github.com/lib/pq"
)

// Sql connection configuration
type Connect_token struct{
	Host    string
	Port    int
	User    string
	Pass    string
	Dbname  string
}

// Connects the target db and returns the handle
func Connect_db(t *Connect_token) (*sql.DB){
	fmt.Println("connect")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", (*t).Host, (*t).Port, (*t).User, (*t).Pass, (*t).Dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err!= nil {
		panic(err)
		}
	fmt.Println("connected")
	return db
}

// Executes insert queries
func Insert_data(db *sql.DB, query string, test bool){

	if !test {
		_ , err := db.Exec(query)
		if err!= nil {
			fmt.Println("##################################################")
			fmt.Println(query)
			fmt.Println("##################################################")
			panic(err)
			}
		} else {
			fmt.Println(query)
			}
}

// Executes insert query for instance table and returns the id allocated
func Insert_datawID(db *sql.DB, query string) int{
	var res		int

	_ , err := db.Exec(query)
	if err!= nil {
		fmt.Println("##################################################")
		fmt.Println(query)
		fmt.Println("##################################################")
		panic(err)
		}
	rows, err := db.Query("SELECT currval('instances_instance_id_seq');")
	if err != nil {
		panic(err)
		}
	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&res); err != nil {
		panic(err)
		}

	return res
}
