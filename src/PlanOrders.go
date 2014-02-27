// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"math"
)

func PlanOrders_Next_Natural(params map[string]int64) (ret func(AssetValue,
		MoneyValue, OrderType) Order, err error) {
	if params["natural_money"] < 1 {
		err = fmt.Errorf("incorrect value of '-natural'")
		return
	}
	size_money := MoneyValue(params["natural_money"])

	target := float64(params["target"])*0.00001
	if target < 0.0001 {
		err = fmt.Errorf("incorrect value of '-target'")
		return
	}

	spread := float64(1.0)
	if params["spread"] != 0 {
		s := float64(params["spread"])*0.00001
		if (s < 0.0) || (s > 0.6) {
			err = fmt.Errorf("incorrect value of '-spread'")
			return
		}
		spread = math.Sqrt(1.0+s)
	}

	ret = func(asset AssetValue, money MoneyValue, t OrderType) Order {
		if t == BUY {
			if money < size_money { return Order {} }
			s := float64(size_money)
			delta := s*s + 4.0*float64(asset)*(float64(money)-s)*target
			newprice := (math.Sqrt(delta)-s)/float64(2*asset)
			if newprice <= 0 { return Order {} }
			buy := (s/newprice)*spread
			if buy < 90.0 { return Order {} }
			return Order {AssetValue(buy), size_money, BUY}
		} else {
			s := float64(size_money)
			delta := s*s + 4.0*float64(asset)*(float64(money)+s)*target
			newprice := (math.Sqrt(delta)+s)/float64(2*asset)
			sell := (s/newprice)/spread
			if sell < 90.0 { return Order {} }
			return Order {AssetValue(sell), size_money, SELL}
		}
	}
	return
}

func PlanOrders_Next_Balance(params map[string]int64) (ret func(AssetValue,
		MoneyValue, OrderType) Order, err error) {
	if params["balance_money"] < 1 {
		err = fmt.Errorf("incorrect value of '-balance'")
		return
	}
	size_money := MoneyValue(params["balance_money"])

	spread := float64(1.0)
	if params["spread"] != 0 {
		s := float64(params["spread"])*0.00001
		if (s < 0.0) || (s > 0.6) {
			err = fmt.Errorf("incorrect value of '-spread'")
			return
		}
		spread = math.Sqrt(1.0+s)
	}

	ret = func(asset AssetValue, money MoneyValue, t OrderType) Order {
		if t == BUY {
			s := float64(size_money)
			if float64(money) < s*2.01 { return Order {} }
			buy := ((float64(asset)*s)/(float64(money)-s*2.0))*spread
			if buy < 90.0 { return Order {} }
			return Order {AssetValue(buy), size_money, BUY}
		} else {
			s := float64(size_money)
			sell := ((float64(asset)*s)/(float64(money)+s*2.0))/spread
			if sell < 90.0 { return Order {} }
			return Order {AssetValue(sell), size_money, SELL}
		}
	}
	return
}

func PlanOrders(params map[string]int64) (OrderPlanner, error) {
	var next func(asset AssetValue, money MoneyValue, t OrderType) Order
	if (params["natural_money"] > 0) || (params["natural_asset"] > 0) {
		var err error
		next, err = PlanOrders_Next_Natural(params)
		if err != nil { return OrderPlanner {}, err }
	} else if (params["balance_money"] > 0) || (params["balance_asset"] > 0) {
		var err error
		next, err = PlanOrders_Next_Balance(params)
		if err != nil { return OrderPlanner {}, err }
	} else {
		return OrderPlanner {}, fmt.Errorf("planning type not specified")
	}

	mask_money := MoneyValue(0)
	if params["mask_money"] > 0 {
		mask_money = MoneyValue(params["mask_money"])
	}

	mask_asset := AssetValue(0)
	if params["mask_asset"] > 0 {
		mask_asset = AssetValue(params["mask_asset"])
	}

	numplace := int(params["place"]/100000)
	if numplace < 1 { numplace = 3 }

	fee := 100000 - AssetValue(params["fee"])
	if (fee < 90000) || (fee > 100000) {
		return OrderPlanner {}, fmt.Errorf("incorrect value of '-fee'")
	}

	targetOrders := func(asset AssetValue, money MoneyValue) []Order {
		if asset < 90 { return nil }
		if money < 90 { return nil }
		ret := make([]Order, 0, 10)
		{
			a, m := asset, money
			for i := 0; i < numplace; i++ {
				if m < 90 { break }
				o := next(a.Subtract(mask_asset), m.Subtract(mask_money), BUY)
				if o.asset < 1 { break }
				if m < o.money+mask_money { break }
				ret = append(ret, Order {o.asset*100000/AssetValue(fee),
						o.money, BUY})
				m -= o.money
				a += o.asset
			}
		}
		{
			a, m := asset, money
			for i := 0; i < numplace; i++ {
				if a < 90 { break }
				o := next(a.Subtract(mask_asset), m.Subtract(mask_money), SELL)
				if o.asset < 1 { break }
				if a < o.asset+mask_asset { break }
				ret = append(ret, Order {o.asset,
						o.money*100000/MoneyValue(fee), SELL})
				m += o.money
				a -= o.asset
			}
		}
		return ret
	}

	if (params["test_money"] > 0) || (params["test_asset"] > 0) {
		m := MoneyValue(params["test_money"])
		a := AssetValue(params["test_asset"])
		{
			pr := float64(1.0)
			if params["target"] > 0 {
				pr = float64(params["target"])*0.00001
			}
			if m < 1 {
				a /= 2
				m = MoneyValue(float64(a)*pr)
			} else if a < 1 {
				m /= 2
				a = AssetValue(float64(m)/pr)
			}
		}
		target := targetOrders(a, m)
		{
			a_, m_ := a, m
			lines := make([]string, 0, 100)
			for _, o := range target {
				if o.t != BUY { continue }
				m_ -= o.money
				a_ += o.asset*AssetValue(fee)/1000
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
			m += o.money*MoneyValue(fee)/1000
			a -= o.asset
			fmt.Println("                <-", o)
			fmt.Println(a, m)
		}
		return OrderPlanner {}, fmt.Errorf("test mode, exiting")
	}

	return OrderPlanner {targetOrders}, nil
}
