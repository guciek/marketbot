// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"decimal"
	"fmt"
	"money"
	"sort"
	"strconv"
	"strings"
)

func planOrdersBalance(params map[string]string) (
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

	return func(am1, am2 money.Money) Order {
		if ((am1.Currency() != pair[0]) || (am2.Currency() != pair[1])) &&
				((am1.Currency() != pair[1]) || (am2.Currency() != pair[0])) {
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

func percentValue(s string) (decimal.Decimal, error) {
	percent := false
	if s[len(s)-1] == '%' {
		percent = true
		s = s[0:len(s)-1]
	}
	v, err := decimal.ParseDecimal(s)
	if err != nil {
		return decimal.Decimal {}, err
	}
	if percent {
		m, _ := decimal.ParseDecimal("0.01")
		v = v.Mult(m)
	}
	return v, nil
}

func PlanOrders(params map[string]string) (OrderPlanner, error) {
	if params["test"] != "" {
		balance := make([]money.Money, 0, 10)
		for _, s := range strings.Split(params["test"], ",") {
			m, err := money.ParseMoney(s)
			if err != nil { return OrderPlanner {}, err }
			balance = append(balance, m)
		}
		delete(params, "test")
		planner, err := PlanOrders(params)
		if err != nil { return OrderPlanner {}, err }
		target := SortOrdersByPriority(planner.TargetOrders(balance))
		sort.Sort(target)
		for _, t := range target {
			fmt.Println(t)
		}
		return OrderPlanner {}, fmt.Errorf("test mode, exiting")
	}

	var next func(am1, am2 money.Money) Order
	if params["balance"] != "" {
		var err error
		next, err = planOrdersBalance(params)
		if err != nil { return OrderPlanner {}, err }
	} else {
		return OrderPlanner {}, fmt.Errorf("planning type not specified")
	}

	var place int64 = 3
	if params["place"] != "" {
		var err error
		place, err = strconv.ParseInt(params["place"], 10, 64)
		if err != nil { return OrderPlanner {}, err }
		if place < 1 {
			return OrderPlanner {}, fmt.Errorf("invalid value of \"-place\"")
		}
		delete(params, "place")
	}

	buy_multiplier := decimal.Value(1)
	if s := params["buy"]; s != "" {
		v, err := percentValue(s)
		if (err != nil) || (! v.Less(decimal.Value(2)) ||
				v.Less(decimal.Value(1))) {
			return OrderPlanner {}, fmt.Errorf("invalid value of \"-buy\"")
		}
		buy_multiplier = v
		delete(params, "buy")
	}

	for p, _ := range params {
		return OrderPlanner {}, fmt.Errorf("unrecognized parameter %q", "-"+p)
	}

	return OrderPlanner {
		TargetOrders: func(balance []money.Money) []Order {
			ret := make([]Order, 0, 12)
			for i := 0; i < len(balance); i++ {
				for j := 0; j < len(balance); j++ {
					if i == j { continue }
					b1, b2 := balance[i], balance[j]
					for k := int64(0); k < place; k++ {
						o := next(b1, b2)
						if o.buy.IsNull() { break }
						if o.sell.IsNull() { break }
						if ! o.sell.LessNotSimilar(b2) { break }
						if ! money.PriceLess(b1, b2, o.buy, o.sell) { break }
						b2 = b2.Sub(o.sell)
						b1 = b1.Add(o.buy)
						o.buy = o.buy.Mult(buy_multiplier)
						ret = append(ret, o)
					}
				}
			}
			return ret
		},
	}, nil
}
