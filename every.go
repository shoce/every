/*
history:
022/1204 v1

go mod init github.com/shoce/every
go get -a -u -v
go mod tidy

GoFmt
GoBuildNull
GoBuild

*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	NL = "\n"
)

func main() {

	var err error

	var Duration time.Duration
	var StopAfter time.Duration

	var cmd string
	var args []string
	var Command *exec.Cmd

	if len(os.Args) < 4 || (os.Args[2] != "--" && os.Args[3] != "--") {
		fmt.Fprintf(os.Stderr, "usage: every duration [stopafter] -- command [args]"+NL)
		os.Exit(1)
	}

	Duration, err = time.ParseDuration(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "time.ParseDuration %s: %v"+NL, os.Args[1], err)
		os.Exit(1)
	}

	if os.Args[2] != "--" {
		StopAfter, err = time.ParseDuration(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "time.ParseDuration %s: %v"+NL, os.Args[2], err)
			os.Exit(1)
		}
	}

	if os.Args[2] == "--" {
		cmd = os.Args[3]
		args = os.Args[4:]
	} else if os.Args[3] == "--" {
		cmd = os.Args[4]
		args = os.Args[5:]
	} else {
		fmt.Fprintf(os.Stderr, "error: there must be '--' before the command"+NL)
		os.Exit(1)
	}

	StartTime := time.Now()

	for {
		Command = exec.Command(cmd, args...)
		Command.Stdin, Command.Stdout, Command.Stderr = os.Stdin, os.Stdout, os.Stderr
		fmt.Fprintf(os.Stderr, NL+"%s:"+NL, Command)
		err = Command.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, NL+"err: %v"+NL, err)
		}

		fmt.Fprintf(os.Stderr, NL+"sleeping %v"+NL, Duration)
		time.Sleep(Duration)

		if StopAfter > 0 && time.Now().Sub(StartTime) > StopAfter {
			fmt.Fprintf(os.Stderr, NL+"stopping after %v"+NL, StopAfter)
			break
		}
	}

}
