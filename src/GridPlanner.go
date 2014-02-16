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

type Grid_el struct  {
	price_b, price_s PriceValue
	am_for_b, am_for_s MoneyValue
}

func GridPlanner() OrderPlanner {
	var grid []Grid_el
	grid = append(grid, Grid_el {255000000, 262000000, 10000000, 0})
	words := func(s string) []string {
		w := strings.Split(s, " ")
		i, j := 0, 0
		for j < len(w) {
			if len(w[j]) > 0 {
				w[i] = w[j]
				i++
			}
			j++
		}
		if i > 0 { return w[0:i] }
		return nil
	}
	return OrderPlanner {
		TargetOrders: func() []Order {
			r := make([]Order, 0, 2*len(grid))
			for _, el := range grid {
				if el.price_b >= el.price_s { panic("price_b >= price_s") }
				if ! el.am_for_b.Similar(0) {
					r = append(r, Order {price: el.price_b,
							cost: el.am_for_b, t: BUY})
				}
				if ! el.am_for_s.Similar(0) {
					r = append(r, Order {price: el.price_s,
							cost: el.am_for_s, t: SELL})
				}
			}
			return r
		},
		OnNewTransaction: func(t Transaction) bool {
			if t.t == BUY {
				for i := range grid {
					if ! t.price.Similar(grid[i].price_b) { continue }
					if t.cost.Similar(grid[i].am_for_b) ||
							(t.cost < grid[i].am_for_b) {
						grid[i].am_for_b = grid[i].am_for_b.Subtract(t.cost)
						grid[i].am_for_s += t.gain
						return true
					}
				}
			} else {
				for i := range grid {
					if ! t.price.Similar(grid[i].price_s) { continue }
					if t.cost.Similar(grid[i].am_for_s) ||
							(t.cost < grid[i].am_for_s) {
						grid[i].am_for_s = grid[i].am_for_s.Subtract(t.cost)
						grid[i].am_for_b += t.gain
						return true
					}
				}
			}
			return false
		},
		LoadData: func(d string) error {
			lines := strings.Split(d, "\n")
			header := true
			newgrid := make([]Grid_el, 0, len(lines))
			for _, line := range lines {
				if line == "" { continue }
				w := words(line)
				if header {
					header = false
					if len(w) != 4 {
						return fmt.Errorf("wrong header: "+line)
					}
					if (w[0] != "buy_p") || (w[1] != "sell_p") ||
							(w[2] != "buy_am") || (w[3] != "sell_am") {
						return fmt.Errorf("wrong header: "+line)
					}
				} else {
					if len(w) != 4 {
						return fmt.Errorf("wrong data format: "+line)
					}
					v1, err1 := strconv.ParseUint(w[0], 10, 64)
					v2, err2 := strconv.ParseUint(w[1], 10, 64)
					v3, err3 := strconv.ParseUint(w[2], 10, 64)
					v4, err4 := strconv.ParseUint(w[3], 10, 64)
					if (err1 != nil) || (err2 != nil) ||
							(err3 != nil) || (err4 != nil) {
						return fmt.Errorf("wrong data format: "+line)
					}
					newgrid = append(newgrid, Grid_el {
						price_b: PriceValue(v1*100000+1000),
						price_s: PriceValue(v2*100000+99000),
						am_for_b: MoneyValue(v3),
						am_for_s: MoneyValue(v4),
					})
				}
			}
			if header {
				return fmt.Errorf("missing header")
			}
			grid = newgrid
			return nil
		},
		SaveData: func() string {
			r := fmt.Sprintf("% 10s % 10s % 10s % 10s\n",
				"buy_p", "sell_p", "buy_am", "sell_am")
			for _, el := range grid {
				r += fmt.Sprintf("% 10d % 10d % 10d % 10d\n",
					el.price_b/100000, el.price_s/100000,
					el.am_for_b, el.am_for_s)
			}
			return r
		},
	}
}
