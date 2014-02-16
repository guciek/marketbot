// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
)

func Example_OrderSimilar() {
	a := Order { PriceValue(200000000), MoneyValue(5000000), SELL }
	b := Order { PriceValue(200000000), MoneyValue(5000000), BUY }
	c := Order { PriceValue(200000000), MoneyValue(5005000), SELL }
	d := Order { PriceValue(200010000), MoneyValue(5000000), SELL }
	e := Order { PriceValue(200001000), MoneyValue(5001000), SELL }
	fmt.Println(a.Similar(b))
	fmt.Println(a.Similar(c))
	fmt.Println(a.Similar(d))
	fmt.Println(a.Similar(e))
	// Output:
	// false
	// false
	// false
	// true
}

