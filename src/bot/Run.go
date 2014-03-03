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
		if ret != "" { ret += ", " }
		ret += o.ShortString()
	}
	if len(m) < len(ords) { ret += ", ..." }
	return ret
}

func runUpdateInfo(market MarketController,
		log func(string)) func() ([]Order, map[string]money.Money) {
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

	var prev_balance map[string]money.Money
	var prev_orders []Order

	orders_updated := false
	last_order_info := ""
	return func() ([]Order, map[string]money.Money) {
		if orders_updated {
			balance_arr, err := market.GetTotalBalance()
			if err != nil { return nil, nil }
			orders_updated = false
			balance_ := prev_balance
			prev_balance = balance2map(balance_arr)
			if balance_ == nil {
				log("Balance: "+balanceInfo(prev_balance))
				return nil, nil
			}
			if ! balance_similar(balance_, prev_balance) {
				log("Balance changed: "+balanceInfo(prev_balance))
				return nil, nil
			}
			return prev_orders, prev_balance
		}
		c, err := market.GetOrders()
		if err != nil { return nil, nil }
		c_ := prev_orders
		prev_orders = c
		if c_ == nil {
			last_order_info = orderInfo(c)
			log("Orders: "+last_order_info)
			return nil, nil
		}
		a, b := DiffOrders(c_, c)
		if (len(a) != 0) || (len(b) != 0) {
			new_info := orderInfo(c)
			if new_info != last_order_info {
				log("Orders changed: "+new_info)
				last_order_info = new_info
			}
			return nil, nil
		}
		orders_updated = true
		return nil, nil
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
				log("Cancelled order: "+o.String())
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
				log("Warning: could not place order: "+o.String())
			}
		}
		return false
	}

	log(fmt.Sprintf(
		"Status: %s, orders: %d/%d",
		balanceInfo(balance), len(target)-len(place), len(target),
	))

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
		if o, b := info_updater(); o != nil {
			runUpdateOrders(market, planner, o, b, log_status, interrupted)
		}

		if interrupted() { return }
		if err := market.Wait(); err != nil {
			log("Error: "+err.Error())
			return
		}
	}
}
