// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"math"
	"money"
	"strconv"
	"strings"
)

func PlanOrders_Natural(params map[string]string) (
		ret func(am1, am2 money.Money) Order, err error) {
	var target money.Price
	target, err = money.ParsePrice(params["natural"])
	if err != nil { return }
	delete(params, "natural")
	OrderPrintPrice(target)

	if params["order"] == "" {
		err = fmt.Errorf("missing parameter \"-order\"")
		return
	}
	var size money.Money
	size, err = money.ParseMoney(params["order"])
	if err != nil { return }
	if (size.Currency() != target.Currency1()) &&
			(size.Currency() != target.Currency2()) {
		err = fmt.Errorf("currency of \"-order\" should be %q or %q",
			target.Currency1(), target.Currency2())
		return
	}
	delete(params, "order")

	match_target := func(a, b money.Money) bool {
		if a.Currency() != target.Currency1() { return false }
		if b.Currency() != target.Currency2() { return false }
		return true
	}

	d2sqrt := func (am1, am2 money.Money) money.Money {
		price2float := func(a money.Price) float64 {
			// TODO: this is bad
			s := strings.Split(a.StringPrecision(16), " ")
			v, err := strconv.ParseFloat(s[0], 64)
			if err != nil { panic(err.Error()) }
			return v
		}
		money2float := func(a money.Money) float64 {
			// TODO: this is bad
			s := strings.Split(a.String(), " ")
			v, err := strconv.ParseFloat(s[0], 64)
			if err != nil { panic(err.Error()) }
			return v
		}
		float2money := func(a float64, currency string) money.Money {
			// TODO: this is bad
			v, err := money.ParseMoney(fmt.Sprintf("%f %s", a, currency))
			if err != nil { panic(err.Error()) }
			return v
		}
		sign := float64(1.0)
		if size.Currency() == am2.Currency() {
			am1, am2 = am2, am1
			sign = -sign
		}
		t := price2float(target)
		if match_target(am2, am1) {
			t = 1.0/t
		}
		s, a1, a2 := money2float(size), money2float(am1), money2float(am2)
		ret := a1+s*sign
		if ret < 0.0 { return money.Money {} }
		ret = math.Sqrt(s*s + 4.0*a2*ret*t)
		return float2money(ret, size.Currency())
	}

	return func(am1, am2 money.Money) Order {
		if (! match_target(am1, am2)) && (! match_target(am2, am1)) {
			return Order {}
		}
		d := d2sqrt(am1, am2)
		if d.IsNull() { return Order {} }
		if size.Currency() == am1.Currency() {
			part := size.Add(size).Div(d.Add(size), 10)
			return Order {buy: size, sell: am2.Mult(part).Round(6)}
		} else {
			if ! size.LessNotSimilar(d) { return Order {} }
			part := size.Add(size).Div(d.Sub(size), 10)
			return Order {buy: am1.Mult(part).Round(6), sell: size}
		}
	}, nil
}
