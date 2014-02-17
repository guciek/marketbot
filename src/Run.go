// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"sort"
)

func Run(market MarketController, saver DataSaver,
		log func(string), interrupted func() bool) {
	defer market.Exit()
	saver = CachedDataSaver(saver)

	validOrder := func(o Order) bool {
		// Warning: hard-coded values
		if o.t == BUY {
			return (o.cost >= 5000500)
		} else {
			return (o.cost >= MoneyValue(5000500).AfterBuy(o.price))
		}
	}

	var prev_orders, prev_orders_alternate OrderList
	has_alternate_orders := false
	{
		a, err := saver.Read("orders")
		if (err == nil) && (len(a) > 0) {
			if err := prev_orders.LoadData(a); err != nil {
				panic("could not parse order data: "+err.Error())
			}
		} else if err != nil {
			log("Warning: "+err.Error())
		}
	}

	planner := GridPlanner()
	{
		a, err := saver.Read("grid")
		if (err == nil) && (len(a) > 0) {
			err = planner.LoadData(a)
			if err != nil {
				panic("could not parse grid data: "+err.Error())
			}
		} else if err != nil {
			panic("could not read grid data: "+err.Error())
		} else {
			panic("grid data is empty")
		}
	}

	var currentTime Time
	{
		oldlog := log
		log = func(m string) {
			oldlog(fmt.Sprintf("[%v] %v", currentTime, m))
		}
	}

	log_status := func() func(string) {
		var last_t Time
		var last_m string
		return func(m string) {
			if (m == last_m) && (last_t+3597 > currentTime) { return }
			last_m = m
			last_t = currentTime
			log(m)
		}
	}()

	update_transactions := func() bool {
		market_orders, err := market.GetOrders()
		if err != nil { return false }

		trans, err := DetectNewTransactions(prev_orders, market_orders)
		if err != nil {
			if has_alternate_orders {
				trans, err = DetectNewTransactions(
					prev_orders_alternate, market_orders)
				if err != nil {
					log_status("Could not recover from error: "+err.Error())
					return false
				} else {
					has_alternate_orders = false
					log("Recovered from error: orders match variation 2")
				}
			} else {
				log_status("Could not match orders: "+err.Error())
				return false
			}
		} else if has_alternate_orders {
			has_alternate_orders = false
			log("Recovered from error: orders match variation 1")
		}

		if len(trans) > 0 {
			sort.Sort(ByCostDecreasing(trans))
			for _, t := range trans {
				log("Transaction "+t.String())
				if ! planner.OnNewTransaction(t) {
					panic("could not match transaction: "+t.String())
				}
			}
		}

		prev_orders = market_orders
		return true
	}

	update_orders := func() bool {
		target := planner.TargetOrders()
		cancel, place := OrderDiff(prev_orders.orders, target)
		changed := false
		if len(cancel) > 0 {
			changed = true
			for _, c := range cancel {
				if market.CancelOrder(c) == nil {
					log("Cancelled order "+c.String())
					if ! prev_orders.RemoveOrder(c) { panic("?") }
				} else {
					log("Error while cancelling order "+c.String())
					prev_orders_alternate = prev_orders.Copy()
					if ! prev_orders_alternate.RemoveOrder(c) { panic("?") }
					has_alternate_orders = true
					return false
				}
			}
		} else if len(place) > 0 {
			am_for_b := prev_orders.am_for_b
			am_for_s := prev_orders.am_for_s
			for _, p := range place {
				if p.t == BUY {
					am_for_b = am_for_b.Subtract(p.cost)
					if am_for_b.Similar(0) { continue }
				} else {
					am_for_s = am_for_s.Subtract(p.cost)
					if am_for_s.Similar(0) { continue }
				}
				if ! validOrder(p) { continue }
				changed = true
				if market.NewOrder(p) == nil {
					log("Placed order "+p.String())
					if ! prev_orders.AddOrder(p) { panic("?") }
				} else {
					log("Error while placing order "+p.String())
					prev_orders_alternate = prev_orders.Copy()
					if ! prev_orders.AddOrder(p) { panic("?") }
					has_alternate_orders = true
					return false
				}
			}
		}
		if ! changed {
			var all_b, all_s MoneyValue = 0, 0
			for _, o := range target {
				if o.t == BUY { all_b += o.cost } else { all_s += o.cost }
			}
			var placed_b, placed_s = all_b, all_s
			var invalid_b, invalid_s MoneyValue = 0, 0
			for _, o := range place {
				if o.t == BUY { placed_b -= o.cost } else { placed_s -= o.cost }
				if ! validOrder(o) {
					if o.t == BUY {
						all_b -= o.cost
						invalid_b += o.cost
					} else {
						all_s -= o.cost
						invalid_s += o.cost
					}
				}
			}
			log_status(fmt.Sprintf("Placed orders: buy %v/%v [+%v], "+
				"sell %v/%v [+%v]", placed_b, all_b, invalid_b,
				placed_s, all_s, invalid_s))
		}
		return true
	}

	for ! interrupted() {
		_, ts, err := market.GetPrice()
		if err != nil { continue }

		if ts < currentTime {
			panic("time going back")
		}
		currentTime = ts

		if update_transactions() {
			pl_data := planner.SaveData()
			or_data := prev_orders.SaveData()
			upd_result := update_orders()
			if upd_result {
				pl_data = planner.SaveData()
				or_data = prev_orders.SaveData()
			}
			if err := saver.Write("grid", pl_data); err != nil {
				panic("could not write grid data: "+err.Error())
			}
			if err := saver.Write("orders", or_data); err != nil {
				panic("could not write order data: "+err.Error())
			}
			if ! upd_result { continue }
		}

		if interrupted() { break }

		if err := market.Wait(); err != nil {
			log("Error: "+err.Error())
			break
		}
	}
}
