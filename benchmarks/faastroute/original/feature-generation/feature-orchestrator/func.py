from rpc import RPC
import base64
import requests
import redis
import json
from functools import partial
from multiprocessing.dummy import Pool as ThreadPool
import os
import random

#redisClient = redis.Redis(host=os.environ['REDIS_URL'], password=os.environ['REDIS_PASSWORD'])

def invoke_lambda(bucket, dest, key, context):
    RPC(context, os.environ["FEATURE_EXTRACTOR"], [str({"input_bucket": bucket, "key": key, "dest": dest})])

def function_handler(context):
    if context["request_type"] != "GRPC":
        params = context["request"]
        bucket = params['bucket']
        dest = str(random.randint(0, 10000000)) + "-" + bucket
        all_keys = []

        for key in ["reviews100mb.csv", "reviews10mb.csv", "reviews20mb.csv", "reviews50mb.csv"]:
            all_keys.append(key)
        print("Number of File : " + str(len(all_keys)))
        print("File : " + str(all_keys))

        pool = ThreadPool(len(all_keys))
        pool.map(partial(invoke_lambda, bucket, dest, context), all_keys)
        pool.close()
        pool.join()

        return RPC(context, os.environ["FEATURE_WAIT"], [str({"num_of_file": str(len(all_keys)), "input_bucket": dest})])[0], 200
    else:
        print("Empty request", flush=True)
        return "{}", 200