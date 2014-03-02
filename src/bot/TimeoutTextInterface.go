// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"time"
)

func TimeoutTextInterface(i TextInterfaceController,
		ms uint32) TextInterfaceController {
	var timeout = time.Duration(ms)*time.Millisecond
	cmd_write := make(chan string, 1)
	ret_write := make(chan error)
	cmd_read := make(chan struct {}, 1)
	ret_read := make(chan struct {
		s string
		e error
	})
	cmd_exit := make(chan struct {}, 1)
	ret_exit := make(chan error)
	go func() {
		for {
			select {
				case s := <-cmd_write:
					ret_write <- i.Writeln(s)
				case <-cmd_read:
					v, err := i.Readln()
					ret_read <- struct {
						s string
						e error
					} {s: v, e: err}
				case <-cmd_exit:
					ret_exit <- i.Exit()
			}
		}
	}()
	var on_error func(error) error
	var has_error func() error
	{
		var glob_error error
		on_error = func(e error) error {
			if e == nil { return e }
			if glob_error == nil {
				cmd_exit <- struct {} {}
			}
			glob_error = e
			return e
		}
		has_error = func() error {
			return glob_error
		}
	}
	return TextInterfaceController {
		Writeln: func(s string) error {
			if has_error() != nil { return has_error() }
			cmd_write <- s
			select {
				case ret := <-ret_write:
					return on_error(ret)
				case <-time.After(timeout):
					return on_error(fmt.Errorf("timeout in Writeln()"))
			}
		},
		Readln: func() (string, error) {
			if has_error() != nil { return "", has_error() }
			cmd_read <- struct {} {}
			select {
				case ret := <-ret_read:
					return ret.s, on_error(ret.e)
				case <-time.After(timeout):
					return "", on_error(fmt.Errorf("timeout in Readln()"))
			}
		},
		Exit: func() error {
			if has_error() != nil { return has_error() }
			on_error(fmt.Errorf("interface is closed"))
			select {
				case ret := <-ret_exit:
					return on_error(ret)
				case <-time.After(timeout):
					return on_error(fmt.Errorf("timeout in Exit()"))
			}
		},
	}
}
