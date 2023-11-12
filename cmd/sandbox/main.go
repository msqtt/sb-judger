package main

import (
	"errors"
	"os"
	"strconv"

	"github.com/msqtt/sb-judger/internal/sandbox"
)

func main() {
	if len(os.Args) < 2 {
		panic(errors.New("not enough args"))
	}
	arg := os.Args[1]
	switch arg {
	case strconv.FormatUint(sandbox.ArgInit, 10):
		err := sandbox.InitEntry()
		if err != nil {
			panic(err)
		}
	case strconv.FormatUint(sandbox.ArgLaunch, 10):
		err := sandbox.LaunchEntry()	
		if err != nil {
			panic(err)
		}
	default:
	panic(errors.New("invalid args"))
	}
}
