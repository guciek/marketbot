// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"money"
)

type SortOrdersByPriority []Order

func (a SortOrdersByPriority) Len() int {
	return len(a)
}

func (a SortOrdersByPriority) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a SortOrdersByPriority) Less(i, j int) bool {
	if a[i].buy.Currency() < a[j].buy.Currency() { return true }
	if a[i].buy.Currency() > a[j].buy.Currency() { return false }
	if a[i].sell.Currency() < a[j].sell.Currency() { return true }
	if a[i].sell.Currency() > a[j].sell.Currency() { return false }
	return money.PriceLess(a[i].buy, a[i].sell, a[j].buy, a[j].sell)
}
