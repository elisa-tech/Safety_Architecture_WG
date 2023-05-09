package main

import (
	"fmt"
)

type successorsT struct {
	es  []entry
	err error
}

type subsysnameT struct {
	subsysN string
	err     error
}

type numT struct {
	num int
	err error
}

type entryT struct {
	e   entry
	err error
}

type sqlMock struct {
	GetExploredSubsystemByNameValues map[string]string
	getSuccessorsByIdValues          map[string]successorsT
	getSubsysFromSymbolNameValues    map[string]subsysnameT
	sym2numValues                    map[string]numT
	symbSubsysValues                 map[string]subsysnameT
	getEntryByIdValues               map[string]entryT
}

func (d *sqlMock) LOADGetExploredSubsystemByNameValues(subs, explored string) {
	d.GetExploredSubsystemByNameValues[subs] = explored
}

func (d *sqlMock) LOADgetSuccessorsByIdValues(symbolId int, instance int, es []entry, err error) {
	key := fmt.Sprintf("%04x%02x", symbolId, instance)
	d.getSuccessorsByIdValues[key] = successorsT{es, err}
}

func (d *sqlMock) LOADgetSubsysFromSymbolNameValues(symbol string, instance int, subsysN string, err error) {
	key := fmt.Sprintf("%s%02x", symbol, instance)
	d.getSubsysFromSymbolNameValues[key] = subsysnameT{subsysN, err}
}

func (d *sqlMock) LOADsym2numValues(symb string, instance int, num int, err error) {
	key := fmt.Sprintf("%s%02x", symb, instance)
	d.sym2numValues[key] = numT{num, err}

}

func (d *sqlMock) LOADsymbSubsysValues(symblist []int, instance int, subsysN string, err error) {
	key := fmt.Sprintf("%04x%02x", murmurHash3(symblist), instance)
	d.symbSubsysValues[key] = subsysnameT{subsysN, err}
}

func (d *sqlMock) LOADgetEntryByIdValues(symbolId int, instance int, e entry, err error) {
	key := fmt.Sprintf("%04x%02x", symbolId, instance)
	d.getEntryByIdValues[key] = entryT{e, err}
}

func (d *sqlMock) init(arg interface{}) (err error) {
	d.GetExploredSubsystemByNameValues = make(map[string]string, 100)
	d.getSuccessorsByIdValues = make(map[string]successorsT, 100)
	d.getSubsysFromSymbolNameValues = make(map[string]subsysnameT, 100)
	d.sym2numValues = make(map[string]numT, 100)
	d.symbSubsysValues = make(map[string]subsysnameT, 100)
	d.getEntryByIdValues = make(map[string]entryT, 100)
	return nil
}

func (d *sqlMock) GetExploredSubsystemByName(subs string) string {
	debugIOPrintln("input subs=", subs)
	app := d.GetExploredSubsystemByNameValues[subs]
	debugIOPrintln("output =", subs)
	return app
}

func (d *sqlMock) getSuccessorsById(symbolId int, instance int) ([]entry, error) {
	debugIOPrintf("input symbolId=%d, instance=%d\n", symbolId, instance)
	key := fmt.Sprintf("%04x%02x", symbolId, instance)
	app1 := d.getSuccessorsByIdValues[key].es
	app2 := d.getSuccessorsByIdValues[key].err
	debugIOPrintf("output []entry=%+v, error=%s\n", app1, app2)
	return app1, app2
}

func (d *sqlMock) getSubsysFromSymbolName(symbol string, instance int) (string, error) {
	debugIOPrintf("input symbol=%s, instance=%d\n", symbol, instance)
	key := fmt.Sprintf("%s%02x", symbol, instance)
	app1 := d.getSubsysFromSymbolNameValues[key].subsysN
	app2 := d.getSubsysFromSymbolNameValues[key].err
	debugIOPrintf("output  string=%s, error=%s\n", app1, app2)
	return app1, app2
}

func (d *sqlMock) sym2num(symb string, instance int) (int, error) {
	debugIOPrintf("input symbol=%s, instance=%d\n", symb, instance)
	key := fmt.Sprintf("%s%02x", symb, instance)
	app1 := d.sym2numValues[key].num
	app2 := d.sym2numValues[key].err
	debugIOPrintf("output int=%d, error=%s\n", app1, app2)
	return app1, app2
}

func (d *sqlMock) symbSubsys(symblist []int, instance int) (string, error) {
	debugIOPrintf("input symblist=%+v, instance=%d\n", symblist, instance)
	key := fmt.Sprintf("%04x%02x", murmurHash3(symblist), instance)
	app1 := d.symbSubsysValues[key].subsysN
	app2 := d.symbSubsysValues[key].err
	debugIOPrintf("output string=%s, error=%s\n", app1, app2)
	return app1, app2
}

func (d *sqlMock) getEntryById(symbolId int, instance int) (entry, error) {
	debugIOPrintf("input symbolId=%d, instance=%d\n", symbolId, instance)
	key := fmt.Sprintf("%04x%02x", symbolId, instance)
	app1 := d.getEntryByIdValues[key].e
	app2 := d.getEntryByIdValues[key].err
	debugIOPrintf("output entry=%+v, error=%s\n", app1, app2)
	return app1, app2
}

func murmurHash3(arr []int) uint32 {
	const (
		seed = 0x9747b28c
		c1   = 0xcc9e2d51
		c2   = 0x1b873593
		r1   = 15
		r2   = 13
		m    = 5
		n    = 0xe6546b64
	)
	var hash uint32 = seed
	for i, val := range arr {
		k := uint32(val)
		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2
		hash ^= k
		hash = (hash << r2) | (hash >> (32 - r2))
		hash = (hash * m) + n + uint32(i*4)
	}
	hash ^= uint32(len(arr) * 4)
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16
	return hash
}
