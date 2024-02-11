from __future__ import print_function
import json
import os
import time
from flask import Flask,request
import requests
import mmap
app = Flask(__name__)


# import pickle
# import sys
# import os

import cv2
import tempfile
import base64

from concurrent.futures import ThreadPoolExecutor

def decode(bytes):
    temp = tempfile.NamedTemporaryFile(suffix=".mp4")
    temp.write(bytes)
    temp.seek(0)

    all_frames = []
    vidcap = cv2.VideoCapture(temp.name)
    for i in range(int(6)):
        success,image = vidcap.read()
        all_frames.append(cv2.imencode('.jpg', image)[1].tobytes())

    return all_frames

def Recognise(frame,frame_number,file_no):
    filename = f"../pv/frame_{file_no}_{frame_number}"
    with open(filename, 'w+b') as f:
        f.write(frame)
    with open(filename, 'r+b') as f:
        mm = mmap.mmap(f.fileno(), 0)
        mm.write(frame)
    url = f"http://{os.environ.get('RECOG_IP')}:{os.environ.get('RECOG_PORT')}/recog"
    response = requests.post(url, data={"file_number": f"frame_{file_no}_{frame_number}"})
    outputfile = response.content.decode('utf-8')
    with open(outputfile, 'rb') as file:
        mm = mmap.mmap(file.fileno(), 0, access=mmap.ACCESS_READ)
        data= mm.read()
    return data.decode()


def processFrames(videoBytes,file_no):
    frames = decode(videoBytes)
    frames = frames[0:6]
    ex = ThreadPoolExecutor(max_workers=6)
    all_result_futures = ex.map(Recognise, frames,range(len(frames)),[file_no]*len(frames))      
    results = ""
    for result in all_result_futures:
        results = results + result + ","
    filename = f"../pv/result__{file_no}"
    with open(filename, 'w+') as f:
        f.write(results)
    with open(filename, 'r+b') as f:
        mm = mmap.mmap(f.fileno(), 0)
        mm.write(results.encode())
    return filename


def Decode(request,file_no):
    videoBytes = request
    results = processFrames(videoBytes,file_no)
    return results


@app.route('/decode', methods=['POST'])
def results():
    file_no= request.form.get("file_number")
    filename = f"../pv/video_data_{file_no}"
    with open(filename, 'rb') as file:
        mm = mmap.mmap(file.fileno(), 0, access=mmap.ACCESS_READ)
        data= mm.read()
        ret = Decode(data,file_no)
    return ret


    

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5003)



