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

func PlanOrders_Balance(params map[string]string) (
		ret func(am1, am2 money.Money) Order, err error) {
	pair := strings.Split(strings.ToUpper(params["balance"]), "/")
	if len(pair) != 2 {
		err = fmt.Errorf("invalid value of \"-balance\"")
		return
	}
	delete(params, "balance")

	if params["order"] == "" {
		err = fmt.Errorf("missing parameter \"-order\"")
		return
	}
	var size money.Money
	size, err = money.ParseMoney(params["order"])
	if err != nil { return }
	if (size.Currency() != pair[0]) && (size.Currency() != pair[1]) {
		err = fmt.Errorf("currency of \"-order\" should be %q or %q",
			pair[0], pair[1])
		return
	}
	delete(params, "order")

	match := func(a, b money.Money) bool {
		if a.Currency() != pair[0] { return false }
		if b.Currency() != pair[1] { return false }
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
