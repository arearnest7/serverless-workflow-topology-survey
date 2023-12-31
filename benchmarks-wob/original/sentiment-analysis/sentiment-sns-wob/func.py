from parliament import Context
from flask import Request
import base64
import requests
import json
import random
import os

def main(context: Context):
    if 'request' in context.keys():
        event = context.request.json
        #response = requests.get(url = 'http://' + OF_Gateway_IP + ':' + OF_Gateway_Port + '/function/sha>
        #    "Subject": 'Negative Review Received',
        #    "Message": 'Review (ID = %i) of %s (ID = %i) received with negative results from sentiment a>
        #    event['reviewType'], int(event['productID']), int(event['customerID']), event['feedback'])
        #})

        response = requests.get(url=os.environ["SENTIMENT_DB"], json={
            'sentiment': event['sentiment'],
            'reviewType': event['reviewType'],
            'reviewID': event['reviewID'],
            'customerID': event['customerID'],
            'productID': event['productID'],
            'feedback': event['feedback']
        })
        return response.text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
