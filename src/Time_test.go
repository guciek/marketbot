// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"fmt"
)

func Example_Time() {
	var t Time
	fmt.Println(t)
	fmt.Println(Time(1392390089))
	// Output:
	// 1970-01-01 00:00:00
	// 2014-02-14 15:01:29
}
