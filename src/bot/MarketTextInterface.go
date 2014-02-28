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
	parseOrders := func(lines []string) ([]Order, error) {
		var ret []Order
		for _, line := range lines {
			w := strings.Split(line, " ")
			if (len(w) == 4) && ((w[0] == "sell") || (w[0] == "buy") &&
					(w[2] == "for")) {
				v1, err := strconv.ParseUint(w[1], 10, 64)
				if err != nil { return nil, err }
				v2, err := strconv.ParseUint(w[3], 10, 64)
				if err != nil { return nil, err }
				if w[0] == "buy" {
					ret = append(ret,
						Order {AssetValue(v1), MoneyValue(v2), BUY})
				} else {
					ret = append(ret,
						Order {AssetValue(v1), MoneyValue(v2), SELL})
				}
				
			} else {
				return nil, fmt.Errorf("order list: %q", line)
			}
		}
		return ret, nil
	}
	return MarketController {
		GetTime: func() (Time, error) {
			market.Writeln("time")
			line := market.Readln()
			words := append(strings.Split(line, " "), "")
			v, err := strconv.ParseUint(words[1], 10, 64)
			if (words[0] != "time") || (err != nil) {
				return 0, fmt.Errorf("get time: %q", line)
			}
			return Time(v), nil
		},
		GetTotalBalance: func() (AssetValue, MoneyValue, error) {
			market.Writeln("totalbalance")
			line := market.Readln()
			words := append(strings.Split(line, " "), "", "")
			if words[0] != "totalbalance" {
				return 0, 0, fmt.Errorf("get balance: %q", line)
			}
			am1, err1 := strconv.ParseUint(words[1], 10, 64)
			am2, err2 := strconv.ParseUint(words[2], 10, 64)
			if (err1 != nil) || (err2 != nil) {
				return 0, 0, fmt.Errorf("get balance: %q", line)
			}
			return AssetValue(am1), MoneyValue(am2), nil
		},
		GetOrders: func() ([]Order, error) {
			market.Writeln("orders")
			line := market.Readln()
			if line != "orders:" {
				return nil, fmt.Errorf("get orders: %q", line)
			}
			var lines []string
			line = market.Readln()
			for line != "." {
				if len(line) < 1 {
					return nil, fmt.Errorf("get orders: end of input")
				}
				lines = append(lines, line)
				line = market.Readln()
			}
			ret, err := parseOrders(lines)
			return ret, err
		},
		NewOrder: func(o Order) error {
			cmd := fmt.Sprintf("%v %d for %d", o.t, o.asset, o.money)
			market.Writeln(cmd)
			if line := market.Readln(); line != ("ok "+cmd) {
				return fmt.Errorf("new order: %q", line)
			}
			return nil
		},
		CancelOrder: func(o Order) error {
			cmd := fmt.Sprintf("%v %d for %d", o.t, o.asset, o.money)
			market.Writeln("cancel "+cmd)
			if line := market.Readln(); line != ("ok cancel "+cmd) {
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
