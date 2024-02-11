# from parliament import Context
import ssl
import mmap
# Disable SSL certificate verification
ssl._create_default_https_context = ssl._create_unverified_context
import time
# from flask import Request
import json
from torchvision import transforms
from PIL import Image
import torch
import torchvision.models as models

import sys
import os

import io


from concurrent import futures

from flask import Flask,request
app = Flask(__name__)

# Load model
model = models.squeezenet1_1(pretrained=True)

labels_fd = open('imagenet_labels.txt', 'r')
labels = []
for i in labels_fd:
    labels.append(i)
labels_fd.close()

def preprocessImage(imageBytes):
    img = Image.open(io.BytesIO(imageBytes))
    transform = transforms.Compose([
        transforms.Resize(256),
        transforms.CenterCrop(224),
        transforms.ToTensor(),
        transforms.Normalize(
            mean=[0.485, 0.456, 0.406],
            std=[0.229, 0.224, 0.225]
        )
    ])

    img_t = transform(img)
    return torch.unsqueeze(img_t, 0)


def infer(batch_t):
    # Set up model to do evaluation
    model.eval()

    # Run inference
    with torch.no_grad():
        out = model(batch_t)

    # Print top 5 for logging
    _, indices = torch.sort(out, descending=True)
    percentages = torch.nn.functional.softmax(out, dim=1)[0] * 100

    # make comma-seperated output of top 100 label
    out = ""
    for idx in indices[0][:100]:
        out = out + labels[idx] + ","
    return out

def Recognise(frame):
    classification = infer(preprocessImage(frame))
    return classification

@app.route('/recog', methods=['POST'])
def trigger():
    file_name= request.form.get("file_number")
    filename = f"../pv/{file_name}"
    with open(filename, "rb") as file:
        mm = mmap.mmap(file.fileno(), 0, access=mmap.ACCESS_READ)
        data= mm.read()
    ret = Recognise(data).encode()
    filename = f"../pv/recog__{file_name}"
    with open(filename, 'w+b') as f:
        f.write(ret)
    with open(filename, 'r+b') as f:
        mm = mmap.mmap(f.fileno(), 0)
        mm.write(ret)
    return filename
    

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5004)


