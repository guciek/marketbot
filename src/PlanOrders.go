// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"math"
)

func PlanOrders(params map[string]int64) (OrderPlanner, error) {
	if _, found := params["natural"]; found {
		return PlanOrders_Natural(params)
	}
	return OrderPlanner {}, fmt.Errorf("planning type is not specified")
}

func PlanOrders_Natural(params map[string]int64) (OrderPlanner, error) {
	size := MoneyValue(params["for"])
	if size < 10000 { return OrderPlanner {},
		fmt.Errorf("missing parameter '-for'") }

	target := float64(params["natural"])/1000.0
	if target < 0.01 { return OrderPlanner {},
		fmt.Errorf("incorrect value for '-natural'") }

	gain := AssetValue(params["gain"])
	if (gain < 1000) || (gain > 10000) { gain = 1012 }

	mask_money := MoneyValue(0)
	if params["mask_money"] > 0 {
		mask_money = MoneyValue(params["mask_money"]*100000)
	}
	mask_asset := AssetValue(params["mask_asset"])
	if params["mask_asset"] > 0 {
		mask_money = MoneyValue(params["mask_asset"]*100000)
	}

	numplace := int(params["place"])
	if numplace < 1 { numplace = 3 }

	next_buy := func(asset AssetValue, money MoneyValue) Order {
		if money < size { return Order {} }
		s := float64(size)
		delta := s*s + 4.0*float64(asset)*(float64(money)-s)*target
		newprice := (math.Sqrt(delta)-s)/float64(2*asset)
		if newprice <= 0 { return Order {} }
		buy := AssetValue(s/newprice)
		return Order {buy*gain/1000, size, BUY}
	}

	next_sell := func(asset AssetValue, money MoneyValue) Order {
		if asset < 100000 { return Order {} }
		s := float64(size)
		delta := s*s + 4.0*float64(asset)*(float64(money)+s)*target
		newprice := (math.Sqrt(delta)+s)/float64(2*asset)
		sell := AssetValue(s/newprice)
		if sell < 100000 { return Order {} }
		return Order {sell*1000/gain, size, SELL}
	}

	targetOrders := func(asset AssetValue, money MoneyValue) []Order {
		ret := make([]Order, 0, 10)
		{
			a, m := asset, money
			for i := 0; i < numplace; i++ {
				o := next_buy(a.Subtract(mask_asset),
					m.Subtract(mask_money))
				if o.asset < 1 { break }
				if o.money > m { break }
				ret = append(ret, o)
				m -= o.money
				a += o.asset
			}
		}
		{
			a, m := asset, money
			for i := 0; i < numplace; i++ {
				o := next_sell(a.Subtract(mask_asset),
					m.Subtract(mask_money))
				if o.asset < 1 { break }
				if o.asset > a { break }
				ret = append(ret, o)
				m += o.money
				a -= o.asset
			}
		}
		return ret
	}

	if params["test"] > 0 {
		m := MoneyValue(params["test"]*50000)
		a := AssetValue(float64(m)/target)
		m += mask_money
		a += mask_asset
		target := targetOrders(a, m)
		{
			a_, m_ := a, m
			lines := make([]string, 0, 100)
			for _, o := range target {
				if o.t != BUY { continue }
				m_ -= o.money
				a_ += o.asset
				lines = append(lines, fmt.Sprint("                <-", o))
				lines = append(lines, fmt.Sprint(a_, m_))
			}
			for i := len(lines)-1; i >= 0; i-- {
				fmt.Println(lines[i])
			}
		}
		fmt.Println(a, m)
		for _, o := range target {
			if o.t != SELL { continue }
			m += o.money
			a -= o.asset
			fmt.Println("                <-", o)
			fmt.Println(a, m)
		}
		return OrderPlanner {}, fmt.Errorf("test mode, exiting")
	}

	return OrderPlanner {targetOrders}, nil
}
