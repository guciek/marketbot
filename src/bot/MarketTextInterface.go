// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"money"
	"strconv"
	"strings"
)

func MarketTextInterface(market TextInterfaceController) MarketController {
	parseOrder := func(line string) (Order, error) {
		var oid string
		var part []string
		{
			p := strings.SplitN(line, " ", 2)
			if len(p) != 2 {
				return Order {}, fmt.Errorf("invalid order description: %q", line)
			}
			oid = p[0]
			part = strings.Split(p[1], " for ")
		}
		if (len(part) != 2) || (len(part[0]) < 5) ||
				(part[0][0:4] != "buy ") {
			return Order {}, fmt.Errorf("invalid order description: %q", line)
		}
		v1, err := money.ParseMoney(part[0][4:len(part[0])])
		if err != nil { return Order {}, err }
		v2, err := money.ParseMoney(part[1])
		if err != nil { return Order {}, err }
		return Order {buy: v1, sell: v2, id: oid}, nil
	}
	check := uint64(1)
	return MarketController {
		GetTime: func() (Time, error) {
			market.Writeln("time")
			var words []string
			{
				line, err := market.Readln()
				if err != nil { return 0, err }
				words = append(strings.Split(line, " "), "")
			}
			if words[0] != "time" {
				return 0, fmt.Errorf("get time: %q", words[0])
			}
			v, err := strconv.ParseUint(words[1], 10, 64)
			return Time(v), err
		},
		GetTotalBalance: func() ([]money.Money, error) {
			market.Writeln("totalbalance")
			{
				line, err := market.Readln()
				if err != nil { return nil, err }
				if line != "totalbalance:" {
					return nil, fmt.Errorf("get balance: %q", line)
				}
			}
			sum := make(map[string]money.Money)
			{
				line, err := market.Readln()
				if err != nil { return nil, err }
				for line != "." {
					if len(line) < 1 {
						return nil, fmt.Errorf("get balance: empty line")
					}
					v, err1 := money.ParseMoney(line)
					if err1 != nil { return nil, err }
					sum[v.Currency()] = sum[v.Currency()].Add(v)
					line, err = market.Readln()
					if err != nil { return nil, err }
				}
			}
			ret := make([]money.Money, 0, len(sum))
			for _, s := range sum {
				if ! s.IsZero() {
					ret = append(ret, s)
				}
			}
			return ret, nil
		},
		GetOrders: func() ([]Order, error) {
			market.Writeln("orders")
			{
				line, err := market.Readln()
				if err != nil { return nil, err }
				if line != "orders:" {
					return nil, fmt.Errorf("get orders: %q", line)
				}
			}
			ret := make([]Order, 0, 10)
			{
				line, err := market.Readln()
				if err != nil { return nil, err }
				for line != "." {
					var o Order
					o, err = parseOrder(line)
					if err != nil { return nil, err }
					ret = append(ret, o)
					line, err = market.Readln()
					if err != nil { return nil, err }
				}
			}
			return ret, nil
		},
		NewOrder: func(o Order) error {
			cmd := fmt.Sprintf("buy %v for %v", o.buy, o.sell)
			market.Writeln(cmd)
			line, err := market.Readln()
			if err != nil { return err }
			if line != "ok buy" {
				return fmt.Errorf("new order: %q", line)
			}
			return nil
		},
		CancelOrder: func(oid string) error {
			market.Writeln("cancel "+oid)
			line, err := market.Readln()
			if err != nil { return err }
			if line != "ok cancel" {
				return fmt.Errorf("cancel order: %q", line)
			}
			return nil
		},
		CheckConnection: func() error {
			echo := fmt.Sprintf("echo CHECK#%d", check)
			check++
			market.Writeln(echo)
			line, err := market.Readln()
			if err != nil { return err }
			if line != echo {
				return fmt.Errorf("connection check failed")
			}
			return nil
		},
		Wait: func() error {
			market.Writeln("wait")
			line, err := market.Readln()
			if err != nil { return err }
			if line != "ok wait" {
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
