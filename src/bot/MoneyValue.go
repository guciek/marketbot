// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
)

type MoneyValue uint64

func (a MoneyValue) Similar(b MoneyValue) bool {
	if a+2000 < b { return false }
	if b+2000 < a { return false }
	return true
}

func (a MoneyValue) Subtract(b MoneyValue) MoneyValue {
	if a > b { return a-b }
	return 0
}

func (p MoneyValue) String() string {
	return fmt.Sprintf("$%.2f", float64(p)*0.00001)
}
