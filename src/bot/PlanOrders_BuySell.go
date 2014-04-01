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

func PlanOrders_Buy(params_str string) (
		ret func(am1, am2 money.Money) Order, err error) {
	params := strings.Split(params_str, ",")
	if len(params) != 2 {
		err = fmt.Errorf("invalid value of \"-buy\"")
		return
	}

	var target money.Price
	target, err = money.ParsePrice(params[0])
	if err != nil { return }
	OrderPrintPrice(target)

	var size money.Money
	size, err = money.ParseMoney(params[1])
	if err != nil { return }

	if (size.Currency() != target.Currency1()) &&
			(size.Currency() != target.Currency2()) {
		err = fmt.Errorf("order size should be in %q or %q",
			target.Currency1(), target.Currency2())
		return
	}

	return func(am1, am2 money.Money) Order {
		if am1.Currency() != target.Currency2() { return Order {} }
		if am2.Currency() != target.Currency1() { return Order {} }
		if size.Currency() == target.Currency2() {
			return Order {buy: size, sell: size.MultPrice(target)}
		} else {
			return Order {buy: size.MultPricePrecision(target.Inverse(), 6),
				sell: size}
		}
	}, nil
}

func PlanOrders_Sell(params_str string) (
		ret func(am1, am2 money.Money) Order, err error) {
	params := strings.Split(params_str, ",")
	if len(params) != 2 {
		err = fmt.Errorf("invalid value of \"-sell\"")
		return
	}

	var target money.Price
	target, err = money.ParsePrice(params[0])
	if err != nil { return }
	OrderPrintPrice(target)

	var size money.Money
	size, err = money.ParseMoney(params[1])
	if err != nil { return }

	if (size.Currency() != target.Currency1()) &&
			(size.Currency() != target.Currency2()) {
		err = fmt.Errorf("order size should be in %q or %q",
			target.Currency1(), target.Currency2())
		return
	}

	return func(am1, am2 money.Money) Order {
		if am1.Currency() != target.Currency1() { return Order {} }
		if am2.Currency() != target.Currency2() { return Order {} }
		if size.Currency() == target.Currency1() {
			return Order {buy: size,
					sell: size.MultPricePrecision(target.Inverse(), 6)}
		} else {
			return Order {buy: size.MultPrice(target), sell: size}
		}
	}, nil
}
