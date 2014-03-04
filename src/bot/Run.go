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
		info += balance[c].Round(6).String()
	}
	if info == "" { return "zero balance" }
	return info
}

func orderInfo(ords_ []Order) string {
	if (ords_ == nil) || (len(ords_) < 1) { return "none" }
	ords := make([]Order, len(ords_))
	copy(ords, ords_)
	sort.Sort(SortOrdersByPriority(ords))
	ret := ""
	m := make(map[string]bool)
	for _, o := range ords {
		n := o.buy.Currency()+" "+o.sell.Currency()
		if _, found := m[n]; found { continue }
		m[n] = true
		add := o.ShortString()
		if add[0] == 'b' {
			if ret != "" { add = add+", " }
			ret = add+ret
		} else {
			if ret != "" { add = ", "+add }
			ret += add
		}
	}
	if len(m) < len(ords) { ret += ", ..." }
	return ret
}

func runUpdateInfo(market MarketController,
		log func(string)) func() ([]Order, map[string]money.Money, string) {
	balance_similar := func(a, b map[string]money.Money) bool {
		for _, aa := range a {
			if !aa.Similar(b[aa.Currency()]) { return false }
		}
		for _, bb := range b {
			if !bb.Similar(a[bb.Currency()]) { return false }
		}
		return true
	}
	balance2map := func(b []money.Money) map[string]money.Money {
		ret := make(map[string]money.Money)
		for _, m := range b {
			ret[m.Currency()] = m
		}
		return ret
	}
	balance_sanity := func(a, b map[string]money.Money) bool {
		for _, aa := range a {
			if b[aa.Currency()].IsNull() { return false }
		}
		for _, aa := range a {
			if aa.LessNotSimilar(b[aa.Currency()]) { return true }
		}
		for _, bb := range b {
			if a[bb.Currency()].IsNull() { return false }
			if bb.LessNotSimilar(a[bb.Currency()]) { return false }
		}
		return true
	}

	var prev_balance map[string]money.Money
	var prev_orders []Order

	next_update_balance := false
	last_orders_matched := false
	return func() ([]Order, map[string]money.Money, string) {
		if next_update_balance {
			var next_balance map[string]money.Money
			{
				a, err := market.GetTotalBalance()
				if err != nil { return nil, nil, "wait" }
				next_balance = balance2map(a)
			}
			next_update_balance = false
			if prev_balance == nil {
				prev_balance = next_balance
				log("Balance: "+balanceInfo(prev_balance))
				return nil, nil, "wait"
			}
			if ! balance_similar(prev_balance, next_balance) {
				log("Balance: "+balanceInfo(next_balance))
				if ! balance_sanity(prev_balance, next_balance) {
					log("Error: Balance changes are not sane")
					return nil, nil, "exit"
				}
				prev_balance = next_balance
				return nil, nil, "wait"
			}
			if last_orders_matched {
				return prev_orders, prev_balance, "wait"
			}
			return nil, nil, "wait"
		} else {
			c, err := market.GetOrders()
			if err != nil { return nil, nil, "wait" }
			next_update_balance = true
			last_orders_matched = false
			c_ := prev_orders
			prev_orders = c
			if c_ == nil { return nil, nil, "ok" }
			a, b := DiffOrders(c_, c)
			if (len(a) != 0) || (len(b) != 0) { return nil, nil, "ok" }
			last_orders_matched = true
			return nil, nil, "ok"
		}
	}
}

func runUpdateOrders(market MarketController, planner OrderPlanner,
		current []Order, balance map[string]money.Money,
		log func(string), interrupted func() bool) bool {
	var target, cancel, place []Order
	{
		b := make([]money.Money, 0, len(balance))
		for _, m := range balance {
			b = append(b, m)
		}
		target = planner.TargetOrders(b)
		cancel, place = DiffOrders(current, target)
	}

	if len(cancel) > 0 {
		for _, o := range cancel {
			if interrupted() { return false }
			if err := market.CancelOrder(o.id); err == nil {
				log("Cancel order: "+o.String())
			} else {
				log("Cancel order: "+o.String()+" -> error")
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
				log("New order: "+o.String()+" -> error")
			}
		}
		return false
	}

	log(fmt.Sprintf("%d orders placed: %s", len(target), orderInfo(current)))

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
		oldlog := log
		var last_m string
		var last_t Time
		newlog := func(status bool) func (string) { return func(m string) {
			if (m == last_m) && status && (last_t+3590 > currentTime) { return }
			last_m = m
			last_t = currentTime
			oldlog(m)
		}}
		log = newlog(false)
		return newlog(true)
	}()

	info_updater := runUpdateInfo(market, log_status)
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

		if interrupted() { return }
		o, b, status := info_updater()
		if (o != nil) && (b != nil) {
			runUpdateOrders(market, planner, o, b, log_status, interrupted)
		}
		if status == "ok" {
		} else if status == "wait" {
			if interrupted() { return }
			if err := market.Wait(); err != nil {
				log("Error: "+err.Error())
				return
			}
		} else {
			break
		}
	}
}
