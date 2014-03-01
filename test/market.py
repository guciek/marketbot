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
	store["plns"] = 500000000
	store["cashout_price"] = 250000000
	store["btcs"] = (store["plns"] * 100000000000
		) // (store["cashout_price"] * 996)
	store["orders_sell"] = dict()
	store["orders_buy"] = dict()
	store["market_price"] = None
	store["market_price_ts"] = None
	def info():
		o_plns = store["plns"]
		for i in store["orders_buy"]:
			o_plns = o_plns + store["orders_buy"][i]
		o_btcs = store["btcs"]
		for i in store["orders_sell"]:
			o_btcs = o_btcs + store["orders_sell"][i]
		sys.stderr.write(
			"[Market] "+
			("%0.2f mBTC, %0.2f PLN, "+
			"cash out %0.2f PLN at %0.2f PLN/BTC\n") %
			(
				o_btcs*0.00001,
				o_plns*0.00001,
				0.00001 * (o_plns + (o_btcs * store["cashout_price"] * 996)
					// 100000000000),
				store["cashout_price"]*0.00001
			)
		)
	def transactionSell(price, marketpr = False):
		btcs = store["orders_sell"][price]
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
		btcs += randint(-10, 10)
		price += randint(-10, 10)
		if btcs < 1: return False
		if price < 1: return False
		if store["btcs"] < btcs: return False
		while price in store["orders_sell"]:
			price = price + 1
		if btcs * price < 500000000000000: return False
		store["btcs"] -= btcs
		store["orders_sell"][price] = btcs
		runTransactions(True)
		return True
	def buy(price, plns):
		plns += randint(-10, 10)
		price += randint(-10, 10)
		if plns < 1: return False
		if price < 1: return False
		if store["plns"] < plns: return False
		while price in store["orders_buy"]:
			price = price - 1
		if plns < 5000000: return False
		store["plns"] -= plns
		store["orders_buy"][price] = plns
		runTransactions(True)
		return True
	def orders():
		ret = []
		for i in store["orders_buy"]:
			descr = "#b%d buy %.8f BTC for %.5f PLN" % (
				i, (store["orders_buy"][i] * 1.0) / i,
				store["orders_buy"][i]*0.00001
			)
			ret.append(descr)
		for i in store["orders_sell"]:
			descr = "#s%d buy %.5f PLN for %.8f BTC" % (
				i, (store["orders_sell"][i] * i)*0.0000000000001,
				store["orders_sell"][i]*0.00000001
			)
			ret.append(descr)
		return ret
	def cancel(descr):
		price = int(descr[2:])
		if (descr[1] == 's') and (price in store["orders_sell"]):
			store["btcs"] += store["orders_sell"][price]
			del store["orders_sell"][price]
			return True
		if (descr[1] == 'b') and (price in store["orders_buy"]):
			store["plns"] += store["orders_buy"][price]
			del store["orders_buy"][price]
			return True
		return False
	def balance():
		o_plns = store["plns"]
		for i in store["orders_buy"]:
			o_plns = o_plns + store["orders_buy"][i]
		o_btcs = store["btcs"]
		for i in store["orders_sell"]:
			o_btcs = o_btcs + store["orders_sell"][i]
		return o_btcs, o_plns
	return onPrice, sell, buy, orders, cancel, balance, info
onPriceChange, apiSell, apiBuy, apiGetOrders, \
	apiCancel, apiTotalBalance, printInfo = fakeMarket()

def passTime(fn, store = dict()):
	store["f"] = open(fn, "rb")
	l = str(store["f"].readline(), "ascii").split(" ")
	store["t"] = int(l[0])-1
	store["next_ts"] = int(l[0])
	store["next_price"] = int(l[1])
	def nextPrice(step):
		store["t"] = store["t"]+step
		while store["t"] >= store["next_ts"]:
			onPriceChange(store["next_ts"], store["next_price"])
			l = str(store["f"].readline(), "ascii")
			if l == "":
				store["f"].close()
				return False
			l = l.split(" ")
			store["next_ts"] = int(l[0])
			store["next_price"] = int(l[1])
		return True
	def getTime():
		return store["t"]
	nextPrice(1)
	return nextPrice, getTime
getTime = None

def cmdLine(line):
	if line == "": return False
	ret = True
	line = line.split(" ")
	if line[0] == "time":
		if randint(1, 100) <= 4:
			print("error")
		else:
			print("time "+str(getTime()))
	elif line[0] == "wait":
		ret = passTime(60)
		print("ok wait" if ret else "exit")
	elif (line[0] == "buy") and (line[3] == "for"):
		if randint(1, 100) <= 7:
			print("error")
		else:
			am1 = float(line[1])*100000.0
			am2 = float(line[4])*100000.0
			if (line[2] == "PLN") and (line[5] == "BTC"):
				pr = int(am1*100000.0/am2)
				if apiSell(int(pr), int(am2*1000.0)) and (randint(1, 100) <= 96):
					print("ok buy")
				else:
					print("error")
			elif (line[2] == "BTC") and (line[5] == "PLN"):
				pr = int(am2*100000.0/am1)
				if apiBuy(int(pr), int(am2)) and (randint(1, 100) <= 96):
					print("ok buy")
				else:
					print("error")
			else:
				sys.stdout.write("[Market] Rejecting an order"+
					" with unsupported currencies")
				print("error")
	elif line[0] == "cancel":
		if randint(1, 100) <= 5:
			print("error")
		else:
			if apiCancel(line[1]) and (randint(1, 10) <= 8):
				print("ok cancel")
			else:
				print("error")
	elif line[0] == "orders":
		if randint(1, 100) <= 3:
			print("error")
		else:
			os = apiGetOrders()
			if os != None:
				print("orders:")
				for o in os: print(o)
				print(".")
			else:
				print("error")
	elif line[0] == "totalbalance":
		if randint(1, 100) <= 2:
			print("error")
		else:
			bbtc, bpln = apiTotalBalance()
			print("totalbalance:")
			print("%.5f PLN"%(bpln*0.00001))
			print("%.8f BTC"%(bbtc*0.00000001))
			print(".")
	elif line[0] == "echo":
		print("echo "+line[1])
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
		printInfo()
		passTime, getTime = passTime(sys.argv[1])
		try:
			for line in sys.stdin:
				if not cmdLine(line.strip()):
					break
				passTime(5)
		except KeyboardInterrupt:
			sys.stderr.write("\n[Market] Interrupted\n")
		except Exception as e:
			sys.stderr.write("\n[Market] Error: "+str(e)+"\n")
		printInfo()
except Exception as e:
	sys.stderr.write("\n[Market] Error: "+str(e)+"\n")
	sys.exit(1)
