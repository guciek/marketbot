// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package money

type Price struct {
	am1, am2 Money
}

func (p1 Price) Less(p2 Price) bool {
	if (p1.am1.currency == "") || (p2.am1.currency == "") {
		panic("comparing null price values")
	}
	if (p1.am1.currency != p2.am1.currency) ||
			(p1.am2.currency != p2.am2.currency) {
		panic("comparing prices in different currencies")
	}
	return p1.am1.v.Mult(p2.am2.v).Less(p2.am1.v.Mult(p1.am2.v))
}

func (p Price) StringPrecision(precision uint32) string {
	if p.am1.currency == "" {
		panic("printing null price value")
	}
	return p.am1.v.Div(p.am2.v, precision).StringPrecision(2)+
		" "+p.am1.currency+"/"+p.am2.currency
}
