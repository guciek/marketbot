#!/usr/bin/python3
# Copyright by Karol Guciek (http://guciek.github.io)
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, version 2 or 3.

import time
import sys
import json
from urllib import request
from random import randint

def fakeMarket(store = dict()):
	store["plns"] = 200000000
	store["cashout_price"] = 250000000
	store["btcs"] = (store["plns"] *
			100000000000) // (store["cashout_price"] * 996)
	store["orders_sell"] = dict()
	store["orders_buy"] = dict()
	store["market_price"] = None
	store["market_price_ts"] = None
	def info():
		o_plns = 0
		for i in store["orders_buy"]:
			o_plns = o_plns + store["orders_buy"][i]
		o_btcs = 0
		for i in store["orders_sell"]:
			o_btcs = o_btcs + store["orders_sell"][i]
		sys.stderr.write(
			"[Market] "+
			("%0.2f + %0.2f mBTC, %0.2f + %0.2f PLN, "+
				"cash out %0.2f PLN at %0.2f PLN/BTC\n") %
			(
				store["btcs"]*0.00001, o_btcs*0.00001,
				store["plns"]*0.00001, o_plns*0.00001,
				0.00001*
					(store["plns"] + o_plns +
					((store["btcs"]+o_btcs) * store["cashout_price"] * 996) // 100000000000),
				store["cashout_price"]*0.00001
			)
		)
	def transactionSell(price, marketpr = False):
		btcs = store["orders_sell"][price]
		if (btcs > 1000000) and (randint(1, 10) <= 2) and (not marketpr):
			btcs = btcs*randint(50, 98)/100
			store["orders_sell"][price] -= btcs
		else:
			del store["orders_sell"][price]
		plns = (btcs * (store["market_price"] if marketpr
			else price) * 996) // 100000000000
		store["plns"] += plns
		sys.stderr.write(("[Market] Sold %.2f mBTC for %.2f PLN "+
			"(order price %.2f PLN/BTC)\n")
			% (btcs*0.00001, plns*0.00001, price*0.00001))
		return True
	def transactionBuy(price, marketpr = False):
		plns = store["orders_buy"][price]
		if (plns > 4000000) and (randint(1, 10) <= 2) and (not marketpr):
			plns = plns*randint(50, 98)/100
			store["orders_buy"][price] -= plns
		else:
			del store["orders_buy"][price]
		btcs = (plns * 99600000000) // (1000 * (store["market_price"]
			if marketpr else price))
		store["btcs"] += btcs
		sys.stderr.write(("[Market] Bought %.2f mBTC for %.2f PLN "+
			"(order price %.2f PLN/BTC)\n")
			% (btcs*0.00001, plns*0.00001, price*0.00001))
		return True
	def runTransactions(marketpr = False):
		if not store["market_price"]: return
		for price in sorted(store["orders_sell"].keys()):
			if store["market_price"] < price: break
			transactionSell(price, marketpr)
		for price in sorted(store["orders_buy"].keys(), reverse=True):
			if store["market_price"] > price: break
			transactionBuy(price, marketpr)
	def onPrice(ts, price):
		store["market_price"] = price
		store["market_price_ts"] = ts
		runTransactions()
	def sell(price, btcs):
		btcs += randint(-400, 400)
		price += randint(-400, 400)
		if btcs < 1: return False
		if price < 10000000: return False
		if store["btcs"] < btcs: return False
		while price in store["orders_sell"]:
			price = price + 1
		if btcs * price < 500000000000000: return False
		store["btcs"] -= btcs
		store["orders_sell"][price] = btcs
		runTransactions(True)
		return True
	def buy(price, plns):
		plns += randint(-400, 400)
		price += randint(-400, 400)
		if plns < 1: return False
		if price < 1000000: return False
		if store["plns"] < plns: return False
		while price in store["orders_buy"]:
			price = price + 1
		if plns < 5000000: return False
		store["plns"] -= plns
		store["orders_buy"][price] = plns
		runTransactions(True)
		return True
	def orders():
		ret = ["available %d %d" % (store["plns"], store["btcs"])]
		for i in store["orders_buy"]:
			ret.append("buy %d %d" % (i, store["orders_buy"][i]))
		for i in store["orders_sell"]:
			ret.append("sell %d %d" % (i, store["orders_sell"][i]))
		return ret
	def cancelSell(price, btcs):
		if price in store["orders_sell"]:
			if btcs != store["orders_sell"][price]:
				return False
			del store["orders_sell"][price]
			store["btcs"] += btcs
			return True
		return False
	def cancelBuy(price, plns):
		if price in store["orders_buy"]:
			if plns != store["orders_buy"][price]:
				return False
			del store["orders_buy"][price]
			store["plns"] += plns
			return True
		return False
	def getPrice():
		return store["market_price"]
	return onPrice, sell, buy, orders, cancelBuy, \
		cancelSell, getPrice, info
