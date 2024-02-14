from rpc import RPC
import base64
import requests
import json
import random
import os

def function_handler(context):
    if context["request_type"] == "GRPC":
        event = json.loads(context["request"])

        results = ""
        if event["reviewType"] == "Product":
            results = RPC(os.environ["SENTIMENT_PRODUCT_SENTIMENT"], [context["request"]], context["workflow_id"])[0]
        elif event["reviewType"] == "Service":
            results = RPC(os.environ["SENTIMENT_SERVICE_SENTIMENT"], [context["request"]], context["workflow_id"])[0]
        else:
            results = RPC(os.environ["SENTIMENT_CFAIL"], [context["request"]], context["workflow_id"])[0]
        return results, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
