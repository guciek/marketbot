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

func PlanOrders(params map[string][]string) (OrderPlanner, error) {
	if params["test"] != nil {
		balance := make([]money.Money, 0, 10)
		for _, param := range params["test"] {
			for _, s := range strings.Split(param, ",") {
				m, err := money.ParseMoney(s)
				if err != nil { return OrderPlanner {}, err }
				balance = append(balance, m)
			}
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

	var planners []func(money.Money, money.Money) Order

	for _, p := range params["balance"] {
		n, err := PlanOrders_Balance(p)
		if err != nil { return OrderPlanner {}, err }
		planners = append(planners, n)
	}
	delete(params, "balance")

	for _, p := range params["natural"] {
		n, err := PlanOrders_Natural(p)
		if err != nil { return OrderPlanner {}, err }
		planners = append(planners, n)
	}
	delete(params, "natural")

	if len(planners) < 1 {
		return OrderPlanner {}, fmt.Errorf("planning type not specified")
	}

	var place = 3
	for _, p := range params["place"] {
		v, err := strconv.ParseInt(p, 10, 32)
		place = int(v)
		if (err != nil) || (place < 1) {
			return OrderPlanner {}, fmt.Errorf("invalid value of \"-place\"")
		}
	}
	delete(params, "place")

	if len(params["fee"]) > 1 {
		return OrderPlanner {}, fmt.Errorf("multiple values of \"-fee\"")
	}
	fee_mult := decimal.Value(1)
	for _, s := range params["fee"] {
		v, err := percentValue(s)
		if (err != nil) || (! v.Add(v).Less(decimal.Value(1))) {
			return OrderPlanner {}, fmt.Errorf("invalid value of \"-fee\"")
		}
		fee_mult = decimal.Value(1).Div(decimal.Value(1).Sub(v), 8)
		if fee_mult.Less(decimal.Value(1)) { panic("assertion failed") }
	}
	delete(params, "fee")

	if len(params["gain"]) > 1 {
		return OrderPlanner {}, fmt.Errorf("multiple values of \"-gain\"")
	}
	var gain decimal.Decimal
	for _, p := range params["gain"] {
		var err error
		gain, err = percentValue(p)
		if (err != nil) || (! gain.Less(decimal.Value(1))) {
			return OrderPlanner {}, fmt.Errorf("invalid value of \"-gain\"")
		}
	}
	delete(params, "gain")

	var cashout []money.Money
	for _, p := range params["cashout"] {
		for _, s := range strings.Split(p, ",") {
			m, err := money.ParseMoney(s)
			if err != nil {
				return OrderPlanner {},
					fmt.Errorf("invalid value of \"-cashout\"")
			}
			cashout = append(cashout, m)
		}
	}
	delete(params, "cashout")

	for p, _ := range params {
		return OrderPlanner {}, fmt.Errorf("unrecognized parameter %q", "-"+p)
	}

	generate := func(b1, b2 money.Money,
			next func(a1, a2 money.Money) Order) ([]Order) {
		sub_cashout := func() func(v *money.Money) {
			co_filled := make([]bool, len(cashout))
			return func(v *money.Money) {
				for i, co := range cashout {
					if co.Currency() != v.Currency() { continue }
					if co_filled[i] { continue }
					if ! co.LessNotSimilar(*v) { break }
					*v = (*v).Sub(co)
					co_filled[i] = true
				}
			}
		}()
		sub_cashout(&b2)
		var ret []Order
		for len(ret) < place+1 {
			sub_cashout(&b1)
			o := next(b1, b2)
			if o.buy.IsNull() { break }
			if o.sell.IsNull() { break }
			if ! o.sell.LessNotSimilar(b2) { break }
			o.buy = o.buy.Add(o.buy.Mult(gain))
			if len(ret) >= 1 {
				prev := &(ret[len(ret)-1])
				prev_pr := prev.buy.DivPrice(prev.sell)
				if prev_pr.Less(o.buy.DivPrice(o.sell)) {
					ret = append(ret, o)
				} else {
					a := o.sell.MultPricePrecision(prev_pr, 8).Sub(o.buy)
					o.buy = o.buy.Add(a)
					prev.buy = prev.buy.Add(o.buy)
					prev.sell = prev.sell.Add(o.sell)
				}
			} else {
				ret = append(ret, o)
			}
			b1 = b1.Add(o.buy)
			b2 = b2.Sub(o.sell)
		}
		if len(ret) >= place {
			ret = ret[0:place]
		}
		for i := range ret {
			ret[i].buy = ret[i].buy.Mult(fee_mult)
		}
		return ret
	}

	return OrderPlanner {
		TargetOrders: func(bal []money.Money) []Order {
			sum := make(map[string]money.Money)
			for _, m := range bal {
				if ! m.IsZero() {
					sum[m.Currency()] = m.Add(sum[m.Currency()])
				}
			}
			ret := make([]Order, 0, place*2)
			for _, m1 := range sum {
				added := false
				for _, m2 := range sum {
					if m1.Currency() == m2.Currency() { continue }
					for _, planner := range planners {
						a := generate(m1, m2, planner)
						if len(a) < 1 { continue }
						if added { panic("planners conflict with each other") }
						added = true
						ret = append(ret, a...)
					}
				}
			}
			return ret
		},
	}, nil
}
