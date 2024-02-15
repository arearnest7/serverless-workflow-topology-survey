from flask import Flask
import os
import random
import shutil
import time
import requests
import mmap
import glob
app = Flask(__name__)

@app.route('/')
def trigger():
        number = random.randint(0, 100000)
        with open("reference/video.mp4", "rb") as videoFile:
            videoFragment = videoFile.read()
        with open(f"../pv/video_data_{number}", 'w+b') as f:
            f.write(videoFragment)
        with open(f"../pv/video_data_{number}", 'r+b') as f:
            mm = mmap.mmap(f.fileno(), 0)
            mm.write(videoFragment)
        url = f"http://{os.environ.get('DECODER_IP')}:{os.environ.get('DECODER_PORT')}/decode"
        response = requests.post(url, data={f"file_number": number})
        filename = response.content.decode('utf-8')
        with open(filename, "rb") as f:
            mm = mmap.mmap(f.fileno(), 0, access=mmap.ACCESS_READ)
            file_content = mm.read().decode()
        mm.close()
        files_to_delete = glob.glob("../pv/*")
        for file in files_to_delete:
             os.remove(file)
        return file_content

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5002)

        
