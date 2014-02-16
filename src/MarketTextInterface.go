// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"strconv"
	"strings"
)

func MarketTextInterface(market TextInterfaceController) MarketController {
	return MarketController {
		GetPrice: func() (PriceValue, Time, error) {
			market.Writeln("price")
			line := market.Readln()
			words := append(strings.Split(line, " "), "", "")
			v1, err1 := strconv.ParseUint(words[1], 10, 64)
			v2, err2 := strconv.ParseUint(words[2], 10, 64)
			if (words[0] != "price") || (err1 != nil) || (err2 != nil) {
				return 0, 0, fmt.Errorf("get price: %q", line)
			}
			return PriceValue(v1), Time(v2), nil
		},
		GetOrders: func() (ret OrderList, err error) {
			market.Writeln("orders")
			line := market.Readln()
			if line != "orders:" {
				err = fmt.Errorf("get orders: %q", line)
				return
			}
			var lines string
			line = market.Readln()
			for line != "." {
				if len(line) < 1 {
					err = fmt.Errorf("get orders: end of input")
					return
				}
				lines += line+"\n"
				line = market.Readln()
			}
			ret.LoadData(lines)
			return
		},
		NewOrder: func(o Order) error {
			cmd := fmt.Sprintf("%v %d %d", o.t, o.price, o.cost)
			market.Writeln(cmd)
			if line := market.Readln(); line != ("ok "+cmd) {
				return fmt.Errorf("new order: %q", line)
			}
			return nil
		},
		CancelOrder: func(o Order) error {
			cmd := fmt.Sprintf("cancel %v %d %d", o.t, o.price, o.cost)
			market.Writeln(cmd)
			if line := market.Readln(); line != ("ok "+cmd) {
				return fmt.Errorf("cancel order: %q", line)
			}
			return nil
		},
		Wait: func() error {
			market.Writeln("wait")
			if line := market.Readln(); line != "ok wait" {
				return fmt.Errorf("wait: %q", line)
			}
			return nil
		},
		Exit: func() error {
			market.Writeln("exit")
			if err := market.Exit(); err != nil { return err }
			return nil
		},
	}
}
