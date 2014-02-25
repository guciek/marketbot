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
	size_money := MoneyValue(params["natural_money"]*100000)

	target := float64(params["target"])/1000.0
	if target < 0.002 {
		err = fmt.Errorf("incorrect value of '-target'")
		return
	}

	gain := AssetValue(params["gain"])
	if gain == 0 { gain = 1012 }
	if (gain < 1000) || (gain > 2000) {
		err = fmt.Errorf("incorrect value of '-gain'")
		return
	}

	ret = func(asset AssetValue, money MoneyValue, t OrderType) Order {
		if t == BUY {
			if money < size_money { return Order {} }
			s := float64(size_money)
			delta := s*s + 4.0*float64(asset)*(float64(money)-s)*target
			newprice := (math.Sqrt(delta)-s)/float64(2*asset)
			if newprice <= 0 { return Order {} }
			buy := AssetValue(s/newprice)
			return Order {buy*gain/1000, size_money, BUY}
		} else {
			s := float64(size_money)
			delta := s*s + 4.0*float64(asset)*(float64(money)+s)*target
			newprice := (math.Sqrt(delta)+s)/float64(2*asset)
			sell := AssetValue(s/newprice)
			if sell < 100000 { return Order {} }
			return Order {sell*1000/gain, size_money, SELL}
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
	size_money := MoneyValue(params["balance_money"]*100000)

	gain := AssetValue(params["gain"])
	if gain == 0 { gain = 1012 }
	if (gain < 1000) || (gain > 2000) {
		err = fmt.Errorf("incorrect value of '-gain'")
		return
	}

	ret = func(asset AssetValue, money MoneyValue, t OrderType) Order {
		if t == BUY {
			if money < size_money*3 { return Order {} }
			s := float64(size_money)
			buy := AssetValue((float64(asset)*s)/(float64(money)-s*2.0))
			return Order {buy*gain/1000, size_money, BUY}
		} else {
			s := float64(size_money)
			sell := AssetValue((float64(asset)*s)/(float64(money)+s*2.0))
			if sell < 100000 { return Order {} }
			return Order {sell*1000/gain, size_money, SELL}
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
		mask_money = MoneyValue(params["mask_money"]*100000)
	}

	mask_asset := AssetValue(0)
	if params["mask_asset"] > 0 {
		mask_asset = AssetValue(params["mask_asset"]*100000)
	}

	numplace := int(params["place"])
	if numplace < 1 { numplace = 3 }

	fee := AssetValue(params["fee"])
	if fee == 0 { fee = 1000 }
	if (fee < 900) || (fee > 1000) {
		return OrderPlanner {}, fmt.Errorf("incorrect value of '-fee'")
	}

	targetOrders := func(asset AssetValue, money MoneyValue) []Order {
		ret := make([]Order, 0, 10)
		{
			a, m := asset, money
			for i := 0; i < numplace; i++ {
				o := next(a.Subtract(mask_asset), m.Subtract(mask_money), BUY)
				if o.asset < 1 { break }
				if m < o.money+mask_money { break }
				ret = append(ret, Order {o.asset*1000/AssetValue(fee),
						o.money, BUY})
				m -= o.money
				a += o.asset
			}
		}
		{
			a, m := asset, money
			for i := 0; i < numplace; i++ {
				o := next(a.Subtract(mask_asset), m.Subtract(mask_money), SELL)
				if o.asset < 1 { break }
				if a < o.asset+mask_asset { break }
				ret = append(ret, Order {o.asset,
						o.money*1000/MoneyValue(fee), SELL})
				m += o.money
				a -= o.asset
			}
		}
		return ret
	}

	if (params["test_money"] > 0) || (params["test_asset"] > 0) {
		m := MoneyValue(params["test_money"]*100000)
		a := AssetValue(params["test_asset"]*100000)
		{
			pr := float64(1.0)
			if params["target"] > 0 {
				pr = float64(params["target"])/1000.0
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
