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

func (o Order) String() string {
	if o.buy.Currency() > o.sell.Currency() {
		return fmt.Sprintf("buy %v for %v (%v)",
			o.buy, o.sell, o.buy.DivPrice(o.sell, 6))
	} else if o.buy.Currency() < o.sell.Currency() {
		return fmt.Sprintf("buy %v for %v (%v)",
			o.buy, o.sell, o.sell.DivPrice(o.buy, 6))
	} else {
		return fmt.Sprintf("buy %v for %v", o.buy, o.sell)
	}
}

func (o1 Order) Similar(o2 Order) bool {
	return o1.buy.Similar(o2.buy) &&
		o1.sell.Similar(o2.sell)
}
