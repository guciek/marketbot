// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

type Transaction struct {
	Order
	gain MoneyValue
}

func (o Transaction) String() string {
	var realpr PriceValue
	if o.t == SELL {
		realpr = o.cost.PriceIfBuyFor(o.gain)
	} else {
		realpr = o.gain.PriceIfBuyFor(o.cost)
	}
	return o.Order.String()+" -> +"+o.gain.String()+
		", real price "+realpr.String()
}

type ByCostDecreasing []Transaction

func (a ByCostDecreasing) Len() int {
	return len(a)
}
func (a ByCostDecreasing) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByCostDecreasing) Less(i, j int) bool {
	if a[i].cost < a[j].cost { return false }
	if a[i].cost > a[j].cost { return true }
	if a[i].t < a[j].t { return true }
	if a[i].t > a[j].t { return false }
	if a[i].price < a[j].price { return a[i].t != BUY }
	if a[i].price > a[j].price { return a[i].t == BUY }
	return a[i].gain < a[j].gain
}
