// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
)

type PriceValue uint64

func (a PriceValue) Similar(b PriceValue) bool {
	if a+2000 < b { return false }
	if b+2000 < a { return false }
	if a*1001 < b*1000 { return false }
	if b*1001 < a*1000 { return false }
	return true
}

func (p PriceValue) String() string {
	return fmt.Sprintf("%.3f", float64(p)*0.00000001)
}
