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
	"strconv"
	"fmt"
	"os"
	"errors"
	"encoding/json"
	"io/ioutil"
	)

var conf_fn = "conf.json"

const app_name	string="App Name: nav"

const app_descr	string="Descr: kernel symbol navigator"

type Arg_func func(*configuration, []string) (error)

// Command line switch elements
type cmd_line_items struct {
	id		int
	Switch		string
	Help_srt	string
	Has_arg		bool
	Needed		bool
	Func		Arg_func
}

// Represents the application configuration
type configuration struct {
	DBURL		string
	DBPort		int
	DBUser		string
	DBPassword	string
	DBTargetDB	string
	Symbol		string
	Instance	int
	Mode		OutMode
	Excluded_before	[]string
	Excluded_after	[]string
	Target_sybsys	[]string
	MaxDepth	int
	Jout		string
	cmdlineNeeds	map[string] bool
}

// Instance of default configuration values
var	Default_config  configuration = configuration{
	DBURL:		"dbs.hqhome163.com",
	DBPort:		5432,
	DBUser:		"alessandro",
	DBPassword:	"<password>",
	DBTargetDB:	"kernel_bin",
	Symbol:		"",
	Instance:	0,
	Mode:		PRINT_SUBSYS,
	Excluded_before:	[]string{},
	Excluded_after:	[]string{},
	Target_sybsys:	[]string{},
	MaxDepth:	0,		//0: no limit
	Jout:		"GraphOnly",
	cmdlineNeeds:	map[string] bool{},
	}

// Inserts a commandline item item, which is composed by:
// * switch string
// * switch descriptio
// * if the switch requires an additiona argument
// * a pointer to the function that manages the switch
// * the configuration that gets updated
func push_cmd_line_item(Switch string, Help_str string, Has_arg bool, Needed bool, Func Arg_func, cmd_line *[]cmd_line_items){
	*cmd_line = append(*cmd_line, cmd_line_items{id: len(*cmd_line)+1, Switch: Switch, Help_srt: Help_str, Has_arg: Has_arg, Needed: Needed, Func: Func})
}

// This function initializes configuration parser subsystem
// Inserts all the commandline switches suppported by the application
func cmd_line_item_init() ([]cmd_line_items){
	var res	[]cmd_line_items

	push_cmd_line_item("-j", "Force Json output with subsystems data",	true,  false,	func_outtype,	&res)
	push_cmd_line_item("-s", "Specifies symbol",				true,  true,	func_symbol,	&res)
	push_cmd_line_item("-i", "Specifies instance",				true,  true,	func_instance,	&res)
	push_cmd_line_item("-f", "Specifies config file",			true,  false,	func_jconf,	&res)
	push_cmd_line_item("-u", "Forces use specified database userid",	true,  false,	func_DBUser,	&res)
	push_cmd_line_item("-p", "Forecs use specified password",		true,  false,	func_DBPass,	&res)
	push_cmd_line_item("-d", "Forecs use specified DBhost",			true,  false,	func_DBHost,	&res)
	push_cmd_line_item("-p", "Forecs use specified DBPort",			true,  false,	func_DBPort,	&res)
	push_cmd_line_item("-m", "Sets display mode 2=subsystems,1=all",	true,  false,	func_Mode,	&res)
	push_cmd_line_item("-x", "Specify Max depth in call flow exploration",	true,  false,	func_depth,	&res)
	push_cmd_line_item("-h", "This Help",					false, false,	func_help,	&res)

	return res
}

func func_help		(conf *configuration,fn []string)		(error){
	return errors.New("Command Help")
}

func func_outtype(conf *configuration, jout []string)			(error){
	(*conf).Jout=jout[0]
	return nil
}

func func_jconf		(conf *configuration,fn []string)		(error){
	jsonFile, err := os.Open(fn[0])
	if err != nil {
		return err
		}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	err=json.Unmarshal(byteValue, conf)
	if err != nil {
		return err
		}
	return nil
}

func func_symbol	(conf *configuration, fn []string)	(error){
	(*conf).Symbol=fn[0]
	return nil
}

func func_DBUser	(conf *configuration, user []string)	(error){
	(*conf).DBUser=user[0]
	return nil
}

func func_DBPass	(conf *configuration, pass []string)	(error){
	(*conf).DBPassword=pass[0]
	return nil
}

func func_DBHost	(conf *configuration, host []string)	(error){
	(*conf).DBURL=host[0]
	return nil
}

func func_DBPort	(conf *configuration, port []string)	(error){
	s, err := strconv.Atoi(port[0])
	if err!=nil {
		return err
		}
	(*conf).DBPort=s
	return nil
}

func func_depth		(conf *configuration, depth []string)	(error){
	s, err := strconv.Atoi(depth[0])
	if err!=nil {
		return err
		}
	if s<0 {
		return errors.New("Depth must be >= 0")
		}
	(*conf).MaxDepth=s
	return nil
}

func func_instance	(conf *configuration, instance []string)    (error){
	s, err := strconv.Atoi(instance[0])
	if err!=nil {
		return err
		}
	(*conf).Instance=s
	return nil
}

func func_Mode		(conf *configuration, mode []string)    (error){
	s, err := strconv.Atoi(mode[0])
	if err!=nil {
		return err
		}
	if OutMode(s)<PRINT_ALL || OutMode(s)>=OutModeLast {
		return errors.New("unsupported mode")
		}
	(*conf).Mode=OutMode(s)
	return nil
}

// Uses commandline args to generate the help string
func print_help(lines []cmd_line_items){

	fmt.Println(app_name)
	fmt.Println(app_descr)
	for _,item := range lines{
		fmt.Printf(
			"\t%s\t%s\t%s\n",
			item.Switch,
			func (a bool)(string){
				if a {
					return "<v>"
					}
				return ""
			}(item.Has_arg),
			item.Help_srt,
			)
		}
}

// Used to parse the command line and generate the command line
func args_parse(lines []cmd_line_items)(configuration, error){
	var	extra		bool=false;
	var	conf		configuration=Default_config
	var 	f		Arg_func

	for _, item := range lines{
		if item.Needed {
			conf.cmdlineNeeds[item.Switch]=false
			}
		}

	for _, os_arg := range os.Args[1:] {
		if !extra {
			for _, arg := range lines{
				if arg.Switch==os_arg {
					if arg.Needed {
						conf.cmdlineNeeds[arg.Switch]=true
						}
					if arg.Has_arg{
						f=arg.Func
						extra=true
						break
						}
					err := arg.Func(&conf, []string{})
					if err != nil {
						return Default_config, err
						}
					}
				}
			continue
			}
		if extra{
			err := f(&conf,[]string{os_arg})
			if err != nil {
				return Default_config, err
				}
			extra=false
			}

		}
	if extra {
		 return  Default_config, errors.New("Missing switch arg")
		}

	res:=true
	for _, element := range conf.cmdlineNeeds {
		res = res && element
		}
	if res {
		return	conf, nil
		}
	return Default_config, errors.New("Missing needed arg")
}
