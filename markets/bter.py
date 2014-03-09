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

class Bter:
	def __init__(self, api_key, api_secret):
		self.__api_key = str(api_key).strip()
		self.__api_secret = bytes(str(api_secret).strip(), "ascii")
		self.__api_url = "https://bter.com/api/1/private/"

	def __q(self, fun, args = {}):
		try:
			args['nonce'] = int(time.time())
			post_data = bytes(parse.urlencode(args), "utf-8")
			sign = hmac.new(self.__api_secret, post_data, hashlib.sha512).hexdigest()
			headers = {
				'KEY' : self.__api_key,
				'SIGN': sign,
				'User-Agent': "Mozilla/4.0 (compatible; Bter bot)"
			}
			url = self.__api_url + fun
			req = request.Request(url, post_data, headers)
			response = request.urlopen(req, None, 5)
			d = str(response.read(), "utf-8")
			try:
				ret = dict(json.loads(d))
			except Exception as e:
				raise Exception("could not parse JSON")
			if (ret["result"] is True) or (str(ret["result"]).lower() == "true"):
				return ret
			raise Exception("result is "+str(ret["result"]))
		except Exception as e:
			raise Exception("API request failed: "+str(e))

	def getFunds(self):
		return self.__q("getfunds")

	def getOrders(self):
		return self.__q("orderlist")

	def placeOrder(self, pair, tpe, amount, price):
		args = {"pair": str(pair), "type": str(tpe),
			"rate": ("%.15f"%float(price)), "amount": str(amount)}
		return self.__q("placeorder", args)

	def cancelOrder(self, oid):
		return self.__q("cancelorder", {"order_id": int(oid)})

def cmdLine(market, line):
	if line == "": return False
	line = line.split(" ")

	trading_pairs = set([
		"ltc_btc","bqc_btc","btb_btc","buk_btc","cdc_btc",
		"cmc_btc","cnc_btc","dgc_btc","doge_btc","dtc_btc",
		"exc_btc","frc_btc","ftc_btc","max_btc","mec_btc",
		"mint_btc","mmc_btc","nec_btc","nmc_btc","nxt_btc",
		"ppc_btc","pts_btc","qrk_btc","src_btc","tag_btc",
		"yac_btc","vtc_btc","wdc_btc","xpm_btc","zcc_btc",
		"zet_btc","bqc_ltc","cent_ltc","cnc_ltc","dvc_ltc",
		"ftc_ltc","frc_ltc","ifc_ltc","net_ltc","nmc_ltc",
		"ppc_ltc","red_ltc","tips_ltc","tix_ltc","wdc_ltc",
		"yac_ltc"
	])

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
		am1 = line[1]
		cur1 = line[2].lower()
		am2 = line[4]
		cur2 = line[5].lower()
		tpe = "BUY"
		if cur2+"_"+cur1 in trading_pairs:
			am1, am2 = am2, am1
			cur1, cur2 = cur2, cur1
			tpe = "SELL"
		if cur1+"_"+cur2 in trading_pairs:
			pr = float(am2)/float(am1)
			market.placeOrder(cur1+"_"+cur2, tpe, am1, pr)
			print("ok buy")
		else:
			print("error unsupported trading pair '%s'"%(cur1+"_"+cur2))
		return True

	if line[0] == "cancel":
		market.cancelOrder(str(line[1]))
		print("ok cancel")
		return True

	if line[0] == "orders":
		r = market.getOrders()
		print("orders:")
		for o in r["orders"]:
			oid = str(o["id"])
			print ("%s buy %s %s for %s %s" % (
				oid,
				str(o["buy_amount"]), str(o["buy_type"]),
				str(o["sell_amount"]), str(o["sell_type"])
			))
		print(".")
		return True

	elif line[0] == "totalbalance":
		r = market.getFunds()
		print("totalbalance:")
		if "available_funds" in r:
			for k in r["available_funds"]:
				print(str(r["available_funds"][k])+" "+str(k))
		if "locked_funds" in r:
			for k in r["locked_funds"]:
				print(str(r["locked_funds"][k])+" "+str(k))
		print(".")
		return True

	if line[0] == "exit":
		print("exit")
		return False

	print("error Unknown command")
	return True

def run():
	b = Bter(sys.argv[1], sys.argv[2])
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
except Exception as e:
	sys.stderr.write("Error: "+str(e)+"\n")
