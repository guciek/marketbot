// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
)

func Example_PriceSimilar() {
	fmt.Println(PriceValue(100000000).Similar(PriceValue(100001500)))
	fmt.Println(PriceValue(100000000).Similar(PriceValue(99998500)))
	fmt.Println(PriceValue(100000000).Similar(PriceValue(100002500)))
	fmt.Println(PriceValue(100000000).Similar(PriceValue(99997500)))
	fmt.Println(PriceValue(1).Similar(PriceValue(0)))
	// Output:
	// true
	// true
	// false
	// false
	// false
}
