// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"money"
	"sort"
)

func balanceInfo(balance map[string]money.Money) string {
	currencies := make([]string, 0, len(balance))
	for c, _ := range balance {
		currencies = append(currencies, c)
	}
	sort.Strings(currencies)
	info := ""
	for _, c := range currencies {
		if info != "" { info += ", " }
		info += balance[c].String();
	}
	if info == "" { return "empty" }
	return info
}

func runUpdate(market MarketController, planner OrderPlanner,
		current []Order, balance map[string]money.Money,
		log, log_status func(string),
		interrupted func() bool) bool {
	var target, cancel, place []Order
	{
		b := make([]money.Money, 0, len(balance))
		for _, m := range balance {
			b = append(b, m)
		}
		target = planner.TargetOrders(b)
		cancel, place = DiffOrders(current, target)
	}
	log_status(fmt.Sprintf(
		"Confirmed orders: %d/%d (+%d), balance: %s",
		len(target)-len(place), len(target), len(cancel), balanceInfo(balance),
	))

	if len(cancel) > 0 {
		for _, o := range cancel {
			if interrupted() { return false }
			if err := market.CancelOrder(o.id); err == nil {
				log("Cancel order: "+o.String())
			} else {
				log("Warning: could not cancel order: "+o.String())
			}
		}
		return false
	}

	if len(place) > 0 {
		for _, o := range place {
			if interrupted() { return false }
			if err := market.NewOrder(o); err == nil {
				log("New order: "+o.String())
			} else {
				log("Warning: could not add order: "+o.String())
			}
		}
		return false
	}

	return true
}

func Run(market MarketController, planner OrderPlanner,
		log func(string), interrupted func() bool) {
	defer market.Exit()

	var currentTime Time
	{
		oldlog := log
		log = func(m string) {
			oldlog(fmt.Sprintf("[%v] %v", currentTime, m))
		}
	}

	log_status := func() func(string) {
		var last_m string
		var last_t Time
		return func(m string) {
			if (m == last_m) && (last_t+3590 > currentTime) { return }
			last_m = m
			last_t = currentTime
			log(m)
		}
	}()

	balance_similar := func(a, b map[string]money.Money) bool {
		for _, aa := range a {
			if !aa.Similar(b[aa.Currency()]) { return false }
		}
		for _, bb := range b {
			if !bb.Similar(a[bb.Currency()]) { return false }
		}
		return true
	}

	prev_balance := make(map[string]money.Money)
	prev_orders := make([]Order, 0)
	orders_updated := false
	for ! interrupted() {
		{
			err := market.CheckConnection()
			if err != nil {
				log("Error: "+err.Error())
				return
			}
		}
		var ts Time
		{
			var err error
			ts, err = market.GetTime()
			if err != nil { continue }
		}

		if ts < currentTime {
			panic("time going back")
		}
		currentTime = ts

		info_updated := false

		if interrupted() { return }
		if orders_updated {
			balance, err := market.GetTotalBalance()
			if err == nil {
				balance_ := prev_balance
				prev_balance = balance
				if balance_similar(balance_, balance) {
					orders_updated = false
					info_updated = true
				} else {
					log("Balance changed: "+balanceInfo(balance))
				}
			}
		} else {
			c, err := market.GetOrders()
			if err == nil {
				c_ := prev_orders
				prev_orders = c
				a, b := DiffOrders(c_, c)
				if (len(a) == 0) && (len(b) == 0) {
					orders_updated = true
				} else {
					log(fmt.Sprintf(
						"Orders changed: total %d (added %d, removed %d)",
						len(c), len(b), len(a),
					))
				}
			}
		}

		if info_updated {
			runUpdate(market, planner, prev_orders, prev_balance,
				log, log_status, interrupted)
		}

		if interrupted() { return }
		if err := market.Wait(); err != nil {
			log("Error: "+err.Error())
			return
		}
	}
}
