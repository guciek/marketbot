// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
)

func Example_AssetSimilar() {
	fmt.Println(AssetValue(100000).Similar(AssetValue(101500)))
	fmt.Println(AssetValue(100000).Similar(AssetValue(98500)))
	fmt.Println(AssetValue(100000).Similar(AssetValue(102500)))
	fmt.Println(AssetValue(100000).Similar(AssetValue(97500)))
	fmt.Println(AssetValue(1800).Similar(AssetValue(0)))
	// Output:
	// true
	// true
	// false
	// false
	// true
}

func Example_AssetSubtract() {
	var a AssetValue
	fmt.Println(a)
	a = AssetValue(10000000)
	fmt.Println(a)
	b := a.Subtract(7000000)
	fmt.Println(a, b)
	fmt.Println(b.Subtract(7000000))
	// Output:
	// @0.00
	// @100.00
	// @100.00 @30.00
	// @0.00
}