onPriceChange, apiSell, apiBuy, apiGetOrders, \
	apiCancelBuy, apiCancelSell, apiGetPrice, printInfo = fakeMarket()

def passTime(fn, store = dict()):
	store["step"] = 600
	store["f"] = open(fn, "rb")
	l = str(store["f"].readline(), "ascii").split(" ")
	store["t"] = int(l[0])-store["step"]
	store["next_ts"] = int(l[0])
	store["next_price"] = int(l[1])
	def nextPrice():
		store["t"] = store["t"]+store["step"]
		while store["t"] >= store["next_ts"]:
			l = str(store["f"].readline(), "ascii")
			if l == "":
				store["f"].close()
				return False
			l = l.split(" ")
			store["next_ts"] = int(l[0])
			store["next_price"] = int(l[1])
		onPriceChange(store["t"], store["next_price"])
		return True
	def getTime():
		return store["t"]
	nextPrice()
	return nextPrice, getTime
getTime = None

def cmdLine(line):
	if line == "": return False
	ret = True
	line = line.split(" ")
	if line[0] == "price":
		if randint(1, 10) <= 2:
			print("error")
		else:
			p = apiGetPrice()
			print("price "+str(p)+" "+str(getTime()) if p else "error")
	elif line[0] == "wait":
		ret = passTime()
		print("ok wait" if ret else "exit")
	elif line[0] == "sell":
		if randint(1, 10) <= 2:
			print("error")
		else:
			if apiSell(int(line[1]), int(line[2])) and (randint(1, 10) <= 9):
				print("ok sell %d %d" % (int(line[1]), int(line[2])))
			else:
				print("error")
	elif line[0] == "buy":
		if randint(1, 10) <= 2:
			print("error")
		else:
			if apiBuy(int(line[1]), int(line[2])) and (randint(1, 10) <= 9):
				print("ok buy %d %d" % (int(line[1]), int(line[2])))
			else:
				print("error")
	elif line[0] == "cancel":
		if randint(1, 10) <= 3:
			print("error")
		else:
			if line[1] == "sell":
				if apiCancelSell(int(line[2]), int(line[3])
						) and (randint(1, 10) <= 8):
					print("ok cancel sell %d %d" % (int(line[2]), int(line[3])))
				else:
					print("error")
			elif line[1] == "buy":
				if apiCancelBuy(int(line[2]), int(line[3])
						) and (randint(1, 10) <= 8):
					print("ok cancel buy %d %d" % (int(line[2]), int(line[3])))
				else:
					print("error")
	elif line[0] == "orders":
		if randint(1, 100) <= 1:
			print("error")
		else:
			os = apiGetOrders()
			if os != None:
				print("orders:")
				for o in os: print(o)
				print(".")
			else:
				print("error")
	elif line[0] == "exit":
		print("exit")
		ret = False
	else:
		print("error Unknown command")
	sys.stdout.flush()
	return ret

try:
	if len(sys.argv) < 2:
		sys.stderr.write("\nUsage:\n\ttest-market.py <log-file>\n\n")
	else:
		passTime, getTime = passTime(sys.argv[1])
		printInfo()
		for line in sys.stdin:
			if not cmdLine(line.strip()):
				break
		printInfo()
except KeyboardInterrupt:
	sys.stderr.write("\n[Market] Interrupted\n")
	sys.exit(1)
except Exception as e:
	sys.stderr.write("\n[Market] "+str(e)+"\n")
	sys.exit(1)
	

