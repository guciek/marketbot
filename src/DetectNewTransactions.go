// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
	"sort"
)

func DetectNewTransactions(prev, cur OrderList) ([]Transaction, error) {
	var ret []Transaction
	onTransaction := func(t Transaction) {
		if t.gain >= 20 {
			t.gain -= 10
		} else if t.gain > 10 {
			t.gain = 10
		}
		ret = append(ret, t)
	}

	prev = prev.Copy()
	cur = cur.Copy()
	sort.Sort(prev)
	sort.Sort(cur)

	var maxSellPrice, minBuyPrice PriceValue = 0, 0

	for p := prev.Len()-1; p >= 0; p-- {
		for c := cur.Len()-1; c >= 0; c-- {
			if (cur.orders[c].price > 0) &&
					prev.orders[p].Similar(cur.orders[c]) {
				if prev.orders[p].t == SELL {
					if (maxSellPrice <= 0) ||
							(prev.orders[p].price < maxSellPrice) {
						maxSellPrice = prev.orders[p].price
					}
				} else {
					if (minBuyPrice <= 0) ||
							(prev.orders[p].price > minBuyPrice) {
						minBuyPrice = prev.orders[p].price
					}
				}
				cur.orders[c].price = 0
				prev.orders[p].price = 0
			}
		}
	}

	for p := prev.Len()-1; p >= 0; p-- {
		for c := cur.Len()-1; c >= 0; c-- {
			if (cur.orders[c].price > 0) &&
					prev.orders[p].price.Similar(cur.orders[c].price) &&
					(prev.orders[p].t == cur.orders[c].t) &&
					(prev.orders[p].cost > cur.orders[c].cost) {
				prev.orders[p].cost -= cur.orders[c].cost
				cur.orders[c].price = 0
			}
		}
	}
	for c := cur.Len()-1; c >= 0; c-- {
		if cur.orders[c].price > 0 {
			if cur.orders[c].t == SELL {
				cur.am_for_s += cur.orders[c].cost
			} else {
				cur.am_for_b += cur.orders[c].cost
			}
			cur.orders[c].price = 0
		}
	}

	err := func() error {
		var expect_gain, real_cost MoneyValue = 0, 0
		for _, p := range prev.orders {
			if (p.t != SELL) || (p.price <= 0) { continue }
			real_cost += p.cost
			expect_gain += p.cost.AfterSell(p.price)
		}
		if expect_gain.Similar(0) {
			if cur.am_for_b.Less(prev.am_for_b) {
				return fmt.Errorf("amount for buying decreased from %v to %v",
					prev.am_for_b, cur.am_for_b)
			}
			return nil
		}
		real_gain := cur.am_for_b.Subtract(prev.am_for_b)
		if real_gain*100 < expect_gain*99 {
			return fmt.Errorf("gain +%v from sell orders, expected +%v",
				real_gain, expect_gain)
		}
		if maxSellPrice > 0 {
			max_real_gain := real_cost.AfterSell(maxSellPrice)
			if (! real_gain.Similar(max_real_gain)) &&
					(real_gain > max_real_gain) {
				return fmt.Errorf("gain +%v from sell orders, "+
					"expected +%v but no more than +%v",
					real_gain, expect_gain, max_real_gain)
			}
		}
		if real_gain.Similar(expect_gain) || (real_gain < expect_gain) {
			for _, p := range prev.orders {
				if (p.t != SELL) || (p.price <= 0) { continue }
				onTransaction(Transaction {p, (real_gain*
					p.cost.AfterSell(p.price))/expect_gain})
			}
		} else {
			for _, p := range prev.orders {
				if (p.t != SELL) || (p.price <= 0) { continue }
				onTransaction(Transaction {p, (real_gain*p.cost)/real_cost})
			}
		}
		return nil
	}()
	if err != nil { return nil, err }

	err = func() error {
		var expect_gain, real_cost MoneyValue = 0, 0
		for _, p := range prev.orders {
			if (p.t != BUY) || (p.price <= 0) { continue }
			real_cost += p.cost
			expect_gain += p.cost.AfterBuy(p.price)
		}
		if expect_gain.Similar(0) {
			if cur.am_for_s.Less(prev.am_for_s) {
				return fmt.Errorf("amount for selling decreased from %v to %v",
					prev.am_for_s, cur.am_for_s)
			}
			return nil
		}
		real_gain := cur.am_for_s.Subtract(prev.am_for_s)
		if real_gain*100 < expect_gain*99 {
			return fmt.Errorf("gain +%v from buy orders, expected +%v",
				real_gain, expect_gain)
		}
		if minBuyPrice > 0 {
			max_real_gain := real_cost.AfterBuy(minBuyPrice)
			if (! real_gain.Similar(max_real_gain)) &&
					(real_gain > max_real_gain) {
				return fmt.Errorf("gain +%v from buy orders, "+
					"expected +%v but no more than +%v",
					real_gain, expect_gain, max_real_gain)
			}
		}
		if real_gain.Similar(expect_gain) || (real_gain < expect_gain) {
			for _, p := range prev.orders {
				if (p.t != BUY) || (p.price <= 0) { continue }
				onTransaction(Transaction {p, (real_gain*
					p.cost.AfterBuy(p.price))/expect_gain})
			}
		} else {
			for _, p := range prev.orders {
				if (p.t != BUY) || (p.price <= 0) { continue }
				onTransaction(Transaction {p, (real_gain*p.cost)/real_cost})
			}
		}
		return nil
	}()
	if err != nil { return nil, err }

	return ret, nil
}
