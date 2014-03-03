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

func (o Order) ShortString() string {
	if order__prCurrencyPairs[o.buy.Currency()+" "+o.sell.Currency()] {
		return fmt.Sprintf("sell at %v",
			o.buy.DivPrice(o.sell).StringPrecision(6))
	}
	return fmt.Sprintf("buy at %v",
		o.sell.DivPrice(o.buy).StringPrecision(6))
}

func (o Order) String() string {
	if order__prCurrencyPairs[o.buy.Currency()+" "+o.sell.Currency()] {
		return fmt.Sprintf("sell %v at %v",
			o.sell, o.buy.DivPrice(o.sell).StringPrecision(6))
	}
	if order__prCurrencyPairs[o.sell.Currency()+" "+o.buy.Currency()] {
		return fmt.Sprintf("buy %v at %v",
			o.buy, o.sell.DivPrice(o.buy).StringPrecision(6))
	}
	return fmt.Sprintf("buy %v for %v", o.buy, o.sell)
}

func (o1 Order) Similar(o2 Order) bool {
	return o1.buy.Similar(o2.buy) &&
		o1.sell.Similar(o2.sell)
}
