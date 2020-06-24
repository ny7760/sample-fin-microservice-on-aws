import json
import pprint
import requests
import time


URL = "http://localhost:8081/trade/order"

payload = {
    "orderDate": 20200607,
    "orderTime": "11:25:00",
    "stockCode": "A002",
    "tradeType": "BUY",
    "orderType": "1",  # 1は指値
    "orderPrice": 2320,
    "orderQuantity": 10
}

# while True:
#     r = requests.post(URL, data=json.dumps(payload))
#     pprint.pprint(r.text)
#     time.sleep(5)

r = requests.post(URL, data=json.dumps(payload))
pprint.pprint(r.text)
