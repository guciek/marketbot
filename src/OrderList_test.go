// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
)

func Example_OrderList() {
	var ol OrderList
	ol.am_for_s = MoneyValue(10000000)
	ol.am_for_b = MoneyValue(20000000)
	fmt.Printf("%q\n", ol.SaveData())
	fmt.Println(ol.AddOrder(
		Order { PriceValue(200000000), MoneyValue(5000000), SELL }))
	s := ol.SaveData()
	fmt.Println(ol.LoadData("available 500 600s\n") == nil)
	fmt.Println(ol.AddOrder(
		Order { PriceValue(150000000), MoneyValue(7000000), SELL }))
	fmt.Println(ol.AddOrder(
		Order { PriceValue(150000000), MoneyValue(4000000), SELL }))
	fmt.Printf("%q\n", ol.SaveData())
	fmt.Println(ol.LoadData(s) == nil)
	fmt.Println(ol.RemoveOrder(
		Order { PriceValue(200000000), MoneyValue(5000000), SELL }))
	ol2 := ol.Copy()
	fmt.Println(ol2.AddOrder(
		Order { PriceValue(300000000), MoneyValue(21000000), BUY }))
	fmt.Println(ol2.AddOrder(
		Order { PriceValue(300000000), MoneyValue(19000000), BUY }))
	fmt.Printf("%q\n", ol.SaveData())
	fmt.Println(ol.LoadData(ol2.SaveData()) == nil)
	fmt.Printf("%q\n", ol.SaveData())
	// Output:
	// "available 20000000 10000000\n"
	// true
	// false
	// false
	// true
	// "available 20000000 1000000\nsell 150000000 4000000\nsell 200000000 5000000\n"
	// true
	// true
	// false
	// true
	// "available 20000000 10000000\n"
	// true
	// "available 1000000 10000000\nbuy 300000000 19000000\n"
}
