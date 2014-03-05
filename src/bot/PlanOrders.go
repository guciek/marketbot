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

func addAmountFromGain(param string) (func(Order) money.Money, error) {
	if v, err := percentValue(param); err == nil {
		if ! v.Less(decimal.Value(1)) {
			return nil, fmt.Errorf("invalid value of \"-gain\"")
		}
		return func(o Order) money.Money {
			return o.buy.Mult(v)
		}, nil
	}

	if pr, err := money.ParsePrice(param); err == nil {
		return func(o Order) money.Money {
			if (o.buy.Currency() == pr.Currency1()) &&
					(o.sell.Currency() == pr.Currency2()) {
				return o.sell.MultPrice(pr, 5)
			}
			if (o.sell.Currency() == pr.Currency1()) &&
					(o.buy.Currency() == pr.Currency2()) {
				m := o.buy.MultPrice(pr, 12)
				if ! m.LessNotSimilar(o.sell) {
					return money.Money {}
				}
				trans_pr := o.buy.DivPrice(o.sell.Sub(m))
				return o.sell.MultPrice(trans_pr, 5).Sub(o.buy)
			}
			panic("minimum gain "+pr.String()+
				" is not applicable to order: "+o.String())
		}, nil
	}

	return nil, fmt.Errorf("invalid value of \"-gain\"")
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

	if len(params["place"]) > 1 {
		return OrderPlanner {}, fmt.Errorf("invalid value of \"-place\"")
	}
	var place int64 = 3
	for _, p := range params["place"] {
		var err error
		place, err = strconv.ParseInt(p, 10, 64)
		if err != nil { return OrderPlanner {}, err }
		if place < 1 {
			return OrderPlanner {}, fmt.Errorf("invalid value of \"-place\"")
		}
	}
	delete(params, "place")

	if len(params["fee"]) > 1 {
		return OrderPlanner {}, fmt.Errorf("invalid value of \"-fee\"")
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

	var gain_adders []func(o Order) money.Money
	for _, p := range params["mingain"] {
		for _, s := range strings.Split(p, ",") {
			adder, err := addAmountFromGain(s)
			if err != nil { return OrderPlanner {}, err }
			gain_adders = append(gain_adders, adder)
		}
	}
	delete(params, "mingain")

	for p, _ := range params {
		return OrderPlanner {}, fmt.Errorf("unrecognized parameter %q", "-"+p)
	}

	generate := func(b1, b2 money.Money,
			next func(a1, a2 money.Money) Order) (ret []Order) {
		for k := int64(0); k < place; k++ {
			o := next(b1, b2)
			if o.buy.IsNull() { break }
			if o.sell.IsNull() { break }
			if ! o.sell.LessNotSimilar(b2) { break }
			{
				max_add := o.buy.Zero()
				for _, adder := range gain_adders {
					add := adder(o)
					if add.IsNull() { return }
					if max_add.LessNotEqual(add) { max_add = add }
				}
				o.buy = o.buy.Add(max_add)
			}
			b2 = b2.Sub(o.sell)
			b1 = b1.Add(o.buy)
			o.buy = o.buy.Mult(fee_mult)
			ret = append(ret, o)
			if len(ret) >= 2 {
				if ! SortOrdersByPriority(ret).Less(len(ret)-2, len(ret)-1) {
					panic("assertion failed")
				}
			}
		}
		return
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
				for _, m2 := range sum {
					if m1.Currency() == m2.Currency() { continue }
					for _, planner := range planners {
						ret = append(ret, generate(m1, m2, planner)...)
					}
				}
			}
			return ret
		},
	}, nil
}
