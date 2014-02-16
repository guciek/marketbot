// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

type Order struct {
	price PriceValue
	cost MoneyValue
	t OrderType
}

func (o Order) String() string {
	if o.t == BUY {
		return "(buy "+o.cost.AfterBuy(o.price).String()+" for "+
			o.cost.String()+", price "+o.price.String()+")"
	} else {
		return "(sell "+o.cost.String()+" for "+
			o.cost.AfterSell(o.price).String()+", price "+o.price.String()+")"
	}
}

func (o1 Order) Similar(o2 Order) bool {
	return (o1.t == o2.t) &&
		o1.price.Similar(o2.price) &&
		o1.cost.Similar(o2.cost)
}
