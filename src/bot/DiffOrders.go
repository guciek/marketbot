// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"money"
	"sort"
)

func DiffOrders(a_, b_ []Order) (aonly []Order, bonly []Order) {
	a := make([]Order, len(a_))
	b := make([]Order, len(b_))

	copy(a, a_)
	copy(b, b_)

	sort.Sort(SortOrdersByPriority(a))
	sort.Sort(SortOrdersByPriority(b))

	for ai := len(a)-1; ai >= 0; ai-- {
		for bi := len(b)-1; bi >= 0; bi-- {
			if (!b[bi].buy.IsNull()) && a[ai].Similar(b[bi]) {
				a[ai].buy = money.Money {}
				b[bi].buy = money.Money {}
				break
			}
		}
	}

	for ai := range a {
		if !a[ai].buy.IsNull() {
			aonly = append(aonly, a[ai])
		}
	}

	for bi := range b {
		if !b[bi].buy.IsNull() {
			bonly = append(bonly, b[bi])
		}
	}

	return
}
