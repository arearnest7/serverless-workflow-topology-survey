from parliament import Context
from flask import Request
import base64
import random
import datetime
import redis

if "LOGGING_NAME" in os.environ:
    loggingClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])

def main(context: Context):
    if 'request' in context.keys():
        if "LOGGING_NAME" in os.environ:
            loggingClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + "0" + "," + "0" + "," + "0" + "," + "kn" + "," + "0" + "\n")
        return "CategorizationFail: Fail: \"Input CSV could not be categorised into 'Product' or 'Service'.\"", 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
