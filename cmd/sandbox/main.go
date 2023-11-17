package main

import (
	"github.com/msqtt/sb-judger/internal/sandbox"
)

func main() {
	err := sandbox.LaunchEntry()	
	if err != nil {
    panic(err)
	}
}
