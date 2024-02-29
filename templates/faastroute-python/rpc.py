from concurrent import futures
import os
import random
import redis
import datetime
import grpc
import mmap
import app_pb2 as pb
import app_pb2_grpc as pb_grpc

MAX_MESSAGE_LENGTH = 1024 * 1024 * 200
opts = [("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),("grpc.max_send_message_length", MAX_MESSAGE_LENGTH)]

def RPC(context, dest, payloads):
    if "LOGGING_NAME" in os.environ:
        redisClient = redis.Redis(host=os.environ['LOGGING_URL'], password=os.environ['LOGGING_PASSWORD'])
    with grpc.insecure_channel(dest, options=opts,) as channel:
        stub = pb_grpc.gRPCFunctionStub(channel)
        tl = []
        pv_paths = []
        request_type = [["gg", "gm"],["mg","mm"]]["RPC_PV" in os.environ]["RPC_DEST_PV" in os.environ]
        print(request_type)
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + context["WorkflowId"] + "," + context["Depth"] + "," + context["Width"] + "," + context["RequestType"] + "," + "10" + "\n")
        with futures.ThreadPoolExecutor(max_workers=len(payloads)) as executor:
            for i in range(len(payloads)):
                payload = payloads[i]
                if request_type == "gg" or request_type == "gm":
                    tl.append(executor.submit(stub.gRPCFunctionHandler, pb.RequestBody(data=payload, workflow_id=context["WorkflowId"], depth=(context["Depth"] + 1), width=i, request_type=request_type)))
                else:
                    pv_path = context["WorkflowId"] + "_" + str(context["Depth"]) + "_" + str(context["Width"]) + "_" + str(random.randint(0, 10000000))
                    pv_paths.append(pv_path)
                    with open(os.environ["RPC_PV"] + "/" + pv_path, "wb") as f:
                        f.write(payload)
                    tl.append(executor.submit(stub.gRPCFunctionHandler, pb.RequestBody(workflow_id=context["WorkflowId"], depth=(context["Depth"] + 1), width=i, request_type=request_type, pv_path=pv_path)))
        print(len(tl))
        results = []
        if "RPC_DEST_PV" not in os.environ:
            results = [t.result().reply for t in tl]
        else:
            for t in tl:
                print(t.result().code)
                print(t.result().reply)
                print(t.result().pv_path)
                reply = b''
                with open(os.environ["RPC_DEST_PV"] + "/" + t.result().pv_path, "rb") as f:
                    mm = mmap.mmap(f.fileno(), 0, access=mmap.ACCESS_READ)
                    reply = mm.read()
                    mm.close()
                print(reply)
                os.remove(os.environ["RPC_DEST_PV"] + "/" + t.result().pv_path)
                results.append(reply.decode())
        if "RPC_PV" in os.environ:
            for pv_path in pv_paths:
                os.remove(os.environ["RPC_PV"] + "/" + pv_path)
        if "LOGGING_NAME" in os.environ:
            redisClient.append(os.environ["LOGGING_NAME"], str(datetime.datetime.now()) + "," + context["WorkflowId"] + "," + context["Depth"] + "," + context["Width"] + "," + context["RequestType"] + "," + "11" + "\n")
        return results
