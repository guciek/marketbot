// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"money"
	"strings"
)

func PlanOrders_Balance(params_str string) (
		ret func(am1, am2 money.Money) Order, err error) {
	params := strings.Split(params_str, ",")
	if len(params) != 2 {
		err = fmt.Errorf("invalid value of \"-balance\"")
		return
	}

	var pair money.Price
	pair, err = money.ParsePrice("1 "+params[0])
	if err != nil {
		err = fmt.Errorf("invalid value of \"-balance\"")
		return
	}
	OrderPrintPrice(pair)

	var size money.Money
	size, err = money.ParseMoney(params[1])
	if err != nil { return }

	if (size.Currency() != pair.Currency1()) &&
			(size.Currency() != pair.Currency2()) {
		err = fmt.Errorf("order size should be in %q or %q",
			pair.Currency1(), pair.Currency2())
		return
	}

	match := func(a, b money.Money) bool {
		if a.Currency() != pair.Currency1() { return false }
		if b.Currency() != pair.Currency2() { return false }
		return true
	}

	return func(am1, am2 money.Money) Order {
		if (! match(am1, am2)) && (! match(am2, am1)) {
			return Order {}
		}
		size2 := size.Add(size)
		if size.Currency() == am1.Currency() {
			part := size.Div(am1.Add(size2), 10)
			return Order {buy: size, sell: am2.Mult(part).Round(6)}
		} else {
			if ! size2.LessNotSimilar(am2) { return Order {} }
			part := size.Div(am2.Sub(size2), 10)
			return Order {buy: am1.Mult(part).Round(6), sell: size}
		}
	}, nil
}
