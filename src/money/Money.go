// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package money

import (
	"bytes"
	"decimal"
	"fmt"
	"strings"
)

type Money struct {
	v decimal.Decimal
	currency string
}

func ParseMoney(s string) (Money, error) {
	s = strings.TrimSpace(s)
	if len(s) < 1 { return Money {}, fmt.Errorf("empty string") }
	var value string
	{
		i := 0
		for i < len(s) {
			if (s[i] != '.') && ((s[i] < '0') || (s[i] > '9')) { break }
			i++
		}
		value = s[0:i]
		s = strings.TrimSpace(s[i:len(s)])
	}
	var currency bytes.Buffer
	for _, c := range s {
		if (c >= 'a') && (c <= 'z') { c += 'A'; c -= 'a' }
		if (c < 'A') || (c > 'Z') {
			return Money {}, fmt.Errorf("invalid currency name")
		}
		currency.WriteByte(byte(c))
	}
	if currency.Len() < 1 {
		return Money {}, fmt.Errorf("missing currency")
	}
	if currency.Len() > 5 {
		return Money {}, fmt.Errorf("currency name too long")
	}
	v_decimal, err := decimal.ParseDecimal(value)
	if err != nil { return Money {}, err }
	return Money {v_decimal, currency.String()}, nil
}

func (a Money) Currency() string {
	if a.currency == "" {
		panic("getting currency of null money value")
	}
	return a.currency
}

func (a Money) Zero() Money {
	if a.currency == "" {
		panic("null money value")
	}
	return Money {decimal.Decimal {}, a.currency}
}

func (a Money) IsZero() bool {
	if a.currency == "" {
		panic("comparing null money value")
	}
	return a.v.IsZero()
}

func (a Money) Similar(b Money) bool {
	if (a.currency == "") && (b.currency == "") {
		panic("comparing null money value")
	}
	if a.currency != b.currency {
		return false
	}
	var zero decimal.Decimal
	if  b.v.Equals(zero) || b.v.Equals(zero) {
		return a.v.Equals(b.v)
	}
	// TODO: this is inefficient
	return a.v.Div(b.v, 3).Equals(decimal.Value(1)) &&
		b.v.Div(a.v, 3).Equals(decimal.Value(1))
}

func (a Money) LessNotSimilar(b Money) bool {
	if (a.currency == "") || (b.currency == "") {
		panic("comparing null money value")
	}
	if a.currency != b.currency {
		panic("comparing values in different currencies")
	}
	if a.Similar(b) { return false }
	return a.v.Less(b.v)
}

func (a Money) LessNotEqual(b Money) bool {
	if (a.currency == "") || (b.currency == "") {
		panic("comparing null money value")
	}
	if a.currency != b.currency {
		panic("comparing values in different currencies")
	}
	return a.v.Less(b.v)
}

func (a Money) Add(b Money) Money {
	if (a.currency == "") && (b.currency == "") {
		panic("adding null money values")
	}
	if a.currency == "" { return b }
	if b.currency == "" { return a }
	if a.currency != b.currency {
		panic("adding values in different currencies")
	}
	return Money {a.v.Add(b.v), a.currency}
}

func (a Money) Sub(b Money) Money {
	if (a.currency == "") || (b.currency == "") {
		panic("subtracting null money value")
	}
	if a.currency != b.currency {
		panic("subtracting values in different currencies")
	}
	return Money {a.v.Sub(b.v), a.currency}
}

func (a Money) Mult(v decimal.Decimal) Money {
	if a.currency == "" {
		panic("multiplying null money value")
	}
	return Money {a.v.Mult(v), a.currency}
}

func (a Money) Div(b Money, precision uint32) decimal.Decimal {
	if (a.currency == "") || (b.currency == "") {
		panic("dividing null money value")
	}
	if a.currency != b.currency {
		panic("dividing values in different currencies")
	}
	return a.v.Div(b.v, precision)
}

func (a Money) DivPrice(b Money) Price {
	if (a.currency == "") || (b.currency == "") {
		panic("dividing null money value")
	}
	if a.currency == b.currency {
		panic("dividing values in the same currency")
	}
	return Price {a, b}
}

func (a Money) MultPrice(p Price, precision uint32) Money {
	if (a.currency == "") || p.IsNull() {
		panic("multiplying null money or price value")
	}
	if a.currency != p.am2.currency {
		panic("multiplying values in different currencies")
	}
	return Money {a.v.Mult(p.am1.v).Div(p.am2.v, precision), p.am1.currency}
}

func (a Money) Round(precision uint32) Money {
	if a.currency == "" {
		panic("rounding null money value")
	}
	return Money {a.v.Div(decimal.Value(1), precision), a.currency}
}

func (a Money) IsNull() bool {
	return a.currency == ""
}

func (a Money) String() string {
	if a.currency == "" {
		panic("printing null money value")
	}
	return a.v.StringPrecision(2)+" "+a.currency
}
