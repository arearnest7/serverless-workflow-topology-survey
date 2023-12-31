from parliament import Context
from flask import Request
import base64
import requests
import redis
import json
from functools import partial
from multiprocessing.dummy import Pool as ThreadPool
import os
import random

redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def invoke_lambda(bucket, dest, key):
    requests.post(url=os.environ["FEATURE_EXTRACTOR"], json={"input_bucket": bucket, "key": key.decode(), "dest": dest})

def main(context: Context):
    if 'request' in context.keys():
        params = context.request.json
        bucket = params['bucket']
        dest = str(random.randint(0, 10000000)) + "-" + bucket
        all_keys = []

        for key in redisClient.scan_iter(bucket + "-*"):
            all_keys.append(key)
        print("Number of File : " + str(len(all_keys)))
        print("File : " + str(all_keys))

        pool = ThreadPool(len(all_keys))
        pool.map(partial(invoke_lambda, bucket, dest), all_keys)
        pool.close()
        pool.join()

        return requests.post(url=os.environ["FEATURE_WAIT"], json={"num_of_file": str(len(all_keys)), "input_bucket": dest}).text, 200
    else:
        print("Empty request", flush=True)
        return "{}", 200
