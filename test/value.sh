#!/bin/bash
# Copyright by Karol Guciek (http://guciek.github.io)
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 2 or 3.

SELLP=$1
if [ -z "$SELLP" ]; then
	SELLP=2500
fi

VALUE=0
VALUE_S=0

O=$(cat orders.txt) || exit 1
while read a b c; do
	case "$a" in
		available)
			let VALUE=$VALUE+$b
			let VALUE_S=$VALUE_S+$c
		;;
		buy)
			if [ "$b" -ge "${SELLP}00000" ]; then
				let VALUE_S=$VALUE_S+$c*100000000/$b
			else
				let VALUE=$VALUE+$c
			fi
		;;
		sell)
			if [ "$b" -le "${SELLP}00000" ]; then
				let VALUE=$VALUE+$b/1000*$c/100000
			else
				let VALUE_S=$VALUE_S+$c
			fi
		;;
	esac
done < <(echo "$O")

let TOTAL=$VALUE+$VALUE_S*$SELLP/1000

P=$(printf "%d.%03d" $(($SELLP / 1000)) $(($SELLP - ($SELLP / 1000)*1000)))
let T=$TOTAL/1000
T=$(printf "%d.%02d" $(($T / 100)) $(($T - ($T / 100)*100)))
let V=$VALUE/1000
V=$(printf "%d.%02d" $(($V / 100)) $(($V - ($V / 100)*100)))
let VS=$VALUE_S/1000
VS=$(printf "%d.%02d" $(($VS / 100)) $(($VS - ($VS / 100)*100)))
echo "$V + sell $VS at $P = $T"
