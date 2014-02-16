// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type OrderList struct {
	orders []Order
	am_for_b, am_for_s MoneyValue
}

// Sorting order: more urgent first

func (a OrderList) Len() int {
	return len(a.orders)
}

func (a OrderList) Swap(i, j int) {
	a.orders[i], a.orders[j] = a.orders[j], a.orders[i]
}

func (a OrderList) Less(i, j int) bool {
	if a.orders[i].t < a.orders[j].t { return true }
	if a.orders[i].t > a.orders[j].t { return false }
	if a.orders[i].price < a.orders[j].price { return a.orders[i].t != BUY }
	if a.orders[i].price > a.orders[j].price { return a.orders[i].t == BUY }
	return a.orders[i].cost < a.orders[j].cost
}

func (ol OrderList) Copy() (ret OrderList) {
	ret = ol
	ret.orders = make([]Order, ol.Len())
	copy(ret.orders, ol.orders)
	return
}

func (ol *OrderList) AddOrder(o Order) bool {
	if o.t == SELL {
		if ol.am_for_s < o.cost { return false }
		ol.am_for_s -= o.cost
	} else {
		if ol.am_for_b < o.cost { return false }
		ol.am_for_b -= o.cost
	}
	ol.orders = append(ol.orders, o)
	return true
}

func (ol *OrderList) RemoveOrder(o Order) bool {
	ll := len(ol.orders)
	for i := 0; i < ll; i++ {
		if o != ol.orders[i] { continue }
		ol.Swap(i, ll-1)
		ol.orders = ol.orders[0:ll-1]
		if o.t == SELL {
			ol.am_for_s += o.cost
		} else {
			ol.am_for_b += o.cost
		}
		return true
	}
	return false
}

func (ol OrderList) SaveData() string {
	ol = ol.Copy()
	sort.Sort(ol)
	data := ""
	data += fmt.Sprintf("available %d %d\n", ol.am_for_b, ol.am_for_s)
	for _, o := range ol.orders {
		data += fmt.Sprintf("%s %d %d\n", o.t.String(), o.price, o.cost)
	}
	return data
}

func (ol *OrderList) LoadData(d string) error {
	ret := make([]Order, 0, 10)
	lines := strings.Split(d, "\n")
	for (len(lines) > 0) && (lines[len(lines)-1] == "") {
		lines = lines[:len(lines)-1]
	}
	parsedAvailable := false
	var av_s, av_b MoneyValue
	for _, line := range lines {
		w := strings.Split(line, " ")
		if (len(w) == 3) && (w[0] == "available") && (! parsedAvailable) {
			parsedAvailable = true
			i, err := strconv.ParseUint(w[1], 10, 64)
			if err != nil { return err }
			av_b = MoneyValue(i)
			i, err = strconv.ParseUint(w[2], 10, 64)
			if err != nil { return err }
			av_s = MoneyValue(i)
		} else if (len(w) == 3) && ((w[0] == "sell") || (w[0] == "buy")) {
			tp := SELL
			if w[0] == "buy" { tp = BUY }
			p, err := strconv.ParseUint(w[1], 10, 64)
			if err != nil { return err }
			c, err := strconv.ParseUint(w[2], 10, 64)
			if err != nil { return err }
			ret = append(ret, Order {price: PriceValue(p),
				cost: MoneyValue(c), t: tp})
		} else {
			return fmt.Errorf("order list: %q", line)
		}
	}
	if ! parsedAvailable {
		return fmt.Errorf("order list: no available balance")
	}
	ol.orders = ret
	ol.am_for_s = av_s
	ol.am_for_b = av_b
	return nil
}
