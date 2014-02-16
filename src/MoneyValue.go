// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
)

type MoneyValue uint64

func (a MoneyValue) Similar(b MoneyValue) bool {
	if a+2000 < b { return false }
	if b+2000 < a { return false }
	return true
}

func (a MoneyValue) Less(b MoneyValue) bool {
	if a.Similar(b) { return false }
	return a < b
}

func (a MoneyValue) AfterSell(price PriceValue) MoneyValue {
	if price < 1 { panic("price < 1") }
	return MoneyValue((uint64(a) * uint64(price)) / 100000000)
}

func (a MoneyValue) AfterBuy(price PriceValue) MoneyValue {
	if price < 1 { panic("price < 1") }
	return MoneyValue((uint64(a) * 100000000) / uint64(price))
}

func (a MoneyValue) PriceIfBuyFor(b MoneyValue) PriceValue {
	if a < 1 { panic("a < 1") }
	return PriceValue((uint64(b) * 100000000) / uint64(a))
}

func (a MoneyValue) Subtract(b MoneyValue) MoneyValue {
	if a > b { return a-b }
	return 0
}

func (p MoneyValue) String() string {
	return fmt.Sprintf("%.2f", float64(p)*0.00001)
}
