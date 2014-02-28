// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"strconv"
)

type OrderType uint8

const (
	BUY = OrderType(0)
	SELL = OrderType(1)
)

func (t OrderType) String() string {
	if uint8(t) > 0 { return "sell" }
	return "buy"
}

type Order struct {
	asset AssetValue
	money MoneyValue
	t OrderType
}

func (o Order) PriceString() string {
	price := float64(o.money)/float64(o.asset)
	digits := 5
	{
		p := price
		for p < 1.0 {
			p *= 10.0
			digits++
		}
	}
	return strconv.FormatFloat(price, 'f', digits, 64)
}

func (o Order) String() string {
	return fmt.Sprintf("(%v %v for %v, price %s)",
		o.t, o.asset, o.money, o.PriceString())
}

func (o1 Order) Similar(o2 Order) bool {
	return (o1.t == o2.t) &&
		o1.asset.Similar(o2.asset) &&
		o1.money.Similar(o2.money)
}
