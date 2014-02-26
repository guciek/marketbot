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

func Run_Update(market MarketController,
		planner func(AssetValue, MoneyValue) []Order,
		log, log_status func(string), interrupted func() bool) bool {
	a, m, err := market.GetTotalBalance()
	if err != nil { return false }

	if interrupted() { return false }
	current, err := market.GetOrders()
	if err != nil { return false }

	target := planner(a, m)
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
		return func(m string) {
			if (m == last_m) { return }
			last_m = m
			log(m)
		}
	}()

	var lastUpdate Time
	for ! interrupted() {
		ts, err := market.GetTime()
		if err != nil { continue }

		if ts < currentTime {
			panic("time going back")
		}
		currentTime = ts

		if Run_Update(market, planner.TargetOrders,
				log, log_status, interrupted) {
			lastUpdate = currentTime
			if interrupted() { break }
			if err := market.Wait(); err != nil {
				log("Error: "+err.Error())
				break
			}
		} else if lastUpdate > 0 {
			if lastUpdate+600 <= currentTime {
				lastUpdate = currentTime
				log("Warning: could not update for last 10 minutes")
			}
		} else {
			lastUpdate = currentTime
			log("Warning: could not update")
		}
	}
}
