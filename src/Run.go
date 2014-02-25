// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"sort"
)

func Run_Update(market MarketController,
		planner func(AssetValue, MoneyValue) []Order,
		log, log_status func(string), interrupted func() bool) bool {
	a, m, err := market.GetTotalBalance()
	if err != nil { return false }

	current, err := market.GetOrders()
	if err != nil { return false }

	target := planner(a, m)
	cancel, place := DiffOrders(current, target)

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

	sort.Sort(SortOrdersByPriority(target))
	firstprice := func(t OrderType) string {
		for _, o := range target {
			if o.t == t { return " "+o.PriceString() }
		}
		return ""
	}
	var target_buy, target_sell = 0, 0
	for _, o := range target {
		if o.t == BUY { target_buy++ } else { target_sell++ }
	}
	log_status(fmt.Sprintf("Placed: %d buy orders, "+
		"%d sell orders, spread%s ~%s", target_buy, target_sell,
		firstprice(BUY), firstprice(SELL)))

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
			if lastUpdate+1200 <= currentTime {
				lastUpdate = currentTime
				log_status("Warning: could not update for 20 minutes")
			}
		}
	}
}
