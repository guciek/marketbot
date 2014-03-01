// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"money"
)

type OrderPlanner struct {
	TargetOrders func([]money.Money) []Order
}

type TextInterfaceController struct {
	Writeln func(string) error
	Readln func() (string, error)
	Exit func() error
}

type MarketController struct {
	GetTime func() (Time, error)
	GetTotalBalance func() (map[string]money.Money, error)
	GetOrders func() ([]Order, error)
	NewOrder func(Order) error
	CancelOrder func(string) error
	CheckConnection func() error
	Wait func() error
	Exit func() error
}
