// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

type SortOrdersByPriority []Order

func (a SortOrdersByPriority) Len() int {
	return len(a)
}

func (a SortOrdersByPriority) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a SortOrdersByPriority) Less(i, j int) bool {
	if a[i].t < a[j].t { return true }
	if a[i].t > a[j].t { return false }
	v := int64(a[i].money)*int64(a[j].asset) -
		int64(a[j].money)*int64(a[i].asset)
	if v < 0 { return a[i].t != BUY }
	if v > 0 { return a[i].t == BUY }
	return a[i].asset < a[j].asset
}
