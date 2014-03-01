// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package decimal

import (
	"fmt"
)

func Example_Decimal_Add() {
	var n1, n2 Decimal
	n2 = n1.Add(Value(500)).Add(n2)
	fmt.Println(n2.Add(Value(7)))
	fmt.Println(n2, n2.exp)
	fmt.Println(n1, n1.exp)
	n1 = n2.Add(Value(7)).Add(Value(2493))
	fmt.Println(n1, n1.exp)
	n1 = Value(11000);
	fmt.Println(n1, n1.exp)
	n1 = n1.Add(n2);
	fmt.Println(n1, n1.exp)
	// Output:
	// 507
	// 500 2
	// 0 0
	// 3000 3
	// 11000 3
	// 11500 2
}

func Example_ParseString() {
	x, err := ParseDecimal("3.14")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("0.00159")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("2060090.0")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("1030500.00000000")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("1030500.00000001")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("0")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("0000000.00000")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal(".00000001")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("100.")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("10a")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("1 ")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("1 2")
	fmt.Println(err == nil, x, x.exp)
	x, err = ParseDecimal("100.00.00")
	fmt.Println(err == nil, x, x.exp)
	// Output:
	// true 3.14 -2
	// true 0.00159 -5
	// true 2060090 1
	// true 1030500 2
	// true 1030500.00000001 -8
	// true 0 0
	// true 0 0
	// false 0 0
	// false 0 0
	// false 0 0
	// false 0 0
	// false 0 0
	// false 0 0
	// false 0 0
}

func Example_Decimal_Fractions() {
	x, _ := ParseDecimal("3.14")
	y, _ := ParseDecimal("0.00159")
	z, _ := ParseDecimal("0.00001")
	fmt.Println(x, "+", y, "=", x.Add(y).Add(z).Sub(z))
	fmt.Println(x, "-", y, "=", x.Sub(y))
	fmt.Println(x, "*", y, "=", x.Mult(y))
	fmt.Println(x, "/", y, "=", x.Div(y, 5), "=", x.Div(y, 7))
	fmt.Println(y, "-", x, "=", y.Sub(x))
	fmt.Println(x, "-", x, "=", x.Sub(x))
	fmt.Println(z, z.exp)
	z = z.Add(y)
	fmt.Println(z, z.exp)
	fmt.Println("45 / 99 =", Value(45).Div(Value(99), 1))
	fmt.Println("45 / 99 =", Value(45).Div(Value(99), 20))
	fmt.Println("45 / 99 =", Value(45).Div(Value(99), 21))
	fmt.Println("45 / 99 =", Value(45).Div(Value(99), 22))
	fmt.Println("45 / 99 =", Value(45).Div(Value(99), 23))
	x = Value(1).Div(Value(3), 6)
	fmt.Println("100 -", x, "=", Value(100).Sub(x))
	// Output:
	// 3.14 + 0.00159 = 3.14159
	// 3.14 - 0.00159 = 3.13841
	// 3.14 * 0.00159 = 0.0049926
	// 3.14 / 0.00159 = 1974.8 = 1974.843
	// 0.00159 - 3.14 = 0
	// 3.14 - 3.14 = 0
	// 0.00001 -5
	// 0.0016 -4
	// 45 / 99 = 0.5
	// 45 / 99 = 0.45454545454545454545
	// 45 / 99 = 0.454545454545454545455
	// 45 / 99 = 0.4545454545454545454545
	// 45 / 99 = 0.45454545454545454545455
	// 100 - 0.333333 = 99.666667
}

func Example_Decimal_Compare() {
	x, _ := ParseDecimal("31415.9")
	y, _ := ParseDecimal("31425.9")
	u, _ := ParseDecimal("31415.92")
	w, _ := ParseDecimal("3141.59")
	z, _ := ParseDecimal("21415.9")
	fmt.Println(x.Less(y), y.Less(x))
	fmt.Println(x.Less(u), u.Less(x))
	fmt.Println(x.Less(w), w.Less(x))
	fmt.Println(x.Less(z), z.Less(x))
	fmt.Println(x.Less(x))
	fmt.Println(x.Equals(y), x.Equals(w))
	fmt.Println(x.Equals(x), x.Equals(y.Sub(Value(10))))
	// Output:
	// true false
	// true false
	// false true
	// false true
	// false
	// false false
	// true true
}

func Example_Decimal_Big() {
	x, _ := ParseDecimal("0.7")
	y := Value(1)
	for i := 0; i < 100; i++ {
		y = y.Mult(x)
	}
	fmt.Println(y)
	// Output:
	// 0.0000000000000003234476509624757991344647769100216810857203198904625400933895331391691459636928060001
}
