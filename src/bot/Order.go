// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"money"
)

type Order struct {
	buy, sell money.Money
	id string
}

var order__prCurrencyPairs map[string]bool

func OrderPrintPrice(p money.Price) {
	if order__prCurrencyPairs == nil {
		order__prCurrencyPairs = make(map[string]bool)
	}
	order__prCurrencyPairs[p.Currency1()+" "+p.Currency2()] = true
}

func (o Order) String() string {
	var pr money.Price
	if order__prCurrencyPairs[o.buy.Currency()+" "+o.sell.Currency()] {
		pr = o.buy.DivPrice(o.sell)
	} else if order__prCurrencyPairs[o.sell.Currency()+" "+o.buy.Currency()] {
		pr = o.sell.DivPrice(o.buy)
	}
	if pr.IsNull() {
		return fmt.Sprintf("buy %v for %v", o.buy, o.sell)
	}
	return fmt.Sprintf("buy %v for %v (%v)", o.buy, o.sell,
		pr.StringPrecision(6))
}

func (o1 Order) Similar(o2 Order) bool {
	return o1.buy.Similar(o2.buy) &&
		o1.sell.Similar(o2.sell)
}
