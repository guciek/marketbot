// Copyright by Karol Guciek (http://guciek.github.io)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 2 or 3.

package main

func CachedDataSaver(ds DataSaver) DataSaver {
	cache := make(map[string]string)
	return DataSaver {
		Write: func(f string, d string) error {
			if v, found := cache[f]; found {
				if v == d { return nil }
			}
			cache[f] = d
			return ds.Write(f, d)
		},
		Read: func(f string) (string, error) {
			if v, found := cache[f]; found { return v, nil }
			v, err := ds.Read(f)
			if err != nil { return "", err }
			cache[f] = v
			return v, nil
		},
	}
}
