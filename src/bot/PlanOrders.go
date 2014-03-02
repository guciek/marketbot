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

	var planner func(money.Money, money.Money) Order

	if params["balance"] != "" {
		n, err := PlanOrders_Balance(params)
		if err != nil { return OrderPlanner {}, err }
		planner = n
	} else if params["natural"] != "" {
		n, err := PlanOrders_Natural(params)
		if err != nil { return OrderPlanner {}, err }
		planner = n
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
	if s := params["gain"]; s != "" {
		v, err := percentValue(s)
		if (err != nil) || (! v.Less(decimal.Value(2)) ||
				v.Less(decimal.Value(1))) {
			return OrderPlanner {}, fmt.Errorf("invalid value of \"-gain\"")
		}
		buy_multiplier = v
		delete(params, "gain")
	}

	for p, _ := range params {
		return OrderPlanner {}, fmt.Errorf("unrecognized parameter %q", "-"+p)
	}

	generate := func(b1, b2 money.Money,
			next func(a1, a2 money.Money) Order) []Order {
		var ret []Order
		for k := int64(0); k < place; k++ {
			o := next(b1, b2)
			if o.buy.IsNull() { break }
			if o.sell.IsNull() { break }
			if ! o.sell.LessNotSimilar(b2) { break }
			b2 = b2.Sub(o.sell)
			b1 = b1.Add(o.buy)
			o.buy = o.buy.Mult(buy_multiplier)
			ret = append(ret, o)
		}
		return ret
	}

	return OrderPlanner {
		TargetOrders: func(bal []money.Money) []Order {
			ret := make([]Order, 0, place*2)
			for i := 0; i < len(bal); i++ {
				for j := 0; j < len(bal); j++ {
					if i == j { continue }
					ret = append(ret, generate(bal[i], bal[j], planner)...)
				}
			}
			return ret
		},
	}, nil
}
