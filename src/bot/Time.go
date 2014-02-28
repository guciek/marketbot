// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

import (
	"time"
)

type Time int64

func (t Time) String() string {
	return time.Unix(int64(t), 0).UTC().Format("2006-01-02 15:04:05")
}
