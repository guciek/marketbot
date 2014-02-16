// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
)

func Example_MoneySimilar() {
	fmt.Println(MoneyValue(100000).Similar(MoneyValue(101500)))
	fmt.Println(MoneyValue(100000).Similar(MoneyValue(98500)))
	fmt.Println(MoneyValue(100000).Similar(MoneyValue(102500)))
	fmt.Println(MoneyValue(100000).Similar(MoneyValue(97500)))
	fmt.Println(MoneyValue(1800).Similar(MoneyValue(0)))
	// Output:
	// true
	// true
	// false
	// false
	// true
}

func Example_MoneySubtract() {
	var a MoneyValue
	fmt.Println(a)
	a = MoneyValue(10000000)
	fmt.Println(a)
	b := a.Subtract(7000000)
	fmt.Println(a, b)
	fmt.Println(b.Subtract(7000000))
	// Output:
	// 0.00
	// 100.00
	// 100.00 30.00
	// 0.00
}

func Example_BuySell() {
	b100 := MoneyValue(10000000)
	p2 := PriceValue(200000000)
	s50 := b100.AfterBuy(p2)
	fmt.Println(b100, s50, p2)
	b100_ := s50.AfterSell(p2)
	fmt.Println(b100_, s50, p2)
	p2_ := s50.PriceIfBuyFor(b100)
	fmt.Println(b100, s50, p2_)
	// Output:
	// 100.00 50.00 2.000
	// 100.00 50.00 2.000
	// 100.00 50.00 2.000
}
