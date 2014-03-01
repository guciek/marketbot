// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package money

import (
	"fmt"
)

func Example_MoneyParse() {
	x, err := ParseMoney("  \n10 USD")
	fmt.Println(err == nil, x)
	y, err := ParseMoney("0.52uSd  ")
	fmt.Println(err == nil, y)
	_, err = ParseMoney("0.0017 mBTC")
	fmt.Println(err == nil)
	_, err = ParseMoney("0.0017 B C")
	fmt.Println(err == nil)
	_, err = ParseMoney("0.0.017 BTC")
	fmt.Println(err == nil)
	_, err = ParseMoney("0.0017 BT$")
	fmt.Println(err == nil)
	fmt.Println(x.Add(y))
	fmt.Println(x.Sub(y))
	// Output:
	// true 10.00 USD
	// true 0.52 USD
	// false
	// false
	// false
	// false
	// 10.52 USD
	// 9.48 USD
}

func Example_MoneySimilar() {
	x, _ := ParseMoney("10 USD")
	y, _ := ParseMoney("10.02 USD")
	z, _ := ParseMoney("10.0004 USD")
	fmt.Println(x.Similar(y), x.Similar(z))
	x, _ = ParseMoney("0.001 BTC")
	y, _ = ParseMoney("0.001002 BTC")
	z, _ = ParseMoney("0.00100004 BTC")
	fmt.Println(x.Similar(y), x.Similar(z))
	// Output:
	// false true
	// false true
}
