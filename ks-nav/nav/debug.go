package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
)

const (
	debugNone                 = 0
	debugIO                   = 1
	debugQuery                = 2
	debugDecorateFunctionName = 14
	debugAddFunctionName      = 15
)

//const DebugLevel uint32 = debugIO | (1<<debugAddFunctionName-1)
const DebugLevel uint32 = debugNone

func debugIOPrintf(format string, a ...interface{}) (int, error) {
	var s, randomSeed string
	var n int
	var err error

	randomSeed = ""
	if DebugLevel&(1<<(debugIO-1)) != 0 {
		if DebugLevel&(1<<(debugAddFunctionName-1)) != 0 {
			pc, _, _, ok := runtime.Caller(1)
			s = "?"
			if ok {
				fn := runtime.FuncForPC(pc)
				if fn != nil {
					s = fn.Name()
				}
			}
			if DebugLevel&(1<<(debugDecorateFunctionName-1)) != 0 {
				randomSeed = "." + strconv.Itoa(rand.Intn(10000))
			}
			newformat := "[" + s + randomSeed + "] " + format
			n, err = fmt.Printf(newformat, a...)
		} else {
			n, err = fmt.Printf(format, a...)
		}
		return n, err
	}
	return 0, nil
}

func debugIOPrintln(a ...interface{}) (int, error) {
	if DebugLevel&(1<<(debugIO-1)) != 0 {
		n, err := fmt.Println(a...)
		return n, err
	}
	return 0, nil
}

func debugQueryPrintf(format string, a ...interface{}) (int, error) {
	if DebugLevel&(1<<(debugQuery-1)) != 0 {
		n, err := fmt.Printf(format, a...)
		return n, err
	}
	return 0, nil
}

func debugQueryPrintln(a ...interface{}) (int, error) {
	if DebugLevel&(1<<(debugQuery-1)) != 0 {
		n, err := fmt.Println(a...)
		return n, err
	}
	return 0, nil
}
