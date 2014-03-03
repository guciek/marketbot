// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package decimal

import (
	"bytes"
	"fmt"
	"math/big"
)

type Decimal struct {
	v *big.Int
	exp int
}

func (p Decimal) String() string {
	return p.StringPrecision(0)
}

func (p Decimal) StringPrecision(prec uint32) string {
	if p.v == nil { return "0" }
	s := p.v.String()
	var ret bytes.Buffer
	if p.exp >= 0 {
		ret.WriteString(s)
		for p := p.exp; p > 0; p-- {
			ret.WriteByte('0')
		}
		if prec > 0 {
			ret.WriteByte('.')
			for p := uint32(0); p < prec; p++ {
				ret.WriteByte('0')
			}
		}
	} else {
		if len(s) > -p.exp {
			ret.WriteString(s[0:len(s)+p.exp])
			ret.WriteByte('.')
			ret.WriteString(s[len(s)+p.exp:len(s)])
		} else {
			ret.WriteByte('0')
			ret.WriteByte('.')
			for p := len(s)+p.exp; p < 0; p++ {
				ret.WriteByte('0')
			}
			ret.WriteString(s)
		}
		for p := int64(prec) + int64(p.exp); p > 0; p-- {
			ret.WriteByte('0')
		}
	}
	return ret.String()
}

func normalizeModify(a *Decimal) {
	if a.v == nil {
		a.exp = 0
		return
	}
	div, mod, ten := new(big.Int), new(big.Int), big.NewInt(10)
	for {
		div.DivMod(a.v, ten, mod)
		if mod.Sign() != 0 { break }
		if div.Sign() == 0 {
			a.v = nil
			a.exp = 0
			return
		}
		a.v.Set(div)
		a.exp++
	}
}

func Value(v uint32) Decimal {
	if v < 1 { return Decimal {} }
	ret := Decimal {big.NewInt(int64(v)), 0}
	normalizeModify(&ret)
	return ret
}

func ParseDecimal(v string) (Decimal, error) {
	if len(v) < 1 {
		return Decimal {}, fmt.Errorf("empty string")
	}
	if v[0] == '.' {
		return Decimal {}, fmt.Errorf("number starts with a dot")
	}
	if v[len(v)-1] == '.' {
		return Decimal {}, fmt.Errorf("number ends with a dot")
	}
	seen_dot := false
	var decimal_places int = 0
	var ret Decimal
	for _, c := range v {
		if c == '.' {
			if seen_dot {
				return Decimal {}, fmt.Errorf("two dots in number")
			}
			seen_dot = true
		} else if (c >= '0') && (c <= '9') {
			if seen_dot {
				decimal_places++
			}
			if ret.v != nil {
				ret.exp++
			}
			ret = ret.Add(Value(uint32(c-'0')))
		} else {
			return Decimal {}, fmt.Errorf("invalid character")
		}
	}
	if ret.v == nil { return Decimal {}, nil }
	ret.exp -= decimal_places
	return ret, nil
}

func (lhs Decimal) Add(rhs Decimal) Decimal {
	if lhs.v == nil { return rhs }
	if rhs.v == nil { return lhs }
	var ret Decimal
	if lhs.exp == rhs.exp {
		ret.v = new(big.Int)
		ret.v.Add(lhs.v, rhs.v)
	} else {
		if lhs.exp < rhs.exp {
			lhs, rhs = rhs, lhs
		}
		ret.v = big.NewInt(10)
		ret.v.Exp(ret.v, big.NewInt(int64(lhs.exp-rhs.exp)), nil)
		ret.v.Mul(ret.v, lhs.v)
		ret.v.Add(ret.v, rhs.v)
	}
	ret.exp = rhs.exp
	normalizeModify(&ret)
	return ret
}

func (lhs Decimal) Sub(rhs Decimal) Decimal {
	if lhs.v == nil { return rhs }
	if rhs.v == nil { return lhs }
	var ret Decimal
	if lhs.exp == rhs.exp {
		if lhs.v.Cmp(rhs.v) <= 0 { return Decimal {} }
		ret.v = new(big.Int)
		ret.v.Sub(lhs.v, rhs.v)
	} else {
		negate := false
		if lhs.exp < rhs.exp {
			lhs, rhs = rhs, lhs
			negate = true
		}
		ret.v = big.NewInt(10)
		ret.v.Exp(ret.v, big.NewInt(int64(lhs.exp-rhs.exp)), nil)
		ret.v.Mul(ret.v, lhs.v)
		ret.v.Sub(ret.v, rhs.v)
		if negate {
			ret.v.Neg(ret.v)
		}
	}
	ret.exp = rhs.exp
	if ret.v.Sign() <= 0 { return Decimal {} }
	normalizeModify(&ret)
	return ret
}

func (lhs Decimal) Mult(rhs Decimal) Decimal {
	if lhs.v == nil { return Decimal {} }
	if rhs.v == nil { return Decimal {} }
	var ret Decimal
	ret.v = new(big.Int)
	ret.v.Mul(lhs.v, rhs.v)
	ret.exp = rhs.exp+lhs.exp
	normalizeModify(&ret)
	return ret
}

func (lhs Decimal) Div(rhs Decimal, precision uint32) Decimal {
	if lhs.v == nil { return Decimal {} }
	if rhs.v == nil { panic("division by 0") }
	if precision < 1 { panic("0 digits of precision") }
	var ret Decimal
	ret.exp = lhs.exp - rhs.exp - 1
	ret.v = new(big.Int)
	ret.v.Set(lhs.v)
	ten := big.NewInt(10)
	for ret.v.Cmp(rhs.v) < 0 {
		ret.v.Mul(ret.v, ten)
		ret.exp--
	}
	ret.v.Mul(ret.v, ten)
	precision_power := new(big.Int)
	precision_power.Exp(ten, big.NewInt(int64(precision)), nil)
	ret.v.Mul(ret.v, precision_power)
	ret.exp = int(int64(ret.exp)-int64(precision))
	{
		halfr := new(big.Int)
		halfr.Div(rhs.v, big.NewInt(2))
		ret.v.Add(ret.v, halfr)
	}
	ret.v.Div(ret.v, rhs.v)
	if ret.v.Cmp(precision_power) >= 0 {
		mod := new(big.Int)
		for ret.v.Cmp(precision_power) >= 0 {
			ret.v.DivMod(ret.v, ten, mod)
			ret.exp++
		}
		if mod.Cmp(big.NewInt(5)) >= 0 {
			ret.v.Add(ret.v, big.NewInt(1))
		}
	}
	normalizeModify(&ret)
	return ret
}

func (a Decimal) IsZero() bool {
	return a.v == nil
}

func (lhs Decimal) Equals(rhs Decimal) bool {
	if (rhs.v == nil) && (lhs.v == nil) { return true }
	if lhs.v == nil { return false }
	if rhs.v == nil { return false }
	if lhs.exp != rhs.exp { return false }
	return lhs.v.Cmp(rhs.v) == 0
}

func (lhs Decimal) Less(rhs Decimal) bool {
	if rhs.v == nil { return false }
	if lhs.v == nil { return true }
	l, r := lhs.v.String(), rhs.v.String()
	d := (len(l)+lhs.exp) - (len(r)+rhs.exp)
	if d > 0 { return false }
	if d < 0 { return true }
	for i := 0; (i < len(l)) && (i < len(r)); i++ {
		if l[i] != r[i] {
			return l[i] < r[i]
		}
	}
	return len(l) < len(r)
}
