#!/usr/bin/python3

import time
import sys
import math
import base64
import json
import hmac
import hashlib
import threading
from urllib import parse, request

class Bitcurex:
	def __init__(self, api_key, api_secret):
		self.__api_key = str(api_key).strip()
		self.__api_secret = base64.b64decode(bytes(str(api_secret).strip(), "utf-8"))
		self.__api_url = "https://pln.bitcurex.com/api/0/"

	def __q(self, fun, args = {}):
		try:
			args['nonce'] = "%f" % time.time()
			post_data = bytes(parse.urlencode(args), "utf-8")
			sign = hmac.new(self.__api_secret, post_data, hashlib.sha512).digest()
			headers = {
				'Content-Type': 'application/x-www-form-urlencoded',
				'Rest-Key' : self.__api_key,
				'Rest-Sign': base64.b64encode(sign)
			}
			url = self.__api_url + fun
			req = request.Request(url, post_data, headers)
			response = request.urlopen(req, None, 5)
			d = str(response.read(), "utf-8")
			try:
				ret = dict(json.loads(d))
			except Exception as e:
				raise Exception("could not parse JSON")
			if "error" in ret:
				raise Exception(str(ret["error"]))
			return ret
		except Exception as e:
			raise Exception("API request failed: "+str(e))

	def getOrders(self):
		r = self.__q("getOrders")
		return r

	def cancelOrder(self, oid, tp):
		r = self.__q("cancelOrder", {'oid': oid, 'type': tp})
		return r

	def sellBTC(self, amount_btc, price):
		price = float(price)
		amount_btc = float(amount_btc)
		if price < 1.0:
			return {"error": "Invalid sell price"}
		r = self.__q("sellBTC", {'amount': "%.12f"%amount_btc, 'price': "%.5f"%price})
		return r

	def buyBTC(self, amount_btc, price):
		price = float(price)
		amount_btc = float(amount_btc)
		if price < 1.0:
			return {"error": "Invalid buy price"}
		r = self.__q("buyBTC", {'amount': "%.12f"%amount_btc, 'price': "%.5f"%price})
		return r

def cmdLine(market, line):
	if line == "": return False
	line = line.split(" ")

	if line[0] == "time":
		print("time "+str(int(time.time())))
		return True

	if line[0] == "echo":
		print("echo "+line[1])
		return True

	if line[0] == "wait":
		try:
			time.sleep(10)
		except:
			pass
		print("ok wait")
		return True

	if (line[0] == "buy") and (line[3] == "for"):
		am1 = float(line[1])
		am2 = float(line[4])
		if (line[2] == "PLN") and (line[5] == "BTC"):
			market.sellBTC(am2, am1/am2)
			print("ok buy")
		elif (line[2] == "BTC") and (line[5] == "PLN"):
			market.buyBTC(am1, am2/am1)
			print("ok buy")
		else:
			print("error unsupported trading pair")
		return True

	if line[0] == "cancel":
		par = line[1].split("#")
		market.cancelOrder(int(par[0]), int(par[1]))
		print("ok cancel")
		return True

	if line[0] == "orders":
		r = market.getOrders()
		print("orders:")
		for o in r["orders"]:
			t = int(o["type"])
			pr = float(o["price"])
			am_btc = float(o["amount"])
			am_pln = am_btc*pr
			oid = str(o["oid"])
			if t == 1:
				print("%s#%d buy %.5f PLN for %.8f BTC" % (oid, t, am_pln, am_btc))
			elif t == 2:
				print("%s#%d buy %.8f BTC for %.5f PLN" % (oid, t, am_btc, am_pln))
			else:
				print("unknown order type")
		print(".")
		return True

	elif line[0] == "totalbalance":
		r = market.getOrders()
		sum_pln = float(r["plns"])
		sum_btc = float(r["btcs"])
		for o in r["orders"]:
			t = int(o["type"])
			pr = float(o["price"])
			am_btc = float(o["amount"])
			am_pln = am_btc*pr
			if t == 1:
				sum_btc += am_btc
			elif t == 2:
				sum_pln += am_pln
			else:
				raise Exception("unknown order type")
		print("totalbalance:")
		print("%.5f PLN" % sum_pln)
		print("%.8f BTC" % sum_btc)
		print(".")
		return True

	if line[0] == "exit":
		print("exit")
		return False

	print("error Unknown command")
	return True

def run():
	b = Bitcurex(sys.argv[1], sys.argv[2])
	for line in sys.stdin:
		try:
			if not cmdLine(b, line.strip()):
				break
			sys.stdout.flush()
			try:
				time.sleep(1)
			except:
				pass
		except IOError:
			raise
		except KeyboardInterrupt:
			raise
		except Exception as e:
			print("error", e)
			sys.stdout.flush()

try:
	run()
except IOError:
	pass
except KeyboardInterrupt:
	pass
except Exception as e:
	sys.stderr.write("Error: "+str(e)+"\n")
	sys.stderr.flush()
