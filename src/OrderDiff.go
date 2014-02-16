// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"sort"
)

func OrderDiff(a, b []Order) (aonly []Order, bonly []Order) {
	al := OrderList{a, 0, 0}.Copy()
	bl := OrderList{b, 0, 0}.Copy()

	sort.Sort(al)
	sort.Sort(bl)

	for ai := len(al.orders)-1; ai >= 0; ai-- {
		for bi := len(bl.orders)-1; bi >= 0; bi-- {
			if (bl.orders[bi].price > 0) &&
					al.orders[ai].Similar(bl.orders[bi]) {
				al.orders[ai].price = 0
				bl.orders[bi].price = 0
				break
			}
		}
	}

	for ai := range al.orders {
		if al.orders[ai].price > 0 {
			aonly = append(aonly, al.orders[ai])
		}
	}

	for bi := range bl.orders {
		if bl.orders[bi].price > 0 {
			bonly = append(bonly, bl.orders[bi])
		}
	}

	return
}
