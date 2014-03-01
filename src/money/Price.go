// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package money

import (
	"decimal"
)

type Price struct {
	pr decimal.Decimal
	currency1, currency2 string
}

func (a Price) String() string {
	if a.currency1 == "" {
		panic("printing null price value")
	}
	return a.pr.StringPrecision(2)+" "+a.currency1+"/"+a.currency2
}
