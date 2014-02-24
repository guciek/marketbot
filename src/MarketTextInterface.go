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
	cancel := make(map[Order]string)
	GetTotalBalance := func() (am_for_s, am_for_b MoneyValue, err error) {
		market.Writeln("totalbalance")
		line := market.Readln()
		words := append(strings.Split(line, " "), "", "")
		if words[0] != "totalbalance" {
			err = fmt.Errorf("get balance: %q", line)
			return
		}
		am1, err1 := strconv.ParseUint(words[1], 10, 64)
		am2, err2 := strconv.ParseUint(words[2], 10, 64)
		if (err1 != nil) || (err2 != nil) {
			err = fmt.Errorf("get balance: %q", line)
		} else {
			am_for_s = MoneyValue(am1)
			am_for_b = MoneyValue(am2)
		}
		return
	}
	parseOrders := func(lines []string) ([]Order, error) {
		var ret []Order
		for _, line := range lines {
			w := strings.Split(line, " ")
			if (len(w) == 4) && ((w[0] == "sell") || (w[0] == "buy") &&
					(w[2] == "for")) {
				asset, err := strconv.ParseUint(w[1], 10, 64)
				if err != nil { return nil, err }
				money, err := strconv.ParseUint(w[3], 10, 64)
				if err != nil { return nil, err }
				var add Order
				add.price = MoneyValue(asset).PriceIfBuyFor(MoneyValue(money))
				if w[0] == "buy" {
					add.cost = MoneyValue(money)
					add.t = BUY
				} else {
					add.cost = MoneyValue(asset)
					add.t = SELL
				}
				cancel[add] = line
				ret = append(ret, add)
			} else {
				return nil, fmt.Errorf("order list: %q", line)
			}
		}
		return ret, nil
	}
	GetOrders := func() ([]Order, error) {
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
		GetOrderList: func() (ret OrderList, err error) {
			ret.am_for_s, ret.am_for_b, err = GetTotalBalance()
			if err != nil { return }
			ret.orders, err = GetOrders()
			if err != nil { return }
			for _, o := range ret.orders {
				if o.t == SELL {
					ret.am_for_s = ret.am_for_s.Subtract(o.cost)
				} else {
					ret.am_for_b = ret.am_for_b.Subtract(o.cost)
				}
			}
			return
		},
		NewOrder: func(o Order) error {
			var cmd string
			if o.t == SELL {
				cmd = fmt.Sprintf("sell %d for %d", o.cost,
					o.cost.AfterSell(o.price))
			} else {
				cmd = fmt.Sprintf("buy %d for %d",
					o.cost.AfterBuy(o.price), o.cost)
			}
			market.Writeln(cmd)
			if line := market.Readln(); line != ("ok "+cmd) {
				return fmt.Errorf("new order: %q", line)
			}
			return nil
		},
		CancelOrder: func(o Order) error {
			cmd := cancel[o]
			if len(cmd) < 1 {
				return fmt.Errorf("cancel order: not found")
			}
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
