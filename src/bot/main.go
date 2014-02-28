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
			"\t"+n+" -balance '$1.23' <market>\n"+
			"\t"+n+" -natural '$1.23' -target 1.234 <market>\n"+
			"\n"+
			"Other options:\n"+
			"\t-fee 0.12%%       Market fee, deducted from transaction gain\n"+
			"\t-spread 1.2%%     Increase spread between buy and sell orders\n"+
			"\t-test '$1234'    Show target orders and exit\n"+
			"\n",
		)
		return
	}

	var planner OrderPlanner
	{
		params := make(map[string]int64)
		parseFloat100000 := func(val string) int64 {
			s := strings.Split(val, ".")
			if len(s) > 2 { panic("could not parse number: "+val) }
			var ret int64 = 0
			{
				v, err := strconv.ParseInt(s[0], 10, 64)
				if err != nil { panic("could not parse number: "+val) }
				ret += v*100000
			}
			if len(s) == 2 {
				if len(s[1]) > 5 { panic("too many decimal places: "+val) }
				for len(s[1]) < 5 { s[1] = s[1]+"0" }
				v, err := strconv.ParseInt(s[1], 10, 64)
				if err != nil { panic("could not parse number: "+val) }
				ret += v
			}
			return ret
		}
		for len(args) >= 3 {
			if args[1][0] != '-' { break }
			name, val := args[1][1:], args[2]
			if len(val) < 1 { panic("invalid value of -"+name) }
			args = args[2:]
			if val[0] == '$' {
				name += "_money"
				val = val[1:]
			}
			if val[0] == '@' {
				name += "_asset"
				val = val[1:]
			}
			if val[len(val)-1] == '%' {
				params[name] = parseFloat100000(val[0:len(val)-1])/100
			} else {
				params[name] = parseFloat100000(val)
			}
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
