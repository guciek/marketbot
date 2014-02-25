// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import(
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}()

	execTextInterface := func(executable string,
			args ...string) (TextInterfaceController, error) {
		var empty TextInterfaceController
		cmd := exec.Command(executable, args...)
		cmd.Stderr = os.Stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil { return empty, err }
		stdout_r := bufio.NewReader(stdout)
		stdin, err := cmd.StdinPipe()
		if err != nil { return empty, err }
		if err := cmd.Start(); err != nil { return empty, err }
		return TextInterfaceController {
			func(line string) {
				_, err := stdin.Write([]byte(line+"\n"))
					if err != nil { panic(err) }
			},
			func() string {
				for {
					line, err := stdout_r.ReadString('\n')
					if err == io.EOF { return "" }
					if err != nil { panic(err) }
					if len(line) >= 2 { return line[:len(line)-1] }
				}
			},
			func() error {
				if stdin.Close() != nil { panic(err) }
				return cmd.Wait()
			},
		}, nil
	}

	args := os.Args
	if len(args) < 2 {
		path := strings.Split(args[0], "/")
		n := path[len(path)-1]
		fmt.Fprintf(
			os.Stderr,
			"\nUsage:\n"+
			"\t"+n+" -natural <target> -for <money> <market-interface>\n"+
			"\n",
		)
		return
	}

	var planner OrderPlanner
	{
		params := make(map[string]int64)
		for len(args) >= 3 {
			if args[1][0] != '-' { break }
			v, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil { panic("could not parse number: "+args[2]) }
			params[args[1][1:]] = v
			args = args[2:]
		}
		var err error
		planner, err = PlanOrders(params)
		if err != nil { panic(err.Error()) }
	}

	if len(args) < 2 { panic("market interface not specified") }

	interrupted := false
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		<-signals
		interrupted = true
	}()

	exec, err := execTextInterface(args[1], args[2:]...)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		return
	}

	Run(
		MarketTextInterface(exec), planner,
		func(msg string) { fmt.Fprintf(os.Stderr, "%s\n", msg) },
		func() bool { return interrupted },
	)
	fmt.Fprint(os.Stderr, "End.\n")
}
