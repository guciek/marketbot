// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package money

import (
	"fmt"
	"strings"
)

type Price struct {
	am1, am2 Money
}

func ParsePrice(s string) (Price, error) {
	s = strings.TrimSpace(s)
	if len(s) < 1 { return Price {}, fmt.Errorf("empty string") }
	parts := strings.Split(s, "/")
	if len(parts) != 2 { return Price {}, fmt.Errorf("invalid price") }
	a1, err1 := ParseMoney(parts[0])
	a2, err2 := ParseMoney("1 "+parts[1])
	if (err1 != nil) || (err2 != nil) {
		return Price {}, fmt.Errorf("invalid price")
	}
	if a1.currency == a2.currency {
		return Price {}, fmt.Errorf("currencies must be different")
	}
	return Price {am1: a1, am2: a2}, nil
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

func (p Price) IsNull() bool {
	return (p.am1.currency == "") || (p.am2.currency == "")
}

func (p Price) Currency1() string {
	return p.am1.Currency()
}

func (p Price) Currency2() string {
	return p.am2.Currency()
}

func (p Price) String() string {
	return p.StringPrecision(6)
}

func (p Price) StringPrecision(precision uint32) string {
	if p.am1.currency == "" {
		panic("printing null price value")
	}
	return p.am1.v.Div(p.am2.v, precision).StringPrecision(2)+
		" "+p.am1.currency+"/"+p.am2.currency
}
