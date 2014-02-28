// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"sort"
	"strings"
)

func Run_Update(market MarketController, planner OrderPlanner,
		current []Order, a AssetValue, m MoneyValue,
		log, log_status func(string),
		interrupted func() bool) bool {
	target := planner.TargetOrders(a, m)
	cancel, place := DiffOrders(current, target)

	{
		prices := func(t OrderType) (ret []string) {
			for _, o := range current {
				if o.t == t { ret = append(ret, o.PriceString()) }
			}
			ret = append(ret, "X")
			return
		}
		sort.Sort(SortOrdersByPriority(current))
		b, s := prices(BUY), prices(SELL)
		if len(s) > 2 { s = s[0:2] }
		if len(b) > 2 { b = b[0:2] }
		if len(b) > 1 { b[0], b[1] = b[1], b[0] }
		log_status(fmt.Sprintf(
			"Orders: %d/%d (+%d), spread: [%v] %v ~ %v [%v]",
			len(target)-len(place), len(target), len(cancel),
			m, strings.Join(b, " "), strings.Join(s, " "), a,
		))
	}

	if len(cancel) > 0 {
		for _, o := range cancel {
			if interrupted() { return false }
			log("Cancelling order "+o.String())
			market.CancelOrder(o)
		}
		return false
	}

	if len(place) > 0 {
		for _, o := range place {
			if interrupted() { return false }
			log("Placing order "+o.String())
			market.NewOrder(o)
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

	var prev_asset AssetValue = 0
	var prev_money MoneyValue = 0
	var prev_orders []Order = make([]Order, 0)
	orders_updated := false
	for ! interrupted() {
		ts, err := market.GetTime()
		if err != nil { continue }

		if ts < currentTime {
			panic("time going back")
		}
		currentTime = ts

		info_updated := false

		if interrupted() { return }
		if orders_updated {
			a, m, err := market.GetTotalBalance()
			if err == nil {
				a_, m_ := prev_asset, prev_money
				prev_asset, prev_money = a, m
				if a.Similar(a_) && (m.Similar(m_)) {
					orders_updated = false
					info_updated = true
				}
			}
		} else {
			c, err := market.GetOrders()
			if err == nil {
				c_ := prev_orders
				prev_orders = c
				a, b := DiffOrders(c, c_)
				if (len(a) == 0) && (len(b) == 0) {
					orders_updated = true
				}
			}
		}

		if info_updated {
			Run_Update(market, planner, prev_orders, prev_asset, prev_money,
				log, log_status, interrupted)
		}

		if interrupted() { return }
		if err := market.Wait(); err != nil {
			log("Error: "+err.Error())
			return
		}
	}
}
