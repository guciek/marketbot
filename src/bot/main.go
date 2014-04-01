// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import(
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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
		cmd := exec.Command(executable, args...)
		cmd.Stderr = os.Stderr
		var write func(line string) error
		var close func() error
		{
			stdin, err := cmd.StdinPipe()
			if err != nil { return TextInterfaceController {}, err }
			write = func(line string) error {
				_, err := stdin.Write([]byte(line+"\n"))
				return err
			}
			close = func() error {
				if err := stdin.Close(); err != nil { return err }
				if err := cmd.Wait(); err != nil { return err }
				return nil
			}
		}
		var read func() (string, error)
		{
			stdout, err := cmd.StdoutPipe()
			if err != nil { return TextInterfaceController {}, err }
			stdout_r := bufio.NewReader(stdout)
			read = func() (string, error) {
				line, err := stdout_r.ReadString('\n')
				if err != nil { return "", err }
				if len(line) >= 2 { return line[:len(line)-1], nil }
				return "", nil
			}
		}
		if err := cmd.Start(); err != nil {
			return TextInterfaceController {}, err
		}
		return TextInterfaceController {write, read, close}, nil
	}

	args := os.Args
	if len(args) < 2 {
		path := strings.Split(args[0], "/")
		n := path[len(path)-1]
		fmt.Fprintf(
			os.Stderr,
			"\nUsage:\n"+
			"\t"+n+" [options] <market>\n"+
			"\t"+n+" [options] -test 123abc,4.5xyz\n"+
			"\n"+
			"Strategy options:\n"+
			"\t-balance abc/xyz,1.23abc <market>\n"+
			"\t-natural 12.3abc/xyz,1.23abc <market>\n"+
			"\t-buy 12.3abc/xyz,1.23abc <market>\n"+
			"\t-sell 23.4abc/xyz,1.23abc <market>\n"+
			"\n"+
			"Other options:\n"+
			"\t-fee 0.12%%              Compensate for transaction fee\n"+
			"\t-gain 1.23%%             Percentage gain for each transaction\n"+
			"\t-cashout 123abc         Subtract this amount if reached\n"+
			"\n",
		)
		return
	}

	var planner OrderPlanner
	{
		params := make(map[string][]string)
		for len(args) >= 2 {
			if args[1][0] != '-' { break }
			if len(args) < 3 { panic("empty value of "+args[1]) }
			name, val := args[1][1:], args[2]
			if len(val) < 1 { panic("invalid value of -"+name) }
			args = args[2:]
			params[name] = append(params[name], val)
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
	exec = TimeoutTextInterface(exec, 60*1000)

	Run(
		MarketTextInterface(exec), planner,
		func(msg string) { fmt.Fprintf(os.Stderr, "%s\n", msg) },
		func() bool { return interrupted },
	)
	fmt.Fprint(os.Stderr, "End.\n")
}
