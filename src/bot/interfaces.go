// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

type OrderPlanner struct {
	TargetOrders func(AssetValue, MoneyValue) []Order
}

type TextInterfaceController struct {
	Writeln func(string)
	Readln func() string
	Exit func() error
}

type MarketController struct {
	GetTime func() (Time, error)
	GetTotalBalance func() (AssetValue, MoneyValue, error)
	GetOrders func() ([]Order, error)
	NewOrder func(o Order) error
	CancelOrder func(o Order) error
	Wait func() error
	Exit func() error
}
