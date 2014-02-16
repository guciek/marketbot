// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

type OrderType uint8

const (
	BUY = OrderType(0)
	SELL = OrderType(1)
)

func (t OrderType) String() string {
	if uint8(t) > 0 { return "sell" }
	return "buy"
}
